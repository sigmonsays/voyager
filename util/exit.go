package util

import (
	"fmt"
	"os"
	"strings"
)

func ExitIfError(err error, s string, args ...interface{}) {
	if err != nil {
		prefix := fmt.Sprintf("ERROR: %s: ", err)
		fmt.Printf(strings.Replace("%", "%%", prefix, -1)+s+"\n", args...)
		os.Exit(1)
	}
}
