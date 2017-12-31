package config

import (
	gologging "github.com/sigmonsays/go-logging"
)

var log gologging.Logger

func init() {
	log = gologging.Register("voy", func(newlog gologging.Logger) { log = newlog })
}
