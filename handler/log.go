package handler

import (
	gologging "github.com/sigmonsays/go-logging"
)

var log gologging.Logger

func init() {
	log = gologging.Register("handler", func(newlog gologging.Logger) { log = newlog })
}
