package api

import (
	"github.com/sigmonsays/voyager/proto/vapi"
	"github.com/sigmonsays/voyager/types"

	"golang.org/x/net/context"
)

func (api *VoyApi) ListFiles(ctx context.Context, in *vapi.ListRequest) (*vapi.ListResponse, error) {

	log.Tracef("listfiles %#v", in)

	err := api.Authenticate(ctx)
	if err != nil {
		return nil, err
	}

	// load the voy file
	req := &types.ListPathRequest{
		User: in.User,
		Path: in.Path,
	}

	voy, err := api.VoyFile.Load(req)
	if err != nil {
		return nil, err
	}

	paths, err := api.VoyFile.ResolvePath(voy, req)
	if err != nil {
		return nil, err
	}

	// load the file contents
	files, err := api.PathLoader.GetFiles(paths.LocalPath)
	if err != nil {
		return nil, err
	}

	log.Tracef("files:%d", len(files))

	res := &vapi.ListResponse{}

	for _, file := range files {
		f := &vapi.File{
			Name: file.Name,
			Size: file.Size,
		}
		res.Files = append(res.Files, f)
	}

	log.Tracef("response %s", res)
	return res, nil
}
