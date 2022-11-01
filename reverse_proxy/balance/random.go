package balance

import (
	"errors"
	"math/rand"
	"strings"
)

type RandomBalance struct {
	currentIndex int
	addresses    []string
	//观察主体
	conf LoadBalanceConf
}

func (r *RandomBalance) Add(params ...string) error {
	if len(params) == 0 {
		return errors.New("param len 1 at least")
	}
	addr := params[0]
	r.addresses = append(r.addresses, addr)
	return nil
}

func (r *RandomBalance) Next() string {
	if len(r.addresses) == 0 {
		return ""
	}
	r.currentIndex = rand.Intn(len(r.addresses))
	return r.addresses[r.currentIndex]
}

func (r *RandomBalance) Get(key string) (string, error) {
	return r.Next(), nil
}

func (r *RandomBalance) SetConf(conf LoadBalanceConf) {
	r.conf = conf
}

func (r *RandomBalance) Update() {
	if conf, ok := r.conf.(*LoadBalanceCheckConf); ok {
		//log.Println("Update get check conf:", conf.GetConf())
		r.addresses = nil
		for _, ip := range conf.GetConf() {
			_ = r.Add(strings.Split(ip, `,`)...)
		}
	}
}
