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

func NewCartographer(orgName string, token string) *Cartographer {
	return &Cartographer{
		client: &http.Client{
			Timeout: time.Second * 10,
		},
		orgName: orgName,
		token:   token,
	}
}

func (c *Cartographer) Modules() ([]modules.Module, error) {
	return modules.Modules(c.client, c.orgName, c.token)
}

func main() {
	c := NewCartographer("thelostsons", os.Getenv("TFTOKEN"))

	mods, err := c.Modules()
	if err != nil {
		log.Fatal(err)
	}

	for _, mod := range mods {
		log.Printf("%s %s %s %s %d %s\n", mod.Name, mod.Source, mod.Version, mod.RegistryType, mod.WorkspaceCount, mod.Workspaces)
	}
}
