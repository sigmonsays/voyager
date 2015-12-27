package api

import (
	gologging "github.com/sigmonsays/go-logging"
)

var log gologging.Logger

func init() {
	log = gologging.Register("api", func(newlog gologging.Logger) { log = newlog })
}
