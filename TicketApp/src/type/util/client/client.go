package client

import (
	"github.com/valyala/fasthttp"
	"time"
)

type Client struct {
	BaseUrl string
}

func (c Client) NewClient() *fasthttp.Client {
	return &fasthttp.Client{
		MaxConnsPerHost:     10,
		MaxIdleConnDuration: time.Second * 10,
		MaxConnDuration:     time.Second * 10,
		MaxConnWaitTimeout:  time.Second * 10,
		RetryIf:             nil,
	}
}

func Get() {
	clientType := Client{BaseUrl: ""}
	client := clientType.NewClient()

	req := &fasthttp.Request{
		Header:        fasthttp.RequestHeader{},
		UseHostHeader: false,
	}

	res := &fasthttp.Response{
		Header:               fasthttp.ResponseHeader{},
		ImmediateHeaderFlush: false,
		SkipBody:             false,
	}
	client.Do(req, res)
}
