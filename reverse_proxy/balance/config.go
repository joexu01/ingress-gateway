package balance

// LoadBalanceConf 配置主体
type LoadBalanceConf interface {
	Attach(o Observer)
	GetConf() []string
	WatchConf()
	UpdateConf(conf []string)
}

type Observer interface {
	Update()
}
