package main

import (
	"errors"
	"fmt"
	// "strconv"
	"strings"
	"sync/atomic"
	"time"

	"github.com/ajmanlove/hyperledger-sandbox/reinsurance_poc/common"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/core/util"
)

var logger = shim.NewLogger("ReinsuranceRequestCC")
var assetManagementCCId = ""
var counter uint64 = 0
var submissionPrefix = "REQ"

type ReinsuranceRequestCC struct {
}

func (t *ReinsuranceRequestCC) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	logger.Debugf("enter Init, function: [%s], args [%s]", function, args)

	switch function {
	case "init":
		if len(args) != 1 {
			return nil, errors.New("Expects chaincode id for asset_management as init arg")
		}
		assetManagementCCId = args[0]
		return nil, nil
	default:
		return nil, errors.New("Unrecognized Init function: " + function)
	}
	return nil, nil
}

func (t *ReinsuranceRequestCC) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	logger.Debugf("enter Invoke, function: [%s], args [%s]", function, args)
	switch function {
	case "submit":
		return t.submit(stub, args)
	default:
		return nil, errors.New("Unrecognized Invoke function: " + function)
	}
}

func (t *ReinsuranceRequestCC) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	logger.Debugf("enter Query, function: [%s], args [%s]", function, args)
	switch function {
	case "get_request":
		if len(args) != 1 {
			return nil, errors.New("Expected 1 arg, asset id")
		}
		return t.get_request(stub, args)
	default:
		return nil, errors.New("Unrecognized Invoke function: " + function)
	}
}

func (t *ReinsuranceRequestCC) get_request(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	requestId := args[0]
	bytes, err := stub.ReadCertAttribute("enrollmentId")
	if err != nil {
		logger.Error(err)
		return nil, errors.New("failed to get enrollmentId attribute")
	}
	enrollmentId := string(bytes)

	invokeArgs := util.ToChaincodeArgs("get_asset_rights", enrollmentId, requestId)
	bytes, err = stub.QueryChaincode(assetManagementCCId, invokeArgs)

	var response common.AssetRightsResponse
	if err := response.Decode(bytes); err != nil {
		logger.Error(err)
		return nil, errors.New("Failed to deserialize AssetRight")
	}

	if response.Contains(common.AVIEWER) {
		return stub.GetState(requestId) // TODO visibility
	} else {
		return nil, errors.New("Insufficient rights to view this asset, enrollment id " + enrollmentId)
	}
}

func (t *ReinsuranceRequestCC) submit(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	logger.Debug("submit()")

	id := t.get_new_submission_id()
	requestees := strings.Split(args[0], ",")
	portfolioSha := args[1]
	portfolioUrl := args[2]
	contractText := args[3]
	schema := args[4]
	schemaVersion := args[5]
	status := "requested"
	bytes, err := stub.ReadCertAttribute("enrollmentId")
	now := get_unix_millisec()

	if err != nil {
		logger.Error(err)
		return nil, errors.New("failed to get enrollmentId attribute")
	}
	requestor := string(bytes)

	rr := common.ReinsuranceRequest{
		Id:           id,
		Requestor:    requestor,
		Requestees:   requestees,
		PortfolioSHA: portfolioSha,
		PortfolioURL: portfolioUrl,
		ContractText: contractText,
		ISQLSchema:   schema,
		ISQLVersion:  schemaVersion,
		Status:       status,
		Created:      now,
		Updated:      now,
	}

	// Submit
	bytes, err = rr.Encode()
	if err != nil {
		logger.Error(err)
		return nil, errors.New("Failed to serialize request")
	}

	err = stub.PutState(id, bytes)
	if err != nil {
		logger.Error(err)
		return nil, errors.New("Failed to submit request")
	}

	// Note with asset management
	invokeArgs := util.ToChaincodeArgs("new_request", id, requestor, strings.Join(requestees, ","), fmt.Sprintf("%d", now))
	response, err := stub.InvokeChaincode(assetManagementCCId, invokeArgs)
	if err != nil {
		logger.Error(err)
		return nil, errors.New("failed to manage new request")
	}

	logger.Debugf("Asset management response is %s", string(response))
	logger.Debugf("Asset management error is %s", err)

	return nil, nil
}

// TODO use stateful batching in case of restart
func (t *ReinsuranceRequestCC) get_new_submission_id() string {
	c := atomic.AddUint64(&counter, 1)
	return fmt.Sprintf("%s-%d", submissionPrefix, c)
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
	err := shim.Start(new(ReinsuranceRequestCC))
	if err != nil {
		fmt.Printf("Error starting ReinsuranceRequestCC: %s", err)
	}
}
