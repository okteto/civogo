package dns

import (
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/google/uuid"
	"github.com/okteto/civogo/pkg/client"
)

var token = os.Getenv("CIVO_TOKEN")
var c *client.Client

func TestMain(m *testing.M) {
	if len(token) == 0 {
		fmt.Println("CIVO_TOKEN is not set")
		os.Exit(1)
	}

	c = client.New(token)
	os.Exit(m.Run())
}
func TestGetRecord(t *testing.T) {
	dname := fmt.Sprintf("%s.example.com", uuid.New().String())
	d, err := NewDomain(c, dname)
	if err != nil {
		t.Fatal(err)
	}

	nd, err := GetDomain(c, dname)
	if err != nil {
		t.Fatal(err)
	}

	if nd.ID != d.ID {
		t.Fatalf("got %+v, expected %+v", nd, d)
	}

	name := uuid.New().String()
	value := uuid.New().String()
	newRecord, err := d.NewRecord(c, TXT, name, value, 10, 600)
	if err != nil {
		t.Fatal(err)
	}

	if len(newRecord.ID) == 0 {
		t.Fatal("ID is empty")
	}

	if newRecord.DomainID != d.ID {
		t.Fatalf("ID is wrong. Got %s, expected %s", newRecord.DomainID, d.ID)
	}

	if newRecord.Name != name {
		t.Fatalf("name is wrong. Got %s, expected %s", newRecord.Name, name)
	}

	if newRecord.Value != value {
		t.Fatalf("value is wrong. Got %s, expected %s", newRecord.Value, value)
	}

	g, err := d.GetRecord(c, newRecord.Name)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(g, newRecord) {
		t.Fatalf("got %+v, expected %+v", g, newRecord)
	}

	if err := g.Delete(c); err != nil {
		t.Fatal(err)
	}

	if err := d.Delete(c); err != nil {
		t.Fatal(err)
	}
}
