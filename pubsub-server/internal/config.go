package internal

import (
	"strings"
	"time"

	"github.com/spf13/viper"
)

const (
	defaultGRPCPort                      = 9000
	defaultHealthCheckGoroutineThreshold = 100
	defaultRetrySendTimeout              = time.Second * 10
	defaultMaxRetries                    = 5
)

func GetViperConfig(env string, conf *viper.Viper) (*viper.Viper, error) {
	conf.SetDefault("environment", env)

	conf.SetDefault("log.level", "debug")

	// App
	conf.SetDefault("app.name", "pubsub")
	conf.SetDefault("app.grpc.port", defaultGRPCPort)
	conf.SetDefault("app.grpc.http_route_prefix", "/v1")

	conf.SetDefault("app.retry_send_timeout", defaultRetrySendTimeout)
	conf.SetDefault("app.max_retries", defaultMaxRetries)

	// Monitoring
	conf.SetDefault("health_check.route.group", "/health")
	conf.SetDefault("health_check.route.live", "/live")
	conf.SetDefault("health_check.route.ready", "/ready")
	conf.SetDefault(
		"health_check.goroutine_threshold", defaultHealthCheckGoroutineThreshold,
	) // Alert when there are more than X Goroutine

	// Prometheus
	conf.SetDefault("prometheus.route", "/metrics")

	// Swagger
	conf.SetDefault("swagger.ui.route.group", "/swaggerui/")
	conf.SetDefault("swagger.json.route.group", "/swagger/")

	conf.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "_", "__")) // APP_DATA__BASE_PASS -> app.data_base.pass
	conf.AutomaticEnv()

	conf.SetConfigName(env)
	conf.AddConfigPath("./configs/" + env)
	conf.ReadInConfig()

	return conf, nil
}

func IsProduction(conf *viper.Viper) bool {
	return conf.GetString("environment") != "production"
}
