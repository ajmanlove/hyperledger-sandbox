package main

import (
	"errors"
	"fmt"
	// "encoding/json"
	// "strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim/crypto/attr"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

var logger = shim.NewLogger("EnrollmentServiceCC")
var enrollmentTable = "Enrollment"

// Enrolls contact information for users
type EnrollmentServiceCC struct {
}

// // TBD
// type EnrollmentRecord struct {
// 	enrollmentId 	string,
// 	contact				string
// }

func (t *EnrollmentServiceCC) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	logger.Debug("Init Chaincode...")

	if len(args) != 0 {
		return nil, errors.New("Init does not support arguments")
	}

	// Create enrollment table
	err := stub.CreateTable(enrollmentTable, []*shim.ColumnDefinition{
		&shim.ColumnDefinition{Name: "enrollmentId", Type: shim.ColumnDefinition_STRING, Key: true},
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
	switch function {
		case "get_contact":
			return t.get_contact(stub, args)
		default:
			return nil, errors.New("Unrecognized function : " + function)
	}
}

func (t *EnrollmentServiceCC) enroll(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	logger.Debug("enroll_0 ...")

	callerCert, err := stub.GetCallerCertificate()
	if err != nil {
		logger.Error(err)
		return nil, errors.New("Failed to get caller cert")
	}
	logger.Debugf("Caller CERT is [ %v ]", callerCert)

	bytes, err := attr.GetValueFrom("enrollmentId", callerCert)
	if err != nil {
		logger.Errorf("Failed to get enrollmentId from cert. error is [ %v ]", err)
		return nil, errors.New("Failed to get enrollmentId from cert")
	}
	id := string(bytes)
	logger.Debugf("Caller enrollmentId is [ %v ]", id)

	bytes, err = attr.GetValueFrom("contact", callerCert)
	if err != nil {
		logger.Errorf("Failed to get contact from cert. error is [ %v ]", err)
		return nil, errors.New("Failed to get contact from cert")
	}
	contact := string(bytes)
	logger.Debugf("Caller CONTACT is [ %v ]", contact)

	ok, err := stub.InsertRow(enrollmentTable, shim.Row{
		Columns: []*shim.Column{
			&shim.Column{Value: &shim.Column_String_{String_: id}},
			&shim.Column{Value: &shim.Column_Bytes{Bytes: []byte(contact)}}},
	})

	if !ok && err == nil {
		fmt.Println("Error inserting row")
		return nil, errors.New("enrollmentId was already enrolled " + id)
	}

	return nil, nil
}

func (t *EnrollmentServiceCC) get_contact(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, errors.New("Expected 1 arg for get_contact")
	}

	enrollId := args[0]

	var columns []shim.Column
	col1 := shim.Column{Value: &shim.Column_String_{String_: enrollId}}
	columns = append(columns, col1)

	row, err := stub.GetRow(enrollmentTable, columns)
	if err != nil {
		return nil, fmt.Errorf("Failed retrieving enrollment [%s]: [%s]", enrollId, err)
	}

	contact := row.Columns[1].GetBytes()

	return contact, nil

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
