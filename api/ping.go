package api

import (
	"github.com/sigmonsays/voyager/proto/vapi"

	"golang.org/x/net/context"
)

func (api *VoyApi) Ping(ctx context.Context, in *vapi.PingRequest) (*vapi.PingResponse, error) {

	err := api.Authenticate(ctx)
	if err != nil {
		return nil, err
	}

	res := &vapi.PingResponse{
		Message: in.Message,
	}
	return res, nil
}
