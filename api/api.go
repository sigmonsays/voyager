package api

import (
	"fmt"

	"github.com/sigmonsays/voyager/config"
	"github.com/sigmonsays/voyager/handler"
	"github.com/sigmonsays/voyager/proto/vapi"

	"golang.org/x/net/context"
	"google.golang.org/grpc/metadata"
)

type VoyApi struct {
	Secret string

	Factory *handler.HandlerFactory
}

func MakeApi(cfg *config.ApplicationConfig) *VoyApi {

	v := &VoyApi{
		Secret: cfg.Rpc.Secret,
	}
	return v
}

func (api *VoyApi) Authenticate(ctx context.Context) error {
	md, ok := metadata.FromContext(ctx)
	if ok == false {
		return fmt.Errorf("authenticate: no metadata")
	}

	secrets, ok := md["request-secret"]
	if ok == false {
		return fmt.Errorf("authenticate: no secret")
	}

	if len(secrets) < 1 {
		return fmt.Errorf("authenticate: no secret")
	}
	secret := secrets[0]
	if secret != "" {
		return fmt.Errorf("authenticate: empty secret")
	}

	if secret != api.Secret {
		return fmt.Errorf("authenticate: invalid secret")
	}
	return nil

	/*
		TODO: Would be nice to use the client and server cert as credentials

		authInfo, ok := credentials.FromContext(ctx)
		if ok == false {
			return fmt.Errorf("not authenticated")
		}

		var authType string
		switch info := authInfo.(type) {
		case credentials.TLSInfo:
			authType = info.AuthType()
			serverName = info.State.ServerName
		default:
			return fmt.Errorf("unknown authInfo type: %s", info.AuthType())
		}

	*/

}

func (api *VoyApi) WithHandlerFactory(x *handler.HandlerFactory) *VoyApi {
	api.Factory = x
	return api
}

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
