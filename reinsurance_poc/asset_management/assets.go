package main

import (
	"errors"

	"github.com/ajmanlove/hyperledger-sandbox/reinsurance_poc/common"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// TODO eventually use account ids
var assetTable = "Assets"

type AssetManager struct {
}

// TODO check exists
func (a *AssetManager) Init(stub shim.ChaincodeStubInterface) error {
	// Create enrollment table
	err := stub.CreateTable(assetTable, []*shim.ColumnDefinition{
		{Name: "AssetId", Type: shim.ColumnDefinition_STRING, Key: true},
		{Name: "Record", Type: shim.ColumnDefinition_BYTES, Key: false},
	})

	if err != nil {
		return errors.New("Failed creating Assets table.")
	}

	return nil
}

func (a *AssetManager) RegisterChaincode(stub shim.ChaincodeStubInterface, cc_id string, cc_name string) (bool, error) {
	exists, err := a.ChaincodeExists(stub, cc_id)
	if err != nil {
		// TODO err
	}

	if exists {
		return stub.InsertRow(assetTable, shim.Row{
			Columns: []*shim.Column{
				{Value: &shim.Column_String_{String_: cc_id}},
				{Value: &shim.Column_Bytes{Bytes: []byte(cc_name)}}},
		})

	} else {
		return stub.ReplaceRow(assetTable, shim.Row{
			Columns: []*shim.Column{
				{Value: &shim.Column_String_{String_: cc_id}},
				{Value: &shim.Column_Bytes{Bytes: []byte(cc_name)}}},
		})
	}
}

func (a *AssetManager) ChaincodeExists(stub shim.ChaincodeStubInterface, cc_id string) (bool, error) {
	return a.AssetExists(stub, cc_id)
}

func (a *AssetManager) GetChaincodeName(stub shim.ChaincodeStubInterface, cc_id string) (string, error) {
	exists, err := a.ChaincodeExists(stub, cc_id)
	if err != nil {
		//TODO
	}
	if exists {
		r, err := a.get_table_row(stub, cc_id)
		if err != nil {
			// TODO
		}
		return string(r.Columns[1].GetBytes()), nil
	} else {
		return "", errors.New("No such chaincode registered with identifier " + cc_id)
	}
}

func (a *AssetManager) AssetExists(stub shim.ChaincodeStubInterface, assetId string) (bool, error) {
	r, err := a.get_table_row(stub, assetId)
	return len(r.Columns) > 0, err
}

func (a *AssetManager) AssignRights(stub shim.ChaincodeStubInterface, assetId string, userId string, rights []common.AssetRight) error {
	record, err := a.get_or_create_record(stub, assetId)
	if err != nil {
		// TODO
	}
	record.AssignUserRights(userId, rights)
	_, err = a.save_record(stub, assetId, record)
	if err != nil {
		// TODO
	}
	return nil
}

func (a *AssetManager) GetAssetRecord(stub shim.ChaincodeStubInterface, assetId string) (common.AssetRecord, error) {
	var r common.AssetRecord
	existing, err := a.get_table_row(stub, assetId)
	// TODO use err
	if len(existing.Columns) > 0 {
		err = r.Decode(existing.Columns[1].GetBytes())
		if err != nil {
			logger.Error(err)
			return r, errors.New("Failed to deserialize asset record: " + assetId)
		}
		return r, nil
	} else {
		return r, errors.New("No such asset record : " + assetId)
	}

}

func (a *AssetManager) GetUserRights(stub shim.ChaincodeStubInterface, assetId string, userId string) ([]common.AssetRight, error) {
	record, err := a.GetAssetRecord(stub, assetId)
	if err != nil {
		return make([]common.AssetRight, 0), err
	}
	return record.Rights[userId], nil
}

func (a *AssetManager) get_or_create_record(stub shim.ChaincodeStubInterface, assetId string) (common.AssetRecord, error) {
	var r common.AssetRecord

	existing, err := a.get_table_row(stub, assetId)
	// TODO use err
	if len(existing.Columns) > 0 {

		err = r.Decode(existing.Columns[1].GetBytes())
		if err != nil {
			logger.Error(err)
			return r, errors.New("Failed to deserialize asset record: " + assetId)
		}
		return r, nil
	} else {
		r.Init()
		return r, nil
	}
}

func (a *AssetManager) get_table_row(stub shim.ChaincodeStubInterface, assetId string) (shim.Row, error) {
	var columns []shim.Column
	col1 := shim.Column{Value: &shim.Column_String_{String_: assetId}}
	columns = append(columns, col1)
	return stub.GetRow(assetTable, columns)
}

func (a *AssetManager) save_record(stub shim.ChaincodeStubInterface, assetId string, record common.AssetRecord) (bool, error) {
	logger.Debug("save_record()")
	logger.Debugf("record is [%s]", record)

	recordBytes, err := record.Encode()
	if err != nil {
		logger.Error(err)
		return false, errors.New("Failed to serialize record")
	}

	exists, err := a.AssetExists(stub, assetId)
	if err != nil {
		// TODO
	}

	if exists {
		return stub.ReplaceRow(assetTable, shim.Row{
			Columns: []*shim.Column{
				{Value: &shim.Column_String_{String_: assetId}},
				{Value: &shim.Column_Bytes{Bytes: []byte(recordBytes)}}},
		})
	} else {
		return stub.InsertRow(assetTable, shim.Row{
			Columns: []*shim.Column{
				{Value: &shim.Column_String_{String_: assetId}},
				{Value: &shim.Column_Bytes{Bytes: []byte(recordBytes)}}},
		})
	}

}
