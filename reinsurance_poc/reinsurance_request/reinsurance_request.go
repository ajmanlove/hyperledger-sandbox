package main

import (
	"errors"
	"fmt"
	"encoding/json"
	// "strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

var logger = shim.NewLogger("SimpleContractChaincode")

type ReinsuranceRequestCC struct {
}

// TBD
type ReinsuranceRequest struct {
	ContractType					string `json:"contractType"`
	ContractSubType				string `json:"contractSubType"`
	AssetType							string `json:"assetType"`
	TotalInsuredValue 		int 	 `json:"totalInsuredValue,string"`
	AggregateLimit				int 	 `json:"aggregateLimit,string"`
	PortfolioHash					string `json:"portfolioHash"`
	PortfolioURL					string `json:"portfolioUrl"`
	InExcessOf						int		 `json:"inExcessOf,string"`
	Status								string `json:"status"`
	Requestor							string `json:"requestor"`
	Requestees						[]string `json:"requestees"`
}

func (t *ReinsuranceRequestCC) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	logger.Debug("Init()")

	switch function {
		case "init":
			if len(args) != 0 {
				return nil, errors.New("Init does not support arguments")
			}
			return nil, nil
		default:
			return nil, errors.New("Unrecognized init function : " + function)
	}
}

func (t *ReinsuranceRequestCC) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	logger.Debug("enter Invoke : " + function)

	switch function {
	case "submit_request":
			return t.submit_request(stub, args)
		default:
			return nil, errors.New("Unknown Invoke function : " + function)
	}

}

func (t *ReinsuranceRequestCC) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	logger.Debug("enter Query : " + function)

	switch function {
		case "get_request":
			id := args[0]
			bytes, err := stub.GetState(id)
			if err != nil {
				logger.Error(err)
				return nil, errors.New("Unable to retrieve contract with id " + args[0])
			}
			return bytes, nil

		default:
			return nil, errors.New("Unknown Query function : " + function)
	}

}

func (t *ReinsuranceRequestCC) submit_request(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	rr := ReinsuranceRequest {
		ContractType : "liability",
		ContractSubType	: "facultative",
		AssetType	: "railroad",
		TotalInsuredValue : 100000000,
		AggregateLimit : 100000000,
		PortfolioHash	: "2e1b1b0cb7bfce4cf47706752a234f29",
		PortfolioURL : "http://mybucket.s3-website-us-east-1.amazonaws.com/",
		InExcessOf : 50000000,
		Status : "open",
		Requestor	: "myusername",
		Requestees	: []string {"someone", "someoneelse"},
	}

	bytes, err := json.Marshal(rr)
	if err != nil {
		logger.Error(err)
		return nil, errors.New("Failed to serialize ReinsuranceRequest object")
	}

	id := "1"

	err = stub.PutState(id, bytes)
	if err != nil {
		logger.Error("err")
		return nil, errors.New("Failed to put request")
	}
	return nil, nil
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
