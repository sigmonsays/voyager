package api

import (
	"github.com/sigmonsays/voyager/proto/vapi"
	"github.com/sigmonsays/voyager/types"

	"golang.org/x/net/context"
)

func (api *VoyApi) GetConfig(ctx context.Context, in *vapi.ConfigRequest) (*vapi.ConfigResponse, error) {

	log.Tracef("GetConfig %#v", in)

	err := api.Authenticate(ctx)
	if err != nil {
		return nil, err
	}

	req := &types.ListPathRequest{
		User: in.User,
	}
	voy, err := api.VoyFile.Load(req)
	if err != nil {
		return nil, err
	}

	res := &vapi.ConfigResponse{
		Allow:   voy.Allow,
		Alias:   voy.Alias,
		Servers: voy.Servers,
	}

	return res, nil
}
