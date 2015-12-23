package crawl

import (
	"os"
	"testing"
)

func DieIf(err error) {
	if err != nil {
		panic(err)
	}
}

func Test_Collection(t *testing.T) {
	root := "/tmp/bbq"
	root = "/data/movies"

	os.MkdirAll(root, 0755)
	c, err := NewCollection(root)
	DieIf(err)

	err = c.Update()
	DieIf(err)
}
