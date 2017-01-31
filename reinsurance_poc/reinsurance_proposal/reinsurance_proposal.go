package main

import (
	"errors"
	"fmt"
	// "encoding/json"
	// "strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"

	"github.com/golang/protobuf/proto"
	pb "github.com/hyperledger/fabric/protos"
)

var logger = shim.NewLogger("ReinsuranceProposalCC")

type ReinsuranceProposalCC struct {
}

// TBD
type ReinsuranceProposal struct {
}

func (t *ReinsuranceProposalCC) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	logger.Debug("Init()")

	bytes, err := stub.GetPayload()
	var spec pb.ChaincodeInput
	err = proto.Unmarshal(bytes, &spec)

	if err != nil {
		logger.Error("ERROR HERE")
		logger.Error(err)

		return nil, fmt.Errorf("Failed to unmarshal payload due to : [%s]", err)
	}
	//if len(args) != 0 {
	//	return nil, errors.New("Init does not support arguments")
	//}
	return nil, nil
}

func (t *ReinsuranceProposalCC) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	logger.Debug("enter Invoke")
	return nil, errors.New("No Invoke Implementation")

}

func (t *ReinsuranceProposalCC) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	return nil, errors.New("No Query Implementation")
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
