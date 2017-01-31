package common

type CanViewResponse struct {
	CanView bool `json:"canView"`
}

type AssetRightsResponse struct {
	Rights []AssetRight
}
