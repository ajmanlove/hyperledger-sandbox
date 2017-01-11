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
		case "remove_contract":
			return t.remove_contract(stub, args)
		default:
			return nil, errors.New("Unknown Invoke function : " + function)
	}

}

func (t *SimpleContractChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	switch function {
		case "get_contract":
			logger.Debug("Query:get_contract : %s", args[0])
			bytes, err := stub.GetState(args[0])
			if err != nil {
				logger.Error(err)
				return nil, errors.New("Unable to retrieve contract with id " + args[0])
			}
			return bytes, nil
		default:
			return nil, errors.New("Unknown function : " + function)
	}
}

// submit a simple contract
func (t *SimpleContractChaincode) submit_contract(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	logger.Debug("submit_contract() %s", args[0])
	sc, err := t.get_contract_template()

	if err != nil {
		logger.Error(err)
		return nil, errors.New("Failed to get contract stub")
	}

	id := args[0]
	bytes, _ := stub.GetState(id)
	if bytes != nil {
		e := errors.New("Record with id " + id + " exists!")
		logger.Error(e)

		return nil, e
	}

	name := args[1]
	createDate, err := strconv.ParseUint(args[2], 10, 64)
	if err != nil {
		logger.Error(err)
		return nil, errors.New("Invalid create date, expected unix timestamp " + args[2])
	}

	tiv, err := strconv.Atoi(args[3])

	if err != nil {
		logger.Error(err)
		return nil, errors.New("Invalid create date, expected unix timestamp " + args[2])
	}

	sc.ItemName = name
	sc.CreateDate = createDate
	sc.TotalInsuredValue = tiv

	return t.save_contract(stub, id, sc)
}

func (t *SimpleContractChaincode) remove_contract(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	id := args[0]

	logger.Debug("remove_contract () id: ", id)
	bytes, err := stub.GetState(id)
	if bytes == nil {
		return nil, errors.New("No contract exists with id : " + id)
	}

	err = stub.DelState(id)
	if err != nil {
		logger.Error(err)
		return nil, errors.New("Failed to remove contract with id : " + id)
	}

	return nil, nil
}

// save contract to the ledger
func (t *SimpleContractChaincode) save_contract(stub shim.ChaincodeStubInterface, id string, sc SimpleContract) ([]byte, error) {

  bytes, err := json.Marshal(sc)
	if err != nil {
		logger.Error(err)
		return nil, errors.New("Failed to serialize SimpleContract json")
	}

	err = stub.PutState(id, bytes)
	if err != nil {
		logger.Error(err)
		return nil, errors.New("Failed to save contract, id : " + id)
	}

	return nil, nil
}


// Get a generic contract struct
func (t *SimpleContractChaincode) get_contract_template() (SimpleContract, error) {
	sc := SimpleContract{ItemName: "UNDEFINED", CreateDate: 0, TotalInsuredValue: 0}
	return sc, nil
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
