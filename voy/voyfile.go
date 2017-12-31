package voy

import (
	"fmt"
	"strings"

	"github.com/sigmonsays/voyager/config"
	"gopkg.in/yaml.v2"
)

type VoyFile struct {
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
	section := "default"

	b, err := config.GetConfigSection(path, section)
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

	return &VoyFile{
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
