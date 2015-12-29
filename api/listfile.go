package api

import (
	"net/url"

	"github.com/sigmonsays/voyager/filetype"
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

	log.Tracef("paths %s", paths)

	// load the file contents
	files, err := api.PathLoader.GetFiles(paths.LocalPath)
	if err != nil {
		return nil, err
	}

	log.Tracef("files:%d", len(files))

	// determine the layout
	layout, err := api.Layout.Resolve(voy, paths.LocalPath, files)
	if err != nil {
		return nil, err
	}

	urlp := &url.URL{}
	urlp.Scheme = "http"
	urlp.Host = api.ServerName
	urlp.Path = paths.UrlPrefix

	res := &vapi.ListResponse{
		Layout:       filetype.TypeToString(layout),
		UrlPrefix:    urlp.String(),
		RelPath:      paths.RelPath,
		LocalPath:    paths.LocalPath,
		RemoteServer: "http://" + api.ServerName,
	}

	for _, file := range files {
		f := &vapi.File{
			IsDir: file.IsDir,
			Name:  file.Name,
			Size:  file.Size,
			Mtime: file.MTime,
		}
		res.Files = append(res.Files, f)
	}

	log.Tracef("response %s", res)
	return res, nil
}
