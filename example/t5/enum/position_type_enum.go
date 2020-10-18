package enum

type (
	// EdgeTypeEnum age
	EdgeTypeEnum int
)

const (
	// EdgeTypeTopLeft 左上端
	EdgeTypeTopLeft EdgeTypeEnum = iota
	// EdgeTypeTop 上端
	EdgeTypeTop
	// EdgeTypeTopRight 右上端
	EdgeTypeTopRight
	// EdgeTypeRight 右端
	EdgeTypeRight
	// EdgeTypeBottomRight 右下端
	EdgeTypeBottomRight
	// EdgeTypeBottom 下端
	EdgeTypeBottom
	// EdgeTypeBottomLeft 左下端
	EdgeTypeBottomLeft
	// EdgeTypeLeft 左端
	EdgeTypeLeft
	// EdgeTypeNotEdge 端ではない
	EdgeTypeNotEdge
)
