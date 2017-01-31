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

	if len(args) != 0 {
		return nil, errors.New("Init does not support arguments")
	}

	// Create enrollment table
	err := stub.CreateTable(assetTable, []*shim.ColumnDefinition{
		{Name: "enrollmentId", Type: shim.ColumnDefinition_STRING, Key: true},
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
	case "new_request":
		// expects ["id", "requestor", "requestees,..", "createDate"]
		if len(args) != 4 {
			return nil, errors.New("Expects three arguments: ['id', 'requestor', 'requestees,..', 'createDate']")
		}
		return t.manage_request(stub, args)
	default:
		return nil, errors.New("Unrecognized Invoke function: " + function)
	}

}

func (t *AssetManagementCC) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	switch function {
	case "get_assets_record":
		bytes, err := stub.ReadCertAttribute("enrollmentId")
		if err != nil {
			logger.Error(err)
			return nil, errors.New("failed to get enrollmentId attribute")
		}
		enrollmentId := string(bytes)
		record, err := t.get_or_create_record(stub, enrollmentId)
		// TODO err

		bytes, err = record.Encode()
		// TODO err

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

	t.give_record_rights(record, requestId, []common.AssetRight{common.AOWNER, common.AVIEWER})

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

		t.give_record_rights(record, requestId, []common.AssetRight{common.AVIEWER})

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

	if len(existing.Columns) == 0 {
		r = common.AssetsRecord{
			AssetRights: make(map[string][]common.AssetRight),
			Submissions: make([]common.SubmissionRecord, 0),
			Requests:    make([]common.RequestRecord, 0),
			Proposals:   make([]common.ProposalRecord, 0),
			Accepted:    make([]common.AcceptedProposal, 0),
			Rejected:    make([]common.RejectedProposal, 0),
			Contracts:   make([]common.SubmissionRecord, 0),
		}

	} else {
		err = r.Decode(existing.Columns[1].GetBytes())
		if err != nil {
			logger.Error(err)
			return r, errors.New("Failed to deserialize asset record: " + enrollmentId)
		}
	}

	return r, nil
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

func (t *AssetManagementCC) give_record_rights(record common.AssetsRecord, assetId string, rights []common.AssetRight) {
	if record.AssetRights[assetId] == nil {
		record.AssetRights[assetId] = rights
	} else {
		for _, e := range rights {
			if record.ContainsRight(assetId, e) {
				record.AssetRights[assetId] = append(record.AssetRights[assetId], e)
			}
		}
	}
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
