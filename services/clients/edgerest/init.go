package edgerest

import (
	"fmt"
	"github.com/NubeDev/flow-eng/services/clients/ffclient/nresty"
	"github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"
)

// Client is used to invoke Form3 Accounts API.
type Client struct {
	client      *resty.Client
	ClientToken string
}

func New(address string, port int) *Client {
	if address == "" {
		address = "0.0.0.0"
	}
	if port == 0 {
		port = 5000
	}
	client := resty.New()
	client.SetDebug(false)
	url := fmt.Sprintf("http://%s:%d", address, port)
	apiURL := url
	client.SetBaseURL(apiURL)
	client.SetError(&nresty.Error{})
	client.SetHeader("Content-Type", "application/json")
	return &Client{client: client}
}

func (*Client) edge28ClientDebugMsg(args ...interface{}) {
	enable := false
	if enable {
		prefix := "Edge28 Client: "
		log.Info(prefix, args)
	}
}

func (*Client) edge28ClientErrorMsg(args ...interface{}) {
	prefix := "Edge28 Client: "
	log.Error(prefix, args)
}
