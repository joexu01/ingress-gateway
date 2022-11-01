package service

type AccessControl struct {
	Available         bool   `json:"available" gorm:"column:open_auth" description:"是否开启权限 1=开启"`
	BlackList         string `json:"black_list" gorm:"column:black_list" description:"黑名单ip	"`
	WhiteList         string `json:"white_list" gorm:"column:white_list" description:"白名单ip	"`
	ClientIPFlowLimit int    `json:"client_ip_flow_limit" gorm:"column:client_ip_flow_limit" description:"客户端ip限流	"`
}

func NewAccessControl() *AccessControl {
	return &AccessControl{
		Available:         true,
		BlackList:         "",
		WhiteList:         "",
		ClientIPFlowLimit: 0,
	}
}
