package balance

type LbType int

const (
	LbRandom LbType = iota
	LbRoundRobin
	LbWeightedRoundRobin
	LbConsistentHash
)

func LoadBalanceFactory(lbType LbType) LoadBalance {
	switch lbType {
	case LbRandom:
		return &RandomBalance{}
	case LbConsistentHash:
		return NewConsistentHashBalance(10, nil)
	case LbRoundRobin:
		return &RoundRobinBalance{}
	case LbWeightedRoundRobin:
		return &WeightedRoundRobinBalance{}
	default:
		return &RandomBalance{}
	}
}

func LoadBalanceFactorWithConf(lbType LbType, mConf LoadBalanceConf) LoadBalance {
	//观察者模式
	switch lbType {
	case LbRandom:
		lb := &RandomBalance{}
		lb.SetConf(mConf)
		mConf.Attach(lb)
		lb.Update()
		return lb
	case LbConsistentHash:
		lb := NewConsistentHashBalance(10, nil)
		lb.SetConf(mConf)
		mConf.Attach(lb)
		lb.Update()
		return lb
	case LbRoundRobin:
		lb := &RoundRobinBalance{}
		lb.SetConf(mConf)
		mConf.Attach(lb)
		lb.Update()
		return lb
	case LbWeightedRoundRobin:
		lb := &WeightedRoundRobinBalance{}
		lb.SetConf(mConf)
		mConf.Attach(lb)
		lb.Update()
		return lb
	default:
		lb := &RandomBalance{}
		lb.SetConf(mConf)
		mConf.Attach(lb)
		lb.Update()
		return lb
	}
}
