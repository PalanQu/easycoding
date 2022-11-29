package app

import (
	"context"
	"easycoding/internal/config"
	"fmt"

	auth_middleware "easycoding/internal/middleware/auth"
	error_middleware "easycoding/internal/middleware/error"
	log_middleware "easycoding/internal/middleware/log"
	otel_middleware "easycoding/internal/middleware/otel"
	prometheus_middleware "easycoding/internal/middleware/prometheus"
	recover_middleware "easycoding/internal/middleware/recover"
	validate_middleware "easycoding/internal/middleware/validate"
	"easycoding/internal/service"
	"easycoding/pkg/db"
	"easycoding/pkg/ent"
	"easycoding/pkg/log"

	pkg_otel "easycoding/pkg/otel"
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
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
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
	log           *logrus.Logger
	gwServer      *http.Server
	grpcServer    *grpc.Server
	swaggerServer *http.Server
	config        *config.Config
	dbClient      *ent.Client
	state         int
	wg            *sync.WaitGroup
	context       cancelContext
	shutdownFns   []func() error
}

func New(configPath string) (*Kernel, error) {
	config := config.LoadConfig(configPath)
	logger := log.New(os.Stderr, config.Log.Level, config.Log.Dir)
	database, err := db.CreateDBClient(config)
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect to db")
	}
	tracer, shutdownTraceFunc, err := pkg_otel.NewTracer()
	if err != nil {
		return nil, errors.Wrap(err, "failed to new tracer")
	}
	// Create a global application context.
	ctx, cancel := context.WithCancel(context.Background())

	gwServer := newGrpcGatewayServer(config, logger)
	grpcServer := newGrpcServer(config, logger, database, tracer)
	swaggerServer := newSwaggerServer(config)

	// Build the Kernel struct with all dependencies.
	app := &Kernel{
		log:           logger,
		config:        config,
		dbClient:      database,
		grpcServer:    grpcServer,
		gwServer:      gwServer,
		swaggerServer: swaggerServer,
		state:         StateStarting,
		shutdownFns: []func() error{
			shutdownTraceFunc,
		},
		wg:      &sync.WaitGroup{},
		context: cancelContext{cancel: cancel, ctx: ctx},
	}

	app.state = StateRunning

	return app, nil
}

func newGrpcServer(
	config *config.Config,
	logger *logrus.Logger,
	db *ent.Client,
	tracer trace.Tracer,
) *grpc.Server {
	opts := []grpc.ServerOption{
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			log_middleware.Interceptor(logger),
			recover_middleware.Interceptor(),
			auth_middleware.Interceptor(),
			validate_middleware.Interceptor(),
			error_middleware.Interceptor(logger),
			prometheus_middleware.Interceptor(),
			otel_middleware.Interceptor(),
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
	service.RegisterServers(grpcServer, logger, db, tracer)
	grpc_prometheus.EnableHandlingTimeHistogram()
	grpc_prometheus.Register(grpcServer)
	reflection.Register(grpcServer)
	return grpcServer
}

func newGrpcGatewayServer(config *config.Config, logger *logrus.Logger) *http.Server {
	gwmux := runtime.NewServeMux(
		runtime.WithMetadata(readFromRequest(logger)),
	)
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
			k.log.Tracef("skipping restarts of server because app is not in running state: state is %d", k.state)
			return
		}

		k.log.Infof("%s started\n", name)
		if err = listen(); err != nil {
			if k.config.Server.RestartOnError {
				k.log.Infof("restart server failed after error on %v", err)
				continue
			}
			k.log.Infof("server failed after error on %v", err)
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
	descrition := fmt.Sprintf("grpc gateway server at %v", k.config.Server.GatewayPort)
	k.listenAndServe(descrition, listen)
}

func (k *Kernel) ListenGrpc() {
	listen := func() error {
		lis, err := net.Listen(
			"tcp", fmt.Sprintf("%s:%s", serveHost, k.config.Server.GrpcPort))
		if err != nil {
			return err
		}
		if err := k.grpcServer.Serve(lis); err != nil {
			return err
		}
		return nil
	}
	descrition := fmt.Sprintf("grpc server at %v", k.config.Server.GrpcPort)
	k.listenAndServe(descrition, listen)
}

func (k *Kernel) ListenSwagger() {
	listen := func() error {
		if err := k.swaggerServer.ListenAndServe(); err != nil {
			return err
		}
		return nil
	}
	descrition := fmt.Sprintf(
		"swagger/metrics server at %v", k.config.Server.SwaggerPort)
	k.listenAndServe(descrition, listen)
}

func (k *Kernel) Shutdown(ctx context.Context) error {
	if k.state != StateRunning {
		k.log.Warn("Application cannot be shutdown since current state is not 'running'")
		return nil
	}

	k.state = StateStopping
	defer func() {
		k.state = StateStopped
	}()

	if k.gwServer != nil {
		if err := k.gwServer.Shutdown(ctx); err != nil {
			k.log.Errorf("server shutdown error: %v\n", err)
		} else {
			k.log.Infoln("gateway server stopped")
		}
	}

	if k.grpcServer != nil {
		k.grpcServer.GracefulStop()
		k.log.Infoln("grpc server stopped")
	}

	if k.swaggerServer != nil {
		if err := k.swaggerServer.Shutdown(ctx); err != nil {
			k.log.Errorf("swagger shutdown error: %v\n", err)
		} else {
			k.log.Infoln("swagger server stopped")
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
			k.log.Errorf("shutdown function returned error: %v\n", shutdownErr)
		}
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-done:
	}

	return k.dbClient.Close()
}

func Version() string {
	return version
}
