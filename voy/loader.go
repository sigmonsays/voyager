package voy

import (
	"fmt"
	"os/user"
	"path/filepath"

	"github.com/sigmonsays/voyager/types"
)

type VoyLoader interface {
	Load(req *types.ListPathRequest) (*VoyFile, error)
}

func NewVoyLoader() *fileLoader {
	l := &fileLoader{}
	return l
}

type fileLoader struct {
}

func (l *fileLoader) Load(req *types.ListPathRequest) (*VoyFile, error) {
	user_ent, err := user.Lookup(req.User)
	if err != nil {
		return nil, fmt.Errorf("user lookup %s: %s", req.User, err)
	}
	homedir := user_ent.HomeDir

	voy := DefaultConfig()
	voyfile := filepath.Join(homedir, ".voyager")

	err = voy.LoadYaml(voyfile)
	if err != nil {
		return nil, fmt.Errorf("load voyfile %s: %s", voyfile, err)
	}

	return voy, nil
}
