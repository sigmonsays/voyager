package vapi

import (
	gologging "github.com/sigmonsays/go-logging"
)

var log gologging.Logger

func init() {
	log = gologging.Register("vapi", func(newlog gologging.Logger) { log = newlog })
}
