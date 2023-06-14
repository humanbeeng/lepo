# Go-LB

This is a toy loadbalancer that I built to better understand how Loadbalancers and Service Discovery works.

### Features

- RoundRobin Load Distribution
- Registration and Deregistration of backends through Service Discovery
- Cancellation propogation and timeouts(Configurable duration)
- Heartbeat and auto-deregistration during backend disconnection
- Supports all HTTP methods
- Handles requests concurrently from single and multiple clients



### Usage

#### Start Loadbalancer
```sh
$ make run
```

#### To enable Service Discovery
1. Import client library
```sh
$ go get github.com/humanbeeng/go-lb/client
```

2. Initialise client and pass the loadbalancer's address(lbAddr) and the server's listenaddress(*addr)
```go
conn, _ := net.Dial("tcp", lbAddr)
c := client.Client{Conn: conn, Addr: *addr}
c.Register()
```



TODO:

- [ ] URL pattern matched balancing
- [ ] Configurable strategy
- [ ] Hot reload config changes
- [ ] Export metrics about healthy servers
- [ ] CLI to add and remove backends and change strategy in realtime
