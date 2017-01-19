package main

import (
	"errors"
	"fmt"
	// "encoding/json"
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
		case "test":
			fmt.Printf("Invoke:Test args : %s", args)
			return nil, nil
		default:
			return nil, errors.New("Unknown Invoke function : " + function)
	}

}

func (t *ReinsuranceRequestCC) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	logger.Debug("enter Query : " + function)

	switch function {
		default:
			return nil, errors.New("Unknown Query function : " + function)
	}

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
