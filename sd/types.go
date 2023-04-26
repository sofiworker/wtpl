package sd

import "fmt"


var (
	ServiceNotFound = fmt.Errorf("the service not found")
)

type ServiceInstance struct {
	Name string
	Host string
	Port string
}

var Services map[string][]*ServiceInstance

func init() {
	Services = make(map[string][]*ServiceInstance, 0)
}

