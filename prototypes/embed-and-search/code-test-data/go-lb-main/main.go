package main

func main() {
	config := LBConfig{ListenAddr: "localhost:3000", ClientExpireDuration: 15, Strategy: RoundRobin, AdminListenAddr: "localhost:8080"}

	lb := NewLoadBalancer(config)
	// lb.Backends = append(lb.Backends, &Backend{Addr: "httpbin.org", Id: 3})

	lb.Start()

}
