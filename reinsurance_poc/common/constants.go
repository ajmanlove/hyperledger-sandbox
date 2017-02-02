package common

type AssetRight int32

const (
	AOWNER  AssetRight = 0
	AVIEWER AssetRight = 1
	AHOLDER AssetRight = 2
)

const (
	request_cc_id  = "reinsurance_request"
	proposal_cc_id = "reinsurance_proposal"
)
