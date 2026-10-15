package healthcheck

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/vucongthanh92/courier-payment-gateway/config"
	"github.com/vucongthanh92/courier-payment-gateway/database"
	"github.com/vucongthanh92/go-base-utils/logger"
	"go.uber.org/zap"
)

func RunHealthCheck(
	ctx context.Context,
	cfg *config.AppConfig,
	readDb database.GormReadDb,
	writeDb database.GormWriteDb,
) func() {
	return func() {
		mux := http.NewServeMux()
		mux.HandleFunc("/live", func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("ok"))
		})
		mux.HandleFunc("/ready", func(w http.ResponseWriter, _ *http.Request) {
			readyCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
			defer cancel()
			if err := pingDatabases(readyCtx, readDb, writeDb); err != nil {
				logger.Warn("healthcheck readiness failed", zap.Error(err))
				http.Error(w, "not ready", http.StatusServiceUnavailable)
				return
			}
			_, _ = w.Write([]byte("ok"))
		})
		logger.Info("Healthcheck server listening on port", zap.String("Port", cfg.Healthcheck.Port))
		if err := http.ListenAndServe(cfg.Healthcheck.Port, mux); err != nil {
			logger.Warn("Healthcheck server", zap.Error(err))
		}
	}
}

func pingDatabases(ctx context.Context, readDb database.GormReadDb, writeDb database.GormWriteDb) error {
	for name, db := range map[string]database.GormReadDb{"read": readDb, "write": database.GormReadDb(writeDb)} {
		if db == nil {
			return fmt.Errorf("%s database is not initialized", name)
		}
		sqlDb, err := (*db).DB()
		if err != nil {
			return err
		}
		if err := sqlDb.PingContext(ctx); err != nil {
			return err
		}
	}
	return nil
}
