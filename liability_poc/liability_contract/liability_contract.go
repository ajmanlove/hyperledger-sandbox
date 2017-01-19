package main

import (
	"errors"
	"fmt"
	// "encoding/json"
	// "strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

var logger = shim.NewLogger("SimpleContractChaincode")

type LiabilityContractCC struct {
}

// TBD
type LiabilityContract struct {
	ItemName 							string 	`json:"itemName"`
	CreateDate 						uint64 	`json:"createDate,string"`
	TotalInsuredValue 		int 		`json:"totalInsuredValue,string"`
}

func (t *LiabilityContractCC) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	return nil, errors.New("No Init Implementation")
}

func (t *LiabilityContractCC) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	logger.Debug("enter Invoke")
	return nil, errors.New("No Invoke Implementation")

}

func (t *LiabilityContractCC) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	return nil, errors.New("No Query Implementation")
}

// ============================================================================================================================
// Main
// ============================================================================================================================
func main() {
	err := shim.Start(new(LiabilityContractCC))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
