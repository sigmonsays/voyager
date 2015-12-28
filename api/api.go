package api

import (
	"fmt"
	"os"

	"github.com/sigmonsays/voyager/config"
	"github.com/sigmonsays/voyager/handler"
	"github.com/sigmonsays/voyager/layout"
	"github.com/sigmonsays/voyager/voy"

	"golang.org/x/net/context"
	"google.golang.org/grpc/metadata"
)

type VoyApi struct {
	Secret     string
	ServerName string

	Factory    *handler.HandlerFactory
	PathLoader handler.PathLoader
	Layout     layout.LayoutResolver
	VoyFile    voy.VoyLoader
}

func MakeApi(cfg *config.ApplicationConfig) *VoyApi {

	v := &VoyApi{
		Secret: cfg.Rpc.Secret,
	}

	var err error
	v.ServerName, err = os.Hostname()
	if err != nil {
		log.Warnf("Hostname: %s", err)
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

	if len(secrets) == 0 {
		return fmt.Errorf("authenticate: no secret")
	}
	secret := secrets[0]
	if secret == "" {
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
