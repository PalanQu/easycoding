package app

import (
	"context"
	"easycoding/internal/config"
	error_middleware "easycoding/internal/middleware/error"
	"fmt"

	log_middleware "easycoding/internal/middleware/log"
	prometheus_middleware "easycoding/internal/middleware/prometheus"
	recover_middleware "easycoding/internal/middleware/recover"
	validate_middleware "easycoding/internal/middleware/validate"
	"easycoding/internal/service"
	"easycoding/pkg/db"
	"easycoding/pkg/log"
	"easycoding/pkg/swagger"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	assetfs "github.com/elazarl/go-bindata-assetfs"
	"github.com/gorilla/mux"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
	"gorm.io/gorm"
)

const (
	version = "v0.0.1"
)

// Possible application states.
const (
	StateStopped int = iota
	StateStarting
	StateRunning
	StateStopping
)

const (
	prometheusPrefix = "/metrics"
	swaggerPrefix    = "/swagger/"

	// TODO(qujiabao): refactor these hard code values into config
	serveHost          = "0.0.0.0"
	maxMsgSize         = 5 * 1024 * 1024
	clientMinWaitPing  = 5 * time.Second
	serverPingDuration = 10 * time.Minute
)

type Kernel struct {
	gwServer      *http.Server
	grpcServer    *grpc.Server
	swaggerServer *http.Server
	Config        *config.Config
	Log           *logrus.Logger
	DB            *gorm.DB
	state         int
	wg            *sync.WaitGroup
	context       cancelContext
	shutdownFns   []func() error
}

func New(configPath string) (*Kernel, error) {
	config := config.LoadConfig(configPath)
	logger := log.New(os.Stderr, config.Log.Level, config.Log.Dir)
	database, err := db.CreateGdb(config, logger)
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect to db")
	}

	// Create a global application context.
	ctx, cancel := context.WithCancel(context.Background())

	gwServer := newGrpcGatewayServer(config)
	grpcServer := newGrpcServer(config, logger, database)
	swaggerServer := newSwaggerServer(config)

	// Build the Kernel struct with all dependencies.
	app := &Kernel{
		Log:           logger,
		Config:        config,
		DB:            database,
		grpcServer:    grpcServer,
		gwServer:      gwServer,
		swaggerServer: swaggerServer,
		state:         StateStarting,
		shutdownFns:   []func() error{},
		wg:            &sync.WaitGroup{},
		context:       cancelContext{cancel: cancel, ctx: ctx},
	}

	app.state = StateRunning

	return app, nil
}

func newGrpcServer(
	config *config.Config,
	logger *logrus.Logger,
	db *gorm.DB,
) *grpc.Server {
	opts := []grpc.ServerOption{
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			prometheus_middleware.Interceptor(),
			log_middleware.Interceptor(logger),
			validate_middleware.Interceptor(),
			recover_middleware.Interceptor(),
			error_middleware.Interceptor(logger),
		)),
		grpc.MaxSendMsgSize(maxMsgSize),
		grpc.MaxRecvMsgSize(maxMsgSize),
		grpc.KeepaliveParams(keepalive.ServerParameters{
			Time: serverPingDuration,
		}),
		grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
			MinTime:             clientMinWaitPing,
			PermitWithoutStream: true,
		}),
	}
	// Create grpc server & register grpc services.
	grpcServer := grpc.NewServer(opts...)
	service.RegisterServers(grpcServer, logger, db)
	grpc_prometheus.EnableHandlingTimeHistogram()
	grpc_prometheus.Register(grpcServer)
	reflection.Register(grpcServer)
	return grpcServer
}

func newGrpcGatewayServer(config *config.Config) *http.Server {
	gwmux := runtime.NewServeMux()
	router := mux.NewRouter()
	router.PathPrefix("/").Handler(gwmux)
	gwServer := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", serveHost, config.Server.GatewayPort),
		Handler: router,
	}
	service.RegisterHandlers(gwmux, fmt.Sprintf("%s:%s", serveHost, config.Server.GrpcPort))
	return gwServer
}

