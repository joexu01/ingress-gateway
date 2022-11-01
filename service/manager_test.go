package service

import (
	"fmt"
	"log"
	"testing"
)

func TestLoadServiceInfo(t *testing.T) {
	m := NewServiceManager()
	services := m.GetAllServices()

	log.Println(services[0].Info.ServiceName)

	log.Printf("%+v\n", services)
}

func TestManager_LoadOnce(t *testing.T) {
	m := NewServiceManager()
	_ = m.LoadOnce()

	fmt.Printf("%+v\n", m.ServiceMap)
}
