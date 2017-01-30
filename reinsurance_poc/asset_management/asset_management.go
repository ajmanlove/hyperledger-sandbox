package main

import (
	"errors"
	"fmt"
	"strings"
	"encoding/json"
	// "strconv"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

var logger = shim.NewLogger("AssetManagementCC")
var assetTable = "Assets"

type AssetManagementCC struct {
}

type AssetsRecord struct {
	Submitted 		[]string	`json:"submitted"`
	Requested			[]string	`json:"submitted"`
	Proposed			[]string 	`json:"proposed"`
	Accepted			[]string	`json:"accepted"`
	Rejected			[]string	`json:"rejected"`
}

func (t *AssetManagementCC) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	logger.Debug("Init Chaincode...")

	if len(args) != 0 {
		return nil, errors.New("Init does not support arguments")
	}

	// Create enrollment table
	err := stub.CreateTable(assetTable, []*shim.ColumnDefinition{
		&shim.ColumnDefinition{Name: "enrollmentId", Type: shim.ColumnDefinition_STRING, Key: true},
		&shim.ColumnDefinition{Name: "Records", Type: shim.ColumnDefinition_BYTES, Key: false},
	})

	if err != nil {
		return nil, errors.New("Failed creating Assets table.")
	}

	logger.Debug("Init Chaincode finished")

	return nil, nil
}

func (t *AssetManagementCC) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	logger.Debug("enter Invoke")
	switch function {
		case "new_request":
			// expects ["id", "requestor", "requestees,.."]
			if len(args) != 3 {
				return nil, errors.New("Expects three arguments: ['id', 'requestor', 'requestees,..']")
			}
			return t.manage_request(stub, args)
		default:
			return nil, errors.New("Unrecognized Invoke function: " + function)
	}

}

func (t *AssetManagementCC) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	switch function {
	case "get_assets_record":
			return t.get_assets_record(stub)
		default:
			return nil, errors.New("Unrecognized function : " + function)
	}
}

func (t *AssetManagementCC) manage_request(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	requestId := args[0]
	requestor := args[1]
	requestees := strings.Split(args[0], ",")

	record, err := t.get_or_create_record(stub, requestor)
	if err != nil {
		return nil, err
	}
	record.Submitted = append(record.Submitted, requestId)
	_, err = t.save_record(stub, record, requestor)
	if err != nil {
		logger.Error(err)
		return nil, errors.New("Failed to save record for id " + requestor)
	}

	for _,element := range requestees {
		record, err := t.get_or_create_record(stub, element)
		if err != nil {
			return nil, err
		}
		record.Requested = append(record.Submitted, requestId)
		_, err = t.save_record(stub, record, element)
		if err != nil {
			logger.Error(err)
			return nil, errors.New("Failed to save record for id " + requestor)
		}
	}

	return nil, err
}

func (t *AssetManagementCC) save_record(stub shim.ChaincodeStubInterface, record AssetsRecord, enrollmentId string) (bool, error) {
	return false, errors.New("save_record() not implemented()")
	recordBytes, err := json.Marshal(record)
	if err != nil {
		logger.Error(err)
		return false, errors.New("Failed to serialize record")
	}

	existing, err := t.get_record(stub, enrollmentId)
	if err != nil {
		logger.Error(err)
		return false, errors.New("Failed to get record for enrollment id : " + enrollmentId)
	}

	if len(existing.Columns) == 0 {
		return stub.InsertRow(assetTable, shim.Row{
			Columns: []*shim.Column{
				&shim.Column{Value: &shim.Column_String_{String_: enrollmentId}},
				&shim.Column{Value: &shim.Column_Bytes{Bytes: []byte(recordBytes)}}},
		})

	} else {
		return stub.ReplaceRow(assetTable, shim.Row{
			Columns: []*shim.Column{
				&shim.Column{Value: &shim.Column_String_{String_: enrollmentId}},
				&shim.Column{Value: &shim.Column_Bytes{Bytes: []byte(recordBytes)}}},
		})
	}
}

func (t *AssetManagementCC) get_or_create_record(stub shim.ChaincodeStubInterface, enrollmentId string) (AssetsRecord, error) {
	var r AssetsRecord

	existing, err := t.get_record(stub, enrollmentId)

	if len(existing.Columns) == 0 {
		r = AssetsRecord {
			Submitted: make([]string, 0),
			Requested: make([]string, 0),
			Proposed: make([]string, 0),
			Accepted: make([]string, 0),
			Rejected: make([]string, 0),
		}

	} else {
		err = json.Unmarshal(existing.Columns[1].GetBytes(), &r)
		if err != nil {
			logger.Error(err)
			return r, errors.New("Failed to deserialize asset record: " + enrollmentId)
		}
	}

	return r, nil
}

func (t *AssetManagementCC) get_assets_record(stub shim.ChaincodeStubInterface) ([]byte, error) {
	bytes, err := stub.ReadCertAttribute("enrollmentId")
	if err != nil {
		logger.Error(err)
		return nil, errors.New("failed to get enrollmentId attribute")
	}
	enrollmentId := string(bytes)

	row, err := t.get_record(stub, enrollmentId)
	if err != nil {
		return nil, err
	}

	if len(row.Columns) == 0 {
		return nil, errors.New("No such record for enrollmentId " + enrollmentId)
	}

	return row.Columns[1].GetBytes(), nil

}

func (t *AssetManagementCC) get_record(stub shim.ChaincodeStubInterface, enrollmentId string) (shim.Row, error) {
	var columns []shim.Column
	col1 := shim.Column{Value: &shim.Column_String_{String_: enrollmentId}}
	columns = append(columns, col1)
	return stub.GetRow(assetTable, columns)
}


// ============================================================================================================================
// Main
// ============================================================================================================================
func main() {
	err := shim.Start(new(AssetManagementCC))
	if err != nil {
		fmt.Printf("Error starting AssetManagementCC: %s", err)
	}
}
