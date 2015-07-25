package server

import (
	gologging "github.com/sigmonsays/go-logging"
)

var log gologging.Logger

func init() {
	log = gologging.Register("server", func(newlog gologging.Logger) { log = newlog })
}
