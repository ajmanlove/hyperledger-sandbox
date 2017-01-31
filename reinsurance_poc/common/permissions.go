package common

type AssetRight int32

const (
	AOWNER  AssetRight = iota
	AVIEWER AssetRight = iota
	AHOLDER AssetRight = iota
)
