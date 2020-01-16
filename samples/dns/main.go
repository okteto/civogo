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
package main

import (
	"github.com/okteto/civogo/client"
	"github.com/okteto/civogo/dns"
	"os"
)

func main() {
	c := client.New(os.Getenv("CIVO_TOKEN"))
	d, err := dns.NewDomain(c, "example.com")
	if err != nil {
		panic(err)
	}

	r, err := d.NewRecord(c, dns.MX, "mail", "10.10.10.1", 10, 600)
	if err != nil {
		panic(err)
	}

	if err := r.Delete(c); err != nil {
		panic(err)
	}

	if err := d.Delete(c); err != nil {
		panic(err)
	}
}
