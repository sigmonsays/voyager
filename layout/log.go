package layout

import (
	gologging "github.com/sigmonsays/go-logging"
)

var log gologging.Logger

func init() {
	log = gologging.Register("layout", func(newlog gologging.Logger) { log = newlog })
}
