package config

import (
	"bytes"
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

type ApplicationConfig struct {
	LogLevel      string
	Hostname      string
	StartupBanner string

	Http HttpConfig
}

type HttpConfig struct {
	BindAddr string

	// maximum number of concurrent requests we'll process (
	MaxConns int
}

func (c *ApplicationConfig) LoadDefault() {
	*c = *GetDefaultConfig()
}

func (c *ApplicationConfig) LoadYaml(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}

	b := bytes.NewBuffer(nil)
	_, err = b.ReadFrom(f)
	if err != nil {
		return err
	}

	if err := c.LoadYamlBuffer(b.Bytes()); err != nil {
		return err
	}

	if err := c.FixupConfig(); err != nil {
		return err
	}

	return nil
}

func (c *ApplicationConfig) LoadYamlBuffer(buf []byte) error {
	err := yaml.Unmarshal(buf, c)
	if err != nil {
		return err
	}
	return nil
}

func (c *ApplicationConfig) PrintYaml() {
	PrintConfig(c)
}

func GetDefaultConfig() *ApplicationConfig {

	hostname, err := os.Hostname()
	if err != nil {
		fmt.Printf("Error getting hostname: %s\n", err)
	}

	return &ApplicationConfig{
		Hostname:      hostname,
		StartupBanner: "Ready",
		Http: HttpConfig{
			BindAddr: ":8181",
			MaxConns: 1000000,
		},
	}
}

// after loading configuration this gives us a spot to "fix up" any configuration
// or abort the loading process
func (c *ApplicationConfig) FixupConfig() error {
	// var emptyConfig ApplicationConfig

	return nil
}

func PrintDefaultConfig() {
	conf := GetDefaultConfig()
	PrintConfig(conf)
}

func PrintConfig(conf *ApplicationConfig) {
	d, err := yaml.Marshal(conf)
	if err != nil {
		fmt.Println("Marshal error", err)
		return
	}
	fmt.Println("-- Configuration --")
	fmt.Println(string(d))
}