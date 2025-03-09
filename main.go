package main

import "github.com/dee-d-dev/go-balancer/loadbalancers"

func main() {
	// servers.RunServers(3)
	loadbalancers.MakeLoadBalancer(3)
}