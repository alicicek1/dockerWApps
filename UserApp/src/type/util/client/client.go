package client

import (
	"UserApp/src/type/util"
	"encoding/json"
	"github.com/valyala/fasthttp"
	"net/http"
	"time"
)

type Client struct {
	BaseUrl     string
	BaseHeaders map[string]string
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

func (c Client) Get(additionalPath string, headers map[string]string) (string, *util.Error) {
	client := c.NewClient()

	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)

	req.SetRequestURI(c.BaseUrl + additionalPath)
	if len(c.BaseHeaders) > 0 {
		for key := range c.BaseHeaders {
			req.Header.Add(key, c.BaseHeaders[key])
		}
	}

	for key := range headers {
		req.Header.Add(key, headers[key])
	}

	err := client.Do(req, resp)
	if err != nil {
		return "false", util.NewError("Client", "ExistById", err.Error(), resp.StatusCode(), 6050)
	}

	bodyBytes := resp.Body()
	bodyStr := string(bodyBytes)

	if resp.StatusCode() == http.StatusOK {
		return bodyStr, nil
	}

	error := util.Error{}
	err = json.Unmarshal(bodyBytes, &error)
	if err != nil {
		return "false", util.NewError("Client", "Unmarshalling", err.Error(), http.StatusBadRequest, 6051)
	}

	if error == (util.Error{}) {
		return "false", util.NewError("Client", "ClientError", string(resp.Body()), resp.StatusCode(), 6051)
	}

	return "false", util.NewError("Client", "ExistById", error.Description, error.StatusCode, 6051)
}
