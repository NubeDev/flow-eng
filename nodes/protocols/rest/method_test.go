package rest

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"testing"
)

func TestHTTP_getClient(t *testing.T) {
	client := resty.New()

	resp, err := client.R().
		Get("http://0.0.0.0:1665/api/flows")
	fmt.Println(err)
	s := resp.String()
	fmt.Println(s)
	fmt.Println(resp.StatusCode())
	fmt.Println(resp.Status())

}
