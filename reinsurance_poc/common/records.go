package common

import "encoding/json"

type Record interface {
	Encode() ([]byte, error)
	Decode([]byte) (Record, error)
	Init()
}

type AssetRecord struct {
	Rights map[string][]AssetRight `json:"assetRights"`
}

func (arr *AssetRecord) UserHasRight(enrollId string, right AssetRight) bool {
	for _, e := range arr.Rights[enrollId] {
		if e == right {
			return true
		}
	}
	return false
}

func (arr *AssetRecord) AssignUserRights(enrollId string, rights []AssetRight) {
	if arr.Rights[enrollId] == nil {
		arr.Rights[enrollId] = rights
	} else {
		for _, e := range rights {
			arr.GiveRight(enrollId, e)
		}
	}
}

func (arr *AssetRecord) GiveRight(enrollId string, right AssetRight) {
	if !arr.UserHasRight(enrollId, right) {
		arr.Rights[enrollId] = append(arr.Rights[enrollId], right)
	}
}

func (r *AssetRecord) Encode() ([]byte, error) {
	return json.Marshal(r)
}

func (r *AssetRecord) Decode(bytes []byte) error {
	return json.Unmarshal(bytes, &r)
}

func (r *AssetRecord) Init() {
	r.Rights = make(map[string][]AssetRight)
}

type UserAssetsRecord struct {
	Submissions map[string]SubmissionRecord `json:"submissions"`
	Requests    map[string]RequestRecord    `json:"requests"`
	Proposals   map[string]ProposalRecord   `json:"proposals"`
	Accepted    map[string]AcceptedProposal `json:"accepted"`
	Rejected    map[string]RejectedProposal `json:"rejected"`
	Contracts   map[string]SubmissionRecord `json:"contracts"`
}

func (r *UserAssetsRecord) Encode() ([]byte, error) {
	return json.Marshal(r)
}

func (r *UserAssetsRecord) Decode(bytes []byte) error {
	return json.Unmarshal(bytes, &r)
}

func (r *UserAssetsRecord) Init() {
	r.Submissions = make(map[string]SubmissionRecord, 0)
	r.Requests = make(map[string]RequestRecord, 0)
	r.Proposals = make(map[string]ProposalRecord, 0)
	r.Accepted = make(map[string]AcceptedProposal, 0)
	r.Rejected = make(map[string]RejectedProposal, 0)
	r.Contracts = make(map[string]SubmissionRecord, 0)
}

type SubmissionRecord struct {
	SubmissionId string   `json:"submissionId"`
	Requestees   []string `json:"requestees"`
	Created      uint64   `json:"created"`
	Updated      uint64   `json:"updated"`
}

type RequestRecord struct {
	SubmissionId string `json:"submissionId"`
	Requestor    string `json:"requestor"`
	Created      uint64 `json:"created"`
	Updated      uint64 `json:"updated"`
}

type ProposalRecord struct {
	SubmissionId string `json:"submissionId"`
	ProposalId   string `json:"proposalId"`
	Created      uint64 `json:"created"`
	Updated      uint64 `json:"updated"`
	UpdatedBy    string `json:"updatedBy"`
}

type AcceptedProposal struct {
	SubmissionId string `json:"submissionId"`
	ProposalId   string `json:"proposalId"`
	Accepted     uint64 `json:"accepted"`
}

type RejectedProposal struct {
	SubmissionId string `json:"submissionId"`
	ProposalId   string `json:"proposalId"`
	Rejected     uint64 `json:"rejected"`
}

type ReinsuranceRequest struct {
	Id           string   `json:"id"`
	PortfolioSHA string   `json:"portfolioSha"`
	PortfolioURL string   `json:"portfolioUrl"`
	Status       string   `json:"status"`
	Requestor    string   `json:"requestor"`
	Requestees   []string `json:"requestees"`
	ContractText string   `json:"contractText"` // TODO needed here?
	ISQLSchema   string   `json:"iSQLSchema"`
	ISQLVersion  string   `json:"iSQLVersion"`
	Created      uint64   `json:"created"`
	Updated      uint64   `json:"updated"`
}

func (r *ReinsuranceRequest) Encode() ([]byte, error) {
	return json.Marshal(r)
}

func (r *ReinsuranceRequest) Decode(bytes []byte) error {
	return json.Unmarshal(bytes, &r)
}

type ReinsuranceBid struct {
	Id           string `json:"id"`
	RequestId    string `json:"requestId"`
	Bidder       string `json:"bidder"`
	ContractText string `json:"contractText"`
	Created      uint64 `json:"created"`
	Updated      uint64 `json:"updated"`
	UpdatedBy    string `json:"updatedBy"`
	Status       string `json:"status"`
}

func (r *ReinsuranceBid) Init() {
	r.Id = ""
	r.RequestId = ""
	r.Bidder = ""
	r.ContractText = ""
	r.Created = 0
	r.Updated = 0
	r.UpdatedBy = ""
	r.Status = ""
}

func (r *ReinsuranceBid) Encode() ([]byte, error) {
	return json.Marshal(r)
}

func (r *ReinsuranceBid) Decode(bytes []byte) error {
	return json.Unmarshal(bytes, &r)
}
