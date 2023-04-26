package sd

import (
	"math/rand"
	"net"
)


func GetServiceByRamdom(name string) (string, error) {
	serviceList, ok := Services[name]; 
	if !ok || len(serviceList) == 0 {
		return "", ServiceNotFound
	}
	index := rand.Int31n(int32(len(serviceList)))
	return net.JoinHostPort(serviceList[index].Host, serviceList[index].Port), nil
}


func GetServiceByPickFirst(name string) (string, error) {
	serviceList, ok := Services[name]; 
	if !ok || len(serviceList) == 0 {
		return "", ServiceNotFound
	}
	return net.JoinHostPort(serviceList[0].Host, serviceList[0].Port), nil
}

