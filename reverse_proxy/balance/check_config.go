package balance

import (
	"fmt"
	"net"
	"reflect"
	"sort"
	"time"
)

const (
	// DefaultCheckMethod default check setting
	DefaultCheckMethod    = 0
	DefaultCheckTimeout   = 5
	DefaultCheckMaxErrNum = 2
	DefaultCheckInterval  = 5
)

func NewBalanceCheckConf(format string, conf map[string]string) (*LoadBalanceCheckConf, error) {
	var activeList []string
	for ip, _ := range conf {
		activeList = append(activeList, ip)
	}
	mConf := &LoadBalanceCheckConf{format: format, activeList: activeList, confIpWeight: conf}
	mConf.WatchConf()
	return mConf, nil
}

type LoadBalanceCheckConf struct {
	observers    []Observer
	confIpWeight map[string]string
	activeList   []string
	format       string
}

func (c *LoadBalanceCheckConf) Attach(o Observer) {
	c.observers = append(c.observers, o)
}

func (c *LoadBalanceCheckConf) NotifyAllObservers() {
	for _, obs := range c.observers {
		obs.Update()
	}
}

func (c *LoadBalanceCheckConf) GetConf() []string {
	var confList []string
	for _, ip := range c.activeList {
		weight, ok := c.confIpWeight[ip]
		if !ok {
			weight = "50" //default to 50
		}
		confList = append(confList, fmt.Sprintf(c.format, ip)+`,`+weight)
	}
	return confList
}

//更新配置时，通知监听者也更新
//主要是用来探测各个IP是否还存活着，重连有次数限制
func (c *LoadBalanceCheckConf) WatchConf() {
	go func() {
		confIpErrNum := make(map[string]int)
		for {
			var changedList []string
			for ip := range c.confIpWeight {
				conn, err := net.DialTimeout(
					"tcp", ip, time.Duration(DefaultCheckTimeout)*time.Second)
				if err == nil {
					_ = conn.Close()
					if _, ok := confIpErrNum[ip]; ok {
						confIpErrNum[ip] = 0
					}
				}
				if err != nil {
					if _, ok := confIpErrNum[ip]; ok {
						confIpErrNum[ip] += 1
					} else {
						confIpErrNum[ip] = 1
					}
				}
				if confIpErrNum[ip] < DefaultCheckMaxErrNum {
					changedList = append(changedList, ip)
				}
			}
			sort.Strings(changedList)
			sort.Strings(c.activeList)
			if !reflect.DeepEqual(changedList, c.activeList) {
				c.UpdateConf(changedList)
			}
			time.Sleep(time.Duration(DefaultCheckInterval) * time.Second)
		}
	}()
}

//更新配置时，通知监听者也更新
func (c *LoadBalanceCheckConf) UpdateConf(conf []string) {
	c.activeList = conf
	for _, obs := range c.observers {
		obs.Update()
	}
}
