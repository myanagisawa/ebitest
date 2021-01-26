package enum

type (
	// SiteTypeEnum age
	SiteTypeEnum int
)

const (
	// SiteTypeBase 拠点
	SiteTypeBase SiteTypeEnum = iota + 1
	// SiteTypeCity 都市
	SiteTypeCity
	// SiteTypeFort 砦
	SiteTypeFort
	// SiteTypePoint ポイント
	SiteTypePoint
	// SiteTypeRemains 遺跡
	SiteTypeRemains
)

// Name FontStyleEnumの区分値表現を返します
func (e SiteTypeEnum) Name() string {
	switch e {
	case SiteTypeBase:
		return "拠点"
	case SiteTypeCity:
		return "都市"
	case SiteTypeFort:
		return "砦"
	case SiteTypePoint:
		return "ポイント"
	case SiteTypeRemains:
		return "遺跡"
	default:
		return "不明"
	}
}
