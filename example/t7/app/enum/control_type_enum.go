package enum

type (
	// ControlTypeEnum コントロールの種類
	ControlTypeEnum int
)

const (
	// ControlTypeDefault デフォルト
	ControlTypeDefault ControlTypeEnum = iota
	// ControlTypeFrame フレーム
	ControlTypeFrame
	// ControlTypeLayer レイヤ
	ControlTypeLayer
)
