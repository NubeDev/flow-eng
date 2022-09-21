package rubixcli

import (
	"context"
	"fmt"
	"github.com/NubeDev/flow-eng/services/clients/rubixcli/nresty"
	"github.com/go-resty/resty/v2"
	"net"
	"net/http"
	"sync"
	"time"
)

var (
	mutex       = &sync.RWMutex{}
	flowClients = map[string]*Client{}
)

type Client struct {
	client *resty.Client
}

// The dialTimeout normally catches: when the server is unreachable and returns i/o timeout within 2 seconds.
// Otherwise, the i/o timeout takes 1.3 minutes on default; which is a very long time for waiting.
// It uses the DialTimeout function of the net package which connects to a server address on a named network before
// a specified timeout.
func dialTimeout(_ context.Context, network, addr string) (net.Conn, error) {
	timeout := 2 * time.Second
	return net.DialTimeout(network, addr, timeout)
}

var transport = http.Transport{
	DialContext: dialTimeout,
}

type Connection struct {
	Ip   string
	Port int
}

func New(conn *Connection) *Client {
	mutex.Lock()
	defer mutex.Unlock()
	ip := conn.Ip
	port := conn.Port
	if ip == "" {
		ip = "0.0.0.0"
	}
	if port == 0 {
		port = 5001
	}

	url := fmt.Sprintf("%s://%s:%d", getSchema(port), ip, port)
	if flowClient, found := flowClients[url]; found {
		return flowClient
	}
	client := resty.New()
	client.SetDebug(false)
	client.SetBaseURL(url)
	client.SetError(&nresty.Error{})
	client.SetTransport(&transport)
	flowClient := &Client{client: client}
	flowClients[url] = flowClient
	return flowClient
}

func getSchema(port int) string {
	if port == 443 {
		return "https"
	}
	return "http"
}
