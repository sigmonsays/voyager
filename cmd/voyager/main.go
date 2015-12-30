package main

import (
	"flag"
	"net"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/sigmonsays/voyager/api"
	"github.com/sigmonsays/voyager/cache"
	"github.com/sigmonsays/voyager/config"
	"github.com/sigmonsays/voyager/handler"
	"github.com/sigmonsays/voyager/http_api"
	"github.com/sigmonsays/voyager/layout"
	"github.com/sigmonsays/voyager/proto/vapi"
	"github.com/sigmonsays/voyager/util"
	"github.com/sigmonsays/voyager/util/devrestarter"
	"github.com/sigmonsays/voyager/voy"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	reload_git "github.com/sigmonsays/git-watch/reload/git"
	gologging "github.com/sigmonsays/go-logging"
)

func Shell(args ...string) error {
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return err
}

func main() {
	cfg := config.GetDefaultConfig()
	flag.StringVar(&cfg.Http.BindAddr, "bind", cfg.Http.BindAddr, "bind address")
	flag.StringVar(&cfg.LogLevel, "log", cfg.LogLevel, "log level")
	flag.Parse()

	var err error

	cfgfile := filepath.Join(os.Getenv("HOME"), ".voyager")
	if util.FileExists(cfgfile) {
		err = cfg.LoadYaml(cfgfile)
		if err != nil {
			log.Errorf("load config %s: %s", cfgfile, err)
			return
		}
	}

	if cfg.LogLevel != "" {
		gologging.SetLogLevel(cfg.LogLevel)
	}
	gologging.SetLogLevels(cfg.LogLevels)

	if log.IsTrace() {
		cfg.PrintYaml()
	}

	if cfg.AutoRestart {
		devrestarter.Init()
	}

	if util.FileExists(cfg.CacheDir) == false {
		err = os.MkdirAll(cfg.CacheDir, 0755)
		if err != nil {
			log.Errorf("cache dir error %s: %s", cfg.CacheDir, err)
			return
		}
	}

	if cfg.AutoUpgrade {
		gw := reload_git.NewGitWatch(".", "master")
		gw.Interval = 30
		gw.OnChange = func(dir, branch, lhash, rhash string) error {
			err = Shell("git", "pull")
			if err != nil {
				log.Warnf("git pull error: %s", err)
			}

			err = Shell("go", "install", "-v", "./...")
			if err != nil {
				log.Warnf("go install error: %s", err)
			}
			return err
		}

		err := gw.Start()
		if err != nil {
			log.Errorf("starting git watch %s", err)
			return
		}
	}

	cache := cache.NewFileCache(cfg.CacheDir)
	if err != nil {
		log.Errorf("cache error: %s %s", cfg.CacheDir, err)
		return
	}

	// core components
	handlerFactory := handler.NewHandlerFactory()
	pathLoader := handler.NewFilesystemPathLoader()

	layoutResolver := layout.NewLayoutResolver()

	voyLoader := voy.NewVoyLoader()

	md := metadata.Pairs("request-secret", cfg.Rpc.Secret)
	ctx := context.Background()
	ctx = metadata.NewContext(ctx, md)

	// the HTTP server
	srv := http_api.NewServer(cfg.Http.BindAddr)
	srv.Conf = cfg
	srv.Cache = cache
	srv.Factory = handlerFactory
	srv.PathLoader = pathLoader
	srv.Layout = layoutResolver
	srv.VoyFile = voyLoader
	srv.Ctx = ctx
	srv.Dashboard = handler.NewDashboardHandler()
	srv.Dashboard.VoyFile = voyLoader
	srv.Dashboard.Username = cfg.Username

	go func() {
		err = srv.Start()
		if err != nil {
			log.Errorf("starting server %s", err)
			return
		}
	}()

	// the RPC server
	lis, err := net.Listen("tcp", cfg.Rpc.BindAddr)
	if err != nil {
		util.ExitIfError(err, "grpc listen %s: %s", cfg.Rpc.BindAddr, err)
	}
	http2_serv := grpc.NewServer()
	http2_api := api.MakeApi(cfg)
	http2_api.WithHandlerFactory(handlerFactory)

	http2_api.PathLoader = pathLoader
	http2_api.Layout = layoutResolver
	http2_api.VoyFile = voyLoader
	http2_api.ServerName = cfg.ServerName

	vapi.RegisterVApiServer(http2_serv, http2_api)

	log.Infof("%s", cfg.StartupBanner)

	err = http2_serv.Serve(lis)
	if err != nil {
		util.ExitIfError(err, "grpc serve %s: %s", cfg.Rpc.BindAddr, err)
	}

}
