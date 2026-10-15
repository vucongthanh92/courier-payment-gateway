package startup

import (
	"context"
	"flag"
	"log"
	"time"

	"github.com/vucongthanh92/courier-payment-gateway/config"
	"github.com/vucongthanh92/courier-payment-gateway/database"
	healthcheck "github.com/vucongthanh92/courier-payment-gateway/helper/healthcheck"
	"github.com/vucongthanh92/courier-payment-gateway/helper/utils"
	"github.com/vucongthanh92/courier-payment-gateway/internal"
	"github.com/vucongthanh92/courier-payment-gateway/redis"
	"github.com/vucongthanh92/go-base-utils/tracing"
)

func Execute() {
	configPath := flag.String("config", "./config/local/config.yaml", "path to config file")
	flag.Parse()
	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatal(err)
	}
	if cfg.Tracing != nil {
		tracing.UseOpenTelemetry(tracing.Config(*cfg.Tracing))
	}

	utils.InitSnowflake(utils.SnowflakeConfig{
		MachineID:   3,
		CustomEpoch: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
	})

	readDB, writeDB := database.GetConnectByGorm(cfg.Database)
	if readDB == nil || writeDB == nil {
		log.Fatal("open payment-gateway database connections")
	}
	container := internal.InitializeContainer(cfg, &readDB, &writeDB, redis.Open(cfg.Redis))
	go healthcheck.RunHealthCheck(context.Background(), cfg, readDB, writeDB)()
	go container.GrpcServer.Run()
	container.HttpServer.Run()
}
