package common

type AssetsRecord struct {
	AssetRights map[string][]string `json:"assetRights"`
	Submissions []SubmissionRecord  `json:"submissions"`
	Requests    []RequestRecord     `json:"requests"`
	Proposals   []ProposalRecord    `json:"proposals"`
	Accepted    []AcceptedProposal  `json:"accepted"`
	Rejected    []RejectedProposal  `json:"rejected"`
	Contracts   []SubmissionRecord  `json:"contracts"`
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
	ContractText string   `json:"contractText"`
	Created      uint64   `json:"created"`
	Updated      uint64   `json:"updated"`
}
