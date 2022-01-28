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
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"x6t.io/circle"
)

func (c *Client) Wechat(ctx context.Context, share circle.WechatShare) error {
	_, err := c.wechat(ctx, share)
	return err
}

func (c *Client) wechat(ctx context.Context, u circle.WechatShare) (circle.Response, error) {
	resps := make(chan result)
	go func() {
		resp, err := c.getWechat(c.URL, u)
		resps <- result{resp, err}
	}()

	select {
	case resp := <-resps:
		return resp.Response, resp.Err
	case <-ctx.Done():
		return nil, ErrUpstreamTimeout
	}
}

func (c *Client) getWechat(u *url.URL, q circle.WechatShare) (circle.Response, error) {
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}
	if c.Token == "" {
		return nil, fmt.Errorf("token must be empty")
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("token", c.Token)

	params := req.URL.Query()
	params.Set("microid", q.Microgrid)
	params.Set("type", q.Type) // default 1

	req.URL.RawQuery = params.Encode()
	hc := &http.Client{}
	hc.Transport = SharedTransport(c.InsecureSkipVerify)
	resp, err := hc.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var response responseType
	decErr := json.NewDecoder(resp.Body).Decode(&response)

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

	if resp.StatusCode != http.StatusOK && response.Err != "" {
		return &response, fmt.Errorf("received status code %d from server",
			resp.StatusCode)
	}

	return &response, nil
}
