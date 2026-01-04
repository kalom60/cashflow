package testutils

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/kalom60/cashflow/initiator"
	"github.com/kalom60/cashflow/internal/constant/model/persistencedb"
	"github.com/kalom60/cashflow/platform/logger"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func SetupTestDB() persistencedb.PersistenceDB {
	configName := "config"
	if os.Getenv("CONFIG_NAME") != "" {
		configName = os.Getenv("CONFIG_NAME")
	}

	logs, err := zap.NewProduction()
	if err != nil {
		log.Fatal("unable to start logger")
	}

	err = initiator.InitConfig(initiator.Config{Name: configName, Path: "../../../config", Logger: logs})
	if err != nil {
		zap.Error(err)
	}

	dsn := viper.GetString("test_db.url")

	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		zap.Error(err)
	}

	idleConnTimeout := viper.GetDuration("database.idle_conn_timeout")
	if idleConnTimeout == 0 {
		idleConnTimeout = 4 * time.Minute
	}

	config.MaxConnIdleTime = idleConnTimeout
	conn, err := pgxpool.ConnectConfig(context.Background(), config)
	if err != nil {
		log.Fatalf("Unable to connect to test database: %v", err)
	}

	CleanTables(conn)

	logger := NewTestLogger()
	persistenceDB := persistencedb.New(conn, logger)

	return persistenceDB
}

func NewTestLogger() logger.Logger {
	logger := initiator.InitLogger()
	return logger
}

func CleanTables(conn *pgxpool.Pool) {
	ctx := context.Background()

	_, err := conn.Exec(
		ctx,
		`TRUNCATE TABLE
			payments, outbox_events
		RESTART IDENTITY CASCADE
	`)
	if err != nil {
		log.Fatalf("failed to clean test tables: %v", err)
	}
}
