package health

import (
	"expvar"
	"fmt"
	"net/http"
	"sync"
	"time"
)

var (
	ExportHealth *expvar.Map
)

func init() {
	ExportHealth = expvar.NewMap("node_health")
}

type HealthMonitor interface {
	IsHealthy(node string) bool
}

func NewHealthCheck() *Health {
	return &Health{
		opts:    DefaultHealthCheckOptions(),
		nodes:   make(map[string]*Node, 0),
		state:   make(map[string]*NodeState, 0),
		control: make(chan int, 1),
	}
}

type Health struct {
	opts    *HealthCheckOptions
	nodes   map[string]*Node
	control chan int
	mx      sync.RWMutex
	state   map[string]*NodeState
}

func DefaultHealthCheckOptions() *HealthCheckOptions {
	opts := &HealthCheckOptions{
		DefaultPort: 8181,
	}
	return opts
}

type HealthCheckOptions struct {
	DefaultPort int
}

type Node struct {
	Name  string
	Addrs []string
}

func (me *Node) AddressKey(addr string) string {
	return me.Name + "/" + addr
}

type NodeState struct {
	Healthy  bool
	HTTPCode int
}

func (me *Health) Start() {
	go me.monitorLoop()
}

func (me *Health) monitorLoop() {

	checkInterval := time.Duration(30) * time.Second
	tick := time.NewTicker(checkInterval)

	for {
		select {
		case <-tick.C:
			me.queueHealthCheck()
		case code := <-me.control:
			if code == 1 {
				go me.runHealthCheck()
			}
		}
	}

}

func (me *Health) runHealthCheck() {
	var httpCode int

	for _, node := range me.nodes {

		for _, addr := range node.Addrs {
			httpCode = 0
			err := me.testHTTP(node.Name, addr, &httpCode)
			healthy := (err == nil)

			key := node.AddressKey(addr)

			// ----- begin lock held
			me.mx.Lock()
			state, ok := me.state[key]
			if ok == false {
				state = &NodeState{}
				me.state[node.Name] = state
			}
			state.Healthy = healthy
			state.HTTPCode = httpCode
			defer me.mx.Unlock()
			// ----- end lock held
		}
	}

}

func (me *Health) testHTTP(node string, host string, httpCode *int) error {

	testUrl := fmt.Sprintf("http://%s:%d/ping", host, me.opts.DefaultPort)

	resp, err := http.Get(testUrl)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if httpCode != nil {
		*httpCode = resp.StatusCode
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("http code %d", resp.StatusCode)
	}

	return nil

}

func (me *Health) queueHealthCheck() {
	select {
	case me.control <- 1:
	default:
		log.Warnf("health check already queued. ignoring")
	}
}

func (me *Health) RegisterNode(node *Node) {
	me.nodes[node.Name] = node
}

func (me *Health) IsHealthy(nodename string) bool {

	me.mx.RLock()
	defer me.mx.RUnlock()

	node, ok := me.nodes[nodename]
	if ok == false {
		return false
	}

	healthy := false

	for _, addr := range node.Addrs {

		key := node.AddressKey(addr)

		state, ok := me.state[key]
		if ok == false {
			continue
		}

		if state.Healthy == true {
			healthy = true
			break
		}
	}
	return healthy
}
