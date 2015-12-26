package api

import (
	"github.com/sigmonsays/voyager/config"
	"github.com/sigmonsays/voyager/proto/vapi"
	"golang.org/x/net/context"
)

type VoyApi struct {
}

func MakeApi(cfg *config.ApplicationConfig) *VoyApi {
	v := &VoyApi{}
	return v
}

func (api *VoyApi) Ping(ctx context.Context, in *vapi.PingRequest) (*vapi.PingResponse, error) {
	res := &vapi.PingResponse{
		Message: in.Message,
	}
	return res, nil
}
