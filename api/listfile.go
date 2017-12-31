package api

import (
	"net/url"
	"os"

	"github.com/sigmonsays/voyager/filetype"
	"github.com/sigmonsays/voyager/proto/vapi"
	"github.com/sigmonsays/voyager/types"

	"golang.org/x/net/context"
)

func (me *VoyApi) ListFiles(ctx context.Context, in *vapi.ListRequest) (*vapi.ListResponse, error) {

	log.Tracef("listfiles %#v", in)

	err := me.Authenticate(ctx)
	if err != nil {
		return nil, err
	}

	// load the voy file
	req := &types.ListPathRequest{
		User: in.User,
		Path: in.Path,
	}

	voy, err := me.VoyFile.Load(req)
	if err != nil {
		return nil, err
	}

	paths, err := me.VoyFile.ResolvePath(voy, req)
	if err != nil {
		return nil, err
	}

	log.Tracef("paths %s", paths)

	// return a error if
	st, err := os.Stat(paths.LocalPath)
	if err != nil {
		return nil, err
	}

	if st.IsDir() == false {
		res := &vapi.ListResponse{
			IsDir:        false,
			RelPath:      paths.RelPath,
			LocalPath:    paths.LocalPath,
			RemoteServer: "http://" + me.ServerName,
		}
		return res, nil
	}

	// load the file contents
	files, err := me.PathLoader.GetFiles(paths.LocalPath)
	if err != nil {
		return nil, err
	}

	log.Tracef("files:%d", len(files))

	// determine the layout
	layout, err := me.Layout.Resolve(voy, paths.LocalPath, files)
	if err != nil {
		return nil, err
	}

	urlp := &url.URL{}
	urlp.Scheme = "http"
	//urlp.Host = me.ServerName
	urlp.Host = voy.ServerName
	urlp.Path = paths.UrlPrefix

	res := &vapi.ListResponse{
		IsDir:        true,
		Layout:       filetype.TypeToString(layout),
		UrlPrefix:    urlp.String(),
		RelPath:      paths.RelPath,
		LocalPath:    paths.LocalPath,
		RemoteServer: "http://" + me.ServerName,
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

	log.Tracef("response layout:%s urlprefix:%s localpath:%s remoteserver:%s files:%d",
		res.Layout, res.UrlPrefix, res.LocalPath, res.RemoteServer, len(res.Files))
	return res, nil
}
