package wrr

import (
	"context"
	"reflect"
	"testing"

	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-kratos/kratos/v2/selector"
	"github.com/go-kratos/kratos/v2/selector/filter"
)

func TestWrr(t *testing.T) {
	//new出来一个selector，一定实现了apply和select，
	wrr := New()
	var nodes []selector.Node
	nodes = append(nodes, selector.NewNode(
		"http",
		"127.0.0.1:8080",
		&registry.ServiceInstance{
			ID:      "127.0.0.1:8080",
			Version: "v2.0.0",
			//配置好权重
			Metadata: map[string]string{"weight": "10"},
		}))
	nodes = append(nodes, selector.NewNode(
		"http",
		"127.0.0.1:9090",
		&registry.ServiceInstance{
			ID:      "127.0.0.1:9090",
			Version: "v2.0.0",
			//配置好权重
			Metadata: map[string]string{"weight": "20"},
		}))
	//把node信息告诉选择器
	wrr.Apply(nodes)
	var count1, count2 int
	for i := 0; i < 90; i++ {
		//不需要其他参数，就直接选择，node信息已经知道了
		n, done, err := wrr.Select(context.Background(), selector.WithNodeFilter(filter.Version("v2.0.0")))
		if err != nil {
			t.Errorf("expect no error, got %v", err)
		}
		if done == nil {
			t.Errorf("expect done callback, got nil")
		}
		if n == nil {
			t.Errorf("expect node, got nil")
		}
		done(context.Background(), selector.DoneInfo{})
		if n.Address() == "127.0.0.1:8080" {
			count1++
		} else if n.Address() == "127.0.0.1:9090" {
			count2++
		}
	}
	//按权重分配
	if !reflect.DeepEqual(count1, 30) {
		t.Errorf("expect 30, got %d", count1)
	}
	if !reflect.DeepEqual(count2, 60) {
		t.Errorf("expect 60, got %d", count2)
	}
}

func TestEmpty(t *testing.T) {
	b := &Balancer{}
	_, _, err := b.Pick(context.Background(), []selector.WeightedNode{})
	if err == nil {
		t.Errorf("expect no error, got %v", err)
	}
}
