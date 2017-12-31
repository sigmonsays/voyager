package config

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

type ApplicationConfig struct {
	LogLevel      string
	LogLevels     map[string]string
	Hostname      string
	ServerName    string
	StartupBanner string

	AutoUpgrade bool
	AutoRestart bool
	CacheDir    string

	Http *HttpConfig
	Rpc  *RpcConfig

	Username string
	Servers  map[string]string

	ACL []string
}

type HttpConfig struct {
	BindAddr string

	// maximum number of concurrent requests we'll process (
	MaxConns int
}

type RpcConfig struct {
	BindAddr string

	Secret string
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
func (me *ApplicationConfig) LoadYamlSection(path, section string) error {
	b, err := GetConfigSection(path, section)
	if err != nil {
		return err
	}
	if err := me.LoadYamlBuffer(b); err != nil {
		return err
	}
	if err := me.FixupConfig(); err != nil {
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
		ServerName:    fmt.Sprintf("%s:%d", hostname, 8181),
		StartupBanner: "Ready",
		AutoUpgrade:   true,
		AutoRestart:   true,
		CacheDir:      "/tmp/voyager",
		LogLevels:     make(map[string]string, 0),
		Http: &HttpConfig{
			BindAddr: ":8181",
			MaxConns: 1000000,
		},
		Rpc: &RpcConfig{
			BindAddr: ":8191",
			Secret:   "changeme",
		},
		Servers: make(map[string]string, 0),
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

// loads configuration from a specific config file in a section by name
// and then returns the yaml stream as []byte to be loaded natively
// allow reading from a section

type partial struct {
	ConfigFile map[string]interface{}
}

func GetConfigSection(path, section string) ([]byte, error) {
	log.Tracef("GetConfigSection path=%q section=%q", path, section)
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	cfg := &partial{
		ConfigFile: make(map[string]interface{}),
	}

	err = yaml.Unmarshal(buf, cfg.ConfigFile)
	if err != nil {
		return nil, err
	}

	conf, ok := cfg.ConfigFile[section]
	if ok == false {
		return nil, fmt.Errorf("No such section %s", section)
	}

	if conf == nil {
		return nil, fmt.Errorf("empty section %s", section)
	}

	blob, err := yaml.Marshal(conf)
	if err != nil {
		return nil, err
	}
	return blob, nil
}
