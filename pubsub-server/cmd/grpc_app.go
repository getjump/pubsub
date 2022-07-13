package main

import (
	"context"
	"net"
	"net/http"
	"os"
	"os/signal"

	grpc_api "github.com/getjump/pubsub/pubsub-server/api/grpc"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"

	// grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/heptiolabs/healthcheck"

	"github.com/getjump/pubsub/pubsub-server/internal"
	"github.com/getjump/pubsub/pubsub-server/pkg/proto"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func main() {
	// Config
	v := viper.GetViper()
	conf, err := internal.GetViperConfig(os.Getenv("BUILD_ENV"), v)
	if err != nil {
		panic(err)
	}

	if !internal.IsProduction(conf) {
		conf.Debug()
	}

	// Prometheus
	http.Handle(conf.GetString("prometheus.route"), promhttp.Handler())

	// Health Check
	healthChecker := healthcheck.NewMetricsHandler(prometheus.DefaultRegisterer, "health_check")
	healthChecker.AddLivenessCheck(
		"Goroutine Threshold",
		healthcheck.GoroutineCountCheck(conf.GetInt("health_check.goroutine_threshold")),
	)

	http.HandleFunc(conf.GetString("health_check.route.group")+conf.GetString("health_check.route.live"), healthChecker.LiveEndpoint)
	http.HandleFunc(conf.GetString("health_check.route.group")+conf.GetString("health_check.route.ready"), healthChecker.ReadyEndpoint)

	// Logger
	zapConfig := zap.NewProductionConfig()
	zapConfig.Level.UnmarshalText([]byte(conf.GetString("log.level")))
	zapConfig.Development = !internal.IsProduction(conf)

	logger, _ := zapConfig.Build()
	defer logger.Sync()

	grpcListener, err := net.Listen("tcp", ":"+conf.GetString("app.grpc.port"))

	if err != nil || grpcListener == nil {
		logger.Sugar().Error("Failed to allocate port %s: %v", conf.GetString("app.grpc.port"), grpcListener)
	}

	grpcServer := grpc.NewServer(
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			grpc_ctxtags.StreamServerInterceptor(),
			grpc_prometheus.StreamServerInterceptor,
			grpc_zap.StreamServerInterceptor(logger),
			// grpc_recovery.StreamServerInterceptor(),
		)),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_ctxtags.UnaryServerInterceptor(),
			grpc_prometheus.UnaryServerInterceptor,
			grpc_zap.UnaryServerInterceptor(logger),
			// grpc_recovery.UnaryServerInterceptor(),
		)),
	)
	defer grpcServer.GracefulStop()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()

	go func() {
		<-ctx.Done()
		logger.Info("Graceful shutdown")
		grpcServer.GracefulStop()
	}()

	grpcServer.RegisterService(&proto.PubSub_ServiceDesc, grpc_api.NewPubSubServer(
		ctx,
		conf.GetDuration("app.retry_send_timeout"),
		conf.GetUint("app.max_retries"),
	))

	logger.Info("Starting gRPC server")

	grpcServer.Serve(grpcListener)
}
