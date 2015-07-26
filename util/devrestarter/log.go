package devrestarter

import (
	gologging "github.com/sigmonsays/go-logging"
)

var log gologging.Logger

func init() {
	log = gologging.Register("devrestarter", func(newlog gologging.Logger) { log = newlog })
}
