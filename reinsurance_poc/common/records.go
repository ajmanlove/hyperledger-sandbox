package common

import "encoding/json"

type Record interface {
	Encode() ([]byte, error)
	Decode([]byte) (Record, error)
	Init()
}

type AssetsRecord struct {
	AssetRights map[string][]AssetRight `json:"assetRights"`
	Submissions []SubmissionRecord      `json:"submissions"`
	Requests    []RequestRecord         `json:"requests"`
	Proposals   []ProposalRecord        `json:"proposals"`
	Accepted    []AcceptedProposal      `json:"accepted"`
	Rejected    []RejectedProposal      `json:"rejected"`
	Contracts   []SubmissionRecord      `json:"contracts"`
}

func (arr *AssetsRecord) ContainsRight(assetId string, right AssetRight) bool {
	for _, e := range arr.AssetRights[assetId] {
		if e == right {
			return true
		}
	}
	return false
}

func (arr *AssetsRecord) GiveRights(assetId string, rights []AssetRight) {
	if arr.AssetRights[assetId] == nil {
		arr.AssetRights[assetId] = rights
	} else {
		for _, e := range rights {
			arr.GiveRight(assetId, e)
		}
	}
}

func (arr *AssetsRecord) GiveRight(assetId string, right AssetRight) {
	if arr.ContainsRight(assetId, right) {
		arr.AssetRights[assetId] = append(arr.AssetRights[assetId], right)
	}
}

func (r *AssetsRecord) Encode() ([]byte, error) {
	return json.Marshal(r)
}

func (r *AssetsRecord) Decode(bytes []byte) error {
	return json.Unmarshal(bytes, &r)
}

func (r *AssetsRecord) Init() {
	r.AssetRights = make(map[string][]AssetRight)
	r.Submissions = make([]SubmissionRecord, 0)
	r.Requests = make([]RequestRecord, 0)
	r.Proposals = make([]ProposalRecord, 0)
	r.Accepted = make([]AcceptedProposal, 0)
	r.Rejected = make([]RejectedProposal, 0)
	r.Contracts = make([]SubmissionRecord, 0)
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
	SubmissionId []string `json:"submissionId"`
	ProposalId   []string `json:"proposalId"`
	Accepted     uint64   `json:"accepted"`
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
