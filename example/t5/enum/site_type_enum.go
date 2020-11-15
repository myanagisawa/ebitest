package enum

type (
	// SiteTypeEnum age
	SiteTypeEnum int
)

const (
	// SiteTypeBase 拠点
	SiteTypeBase SiteTypeEnum = iota
	// SiteTypeCity 都市
	SiteTypeCity
	// SiteTypeFort 砦
	SiteTypeFort
	// SiteTypePoint ポイント
	SiteTypePoint
	// SiteTypeRemains 遺跡
	SiteTypeRemains
)
