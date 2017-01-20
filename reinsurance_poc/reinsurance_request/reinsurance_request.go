package main

import (
	"errors"
	"fmt"
	"encoding/json"
	"strconv"
	"strings"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

var logger = shim.NewLogger("ReinsuranceRequestCC")

type ReinsuranceRequestCC struct {
}

// TBD
type ReinsuranceRequest struct {
	ContractType					string 		`json:"contractType"`
	ContractSubType				string 		`json:"contractSubType"`
	AssetType							string 		`json:"assetType"`
	TotalInsuredValue 		int 	 		`json:"totalInsuredValue,string"`
	AggregateLimit				int 	 		`json:"aggregateLimit,string"`
	PortfolioHash					string 		`json:"portfolioHash"`
	PortfolioURL					string 		`json:"portfolioUrl"`
	InExcessOf						int		 		`json:"inExcessOf,string"`
	Status								string 		`json:"status"`
	Requestor							string 		`json:"requestor"`
	Requestees						[]string 	`json:"requestees"`
}

type RequestEvent struct {
	RequestId 						string 			`json:"requestId"`
	RequestorId						string			`json:"requestorId"`
	Recipients						[]Recipient	`json:"recipients"`
}

type Recipient struct {
	RecipientId 					string	`json:"recipientId"`
	RecipientContact			string	`json:"recipientContact"`
}

func (t *ReinsuranceRequestCC) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	logger.Debug("Init()")

	switch function {
		case "init":
			if len(args) != 0 {
				return nil, errors.New("Init does not support arguments")
			}
			return nil, nil
		default:
			return nil, errors.New("Unrecognized init function : " + function)
	}
}

func (t *ReinsuranceRequestCC) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	logger.Debug("enter Invoke : " + function)

	switch function {
	case "submit_request":
			return t.submit_request(stub, args)
		default:
			return nil, errors.New("Unknown Invoke function : " + function)
	}

}

func (t *ReinsuranceRequestCC) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	logger.Debug("enter Query : " + function)

	switch function {
		case "get_request":
			id := args[0]
			bytes, err := stub.GetState(id)
			if err != nil {
				logger.Error(err)
				return nil, errors.New("Unable to retrieve contract with id " + args[0])
			}
			return bytes, nil

		default:
			return nil, errors.New("Unknown Query function : " + function)
	}

}

func (t *ReinsuranceRequestCC) submit_request(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	tiv, err := strconv.Atoi(args[3])
	if err != nil {
		// do error
	}
	agg_lim, err := strconv.Atoi(args[4])
	if err != nil {
		// do error
	}

	ieo, err := strconv.Atoi(args[7])
	if err != nil {
		// do error
	}

	bytes, err := stub.ReadCertAttribute("enrollmentId")
	if err != nil {
		logger.Error(err)
		return nil, errors.New("failed to get enrollmentId attribute")
	}
	requestor := string(bytes)
	requestees := strings.Split(args[8], ",")

	rr := ReinsuranceRequest {
		ContractType : args[0], //"liability",
		ContractSubType	: args[1], //"facultative",
		AssetType	: args[2], //"railroad",
		TotalInsuredValue : tiv, //100000000,
		AggregateLimit : agg_lim, //100000000,
		PortfolioHash	: args[5], //"2e1b1b0cb7bfce4cf47706752a234f29",
		PortfolioURL : args[6], //"http://mybucket.s3-website-us-east-1.amazonaws.com/",
		InExcessOf : ieo, //50000000,
		Status : "open",
		Requestor	: requestor, //"myusername", // TODO
		Requestees	: requestees, //[]string {"someone", "someoneelse"},
	}

	bytes, err = json.Marshal(rr)
	if err != nil {
		logger.Error(err)
		return nil, errors.New("Failed to serialize ReinsuranceRequest object")
	}

	id := "1"

	err = stub.PutState(id, bytes)
	if err != nil {
		logger.Error("err")
		return nil, errors.New("Failed to put request")
	}

	var recipients []Recipient

	for i := 0; i < len(requestees); i++ {
		recipientId := requestees[i]
		recipientContact := "foo"
		recipient := Recipient {
			RecipientId: recipientId,
			RecipientContact: recipientContact,
		}

		recipients = append(recipients, recipient)
  }

	event := RequestEvent {
		RequestId: id,
		RequestorId: requestor,
		Recipients: recipients}

	bytes, err = json.Marshal(event)
	logger.Debugf("Sending event [ %s ]", bytes)

	if err != nil {
		logger.Error(err)
		return nil, errors.New("Failed to serialize RequestEvent object")
	}

	err = stub.SetEvent("reinsurance_request_event", []byte(bytes))

	if err != nil {
		logger.Error(err)
		return nil, errors.New("Failed to set event, id : " + id)
	}

	return nil, nil
}

// ============================================================================================================================
// Main
// ============================================================================================================================
func main() {
	err := shim.Start(new(ReinsuranceRequestCC))
	if err != nil {
		fmt.Printf("Error starting ReinsuranceRequestCC: %s", err)
	}
}
