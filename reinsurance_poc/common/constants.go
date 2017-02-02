package common

type AssetRight int32

const (
	AOWNER    AssetRight = 0
	AVIEWER   AssetRight = 1
	AAPPROVAL AssetRight = 2
)

const (
	request_cc_id  = "reinsurance_request"
	proposal_cc_id = "reinsurance_proposal"
)

// Chaincode args
const (
	INIT_ARG              = "init"
	AM_REGISTER_CC_ARG    = "register_chaincode"
	AM_NEW_REQ_ARG        = "new_request"
	AM_NEW_BID_ARG        = "new_proposal"
	AM_GET_CC_NAME_ARG    = "get_cc_name"
	AM_GET_U_ASST_ARG     = "get_user_assets"
	AM_GET_AST_RIGHTS_ARG = "get_asset_rights"

	RR_SUBMIT_ARG  = "submit"
	RR_GET_REQ_ARG = "get_request"

	RP_PROPOSE_ARG = "propose"
	RP_COUNTER_ARG = "counter"
	RP_ACCEPT_ARG  = "accept"
	RP_REJECT_ARG  = "reject"
	RP_GET_BID_ARG = "get_proposal"
)
