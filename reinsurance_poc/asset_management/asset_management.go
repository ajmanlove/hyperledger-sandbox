package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/ajmanlove/hyperledger-sandbox/reinsurance_poc/common"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

var logger = shim.NewLogger("AssetManagementCC")
var assetTable = "Assets"

type AssetManagementCC struct {
}

func (t *AssetManagementCC) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	logger.Debug("Init Chaincode...")

	am := AssetManager{}
	am.init()

	um := UserManager{}
	um.init()

	if len(args) != 0 {
		return nil, errors.New("Init does not support arguments")
	}

	// Create enrollment table
	err := stub.CreateTable(assetTable, []*shim.ColumnDefinition{
		{Name: "id", Type: shim.ColumnDefinition_STRING, Key: true},
		{Name: "Records", Type: shim.ColumnDefinition_BYTES, Key: false},
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
	case "register_chaincode":
		if len(args) != 2 {
			return nil, errors.New("Expects 2 args: ['chaincode_name', 'chaincode_identifier']")
		}
		cc_name := args[0]
		cc_id := args[1]

		_, err := t.register_chaincode(stub, cc_name, cc_id)

		return nil, err

	case "new_request":
		// expects ["id", "requestor", "requestees,..", "createDate"]
		if len(args) != 4 {
			return nil, errors.New("Expects three arguments: ['id', 'requestor', 'requestees,..', 'createDate']")
		}
		return t.manage_request(stub, args)
	case "new_proposal":
		return nil, errors.New("new_proposal not implemented")

	default:
		return nil, errors.New("Unrecognized Invoke function: " + function)
	}

}

func (t *AssetManagementCC) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	switch function {
	case "get_cc_name":
		if len(args) != 1 {
			return nil, errors.New("Expects 2 arguments ['chaincode_id']")
		}

		row, err := t.get_record(stub, args[0])
		if err != nil {
			//TODO err
		}

		if len(row.Columns) == 0 {
			return nil, errors.New("No such chaincode record for identifier " + args[0])
		}

		ccn := common.CCNameResponse{Name: string(row.Columns[1].GetBytes())}
		r, err := ccn.Encode()
		if err != nil {
			logger.Error(err)
			return nil, errors.New("failed to encode get_cc_name response")
		}
		return r, nil

	case "get_assets_record":
		bytes, err := stub.ReadCertAttribute("enrollmentId")
		if err != nil {
			logger.Error(err)
			return nil, errors.New("failed to get enrollmentId attribute")
		}
		enrollmentId := string(bytes)
		record, err := t.get_or_create_record(stub, enrollmentId)
		if err != nil {
			return nil, err
		}

		bytes, err = record.Encode()
		if err != nil {
			logger.Error(err)
			return nil, errors.New("Failed to serialize record response")
		}

		return bytes, nil

	case "get_asset_rights":
		// TODO cert attribute ?
		if len(args) != 2 {
			return nil, errors.New("Expects 2 arguments ['enrollmentId', 'assetId']")
		}
		enrollmentId := args[0]
		assetId := args[1]

		rights, err := t.get_asset_rights(stub, enrollmentId, assetId)
		if err != nil {
			logger.Error(err)
			return nil, errors.New("Failed to get asset rights") // TODO better message
		}
		response := common.BuildArr(rights)

		return response.Encode()

	default:
		return nil, errors.New("Unrecognized function : " + function)
	}
}

// TODO think on this
func (t *AssetManagementCC) register_chaincode(stub shim.ChaincodeStubInterface, cc_name string, cc_id string) (bool, error) {
	existing, err := t.get_record(stub, cc_id)
	if err != nil {
		// TODO err
	}

	if len(existing.Columns) == 0 {
		return stub.InsertRow(assetTable, shim.Row{
			Columns: []*shim.Column{
				&shim.Column{Value: &shim.Column_String_{String_: cc_id}},
				&shim.Column{Value: &shim.Column_Bytes{Bytes: []byte(cc_name)}}},
		})

	} else {
		return stub.ReplaceRow(assetTable, shim.Row{
			Columns: []*shim.Column{
				&shim.Column{Value: &shim.Column_String_{String_: cc_id}},
				&shim.Column{Value: &shim.Column_Bytes{Bytes: []byte(cc_name)}}},
		})
	}
}

func (t *AssetManagementCC) manage_request(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	requestId := args[0]
	requestor := args[1]
	requestees := strings.Split(args[2], ",")
	createDate, err := strconv.ParseUint(args[3], 10, 64)
	// TODO parse err

	record, err := t.get_or_create_record(stub, requestor)
	if err != nil {
		return nil, err
	}

	record.GiveRights(requestId, []common.AssetRight{common.AOWNER, common.AVIEWER})

	record.Submissions = append(
		record.Submissions,
		common.SubmissionRecord{
			SubmissionId: requestId,
			Requestees:   requestees,
			Created:      createDate,
			Updated:      createDate,
		})
	_, err = t.save_record(stub, record, requestor)
	if err != nil {
		logger.Error(err)
		return nil, errors.New("Failed to save record for id " + requestor)
	}

	for _, element := range requestees {
		record, err := t.get_or_create_record(stub, element)
		if err != nil {
			return nil, err
		}

		record.Requests = append(
			record.Requests,
			common.RequestRecord{
				SubmissionId: requestId,
				Requestor:    requestor,
				Created:      createDate,
				Updated:      createDate,
			})

		record.GiveRight(requestId, common.AVIEWER)

		_, err = t.save_record(stub, record, element)
		if err != nil {
			logger.Error(err)
			return nil, errors.New("Failed to save record for id " + element)
		}
	}

	return nil, err
}

func (t *AssetManagementCC) save_record(stub shim.ChaincodeStubInterface, record common.AssetsRecord, enrollmentId string) (bool, error) {
	logger.Debug("save_record()")
	logger.Debugf("record is [%s]", record)

	recordBytes, err := record.Encode()
	if err != nil {
		logger.Error(err)
		return false, errors.New("Failed to serialize record")
	}

	existing, err := t.get_record(stub, enrollmentId)
	if err != nil {
		logger.Error(err)
		return false, errors.New("Failed to get record for enrollment id : " + enrollmentId)
	}

	logger.Debugf("marshalled is [%s]", recordBytes)

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

func (t *AssetManagementCC) get_or_create_record(stub shim.ChaincodeStubInterface, enrollmentId string) (common.AssetsRecord, error) {
	var r common.AssetsRecord

	existing, err := t.get_record(stub, enrollmentId)
	// TODO use err
	if len(existing.Columns) > 0 {

		err = r.Decode(existing.Columns[1].GetBytes())
		if err != nil {
			logger.Error(err)
			return r, errors.New("Failed to deserialize asset record: " + enrollmentId)
		}
		return r, nil
	} else {
		r.Init()
		return r, nil
	}
}

func (t *AssetManagementCC) get_record(stub shim.ChaincodeStubInterface, enrollmentId string) (shim.Row, error) {
	var columns []shim.Column
	col1 := shim.Column{Value: &shim.Column_String_{String_: enrollmentId}}
	columns = append(columns, col1)
	return stub.GetRow(assetTable, columns)
}

func (t *AssetManagementCC) get_asset_rights(stub shim.ChaincodeStubInterface, enrollmentId string, assetId string) ([]common.AssetRight, error) {
	record, _ := t.get_or_create_record(stub, enrollmentId)
	// TODO err

	return record.AssetRights[assetId], nil
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
