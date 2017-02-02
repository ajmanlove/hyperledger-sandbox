package main

import (
	"errors"
	"fmt"
	// "encoding/json"
	// "strconv"

	"sync/atomic"
	"time"

	"github.com/ajmanlove/hyperledger-sandbox/reinsurance_poc/common"
	"github.com/hyperledger/fabric/core/chaincode/shim"
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

func (t *ReinsuranceProposalCC) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	logger.Debug("Init()")

	if len(args) != 1 {
		return nil, errors.New("Init expects expects asset management cc id as arg")
	}
	assetManagementCCId = args[0]

	return nil, nil
}

func (t *ReinsuranceProposalCC) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	logger.Debug("enter Invoke")
	switch function {
	case "propose":
		return t.propose(stub, args)
	case "counter":
		return nil, errors.New("counter not implemented")
	case "accept":
		return nil, errors.New("accept not implemented")
	case "reject":
		return nil, errors.New("reject not implemented")
	default:
		return nil, errors.New("Unrecognized Invoke function : " + function)
	}
}

func (t *ReinsuranceProposalCC) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	switch function {
	case "get_proposal":
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

	if len(args) != 2 {
		return nil, errors.New("Requires 2 args: ['requestId', 'contractText']")
	}

	requestId := args[0]
	contractText := args[1]
	now := get_unix_millisec()
	bytes, err := stub.ReadCertAttribute("enrollmentId")
	if err != nil {
		logger.Error(err)
		return nil, errors.New("failed to get enrollmentId attribute")
	}
	enrollmentId := string(bytes)
	// TODO does exist and have rights

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
		// TODO
	}

	err = stub.PutState(id, encoded)
	if err != nil {
		// TODO
	}

	// TODO AM rights and management

	return nil, nil
}

func (t *ReinsuranceProposalCC) get_proposal(stub shim.ChaincodeStubInterface, propId string) (common.ReinsuranceBid, error) {
	existing, err := stub.GetState(propId)
	if err != nil {
		// TODO
	}
	var r common.ReinsuranceBid
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
