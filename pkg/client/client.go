// Copyright 2020 Okteto
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type result struct {
	Result string
}

// A Client to make requests to the Civo API (https://www.civo.com/api)
type Client struct {
	token  string
	client *http.Client
}

// NewClient returns a new Civo API client
func New(token string) *Client {
	c := &http.Client{
		Timeout: 30 * time.Second,
	}

	rt := withHeader(c.Transport)
	rt.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	c.Transport = rt

	return &Client{token: token, client: c}
}

// Get issues a GET to the specified path in the civo API
func (c *Client) Get(path string) (*http.Response, error) {
	url := "https://api.civo.com/v2/" + path
	r, err := c.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get %s: %w", url, err)
	}

	return r, nil
}

// Post issues a POST to the specified path in the civo API
func (c *Client) Post(path string, v interface{}) (*http.Response, error) {
	url := "https://api.civo.com/v2/" + path

	b, err := json.Marshal(v)
	if err != nil {
		return nil, fmt.Errorf("failed to post %s: %w", url, err)
	}

	r, err := c.client.Post(url, "application/json", bytes.NewBuffer(b))
	if err != nil {
		return nil, fmt.Errorf("failed to post %s: %w", url, err)
	}

	return r, nil
}

// Delete issues a DELETE to the specified path in the civo API
func (c *Client) Delete(path string) error {
	url := "https://api.civo.com/v2/" + path
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("failed to delete %s: %w", url, err)
	}

	r, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to delete %s: %w", url, err)
	}

	defer r.Body.Close()

	if r.StatusCode >= 300 {
		return fmt.Errorf("failed to get %s: StatusCode=%d", url, r.StatusCode)
	}

	var s = &result{}
	if err := json.NewDecoder(r.Body).Decode(s); err != nil {
		return fmt.Errorf("failed to delete %s: %w", url, err)
	}

	if s.Result != "success" {
		return fmt.Errorf("failed to delete %s: result=%s", url, s.Result)
	}

	return nil
}

type header struct {
	http.Header
	rt http.RoundTripper
}

func withHeader(rt http.RoundTripper) header {
	if rt == nil {
		rt = http.DefaultTransport
	}

	return header{Header: make(http.Header), rt: rt}
}

func (h header) RoundTrip(req *http.Request) (*http.Response, error) {
	for k, v := range h.Header {
		req.Header[k] = v
	}

	return h.rt.RoundTrip(req)
}
