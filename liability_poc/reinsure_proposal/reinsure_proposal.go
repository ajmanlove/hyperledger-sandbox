package main

import (
	"errors"
	"fmt"
	// "encoding/json"
	// "strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

var logger = shim.NewLogger("SimpleContractChaincode")

type ReinsureProposalCC struct {
}

// TBD
type ReinsureProposal struct {
	ItemName 							string 	`json:"itemName"`
	CreateDate 						uint64 	`json:"createDate,string"`
	TotalInsuredValue 		int 		`json:"totalInsuredValue,string"`
}

func (t *ReinsureProposalCC) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	return nil, errors.New("No Init Implementation")
}

func (t *ReinsureProposalCC) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	logger.Debug("enter Invoke")
	return nil, errors.New("No Invoke Implementation")

}

func (t *ReinsureProposalCC) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	return nil, errors.New("No Query Implementation")
}

// ============================================================================================================================
// Main
// ============================================================================================================================
func main() {
	err := shim.Start(new(ReinsureProposalCC))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
