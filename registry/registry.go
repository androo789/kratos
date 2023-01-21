package registry

import (
	"context"
	"fmt"
	"sort"
)

// Registrar is service registrar.
// 就是我把自己上报上去
type Registrar interface {
	// Register the registration.
	Register(ctx context.Context, service *ServiceInstance) error
	// Deregister the registration.
	// 反注册就是下线
	Deregister(ctx context.Context, service *ServiceInstance) error
}

// Discovery is service discovery.
// 就是获取其他实例的情况
type Discovery interface {
	// GetService return the service instances in memory according to the service name.
	GetService(ctx context.Context, serviceName string) ([]*ServiceInstance, error)
	// Watch creates a watcher according to the service name.
	Watch(ctx context.Context, serviceName string) (Watcher, error)
}

// Watcher is service watcher.
type Watcher interface {
	// Next returns services in the following two cases:
	// 1.the first time to watch and the service instance list is not empty.
	// 2.any service instance changes found.
	// if the above two conditions are not met, it will block until context deadline exceeded or canceled
	Next() ([]*ServiceInstance, error)
	// Stop close the watcher.
	Stop() error
}

// ServiceInstance is an instance of a service in a discovery system.
// 一个基础的实例信息包括以下数据，
type ServiceInstance struct {
	// ID is the unique instance ID as registered.
	ID string `json:"id"`
	// Name is the service name as registered.
	Name string `json:"name"`
	// Version is the version of the compiled.
	Version string `json:"version"`
	// Metadata is the kv pair metadata associated with the service instance.
	Metadata map[string]string `json:"metadata"`
	// Endpoints are endpoint addresses of the service instance.
	// 之所以是数组是因为不同的协议
	// schema:
	//   http://127.0.0.1:8000?isSecure=false
	//   grpc://127.0.0.1:9000?isSecure=false
	Endpoints []string `json:"endpoints"`
}

// 参考公司的例子，确实就是项目名字+唯一ID的组合
func (i *ServiceInstance) String() string {
	return fmt.Sprintf("%s-%s", i.Name, i.ID)
}

// Equal returns whether i and o are equivalent.
func (i *ServiceInstance) Equal(o interface{}) bool {
	if i == nil && o == nil {
		return true
	}

	if i == nil || o == nil {
		return false
	}

	t, ok := o.(*ServiceInstance)
	if !ok {
		return false
	}

	if len(i.Endpoints) != len(t.Endpoints) {
		return false
	}

	sort.Strings(i.Endpoints)
	sort.Strings(t.Endpoints)
	for j := 0; j < len(i.Endpoints); j++ {
		if i.Endpoints[j] != t.Endpoints[j] {
			return false
		}
	}

	if len(i.Metadata) != len(t.Metadata) {
		return false
	}

	for k, v := range i.Metadata {
		if v != t.Metadata[k] {
			return false
		}
	}

	return i.ID == t.ID && i.Name == t.Name && i.Version == t.Version
}
