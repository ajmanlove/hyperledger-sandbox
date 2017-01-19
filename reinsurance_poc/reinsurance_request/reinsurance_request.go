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
	if len(args) != 0 {
		nil, errors.New("Init does not support arguments")
	}
	return nil, nil
}

func (t *ReinsuranceRequestCC) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	logger.Debug("enter Invoke")
	return nil, errors.New("No Invoke Implementation")

}

func (t *ReinsuranceRequestCC) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	return nil, errors.New("No Query Implementation")
}

// ============================================================================================================================
// Main
// ============================================================================================================================
func main() {
	err := shim.Start(new(ReinsuranceRequestCC))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
