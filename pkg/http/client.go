package http

import (
	"net/http"
	"net/http/httputil"
)

type Client struct {
	http.Client
}

func (c *Client) PostRequest(req *http.Request) *http.Response {
	_, err := httputil.DumpRequestOut(req, true)
	if err != nil {
		panic(err)
	}

	resp, err := c.Do(req)
	if err != nil {
		panic(err)
	}

	return resp
}
