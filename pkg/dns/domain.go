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

import (
	"encoding/json"
	"fmt"

	"github.com/okteto/civogo/pkg/client"
)

// A Domain registered in Civo
type Domain struct {
	// The ID of the domain
	ID string `json:"id"`

	// The Name of the domain
	Name string `json:"name"`

	// The Result of the operation
	Result string `json:"result"`
}

// NewDomain creates a domain in Civo
func NewDomain(c *client.Client, name string) (*Domain, error) {
	d := Domain{
		Name: name,
	}

	r, err := c.Post("dns", d)
	if err != nil {
		return nil, fmt.Errorf("internal-server-error: %w", err)
	}

	defer r.Body.Close()
	var v = &Domain{}
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		return nil, fmt.Errorf("internal-server-error: %w", err)
	}

	if v.Result != "success" {
		return nil, fmt.Errorf("internal-server-error: Result=%s", v.Result)
	}

	return v, nil
}

// GetDomain retrieves a domain from Civo
func GetDomain(c *client.Client, domain string) (*Domain, error) {
	r, err := c.Get("dns")
	if err != nil {
		return nil, err
	}

	var ds = make([]Domain, 0)
	if err := json.NewDecoder(r.Body).Decode(&ds); err != nil {
		return nil, fmt.Errorf("internal-server-error: %w", err)
	}

	defer r.Body.Close()

	for _, d := range ds {
		if d.Name == domain {
			return &d, nil
		}
	}
	return nil, fmt.Errorf("not-found")
}

// Delete deletes the domain record in Civo
func (d *Domain) Delete(c *client.Client) error {
	if err := c.Delete(fmt.Sprintf("dns/%s", d.ID)); err != nil {
		return fmt.Errorf("internal-server-error: %w", err)
	}

	return nil
}
