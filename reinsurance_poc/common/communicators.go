package common

import (
	"errors"

	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/core/util"
)

type AssetManagementCommunicator struct {
	CCName string
}

func (a *AssetManagementCommunicator) AssertHasAssetRights(stub shim.ChaincodeStubInterface, assetId string, rights []AssetRight) error {
	enrollmentId, err := a.GetEnrollmentAttr(stub)
	if err != nil {
		return err
	}

	invokeArgs := util.ToChaincodeArgs(AM_GET_AST_RIGHTS_ARG, enrollmentId, assetId)
	bytes, err := stub.QueryChaincode(a.CCName, invokeArgs)
	var response AssetRightsResponse
	if err := response.Decode(bytes); err != nil {
		return fmt.Errorf("Failed to deserialize AssetRightsRespnse due to %s", err)
	}

	if !response.Exists {
		return errors.New("No such asset id " + assetId)
	}

	for _, right := range rights {
		if !response.Contains(right) {
			return fmt.Errorf("Insuffienct rights on asset %s. Missing %d", assetId, right)
		}
	}

	return nil

}

func (a *AssetManagementCommunicator) GetEnrollmentAttr(stub shim.ChaincodeStubInterface) (string, error) {
	bytes, err := stub.ReadCertAttribute("enrollmentId")
	if err != nil {
		return "", fmt.Errorf("failed to get enrollmentId attribute due to : %s", err)
	}
	return string(bytes), nil
}
