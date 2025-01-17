package service

import (
	"fmt"
	"github.com/joexu01/ingress-gateway/public"
	"github.com/joexu01/ingress-gateway/reverse_proxy/balance"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
)

var (
	LoadBalanceHandler *LoadBalancer
	TransporterHandler *Transporter
)

func init() {
	LoadBalanceHandler = NewLoadBalancer()
	TransporterHandler = NewTransporter()
}

type LoadBalance struct {
	CheckMethod   int    `json:"check_method" gorm:"column:check_method" description:"检查方法 tcpchk=检测端口是否握手成功	"`
	CheckTimeout  int    `json:"check_timeout" gorm:"column:check_timeout" description:"check超时时间	"`
	CheckInterval int    `json:"check_interval" gorm:"column:check_interval" description:"检查间隔, 单位s		"`
	RoundType     int    `json:"round_type" gorm:"column:round_type" description:"轮询方式 round/weight_round/random/ip_hash"`
	IpList        string `json:"ip_list" gorm:"column:ip_list" description:"ip列表"`
	WeightList    string `json:"weight_list" gorm:"column:weight_list" description:"权重列表"`
	ForbidList    string `json:"forbid_list" gorm:"column:forbid_list" description:"禁用ip列表"`

	UpstreamConnectTimeout int `json:"upstream_connect_timeout" gorm:"column:upstream_connect_timeout" description:"下游建立连接超时, 单位s"`
	UpstreamHeaderTimeout  int `json:"upstream_header_timeout" gorm:"column:upstream_header_timeout" description:"下游获取header超时, 单位s	"`
	UpstreamIdleTimeout    int `json:"upstream_idle_timeout" gorm:"column:upstream_idle_timeout" description:"下游链接最大空闲时间, 单位s	"`
	UpstreamMaxIdle        int `json:"upstream_max_idle" gorm:"column:upstream_max_idle" description:"下游最大空闲链接数"`
}

func NewLoadBalance() *LoadBalance {
	return &LoadBalance{
		CheckMethod:            0,
		CheckTimeout:           0,
		CheckInterval:          0,
		RoundType:              0,
		IpList:                 "",
		WeightList:             "",
		ForbidList:             "",
		UpstreamConnectTimeout: 0,
		UpstreamHeaderTimeout:  0,
		UpstreamIdleTimeout:    0,
		UpstreamMaxIdle:        0,
	}
}

func (t *LoadBalance) GetIPListByModel() []string {
	return strings.Split(t.IpList, ",")
}

func (t *LoadBalance) GetWeightListByModel() []string {
	return strings.Split(t.WeightList, ",")
}

type LoadBalancer struct {
	LoadBalanceMap   map[string]*LoadBalancerItem
	LoadBalanceSlice []*LoadBalancerItem
	Mutex            sync.RWMutex
}

type LoadBalancerItem struct {
	LoadBalance balance.LoadBalance
	ServiceName string
}

func NewLoadBalancer() *LoadBalancer {
	return &LoadBalancer{
		LoadBalanceMap:   make(map[string]*LoadBalancerItem),
		LoadBalanceSlice: []*LoadBalancerItem{},
		Mutex:            sync.RWMutex{},
	}
}

func (b *LoadBalancer) GetLoadBalancer(service *Detail) (balance.LoadBalance, error) {
	for _, lbrItem := range b.LoadBalanceSlice {
		if lbrItem.ServiceName == service.Info.ServiceName {
			return lbrItem.LoadBalance, nil
		}
	}

	schema := "http://"
	if service.HTTPRule.NeedHttps {
		schema = "https://"
	}
	if service.Info.LoadType == public.LoadTypeTCP || service.Info.LoadType == public.LoadTypeGRPC {
		schema = ""
	}

	//prefix := ""
	//if service.HTTPRule.RuleType == public.HTTPRuleTypePrefixURL {
	//	prefix = service.HTTPRule.Rule
	//}

	ipList := service.LoadBalance.GetIPListByModel()
	weightList := service.LoadBalance.GetWeightListByModel()
	ipConf := make(map[string]string)

	//fmt.Println("ipConf", ipConf)

	for index, ip := range ipList {
		ipConf[ip] = weightList[index]
	}
	//fmt.Println("load balance: ipConf", ipConf)

	mConf, err := balance.NewBalanceCheckConf(
		fmt.Sprintf("%s%s", schema, "%s"), ipConf)
	if err != nil {
		return nil, err
	}
	lb := balance.LoadBalanceFactorWithConf(
		balance.LbType(service.LoadBalance.RoundType), mConf)

	//save to map and slice
	lbItem := &LoadBalancerItem{
		LoadBalance: lb,
		ServiceName: service.Info.ServiceName,
	}
	b.LoadBalanceSlice = append(b.LoadBalanceSlice, lbItem)

	b.Mutex.Lock()
	defer b.Mutex.Unlock()
	b.LoadBalanceMap[service.Info.ServiceName] = lbItem
	return lb, nil
}

type Transporter struct {
	TransportMap   map[string]*TransportItem
	TransportSlice []*TransportItem
	Mutex          sync.RWMutex
}

type TransportItem struct {
	Trans       *http.Transport
	ServiceName string
}

func NewTransporter() *Transporter {
	return &Transporter{
		TransportMap:   make(map[string]*TransportItem),
		TransportSlice: []*TransportItem{},
		Mutex:          sync.RWMutex{},
	}
}

func (t *Transporter) GetTrans(service *Detail) (*http.Transport, error) {
	for _, transItem := range t.TransportSlice {
		if transItem.ServiceName == service.Info.ServiceName {
			return transItem.Trans, nil
		}
	}

	if service.LoadBalance.UpstreamConnectTimeout == 0 {
		service.LoadBalance.UpstreamConnectTimeout = 30
	}
	if service.LoadBalance.UpstreamMaxIdle == 0 {
		service.LoadBalance.UpstreamMaxIdle = 100
	}
	if service.LoadBalance.UpstreamIdleTimeout == 0 {
		service.LoadBalance.UpstreamIdleTimeout = 90
	}
	if service.LoadBalance.UpstreamHeaderTimeout == 0 {
		service.LoadBalance.UpstreamHeaderTimeout = 30
	}
	trans := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   time.Duration(service.LoadBalance.UpstreamConnectTimeout) * time.Second,
			KeepAlive: 30 * time.Second,
			//DualStack: true,
		}).DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          service.LoadBalance.UpstreamMaxIdle,
		IdleConnTimeout:       time.Duration(service.LoadBalance.UpstreamIdleTimeout) * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ResponseHeaderTimeout: time.Duration(service.LoadBalance.UpstreamHeaderTimeout) * time.Second,
	}

	//save to map and slice
	transItem := &TransportItem{
		Trans:       trans,
		ServiceName: service.Info.ServiceName,
	}
	t.TransportSlice = append(t.TransportSlice, transItem)
	t.Mutex.Lock()
	defer t.Mutex.Unlock()
	t.TransportMap[service.Info.ServiceName] = transItem
	return trans, nil
}
