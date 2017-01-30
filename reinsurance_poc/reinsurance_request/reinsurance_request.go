package main

import (
	"errors"
	"fmt"
	"encoding/json"
	// "strconv"
	"strings"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/core/util"

)

var logger = shim.NewLogger("ReinsuranceRequestCC")
var assetManagementCCId = ""

type ReinsuranceRequestCC struct {
}

// TBD
type ReinsuranceRequest struct {
	Id 										string 		`json:"id"`
	PortfolioSHA					string 		`json:"portfolioSha"`
	PortfolioURL					string 		`json:"portfolioUrl"`
	Status								string 		`json:"status"`
	Requestor							string 		`json:"requestor"`
	Requestees						[]string 	`json:"requestees"`
	ContractText					string		`json:"contractText"`
	Created								uint64 		`json:"created"`
	Updated								uint64		`json:"updated"`
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
				return nil, errors.New("Expected 1 arg, id") // TODO temporary
			}
			return stub.GetState(args[0])
		default:
			return nil, errors.New("Unrecognized Invoke function: " + function)
	}
}

func (t *ReinsuranceRequestCC) submit(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	logger.Debug("submit()")

	id := "1" // TODO
	requestees := strings.Split(args[0], ",")
	portfolioSha := args[1]
	portfolioUrl := args[2]
	contractText := args[3]
	status := "requested"
	bytes, err := stub.ReadCertAttribute("enrollmentId")
	now := get_unix_millisec()

	if err != nil {
		logger.Error(err)
		return nil, errors.New("failed to get enrollmentId attribute")
	}
	requestor := string(bytes)

	rr := ReinsuranceRequest {
		Id: id,
		Requestor: requestor,
		Requestees: requestees,
		PortfolioSHA: portfolioSha,
		PortfolioURL: portfolioUrl,
		ContractText: contractText,
		Status: status,
		Created: now,
		Updated: now,
	}

	// Submit
	bytes, err = json.Marshal(rr)
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
	invokeArgs := util.ToChaincodeArgs("new_request", id, requestor, strings.Join(requestees,","), fmt.Sprintf("%d", now))
	response, err := stub.InvokeChaincode(assetManagementCCId, invokeArgs)
	if err != nil {
		logger.Error(err)
		return nil, errors.New("failed to manage new request")
	}

	logger.Debugf("Asset management response is %s", string(response))
	logger.Debugf("Asset management error is %s", err)

	return nil, nil
}


func get_unix_millisec() (uint64) {
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
