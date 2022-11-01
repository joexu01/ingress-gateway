package service

import (
	"time"
)

type Info struct {
	ID          int64     `json:"id" gorm:"primary_key"`
	LoadType    int       `json:"load_type" gorm:"column:load_type" description:"负载类型 0=http 1=tcp 2=grpc"`
	ServiceName string    `json:"service_name" gorm:"column:service_name" description:"服务名称"`
	ServiceDesc string    `json:"service_desc" gorm:"column:service_desc" description:"服务描述"`
	CreatedAt   time.Time `json:"create_at" gorm:"column:create_at" description:"更新时间"`
}

func NewInfo() *Info {
	return &Info{
		ID:          0,
		LoadType:    0,
		ServiceName: "",
		ServiceDesc: "",
		CreatedAt:   time.Time{},
	}
}

//func (i *Info) GetServiceInfoList() (total int, serviceList []Info, err error) {
//	fi, _ := os.Open(public.ServiceRuleFile)
//	bytes, err := ioutil.ReadAll(fi)
//	if err != nil {
//		log.Fatal(err)
//	}
//	return 0, nil, nil
//}
