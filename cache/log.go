package cache

import (
	gologging "github.com/sigmonsays/go-logging"
)

var log gologging.Logger

func init() {
	log = gologging.Register("cache", func(newlog gologging.Logger) { log = newlog })
}
