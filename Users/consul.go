package main

import (
	"fmt"
	"os"

	consulapi "github.com/hashicorp/consul/api"
)

func consulClient() (*consulapi.Client, error) {
	conf := consulapi.DefaultConfig()
	conf.Address = "host.docker.internal:8500"
	client, err := consulapi.NewClient(conf)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func StartConsul(name string, port int) {
	c, err := consulClient()
	if err != nil {
		panic(err)
	}

	registration := new(consulapi.AgentServiceRegistration)
	registration.ID = name
	registration.Name = name
	addr, _ := os.Hostname()
	registration.Address = addr
	registration.Port = port

	registration.Check = new(consulapi.AgentServiceCheck)
	registration.Check.HTTP = fmt.Sprintf("http://%s:%v/healthcheck", addr, port)
	registration.Check.Interval = "5s"
	registration.Check.Timeout = "3s"
	c.Agent().ServiceRegister(registration)
}
