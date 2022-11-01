package service

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/joexu01/ingress-gateway/public"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"sync"
)

var ManagerHandler *Manager

func init() {
	ManagerHandler = NewServiceManager()
}

type Manager struct {
	ServiceMap   map[string]*Detail
	ServiceSlice []*Detail
	Mutex        sync.Mutex
	init         sync.Once
	err          error
}

func NewServiceManager() *Manager {
	return &Manager{
		ServiceMap:   map[string]*Detail{},
		ServiceSlice: []*Detail{},
		Mutex:        sync.Mutex{},
		init:         sync.Once{},
		err:          nil,
	}
}

type Detail struct {
	Info          Info          `json:"info" description:"基本信息"`
	HTTPRule      HttpRule      `json:"http_rule" description:"http_rule"`
	LoadBalance   LoadBalance   `json:"load_balance" description:"balance"`
	AccessControl AccessControl `json:"access_control" description:"access_control"`
}

func (m *Manager) LoadOnce() error {
	m.init.Do(
		func() {
			list := m.GetAllServices()
			m.Mutex.Lock()
			defer m.Mutex.Unlock()
			for _, item := range list {
				temp := item
				m.ServiceMap[temp.Info.ServiceName] = &temp
				m.ServiceSlice = append(m.ServiceSlice, &temp)
			}
			log.Printf("加载的反向代理服务有 %+v\n", *m.ServiceSlice[0])
		},
	)
	return m.err
}

func (m *Manager) HTTPAccessMode(c *gin.Context) (*Detail, error) {
	//1. 前缀匹配    /abc ==> serviceSlice.rule
	//2. 域名匹配    /www.test.com ==> serviceSlice.rule
	log.Printf("HTTPAccessMode: request - %s\n", c.Request.Host)
	//host c.Request.Host
	//path c.Request.Url.Path
	host := c.Request.Host
	host = host[0:strings.Index(host, `:`)]
	//log.Println("host:", host)

	path := c.Request.URL.Path

	for _, serviceItem := range m.ServiceSlice {
		if serviceItem.Info.LoadType != public.LoadTypeHTTP {
			continue
		}
		if serviceItem.HTTPRule.RuleType == public.HTTPRuleTypeDomain {
			if serviceItem.HTTPRule.Rule == host {
				return serviceItem, nil
			}
		}
		if serviceItem.HTTPRule.RuleType == public.HTTPRuleTypePrefixURL {
			if strings.HasPrefix(path, serviceItem.HTTPRule.Rule) {
				return serviceItem, nil
			}
		}
	}
	return nil, errors.New("no service matched")
}

func (m *Manager) GetAllServices() []Detail {
	fi, _ := os.Open(public.ServiceRuleFile)
	bytes, err := ioutil.ReadAll(fi)
	if err != nil {
		log.Fatal(err)
	}

	var details []Detail

	_ = json.Unmarshal(bytes, &details)

	return details
}
