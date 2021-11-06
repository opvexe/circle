/*
Copyright 2021 The SHUMIN Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package queries

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"time"
	"x6t.io/circle"
)

var _ circle.User = &Client{}
var _ circle.Fetcher = &Client{}
var _ circle.Share = &Client{}

// Shared transports for all clients to prevent leaking connections
var skipVerifyTransport *http.Transport
var defaultTransport *http.Transport

// CreateTransport create a new transport
func CreateTransport(skipVerify bool) *http.Transport {
	var transport *http.Transport
	if cloneable, ok := http.DefaultTransport.(interface{ Clone() *http.Transport }); ok {
		transport = cloneable.Clone() // available since go1.13
	} else {
		// This uses the same values as http.DefaultTransport
		transport = &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
				DualStack: true,
			}).DialContext,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		}
	}
	if skipVerify {
		if transport.TLSClientConfig == nil {
			transport.TLSClientConfig = &tls.Config{}
		}
		transport.TLSClientConfig.InsecureSkipVerify = true
	}
	return transport
}

func SharedTransport(skipVerify bool) *http.Transport {
	if skipVerify {
		return skipVerifyTransport
	}
	return defaultTransport
}

func init() {
	skipVerifyTransport = CreateTransport(true)
	defaultTransport = CreateTransport(false)
}

type Client struct {
	URL                *url.URL
	InsecureSkipVerify bool
}

func (c *Client) Connect(ctx context.Context, src *circle.Source) error {
	u, err := url.Parse(src.URL)
	if err != nil {
		return err
	}

	// Only allow acceptance of all certs if the scheme is https AND the user opted into to the setting.
	if u.Scheme == "https" && src.InsecureSkipVerify {
		c.InsecureSkipVerify = src.InsecureSkipVerify
	}

	c.URL = u
	return nil
}

type result struct {
	Response circle.Response
	Err      error
}

func (c *Client) Login(ctx context.Context, u circle.Source) (circle.Response, error)   {
	res,err := c.login(ctx,u)
	if err!=nil{
		return nil, err
	}

	octets, err := res.MarshalJSON()
	if err != nil {
		return nil, err
	}

	fmt.Println(octets)
	return nil,nil
}


func (c *Client) login(ctx context.Context, u circle.Source) (circle.Response, error) {
	resps := make(chan result)
	go func() {
		resp, err := c.Get(c.URL, u)
		resps <- result{resp, err}
	}()

	select {
	case resp := <-resps:
		return resp.Response, resp.Err
	case <-ctx.Done():
		return nil, circle.ErrUpstreamTimeout
	}
}

type responseType struct {
	Results json.RawMessage
	Err     int    `json:"code,omitempty"` // code status
	V2Err   string `json:"msg,omitempty"`  // error message
}

// MarshalJSON returns the raw results bytes from the response
func (r responseType) MarshalJSON() ([]byte, error) {
	return r.Results, nil
}

func (r *responseType) Error() string {
	if r.Err != 1 {
		return r.V2Err
	}
	return ""
}

func (c *Client) Get(u *url.URL, q circle.Source) (circle.Response, error) {
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	params := req.URL.Query()
	params.Set("account", q.Account)
	params.Set("password", q.Password)
	params.Set("tuisongclientid", q.Tuisongclientid)

	req.URL.RawQuery = params.Encode()
	hc := &http.Client{}
	hc.Transport = SharedTransport(c.InsecureSkipVerify)
	resp, err := hc.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var response responseType
	dec := json.NewDecoder(resp.Body)
	decErr := dec.Decode(&response)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received status code %d from server: err: %s", resp.StatusCode, response.Error())
	}
	// ignore this error if we got an invalid status code
	if decErr != nil && decErr.Error() == "EOF" && resp.StatusCode != http.StatusOK {
		decErr = nil
	}

	// If we got a valid decode error, send that back
	if decErr != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK && response.V2Err != "" {
		return &response, fmt.Errorf("received status code %d from server",
			resp.StatusCode)
	}

	return &response, nil
}