func newSwaggerServer(config *config.Config) *http.Server {
	router := mux.NewRouter()
	fileServer := http.FileServer(&assetfs.AssetFS{
		Asset:     swagger.Asset,
		AssetDir:  swagger.AssetDir,
		AssetInfo: swagger.AssetInfo,
		// The swagger-ui is built from npm package swagger-ui-dist.
		Prefix: "dist",
	})
	swaggerJSON, _ := filepath.Abs("api/api.swagger.json")
	jsonPrefix := fmt.Sprintf("%s%s", swaggerPrefix, "api.swagger.json")
	router.PathPrefix(jsonPrefix).Handler(
		http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			f, err := ioutil.ReadFile(swaggerJSON)
			if err != nil {
				panic(err)
			}
			w.Write(f)
		}))

	router.PathPrefix(swaggerPrefix).Handler(
		http.StripPrefix(swaggerPrefix, fileServer))
	router.PathPrefix(prometheusPrefix).Handler(promhttp.Handler())

	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", serveHost, config.Server.SwaggerPort),
		Handler: router,
	}
	return server
}

// cancelContext is a context with a cancel function.
type cancelContext struct {
	cancel context.CancelFunc
	ctx    context.Context
}

func (k *Kernel) listenAndServe(name string, listen func() error) {
	var err error
	for {
		// If the app is stopped or stopping, don't retry to start the server.
		if k.state == StateStopping || k.state == StateStopped {
			k.Log.Tracef("skipping restarts of server because app is not in running state: state is %d", k.state)
			return
		}

		k.Log.Infof("%s started\n", name)
		if err = listen(); err != nil {
			if k.Config.Server.RestartOnError {
				k.Log.Infof("restart server failed after error on %v", err)
				continue
			}
			k.Log.Infof("server failed after error on %v", err)
		}
		return
	}

}

func (k *Kernel) ListenGrpcGateway() {
	listen := func() error {
		if err := k.gwServer.ListenAndServe(); err != nil {
			return err
		}
		return nil
	}
	k.listenAndServe("grpc gateway server", listen)
}

func (k *Kernel) ListenGrpc() {
	listen := func() error {
		lis, err := net.Listen(
			"tcp", fmt.Sprintf("%s:%s", serveHost, k.Config.Server.GrpcPort))
		if err != nil {
			return err
		}
		if err := k.grpcServer.Serve(lis); err != nil {
			return err
		}
		return nil
	}
	k.listenAndServe("grpc server", listen)
}

func (k *Kernel) ListenSwagger() {
	listen := func() error {
		if err := k.swaggerServer.ListenAndServe(); err != nil {
			return err
		}
		return nil
	}
	k.listenAndServe("swagger server", listen)
}

func (k *Kernel) Shutdown(ctx context.Context) error {
	if k.state != StateRunning {
		k.Log.Warn("Application cannot be shutdown since current state is not 'running'")
		return nil
	}

	k.state = StateStopping
	defer func() {
		k.state = StateStopped
	}()

	if k.gwServer != nil {
		if err := k.gwServer.Shutdown(ctx); err != nil {
			k.Log.Errorf("server shutdown error: %v\n", err)
		} else {
			k.Log.Infoln("gateway server stopped")
		}
	}

	if k.grpcServer != nil {
		k.grpcServer.GracefulStop()
		k.Log.Infoln("grpc server stopped")
	}

	if k.swaggerServer != nil {
		if err := k.swaggerServer.Shutdown(ctx); err != nil {
			k.Log.Errorf("swagger shutdown error: %v\n", err)
		} else {
			k.Log.Infoln("swagger server stopped")
		}
	}

	// Cancel global context, then wait for all processes to quit.
	k.context.cancel()
	done := make(chan struct{})
	go func() {
		k.wg.Wait()
		close(done)
	}()

	// Run shutdown functions.
	for _, fn := range k.shutdownFns {
		shutdownErr := fn()
		if shutdownErr != nil {
			k.Log.Errorf("shutdown function returned error: %v\n", shutdownErr)
		}
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-done:
	}

	sqlDB, err := k.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

func Version() string {
	return version
}
