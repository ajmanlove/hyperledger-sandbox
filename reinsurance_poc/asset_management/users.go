package main

import (
	"errors"

	"github.com/ajmanlove/hyperledger-sandbox/reinsurance_poc/common"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// TODO eventually use account ids
type UserManager struct{}

var userAssetsTable = "UserAssets"

// TODO check exists
func (a *UserManager) Init(stub shim.ChaincodeStubInterface) error {
	// Create enrollment table
	err := stub.CreateTable(userAssetsTable, []*shim.ColumnDefinition{
		{Name: "UserId", Type: shim.ColumnDefinition_STRING, Key: true},
		{Name: "Records", Type: shim.ColumnDefinition_BYTES, Key: false},
	})

	if err != nil {
		return errors.New("Failed creating UserAssets table.")
	}

	return nil
}

func (a *UserManager) UserRecordExists(stub shim.ChaincodeStubInterface, userId string) (bool, error) {
	existing, err := a.get_record(stub, userId)
	exists := len(existing.Columns) > 0

	logger.Debugf("UserRecordExists() id: %s, %s", userId, exists)
	return exists, err
}

func (a *UserManager) GetUserAssetRecord(stub shim.ChaincodeStubInterface, userId string) (common.UserAssetsRecord, error) {
	return a.get_or_create_record(stub, userId)
}

func (a *UserManager) SaveUserAssetRecord(stub shim.ChaincodeStubInterface, userId string, record common.UserAssetsRecord) (bool, error) {
	logger.Debug("save_record()")
	logger.Debugf("record is [%s]", record)

	recordBytes, err := record.Encode()
	if err != nil {
		logger.Error(err)
		return false, errors.New("Failed to serialize record")
	}

	existing, err := a.UserRecordExists(stub, userId)
	if err != nil {
		logger.Error(err)
		return false, errors.New("Failed to get record for enrollment id : " + userId)
	}

	logger.Debugf("marshalled is [%s]", recordBytes)

	if existing {
		return stub.InsertRow(userAssetsTable, shim.Row{
			Columns: []*shim.Column{
				{Value: &shim.Column_String_{String_: userId}},
				{Value: &shim.Column_Bytes{Bytes: []byte(recordBytes)}}},
		})

	} else {
		return stub.ReplaceRow(userAssetsTable, shim.Row{
			Columns: []*shim.Column{
				{Value: &shim.Column_String_{String_: userId}},
				{Value: &shim.Column_Bytes{Bytes: []byte(recordBytes)}}},
		})
	}
}

func (a *UserManager) get_or_create_record(stub shim.ChaincodeStubInterface, userId string) (common.UserAssetsRecord, error) {
	var r common.UserAssetsRecord

	existing, err := a.get_record(stub, userId)
	// TODO use err
	if len(existing.Columns) > 0 {
		logger.Debug("Record exists : " + userId)
		err = r.Decode(existing.Columns[1].GetBytes())
		if err != nil {
			logger.Error(err)
			return r, errors.New("Failed to deserialize user assets record: " + userId)
		}
		return r, nil
	} else {
		logger.Debug("Creating non-existant record : " + userId)
		r.Init()
		return r, nil
	}
}

func (a *UserManager) get_record(stub shim.ChaincodeStubInterface, userId string) (shim.Row, error) {
	var columns []shim.Column
	col1 := shim.Column{Value: &shim.Column_String_{String_: userId}}
	columns = append(columns, col1)
	return stub.GetRow(userAssetsTable, columns)
}
