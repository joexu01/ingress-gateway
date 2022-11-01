package balance

import (
	"errors"
	"strings"
)

type RoundRobinBalance struct {
	currentIndex int
	addresses    []string

	conf LoadBalanceConf
}

func (r *RoundRobinBalance) Add(params ...string) error {
	if len(params) == 0 {
		return errors.New("param len 1 at least")
	}
	addr := params[0]
	r.addresses = append(r.addresses, addr)
	return nil
}

func (r *RoundRobinBalance) Get(_ string) (string, error) {
	return r.Next(), nil
}

func (r *RoundRobinBalance) Update() {
	if conf, ok := r.conf.(*LoadBalanceCheckConf); ok {
		//fmt.Println("Update get check conf:", conf.GetConf())
		r.addresses = nil
		for _, ip := range conf.GetConf() {
			_ = r.Add(strings.Split(ip, ",")...)
		}
	}
}

func (r *RoundRobinBalance) Next() string {
	if len(r.addresses) == 0 {
		return ""
	}
	num := len(r.addresses)
	if r.currentIndex >= num {
		r.currentIndex = 0
	}
	currentAddr := r.addresses[r.currentIndex]
	r.currentIndex = (r.currentIndex + 1) % num
	return currentAddr
}

func (r *RoundRobinBalance) SetConf(conf LoadBalanceConf) {
	r.conf = conf
}

