package service

type HttpRule struct {
	RuleType     string `json:"rule_type" gorm:"column:rule_type" description:"匹配类型 domain=域名, url_prefix=url前缀"`
	Rule         string `json:"rule" gorm:"column:rule" description:"type=domain表示域名，type=url_prefix时表示url前缀"`
	NeedHttps    bool   `json:"need_https" gorm:"column:need_https" description:"type=支持https 1=支持"`
	NeedStripUri bool   `json:"need_strip_uri" gorm:"column:need_strip_uri" description:"启用strip_uri 1=启用"`
}

func NewHttpRule() *HttpRule {
	return &HttpRule{
		RuleType:     "domain",
		Rule:         "",
		NeedHttps:    false,
		NeedStripUri: true,
	}
}
