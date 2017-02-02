package common

import "encoding/json"

type Response interface {
	Encode() ([]byte, error)
	Decode([]byte) (Response, error)
}

type AssetRightsResponse struct {
	Exists bool
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

func BuildArr(exists bool, rights []AssetRight) AssetRightsResponse {
	return AssetRightsResponse{Exists: exists, Rights: rights}
}

type CCNameResponse struct {
	Name string
}

func (ccn *CCNameResponse) Encode() ([]byte, error) {
	return json.Marshal(ccn)
}

func (ccn *CCNameResponse) Decode(bytes []byte) error {
	return json.Unmarshal(bytes, &ccn)
}
