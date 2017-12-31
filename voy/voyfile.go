package voy

import (
	"fmt"
	"os"
	"strings"

	"github.com/sigmonsays/voyager/config"
	"gopkg.in/yaml.v2"
)

type VoyFile struct {
	section string

	Allow   []string
	Alias   map[string]string
	Layouts map[string]string
	Servers map[string]string
}

func (c *VoyFile) Allowed(path string) bool {
	allowed := false
	for _, p := range c.Allow {
		if strings.HasPrefix(path, p) {
			allowed = true
			break
		}
	}
	return allowed
}

func (c *VoyFile) LoadDefault() {
	*c = *DefaultConfig()
}

func (c *VoyFile) LoadYaml(path string) error {

	b, err := config.GetConfigSection(path, c.section)
	if err != nil {
		return err
	}
	if err := c.LoadYamlBuffer(b); err != nil {
		return err
	}
	if err := c.FixupConfig(); err != nil {
		return err
	}
	return nil
}

func (c *VoyFile) LoadYamlBuffer(buf []byte) error {
	err := yaml.Unmarshal(buf, c)
	if err != nil {
		return err
	}
	return nil
}

func (c *VoyFile) PrintYaml() {
	PrintConfig(c)
}

func DefaultConfig() *VoyFile {

	hostname, err := os.Hostname()
	if err != nil {
		log.Warnf("Hostname: %s", err)
	}
	var section string
	if hostname == "" {
		section = "noname"
	}
	tmp := strings.Split(hostname, ".")
	if len(tmp) > 0 {
		section = tmp[0]
	}

	return &VoyFile{
		section: section,
		Allow:   make([]string, 0),
		Alias:   make(map[string]string, 0),
		Layouts: make(map[string]string, 0),
		Servers: make(map[string]string, 0),
	}
}

// after loading configuration this gives us a spot to "fix up" any configuration
// or abort the loading process
func (c *VoyFile) FixupConfig() error {
	// var emptyConfig VoyFile

	return nil
}

func PrintConfig(conf *VoyFile) {
	d, err := yaml.Marshal(conf)
	if err != nil {
		fmt.Println("Marshal error", err)
		return
	}
	fmt.Println("-- Configuration --")
	fmt.Println(string(d))
}
