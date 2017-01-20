package main

import (
	"errors"
	"fmt"
	// "encoding/json"
	// "strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

var logger = shim.NewLogger("EnrollmentServiceCC")

type EnrollmentServiceCC struct {
}

// TBD
type EnrollmentService struct {
}

func (t *EnrollmentServiceCC) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	logger.Debug("Init()")
	if len(args) != 0 {
		return nil, errors.New("Init does not support arguments")
	}
	return nil, nil
}

func (t *EnrollmentServiceCC) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	logger.Debug("enter Invoke")
	return nil, errors.New("No Invoke Implementation")

}

func (t *EnrollmentServiceCC) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	return nil, errors.New("No Query Implementation")
}

// ============================================================================================================================
// Main
// ============================================================================================================================
func main() {
	err := shim.Start(new(EnrollmentServiceCC))
	if err != nil {
		fmt.Printf("Error starting ReinsuranceProposalCC: %s", err)
	}
}
