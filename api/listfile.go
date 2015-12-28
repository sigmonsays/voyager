package api

import (
	"github.com/sigmonsays/voyager/proto/vapi"

	"golang.org/x/net/context"
)

func (api *VoyApi) ListFiles(ctx context.Context, in *vapi.ListRequest) (*vapi.ListResponse, error) {

	err := api.Authenticate(ctx)
	if err != nil {
		return nil, err
	}

	res := &vapi.ListResponse{}
	return res, nil
}
