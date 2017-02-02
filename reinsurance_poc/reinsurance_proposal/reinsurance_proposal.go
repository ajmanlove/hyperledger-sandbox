package main

import (
	"errors"
	"fmt"
	// "encoding/json"
	// "strconv"

	"sync/atomic"
	"time"

	"strings"

	"github.com/ajmanlove/hyperledger-sandbox/reinsurance_poc/common"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/core/util"
)

var logger = shim.NewLogger("ReinsuranceProposalCC")
var assetManagementCCId = ""
var counter uint64 = 0
var proposalPrefix = "BID"

type ReinsuranceProposalCC struct {
}

// TBD
type ReinsuranceProposal struct {
}

var amComm = common.AssetManagementCommunicator{}

func (t *ReinsuranceProposalCC) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	logger.Debug("Init()")

	if len(args) != 1 {
		return nil, errors.New("Init expects expects asset management cc id as arg")
	}
	assetManagementCCId = args[0]
	amComm.CCName = assetManagementCCId

	return nil, nil
}

func (t *ReinsuranceProposalCC) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	logger.Debug("enter Invoke")
	switch function {
	case common.RP_PROPOSE_ARG:
		return t.propose(stub, args)
	case common.RP_COUNTER_ARG:
		return nil, errors.New("counter not implemented")
	case common.RP_ACCEPT_ARG:
		return nil, errors.New("accept not implemented")
	case common.RP_REJECT_ARG:
		return nil, errors.New("reject not implemented")
	default:
		return nil, errors.New("Unrecognized Invoke function : " + function)
	}
}

func (t *ReinsuranceProposalCC) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	switch function {
	case common.RP_GET_BID_ARG:
		if len(args) != 1 {
			return nil, errors.New("get_proposal requires 1 arg ['proposalId']")
		}

		proposal, err := t.get_proposal(stub, args[0])
		if err != nil {
			logger.Error(err)
			return nil, err
		}
		return proposal.Encode()

	default:
		return nil, errors.New("Unrecognized Query function : " + function)
	}
}

func (t *ReinsuranceProposalCC) propose(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	logger.Debug("propose() args: " + strings.Join(args, ","))
	if len(args) != 2 {
		return nil, errors.New("Requires 2 args: ['requestId', 'contractText']")
	}

	requestId := args[0]
	contractText := args[1]
	now := get_unix_millisec()

	logger.Debug()

	logger.Debug("get enrollment id")
	enrollmentId, err := amComm.GetEnrollmentAttr(stub)
	if err != nil {
		return nil, err
	}

	logger.Debug("Asserting rights ...")
	err = amComm.AssertHasAssetRights(stub, requestId, []common.AssetRight{common.AVIEWER})
	if err != nil {
		return nil, err
	}

	logger.Debug("Creating record...")
	id := t.create_prop_id(requestId)
	var record common.ReinsuranceBid
	record.Init()

	record.Id = id
	record.RequestId = requestId
	record.Bidder = enrollmentId
	record.ContractText = contractText
	record.Created = now
	record.Updated = now
	record.UpdatedBy = enrollmentId
	record.Status = "bid" // TODO

	encoded, err := record.Encode()
	if err != nil {
		logger.Error(err)
		return nil, errors.New("Failed to encode ReinsuranceBid record")
	}

	err = stub.PutState(id, encoded)
	if err != nil {
		logger.Error(err)
		return nil, errors.New("Failed to put ReinsuranceBid record")
	}

	invokeArgs := util.ToChaincodeArgs(common.AM_NEW_BID_ARG, id, requestId, enrollmentId, fmt.Sprintf("%d", now))
	bytes, err := stub.InvokeChaincode(assetManagementCCId, invokeArgs)

	if err != nil {
		logger.Error(err)
		return nil, errors.New("Failed to manage new proposal asset " + id)
	}

	logger.Debugf("AM RESPONSE is %s", string(bytes))
	return nil, nil
}

func (t *ReinsuranceProposalCC) get_proposal(stub shim.ChaincodeStubInterface, propId string) (common.ReinsuranceBid, error) {
	// Rights
	var r common.ReinsuranceBid

	err := amComm.AssertHasAssetRights(stub, propId, []common.AssetRight{common.AVIEWER})
	if err != nil {
		return r, err
	}

	existing, err := stub.GetState(propId)
	if err != nil {
		// TODO
	}

	if existing != nil {
		err = r.Decode(existing)
		if err != nil {
			// TODO
		}
		return r, nil
	} else {
		return r, errors.New("No such proposal : " + propId)
	}
}

// TODO use stateful batching in case of restart
// TODO id by enrollment id ? BID-[enrollId]-[requestId] ?
func (t *ReinsuranceProposalCC) create_prop_id(requestId string) string {
	c := atomic.AddUint64(&counter, 1)
	return fmt.Sprintf("%s-%s-%d", proposalPrefix, requestId, c)
}

func get_unix_millisec() uint64 {
	now := time.Now()
	nanos := now.UnixNano()
	return uint64(nanos / 1000000)
}

// ============================================================================================================================
// Main
// ============================================================================================================================
func main() {
	err := shim.Start(new(ReinsuranceProposalCC))
	if err != nil {
		fmt.Printf("Error starting ReinsuranceProposalCC: %s", err)
	}
}
