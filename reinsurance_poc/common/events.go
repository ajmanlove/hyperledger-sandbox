package common

type RequestEvent struct {
	RequestId 						string 			`json:"requestId"`
	RequestorId						string			`json:"requestorId"`
	RequestorContact			string			`json:"requestorContact"`
	Recipients						[]Recipient	`json:"recipients"`
}

type Recipient struct {
	RecipientId 					string	`json:"recipientId"`
	RecipientContact			string	`json:"recipientContact"`
}
