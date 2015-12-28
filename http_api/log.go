package http_api

import (
	gologging "github.com/sigmonsays/go-logging"
)

var log gologging.Logger

func init() {
	log = gologging.Register("http_api", func(newlog gologging.Logger) { log = newlog })
}
