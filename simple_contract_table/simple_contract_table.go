package main

import (
	"errors"
	"fmt"
	// "encoding/json"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

var logger = shim.NewLogger("SimpleContractTableChaincode")
var contractTable = "SimpleContractTable"

type SimpleContractTableChaincode struct {
}

type SimpleContract struct {
	ItemName 							string 	`json:"itemName"`
	CreateDate 						uint64 	`json:"createDate,string"`
	TotalInsuredValue 		int32 		`json:"totalInsuredValue,string"`
}

func (t *SimpleContractTableChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	if len(args) != 0 {
		return nil, errors.New("Unexpected arguments for Init")
	}

	err := t.create_contract_table(stub)
	if err != nil {
		logger.Error(err)
		return nil, errors.New("Unable to init contract table")
	}

	return nil, nil
}

func (t *SimpleContractTableChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	logger.Debug("enter Invoke")
	switch function {
		case "submit_contract":
			return t.submit_contract(stub, args)
		// case "update_total_insured_value":
		// 	return t.update_tiv(stub, args)
		// case "remove_contract":
		// 	return t.remove_contract(stub, args)
		default:
			return nil, errors.New("Unknown Invoke function : " + function)
	}
	return nil, nil

}

func (t *SimpleContractTableChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	switch function {
		case "get_contract":
			logger.Debug("Query:get_contract : %s", args[0])
			sc, _ := t.get_contract(stub, args[0])
			// bytes, err := stub.GetState(args[0])
			// if err != nil {
			// 	logger.Error(err)
			// 	return nil, errors.New("Unable to retrieve contract with id " + args[0])
			// }
			return []byte(sc), nil
		default:
			return nil, errors.New("Unknown function : " + function)
	}
	return nil, nil
}

func (t *SimpleContractTableChaincode) get_contract(stub shim.ChaincodeStubInterface, id string) (string, error) {
	var key []shim.Column
	key0 := shim.Column{Value: &shim.Column_String_{String_: id}}
	key = append(key, key0)

	row, err := stub.GetRow(contractTable, key)
	if err != nil {
		logger.Error(err)
		return "", errors.New("Failed to get row with key " + id)
	}

	rowString := fmt.Sprintf("%s", row)
	return rowString, nil
}

// submit a simple contract
func (t *SimpleContractTableChaincode) submit_contract(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
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

	i, err := strconv.ParseInt(args[3], 10, 32)

	if err != nil {
		logger.Error(err)
		return nil, errors.New("Invalid create date, expected unix timestamp " + args[2])
	}
	tiv := int32(i)

	sc.ItemName = name
	sc.CreateDate = createDate
	sc.TotalInsuredValue = tiv

	return t.save_contract(stub, id, sc)
}

// Get a generic contract struct
func (t *SimpleContractTableChaincode) get_contract_template() (SimpleContract, error) {
	sc := SimpleContract{ItemName: "UNDEFINED", CreateDate: 0, TotalInsuredValue: 0}
	return sc, nil
}

// save contract to the ledger
func (t *SimpleContractTableChaincode) save_contract(stub shim.ChaincodeStubInterface, id string, sc SimpleContract) ([]byte, error) {

	var columns []*shim.Column
	col1 := shim.Column{Value: &shim.Column_String_{String_: id}}
	col2 := shim.Column{Value: &shim.Column_String_{String_: sc.ItemName}}
	col3 := shim.Column{Value: &shim.Column_Uint64{Uint64: sc.CreateDate}}
	col4 := shim.Column{Value: &shim.Column_Int32{Int32: sc.TotalInsuredValue}}
	columns = append(columns, &col1)
	columns = append(columns, &col2)
	columns = append(columns, &col3)
	columns = append(columns, &col4)

	row := shim.Row{Columns: columns}

	r, err := stub.InsertRow(contractTable, row)
	if err != nil {
		logger.Error(err)
		return nil, errors.New("Unable to save contract row, id : " + id)
	}
	if !r {
		return nil, errors.New("Row did not save! id : " + id)
	}


	return nil, nil
}

// Create contracts storage as table
func (t *SimpleContractTableChaincode) create_contract_table(stub shim.ChaincodeStubInterface) error {
	// Create table one
	var columnDefs []*shim.ColumnDefinition
	c1 := shim.ColumnDefinition{Name: "id",
		Type: shim.ColumnDefinition_STRING, Key: true}
	c2 := shim.ColumnDefinition{Name: "itemName",
		Type: shim.ColumnDefinition_STRING, Key: false}
	c3 := shim.ColumnDefinition{Name: "createDate",
		Type: shim.ColumnDefinition_UINT64, Key: false}
	c4 := shim.ColumnDefinition{Name: "totalInsuredValue",
		Type: shim.ColumnDefinition_INT32, Key: false}
	columnDefs = append(columnDefs, &c1)
	columnDefs = append(columnDefs, &c2)
	columnDefs = append(columnDefs, &c3)
	columnDefs = append(columnDefs, &c4)
	return stub.CreateTable(contractTable, columnDefs)
}

// ============================================================================================================================
// Main
// ============================================================================================================================
func main() {
	err := shim.Start(new(SimpleContractTableChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
