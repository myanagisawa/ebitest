package enum

type (
	// RouteTypeEnum age
	RouteTypeEnum int
)

const (
	// RouteTypePlain 平地
	RouteTypePlain RouteTypeEnum = iota
	// RouteTypeMountain 山道
	RouteTypeMountain
	// RouteTypeSea 海路
	RouteTypeSea
)
