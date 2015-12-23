package crawl

import (
	gologging "github.com/sigmonsays/go-logging"
)

var log gologging.Logger

func init() {
	log = gologging.Register("crawl", func(newlog gologging.Logger) { log = newlog })
}
