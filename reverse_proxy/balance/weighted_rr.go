package balance

import (
	"errors"
	"strconv"
	"strings"
)

type WeightedRoundRobinBalance struct {
	curIndex    int
	nodes       []*WeightNode
	nodeWeights []int
	//观察主体
	conf LoadBalanceConf
}

type WeightNode struct {
	addr            string
	weight          int
	currentWeight   int
	effectiveWeight int
}

func (w *WeightedRoundRobinBalance) Add(params ...string) error {
	if len(params) != 2 {
		return errors.New("param len need 2")
	}
	parInt, err := strconv.ParseInt(params[1], 10, 64)
	if err != nil {
		return err
	}
	node := &WeightNode{addr: params[0], weight: int(parInt)}
	node.effectiveWeight = node.weight
	w.nodes = append(w.nodes, node)
	return nil
}

func (w *WeightedRoundRobinBalance) Get(_ string) (string, error) {
	return w.Next(), nil
}

func (w *WeightedRoundRobinBalance) Update() {
	if conf, ok := w.conf.(*LoadBalanceCheckConf); ok {
		//fmt.Println("WeightedRoundRobinBalance get check conf:", conf.GetConf())
		w.nodes = nil
		for _, ip := range conf.GetConf() {
			_ = w.Add(strings.Split(ip, ",")...)
		}
	}
}

func (w *WeightedRoundRobinBalance) SetConf(conf LoadBalanceConf) {
	w.conf = conf
}

func (w *WeightedRoundRobinBalance) Next() string {
	total := 0
	var best *WeightNode
	for i := 0; i < len(w.nodes); i++ {
		w := w.nodes[i]
		//step 1 统计所有有效权重之和
		total += w.effectiveWeight

		//step 2 变更节点临时权重为的节点临时权重+节点有效权重
		w.currentWeight += w.effectiveWeight

		//step 3 有效权重默认与权重相同，通讯异常时-1, 通讯成功+1，直到恢复到weight大小
		if w.effectiveWeight < w.weight {
			w.effectiveWeight++
		}
		//step 4 选择最大临时权重点节点
		if best == nil || w.currentWeight > best.currentWeight {
			best = w
		}
	}
	if best == nil {
		return ""
	}
	//step 5 变更临时权重为 临时权重-有效权重之和
	best.currentWeight -= total
	return best.addr
}
