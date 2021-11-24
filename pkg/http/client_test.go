package http

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestHttpReq(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		fmt.Printf("HTTP handler: %q\n", req.RequestURI)
		_, _ = resp.Write([]byte(req.RequestURI))
	}))
	defer func() { server.Close() }()
	const expect = "/hello_world"

	client := Client{
		Client: http.Client{Timeout: time.Second * 3},
	}

	req, err := http.NewRequest(http.MethodPost, server.URL+expect, nil)
	assert.NoError(t, err)

	got := client.PostRequest(req)
	body, err := io.ReadAll(got.Body)
	assert.NoError(t, err)
	assert.Equal(t, expect, string(body))
}
