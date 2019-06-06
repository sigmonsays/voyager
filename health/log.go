package health

import (
	gologging "github.com/sigmonsays/go-logging"
)

var log gologging.Logger

func init() {
	log = gologging.Register("health", func(newlog gologging.Logger) { log = newlog })
}
