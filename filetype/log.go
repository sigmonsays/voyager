package filetype

import (
	gologging "github.com/sigmonsays/go-logging"
)

var log gologging.Logger

func init() {
	log = gologging.Register("filetype", func(newlog gologging.Logger) { log = newlog })
}
