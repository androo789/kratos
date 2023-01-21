package wrr

import (
	"context"
	"sync"

	"github.com/go-kratos/kratos/v2/selector"
	"github.com/go-kratos/kratos/v2/selector/node/direct"
)

const (
	// Name is wrr balancer name
	Name = "wrr"
)

var _ selector.Balancer = (*Balancer)(nil) // Name is balancer name

// Option is random builder option.
type Option func(o *options)

// options is random builder options
type options struct{}

// Balancer is a random balancer.
// 实现了Balancer接口
type Balancer struct {
	mu sync.Mutex
	//key是节点地址，value是权重
	currentWeight map[string]float64
}

// New random a selector.
func New(opts ...Option) selector.Selector {
	return NewBuilder(opts...).Build()
}

// Pick is pick a weighted node.
// 实现了Balancer接口
func (p *Balancer) Pick(_ context.Context, nodes []selector.WeightedNode) (selector.WeightedNode, selector.DoneFunc, error) {
	if len(nodes) == 0 {
		return nil, nil, selector.ErrNoAvailable
	}
	var totalWeight float64
	var selected selector.WeightedNode
	var selectWeight float64

	// nginx wrr load balancing algorithm: http://blog.csdn.net/zhangskd/article/details/50194069
	// 具体算法实现看看
	// 有例子证明，
	p.mu.Lock()
	for _, node := range nodes {
		totalWeight += node.Weight()
		cwt := p.currentWeight[node.Address()]
		// current += effectiveWeight
		cwt += node.Weight()
		p.currentWeight[node.Address()] = cwt
		if selected == nil || selectWeight < cwt {
			selectWeight = cwt
			selected = node
		}
	}
	p.currentWeight[selected.Address()] = selectWeight - totalWeight
	p.mu.Unlock()

	d := selected.Pick()
	return selected, d, nil
}

// NewBuilder returns a selector builder with wrr balancer
// 通过builder实现了
func NewBuilder(opts ...Option) selector.Builder {
	var option options
	for _, opt := range opts {
		opt(&option)
	}
	return &selector.DefaultBuilder{
		Balancer: &Builder{},
		Node:     &direct.Builder{},
	}
}

// Builder is wrr builder
type Builder struct{}

// Build creates Balancer
func (b *Builder) Build() selector.Balancer {
	return &Balancer{currentWeight: make(map[string]float64)}
}
