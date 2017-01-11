package main

import (
	"errors"
	"fmt"
	"encoding/json"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

var logger = shim.NewLogger("SimpleContractChaincode")

type SimpleContractChaincode struct {
}

type SimpleContract struct {
	ItemName 							string 	`json:"itemName"`
	CreateDate 						uint64 	`json:"createDate,string"`
	TotalInsuredValue 		int 		`json:"totalInsuredValue,string"`
}

func (t *SimpleContractChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	return nil, nil
}

func (t *SimpleContractChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	logger.Debug("enter Invoke")
	switch function {
		case "submit_contract":
			return t.submit_contract(stub, args)

		default:
			return nil, errors.New("Unknown Invoke function : " + function)
	}

}

func (t *SimpleContractChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	switch function {
		case "get_contract":
			bytes, err := stub.GetState(args[0])
			if err != nil {
				return nil, errors.New("Unable to retrieve contract with id " + args[0])
			}
			return bytes, nil
		default:
			return nil, errors.New("Unknown function : " + function)
	}
}

// submit a simple contract
func (t *SimpleContractChaincode) submit_contract(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	sc, err := t.get_contract_template()

	if err != nil {
		logger.Error(err)
		return nil, errors.New("Failed to get contract stub")
	}

	id := args[0]
	name := args[1]
	createDate, err := strconv.ParseUint(args[2], 10, 64)
	if err != nil {
		return nil, errors.New("Invalid create date, expected unix timestamp " + args[2])
	}

	tiv, err := strconv.Atoi(args[2])
	if err != nil {
		return nil, errors.New("Invalid create date, expected unix timestamp " + args[2])
	}

	sc.ItemName = name
	sc.CreateDate = createDate
	sc.TotalInsuredValue = tiv

	return t.save_contract(stub, id, sc)
}

// save contract to the ledger
func (t *SimpleContractChaincode) save_contract(stub shim.ChaincodeStubInterface, id string, sc SimpleContract) ([]byte, error) {

  bytes, err := json.Marshal(sc)
	if err != nil {
		return nil, errors.New("Failed to serialize SimpleContract json")
	}

	err = stub.PutState(id, bytes)
	if err != nil {
		return nil, errors.New("Failed to save contract, id : " + id)
	}

	return nil, nil
}


// Get a generic contract struct
func (t *SimpleContractChaincode) get_contract_template() (SimpleContract, error) {
	var sc SimpleContract

	// just do structure instantiation...
	itemName         		:= "\"itemName\":\"UNDEFINED\", "							// Variables to define the JSON
	createDate       		:= "\"createDate\":\"UNDEFINED\", "
	totalInsuredValue   := "\"totalInsuredValue\":\"UNDEFINED\" "

	sc_json := "{"+itemName+createDate+totalInsuredValue+"}"

	err := json.Unmarshal([]byte(sc_json), &sc)

	return sc, err
}


// ============================================================================================================================
// Main
// ============================================================================================================================
func main() {
	err := shim.Start(new(SimpleContractChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
