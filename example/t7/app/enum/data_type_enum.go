package enum

type (
	// DataTypeEnum データ区分
	DataTypeEnum int
)

const (
	// DataTypeSite 拠点
	DataTypeSite DataTypeEnum = iota + 1
	// DataTypeRoute 経路
	DataTypeRoute
)
