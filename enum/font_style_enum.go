package enum

type (
	// FontStyleEnum フォント種別のEnum
	FontStyleEnum int
)

const (
	// FontStyleGenShinGothicBold 源真ゴシック-bold
	FontStyleGenShinGothicBold FontStyleEnum = iota
	// FontStyleGenShinGothicExtraLight 源真ゴシック-extra light
	FontStyleGenShinGothicExtraLight
	// FontStyleGenShinGothicLight 源真ゴシック-light
	FontStyleGenShinGothicLight
	// FontStyleGenShinGothicMedium 源真ゴシック-Medium
	FontStyleGenShinGothicMedium
	// FontStyleGenShinGothicNormal 源真ゴシック-Normal
	FontStyleGenShinGothicNormal
	// FontStyleGenShinGothicRegular 源真ゴシック-Regular
	FontStyleGenShinGothicRegular
)

// Name FontStyleEnumの区分値表現を返します
func (e FontStyleEnum) Name() string {
	switch e {
	case FontStyleGenShinGothicBold:
		return "GenShinGothic-Bold.ttf"
	case FontStyleGenShinGothicExtraLight:
		return "GenShinGothic-ExtraLight.ttf"
	case FontStyleGenShinGothicLight:
		return "GenShinGothic-Light.ttf"
	case FontStyleGenShinGothicMedium:
		return "GenShinGothic-Medium.ttf"
	case FontStyleGenShinGothicNormal:
		return "GenShinGothic-Normal.ttf"
	case FontStyleGenShinGothicRegular:
		return "GenShinGothic-Regular.ttf"
	default:
		return "不明"
	}
}
