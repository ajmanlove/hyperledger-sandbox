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

var am = AssetManager{}
var um = UserManager{}

type AssetManagementCC struct {
}

func (t *AssetManagementCC) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	logger.Debug("Init Chaincode...")

	am.Init(stub)
	um.Init(stub)

	if len(args) != 0 {
		return nil, errors.New("Init does not support arguments")
	}

	logger.Debug("Init Chaincode finished")

	return nil, nil
}

func (t *AssetManagementCC) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	logger.Debug("enter Invoke")
	switch function {
	case common.AM_REGISTER_CC_ARG:
		if len(args) != 2 {
			return nil, errors.New("Expects 2 args: ['chaincode_name', 'chaincode_identifier']")
		}
		cc_name := args[0]
		cc_id := args[1]

		// TODO use bool
		_, err := am.RegisterChaincode(stub, cc_name, cc_id)

		return nil, err

	case common.AM_NEW_REQ_ARG:
		// expects ["id", "requestor", "requestees,..", "createDate"]
		if len(args) != 4 {
			return nil, errors.New("Expects three arguments: ['id', 'requestor', 'requestees,..', 'createDate']")
		}
		return t.manage_request(stub, args)
	case common.AM_NEW_BID_ARG:
		return t.manage_proposal(stub, args)

	default:
		return nil, errors.New("Unrecognized Invoke function: " + function)
	}

}

func (t *AssetManagementCC) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	switch function {
	case common.AM_GET_CC_NAME_ARG:
		if len(args) != 1 {
			return nil, errors.New("Expects 2 arguments ['chaincode_id']")
		}
		name, err := am.GetChaincodeName(stub, args[0])
		if err != nil {
			//TODO err
		}

		ccn := common.CCNameResponse{Name: name}
		r, err := ccn.Encode()
		if err != nil {
			logger.Error(err)
			return nil, errors.New("failed to encode get_cc_name response")
		}
		return r, nil

	case common.AM_GET_U_ASST_ARG:
		bytes, err := stub.ReadCertAttribute("enrollmentId")
		if err != nil {
			logger.Error(err)
			return nil, errors.New("failed to get enrollmentId attribute")
		}
		enrollmentId := string(bytes)
		record, err := um.GetUserAssetRecord(stub, enrollmentId)
		if err != nil {
			return nil, err
		}

		bytes, err = record.Encode()
		if err != nil {
			logger.Error(err)
			return nil, errors.New("Failed to serialize record response")
		}

		return bytes, nil

	case common.AM_GET_AST_RIGHTS_ARG:
		// TODO only admin access to this method?
		// TODO cert attribute ?
		if len(args) != 2 {
			return nil, errors.New("Expects 2 arguments ['enrollmentId', 'assetId']")
		}

		enrollmentId := args[0]
		assetId := args[1]

		exists, err := am.AssetExists(stub, assetId)
		if err != nil {
			// TODO
		}

		var rights []common.AssetRight
		var response common.AssetRightsResponse
		if exists {
			rights, err = am.GetUserRights(stub, assetId, enrollmentId)
			if err != nil {
				// TODO
			}
			// TODO add exists
			response = common.BuildArr(exists, rights)

		} else {
			rights = make([]common.AssetRight, 0)
			response = common.BuildArr(false, rights)
		}

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

	record, err := um.GetUserAssetRecord(stub, requestor)
	if err != nil {
		return nil, err
	}

	// TODO err
	err = am.AssignRights(stub, requestId, requestor, []common.AssetRight{common.AOWNER, common.AVIEWER})

	record.Submissions = append(
		record.Submissions,
		common.SubmissionRecord{
			SubmissionId: requestId,
			Requestees:   requestees,
			Created:      createDate,
			Updated:      createDate,
		})
	_, err = um.SaveUserAssetRecord(stub, requestor, record)
	if err != nil {
		logger.Error(err)
		return nil, errors.New("Failed to save record for id " + requestor)
	}

	for _, requestee := range requestees {
		record, err := um.GetUserAssetRecord(stub, requestee)
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

		// TODO err
		err = am.AssignRights(stub, requestId, requestee, []common.AssetRight{common.AVIEWER})

		_, err = um.SaveUserAssetRecord(stub, requestee, record)
		if err != nil {
			logger.Error(err)
			return nil, errors.New("Failed to save record for id " + requestee)
		}
	}

	return nil, err
}

func (t *AssetManagementCC) manage_proposal(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 4 {
		return nil, errors.New("Expects 4 args ['proposalId', 'requestId', 'bidder', 'createDate']")
	}
	proposalId := args[0]
	requestId := args[1]
	bidder := args[2]
	// TODO parse err
	createDate, err := strconv.ParseUint(args[3], 10, 64)

	record, err := um.GetUserAssetRecord(stub, bidder)
	if err != nil {
		return nil, err
	}

	var originalReq common.RequestRecord
	found := false
	for _, request := range record.Requests {
		if request.SubmissionId == requestId {
			originalReq = request
			found = true
			break
		}
	}

	if !found {
		return nil, fmt.Errorf("IllegalState user %s has no request asset %s", bidder, requestId)
	}

	record.Proposals = append(
		record.Proposals,
		common.ProposalRecord{
			SubmissionId: requestId,
			ProposalId:   proposalId,
			Created:      createDate,
			Updated:      createDate,
			UpdatedBy:    bidder,
		})
	_, err = um.SaveUserAssetRecord(stub, bidder, record)
	if err != nil {
		logger.Error(err)
		return nil, errors.New("Failed to save record for id " + bidder)
	}

	err = am.AssignRights(stub, proposalId, bidder, []common.AssetRight{common.AOWNER, common.AVIEWER, common.AUPDATER})
	if err != nil {
		logger.Error(err)
		return nil, errors.New("Failed to assign rights to id " + bidder)
	}

	record, err = um.GetUserAssetRecord(stub, originalReq.Requestor)
	if err != nil {
		return nil, err
	}

	record.Proposals = append(
		record.Proposals,
		common.ProposalRecord{
			SubmissionId: requestId,
			ProposalId:   proposalId,
			Created:      createDate,
			Updated:      createDate,
			UpdatedBy:    bidder,
		})
	_, err = um.SaveUserAssetRecord(stub, originalReq.Requestor, record)
	if err != nil {
		logger.Error(err)
		return nil, errors.New("Failed to save record for id " + bidder)
	}

	err = am.AssignRights(stub, proposalId, originalReq.Requestor, []common.AssetRight{common.AVIEWER, common.AAPPROVAL, common.AUPDATER})
	if err != nil {
		logger.Error(err)
		return nil, errors.New("Failed to assign rights to id " + bidder)
	}

	return nil, nil
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
