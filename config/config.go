package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"

	"github.com/vucongthanh92/go-base-utils/models"
)

// AppConfig mirrors the common service configuration shape in chat-service.
type AppConfig struct {
	ServiceName string                    `mapstructure:"serviceName"`
	Development bool                      `mapstructure:"development"`
	Logger      *models.LoggerConfig      `mapstructure:"logger"`
	Http        *models.HttpConfig        `mapstructure:"http"`
	GRPC        *models.GrpcConfig        `mapstructure:"grpc"`
	Client      *models.ClientConfig      `mapstructure:"client"`
	Database    *models.DatabaseConfig    `mapstructure:"database"`
	Tracing     *models.TracingConfig     `mapstructure:"tracing"`
	Redis       *models.RedisConfig       `mapstructure:"redis"`
	Kafka       *models.KafkaConfigV2     `mapstructure:"kafka"`
	Healthcheck *models.HealthcheckConfig `mapstructure:"healthcheck"`
	Metrics     *models.MetricsConfig     `mapstructure:"metrics"`
	CronJob     *models.CronJob           `mapstructure:"cronjob"`
	Loki        *models.LokiConfig        `mapstructure:"loki"`
	SePay       models.SePayConfig        `mapstructure:"sepay"`
}

func Load(path string) (*AppConfig, error) {
	v := viper.New()
	v.SetConfigFile(path)
	v.SetEnvPrefix("PAYMENT_GATEWAY")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()
	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}

	var cfg AppConfig
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("decode config: %w", err)
	}

	if cfg.Database != nil {
		if cfg.Database.ReadDbCfg != nil {
			cfg.Database.ReadDbCfg.ConnectionString = os.ExpandEnv(cfg.Database.ReadDbCfg.ConnectionString)
		}
		if cfg.Database.WriteDbCfg != nil {
			cfg.Database.WriteDbCfg.ConnectionString = os.ExpandEnv(cfg.Database.WriteDbCfg.ConnectionString)
		}
	}

	cfg.SePay.MerchantID = os.ExpandEnv(cfg.SePay.MerchantID)
	cfg.SePay.SecretKey = os.ExpandEnv(cfg.SePay.SecretKey)
	cfg.SePay.CheckoutURL = os.ExpandEnv(cfg.SePay.CheckoutURL)
	cfg.SePay.SuccessURL = os.ExpandEnv(cfg.SePay.SuccessURL)
	cfg.SePay.ErrorURL = os.ExpandEnv(cfg.SePay.ErrorURL)
	cfg.SePay.CancelURL = os.ExpandEnv(cfg.SePay.CancelURL)
	cfg.SePay.WebhookSecret = os.ExpandEnv(cfg.SePay.WebhookSecret)

	for i := range cfg.SePay.ReceivingAccountNumbers {
		cfg.SePay.ReceivingAccountNumbers[i] = os.ExpandEnv(cfg.SePay.ReceivingAccountNumbers[i])
	}

	if cfg.SePay.TimestampToleranceSec == 0 {
		cfg.SePay.TimestampToleranceSec = 300
	}

	if cfg.ServiceName == "" ||
		cfg.Http == nil ||
		cfg.Http.Port == "" ||
		cfg.Database == nil ||
		cfg.Database.ReadDbCfg == nil ||
		cfg.Database.WriteDbCfg == nil ||
		cfg.Database.ReadDbCfg.ConnectionString == "" ||
		cfg.Database.WriteDbCfg.ConnectionString == "" {
		return nil, fmt.Errorf("serviceName, http.port, database.readDb and database.writeDb are required")
	}

	if cfg.SePay.Environment != "sandbox" {
		return nil, fmt.Errorf("only sepay sandbox is allowed in this phase")
	}

	return &cfg, nil
}
