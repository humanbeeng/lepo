package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/humanbeeng/go-lb/client"
)

type LBConfig struct {
	ListenAddr           string
	AdminListenAddr      string
	ClientExpireDuration int
	Strategy             Strategy
	Backends             []*Backend
	LastServed           int
}

type LoadBalancer struct {
	LBConfig
	servers map[*client.Client]struct{}
}

func NewLoadBalancer(config LBConfig) *LoadBalancer {
	return &LoadBalancer{
		LBConfig: config,
		servers:  make(map[*client.Client]struct{}),
	}
}

type Backend struct {
	Addr string
	Id   int
	Conn net.Conn
}

func (lb *LoadBalancer) getNextBackend() (*Backend, error) {
	numBackends := len(lb.Backends)
	if numBackends == 0 {
		return nil, fmt.Errorf("no backends available")
	}
	nextBackendIdx := (lb.LastServed + 1) % numBackends
	lb.LastServed = nextBackendIdx

	return lb.Backends[nextBackendIdx], nil
}


// This is a multiline comment
// this is just another line
func (lb *LoadBalancer) registerNewBackend(conn net.Conn) {

	cmdReg := ParseRegisterCommand(conn)
	log.Println("Register command called from", string(cmdReg.Addr))
	b := &Backend{
		Addr: string(cmdReg.Addr),
		Id:   len(lb.Backends) + 1,
		Conn: conn,
	}
	lb.Backends = append(lb.Backends, b)
	log.Printf("Registered %v\n", string(cmdReg.Addr))
}


// This is a comment
func (lb *LoadBalancer) deRegisterBackend(conn net.Conn) {

	log.Println("Deregister command called from", conn.RemoteAddr().String())
	cmdDereg := ParseDeRegisterCommand(conn)

	for i, val := range lb.Backends {
		if val.Addr == string(cmdDereg.Addr) {
			lb.Backends = append(lb.Backends[:i], lb.Backends[i+1:]...)
			break
		}
	}
	log.Printf("Deregistered %v\n", string(cmdDereg.Addr))
}

func proxy(w http.ResponseWriter, r *http.Request, b *Backend) {
	log.Printf("Redirecting %v to %v", r.URL.Path, b.Addr)
	req, err := http.NewRequest(r.Method, "http://"+b.Addr, r.Body)
	if err != nil {
		log.Println("error request", err)
	}
	ctx, cancelFunc := context.WithTimeout(req.Context(), time.Second*TIMEOUT_LIMIT_IN_SECONDS)
	defer cancelFunc()

	req.Header = r.Header
	req.URL.Path = r.URL.Path
	req = req.WithContext(ctx)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			w.WriteHeader(http.StatusGatewayTimeout)
			w.Write([]byte("Request Timedout"))
			log.Println("Request Timedout")
		}
		return
	}

	w.WriteHeader(res.StatusCode)

	// Copy all the headers from the backend's response
	for k, v := range res.Header {
		w.Header().Add(k, v[0])
	}

	io.Copy(w, res.Body)

}

func (lb *LoadBalancer) handleRequest(w http.ResponseWriter, r *http.Request) {

	b, err := lb.getNextBackend()
	if err != nil {
		log.Printf("Error: %v", err)
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte("Service Unavailable"))
		return
	}
	proxy(w, r, b)
}

func (lb *LoadBalancer) handleAdminConnection(conn net.Conn) error {
	defer conn.Close()

	for {
		cmd, err := ParseAdminCommand(conn)
		if err != nil {
			if err == io.EOF {
				conn.Close()
				return nil
			}
			return err
		}
		switch cmd {
		case CmdReg:
			lb.registerNewBackend(conn)

		case CmdDereg:
			lb.deRegisterBackend(conn)

		}

	}
}

func (lb *LoadBalancer) heartBeat(wg *sync.WaitGroup) {
	for range time.Tick(time.Second * time.Duration(lb.ClientExpireDuration)) {
		for i, b := range lb.Backends {
			pingAddr := "http://" + b.Addr + RESERVED_HEALTH_CHECK_PATH
			_, err := http.Get(pingAddr)
			if err != nil {
				log.Printf("Deregistering: %v from pool\n", b.Addr)
				lb.Backends = append(lb.Backends[:i], lb.Backends[i+1:]...)
			}
		}
	}
}

func (lb *LoadBalancer) startClientServer(wg *sync.WaitGroup) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", http.HandlerFunc(lb.handleRequest))

	go func() {
		defer wg.Done()
		log.Printf("Client facing server started to listen on %v\n", lb.ListenAddr)
		http.ListenAndServe(lb.ListenAddr, mux)
	}()
}

func (lb *LoadBalancer) startAdminServer(wg *sync.WaitGroup) {
	adminLn, err := net.Listen("tcp", lb.AdminListenAddr)
	if err != nil {
		log.Fatal(err)
	}

	defer wg.Done()
	log.Printf("Admin server started to listen on %v\n", lb.AdminListenAddr)

	for {
		adminConn, err := adminLn.Accept()
		if err != nil {
			if err == io.EOF {
				return
			}
			log.Printf("accept error %v", err)
			return
		}
		go lb.handleAdminConnection(adminConn)
	}
}

func (lb *LoadBalancer) Start() {
	wg := sync.WaitGroup{}
	wg.Add(3)

	go lb.startClientServer(&wg)

	go lb.startAdminServer(&wg)

	go lb.heartBeat(&wg)

	wg.Wait()
}
