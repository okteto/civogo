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

package dns

import "github.com/okteto/civogo/pkg/client"

import "fmt"

import "encoding/json"

const (
	// A record type
	A = "a"

	// CName record type
	CName = "cname"

	// MX record type
	MX = "mx"

	// TXT record type
	TXT = "txt"
)

// A RecordType is the allowed record types: a, cname, mx or txt
type RecordType string

// A Record is a DNS Record managed by Civo
type Record struct {
	ID       string     `json:"id"`
	DomainID string     `json:"domain_id"`
	Name     string     `json:"name"`
	Value    string     `json:"value"`
	Type     RecordType `json:"type"`
	Priority int        `json:"priority"`
	TTL      int        `json:"ttl"`
}

func (d *Domain) getRecordsURL() string {
	return fmt.Sprintf("dns/%s/records", d.ID)
}

// NewRecord creates a new DNS record
func (d *Domain) NewRecord(c *client.Client, t RecordType, name, value string, priority int, ttl int) (*Record, error) {
	if ttl == 0 {
		ttl = 600
	}

	r := Record{
		Name:     name,
		Value:    value,
		Type:     t,
		Priority: priority,
		TTL:      ttl,
	}

	url := d.getRecordsURL()
	resp, err := c.Post(url, r)
	if err != nil {
		return nil, fmt.Errorf("internal-server-error: %w", err)
	}

	var v = &Record{}
	if err := json.NewDecoder(resp.Body).Decode(v); err != nil {
		return nil, fmt.Errorf("internal-server-error: %w", err)
	}

	return v, nil
}

// GetRecord returns a matching DNS record for the domain or an error
func (d *Domain) GetRecord(c *client.Client, name string) (*Record, error) {
	resp, err := c.Get(d.getRecordsURL())
	if err != nil {
		return nil, fmt.Errorf("internal-server-error: %w", err)
	}

	var rs = make([]Record, 0)
	if err := json.NewDecoder(resp.Body).Decode(&rs); err != nil {
		return nil, fmt.Errorf("internal-server-error: %w", err)
	}

	for _, r := range rs {
		if r.Name == name {
			return &r, nil
		}
	}

	return nil, fmt.Errorf("not-found")
}

// Delete deletes the DNS record
func (r *Record) Delete(c *client.Client) error {
	url := fmt.Sprintf("dns/%s/records/%s", r.DomainID, r.ID)
	if err := c.Delete(url); err != nil {
		return fmt.Errorf("internal-server-error: %w", err)
	}

	return nil
}
