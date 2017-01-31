package common

import "encoding/json"

type Response interface {
	Encode() ([]byte, error)
	Decode([]byte) (Response, error)
}

type CanViewResponse struct {
	CanView bool `json:"canView"`
}

type AssetRightsResponse struct {
	Rights []AssetRight
}

func (arr *AssetRightsResponse) Encode() ([]byte, error) {
	return json.Marshal(arr)
}

func (arr *AssetRightsResponse) Decode(bytes []byte) error {
	return json.Unmarshal(bytes, &arr)
}

func (arr *AssetRightsResponse) Contains(right AssetRight) bool {
	for _, e := range arr.Rights {
		if e == right {
			return true
		}
	}
	return false
}

func BuildArr(rights []AssetRight) AssetRightsResponse {
	return AssetRightsResponse{Rights: rights}
}
