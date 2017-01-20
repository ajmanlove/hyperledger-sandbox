package main

import (
	"errors"
	"fmt"
	// "encoding/json"
	// "strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

var logger = shim.NewLogger("EnrollmentServiceCC")
var enrollmentTable = "Enrollment"

type EnrollmentServiceCC struct {
}

// TBD
type EnrollmentService struct {
}

func (t *EnrollmentServiceCC) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	logger.Debug("Init Chaincode...")

	if len(args) != 0 {
		return nil, errors.New("Init does not support arguments")
	}

	// Create enrollment table
	err := stub.CreateTable(enrollmentTable, []*shim.ColumnDefinition{
		&shim.ColumnDefinition{Name: "Id", Type: shim.ColumnDefinition_STRING, Key: true},
		&shim.ColumnDefinition{Name: "Enrollee", Type: shim.ColumnDefinition_BYTES, Key: false},
	})

	if err != nil {
		return nil, errors.New("Failed creating Enrollment table.")
	}

	logger.Debug("Init Chaincode finished")

	return nil, nil
}

func (t *EnrollmentServiceCC) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	logger.Debug("enter Invoke")
	switch function {
		case "enroll":
			return t.enroll(stub, args)
		default:
			return nil, errors.New("Unrecognized Invoke function: " + function)
	}

}

func (t *EnrollmentServiceCC) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	return nil, errors.New("No Query Implementation")
}

func (t *EnrollmentServiceCC) enroll(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	logger.Debug("Read cert attributes ...")

	callerCert, err := stub.GetCallerCertificate()
	if err != nil {
		logger.Error(err)
		return nil, errors.New("Failed to get caller cert")
	}
	logger.Debugf("Caller CERT is [ %v ]", callerCert)

	callerId, err := stub.ReadCertAttribute("id")
	if err != nil {
		logger.Error(err)
		return nil, errors.New("Failed to read role attribute")
	}
	logger.Debugf("caller ID is [ %v ]", callerId)

	CallerContact, err := stub.ReadCertAttribute("contact")
	if err != nil {
		logger.Error(err)
		return nil, errors.New("Failed to read contact attribute")
	}
	logger.Debugf("Caller CONTACT is [ %v ]", CallerContact)

	id := string(callerId)

	ok, err := stub.InsertRow(enrollmentTable, shim.Row{
		Columns: []*shim.Column{
			&shim.Column{Value: &shim.Column_String_{String_: id}},
			&shim.Column{Value: &shim.Column_Bytes{Bytes: callerCert}}},
	})

	if !ok && err == nil {
		fmt.Println("Error inserting row")
		return nil, errors.New("Id was already enrolled " + id)
	}

	return nil, nil
}
// ============================================================================================================================
// Main
// ============================================================================================================================
func main() {
	err := shim.Start(new(EnrollmentServiceCC))
	if err != nil {
		fmt.Printf("Error starting ReinsuranceProposalCC: %s", err)
	}
}
