package main

import (
	"github.com/majesticbeast/cartographer/modules"
	"log"
	"net/http"
	"os"
	"time"
)

type Cartographer struct {
	client  *http.Client
	orgName string
	token   string
}

type Filter struct {
	Type     string
	Operator string
	Value    string
}

func (c *Cartographer) Client() *http.Client {
	return c.client
}

func (c *Cartographer) OrgName() string {
	return c.orgName
}

func (c *Cartographer) Token() string {
	return c.token
}

func NewCartographer(orgName string, token string) *Cartographer {
	return &Cartographer{
		client: &http.Client{
			Timeout: time.Second * 10,
		},
		orgName: orgName,
		token:   token,
	}
}

func (c *Cartographer) Modules(filters []modules.Filter) (modules.ModuleList, error) {
	return modules.Modules(c, filters)
}

func main() {
	c := NewCartographer("thelostsons", os.Getenv("TFTOKEN"))

	var filters []modules.Filter
	filters = append(filters, modules.Filter{
		Type:     "name",
		Operator: "contains",
		Value:    "iam",
	})

	mods, err := c.Modules(filters)
	if err != nil {
		log.Fatal(err)
	}

	for _, mod := range mods {
		log.Println(mod)
	}

}
