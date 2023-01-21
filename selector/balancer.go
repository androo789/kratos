package selector

import (
	"context"
	"time"
)

// Balancer is balancer interface
type Balancer interface {
	//这里就是输入节点，输出节点
	Pick(ctx context.Context, nodes []WeightedNode) (selected WeightedNode, done DoneFunc, err error)
}

// BalancerBuilder build balancer
type BalancerBuilder interface {
	Build() Balancer
}

// WeightedNode calculates scheduling weight in real time
// 在节点的基础上再新加几个函数，要实现
type WeightedNode interface {
	Node

	// Raw returns the original node
	Raw() Node

	// Weight is the runtime calculated weight
	Weight() float64

	// Pick the node
	Pick() DoneFunc

	// PickElapsed is time elapsed since the latest pick
	PickElapsed() time.Duration
}

// WeightedNodeBuilder is WeightedNode Builder
type WeightedNodeBuilder interface {
	Build(Node) WeightedNode
}
