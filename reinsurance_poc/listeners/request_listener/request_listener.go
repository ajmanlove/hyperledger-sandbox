package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/hyperledger/fabric/events/consumer"
	pb "github.com/hyperledger/fabric/protos"

	"encoding/json"

	"github.com/ajmanlove/hyperledger-sandbox/reinsurance_poc/common"

	"net/smtp"
)

type adapter struct {
	notfy              chan *pb.Event_Block
	rejected           chan *pb.Event_Rejection
	cEvent             chan *pb.Event_ChaincodeEvent
	listenToRejections bool
	chaincodeID        string
}

//GetInterestedEvents implements consumer.EventAdapter interface for registering interested events
func (a *adapter) GetInterestedEvents() ([]*pb.Interest, error) {
	if a.chaincodeID != "" {
		return []*pb.Interest{
			{EventType: pb.EventType_BLOCK},
			{EventType: pb.EventType_REJECTION},
			{EventType: pb.EventType_CHAINCODE,
				RegInfo: &pb.Interest_ChaincodeRegInfo{
					ChaincodeRegInfo: &pb.ChaincodeReg{
						ChaincodeID: a.chaincodeID,
						EventName:   ""}}}}, nil
	}
	return []*pb.Interest{{EventType: pb.EventType_BLOCK}, {EventType: pb.EventType_REJECTION}}, nil
}

//Recv implements consumer.EventAdapter interface for receiving events
func (a *adapter) Recv(msg *pb.Event) (bool, error) {
	if o, e := msg.Event.(*pb.Event_Block); e {
		a.notfy <- o
		return true, nil
	}
	if o, e := msg.Event.(*pb.Event_Rejection); e {
		if a.listenToRejections {
			a.rejected <- o
		}
		return true, nil
	}
	if o, e := msg.Event.(*pb.Event_ChaincodeEvent); e {
		a.cEvent <- o
		return true, nil
	}
	return false, fmt.Errorf("Receive unkown type event: %v", msg)
}

//Disconnected implements consumer.EventAdapter interface for disconnecting
func (a *adapter) Disconnected(err error) {
	fmt.Printf("Disconnected...exiting\n")
	os.Exit(1)
}

func createEventClient(eventAddress string, listenToRejections bool, cid string) *adapter {
	var obcEHClient *consumer.EventsClient

	done := make(chan *pb.Event_Block)
	reject := make(chan *pb.Event_Rejection)
	adapter := &adapter{notfy: done, rejected: reject, listenToRejections: listenToRejections, chaincodeID: cid, cEvent: make(chan *pb.Event_ChaincodeEvent)}
	obcEHClient, _ = consumer.NewEventsClient(eventAddress, 5, adapter)
	if err := obcEHClient.Start(); err != nil {
		fmt.Printf("could not start chat %s\n", err)
		obcEHClient.Stop()
		return nil
	}

	return adapter
}

var emailTemplate = `

	Hello %s
	You have a new reinsurance submission request!

	Requestor: %s
	Requestor Email: %s
	Request Id: %s

	Thanks
`

func main() {
	var eventAddress string
	var listenToRejections bool
	var chaincodeID string
	var senderEmail string
	var senderPass string
	flag.StringVar(&eventAddress, "events-address", "0.0.0.0:7053", "address of events server")
	flag.BoolVar(&listenToRejections, "listen-to-rejections", false, "whether to listen to rejection events")
	flag.StringVar(&chaincodeID, "events-from-chaincode", "", "listen to events from given chaincode")
	flag.StringVar(&senderEmail, "sender-email", "", "email address of the smtp sender")
	flag.StringVar(&senderPass, "sender-password", "", "email password of the smtp sender")

	flag.Parse()

	fmt.Printf("Event Address: %s\n", eventAddress)

	a := createEventClient(eventAddress, listenToRejections, chaincodeID)
	if a == nil {
		fmt.Printf("Error creating event client\n")
		return
	}

	//Set up authentication information.
  auth := smtp.PlainAuth(
      "",
      senderEmail,
      senderPass,
      "smtp.gmail.com",
  )

	for {
		select {
		case b := <-a.notfy:
			fmt.Printf("\n")
			fmt.Printf("\n")
			fmt.Printf("Received block\n")
			fmt.Printf("--------------\n")
			for _, r := range b.Block.Transactions {
				fmt.Printf("Transaction:\n\t[%v]\n", r)
			}
		case r := <-a.rejected:
			fmt.Printf("\n")
			fmt.Printf("\n")
			fmt.Printf("Received rejected transaction\n")
			fmt.Printf("--------------\n")
			fmt.Printf("Transaction error:\n%s\t%s\n", r.Rejection.Tx.Txid, r.Rejection.ErrorMsg)
		case ce := <-a.cEvent:
			fmt.Printf("\n")
			fmt.Printf("\n")
			fmt.Printf("Received chaincode event\n")
			fmt.Printf("------------------------\n")
			fmt.Printf("Chaincode Event:%v\n", ce)

			eventName := string(ce.ChaincodeEvent.EventName)
			switch eventName {
				case "reinsurance_request_event":
					fmt.Printf("GOT event name : %s \n", eventName)

					var rEvent common.RequestEvent
					err := json.Unmarshal(ce.ChaincodeEvent.Payload, &rEvent)
					if err != nil {
						fmt.Printf("FAILED to unmarshal payload due to : %s\n", err)
					}

					for _, r := range rEvent.Recipients {
						fmt.Printf("Sending email to recipient %s (%s)\n", r.RecipientId, r.RecipientContact)
						s := fmt.Sprintf(emailTemplate,
							r.RecipientId,
							rEvent.RequestorId,
							rEvent.RequestorContact,
							rEvent.RequestId,
						)

						err = smtp.SendMail(
						    "smtp.gmail.com:587",
						    auth,
						    rEvent.RequestorContact,
						    []string{r.RecipientContact},
						    []byte(s),
						)

						if err != nil {
								fmt.Printf("FAILED to send email due to : %s\n", err)
						}

					}

				default:
					fmt.Printf("Unrecognized event name : %s \n", eventName)
			}
		}
	}
}
