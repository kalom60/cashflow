package initiator

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kalom60/cashflow/docs"
	"github.com/kalom60/cashflow/internal/constant/model/persistencedb"
	"github.com/kalom60/cashflow/platform/messaging"
	"github.com/kalom60/cashflow/platform/workerpool"
	"github.com/labstack/echo/v4"

	"github.com/spf13/viper"
	echoSwagger "github.com/swaggo/echo-swagger"
	"go.uber.org/zap"
)

func Initiate() {
	docs.SwaggerInfo.Title = "Cashflow API"
	docs.SwaggerInfo.Description = "API documentation for Cashflow"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.BasePath = "/"

	ctx := context.Background()

	log, err := zap.NewProduction()
	if err != nil {
		log.Fatal("unable to start logging")
	}

	configName := "config"
	if os.Getenv("CONFIG_NAME") != "" {
		configName = os.Getenv("CONFIG_NAME")
	}

	err = InitConfig(Config{Name: configName, Path: "config", Logger: log})
	if err != nil {
		log.Fatal("unable to start config", zap.Error(err))
	}

	log.Info("initializing config completed")

	logger := InitLogger()
	log.Info("initializing logger completed")

	// initailizing database connection
	log.Info("initializing database connect")
	pgxPool := initDB("cashflow", logger)
	log.Info("database connection initialized")

	// initializing migration
	logger.Info(ctx, "initializing migration")
	InitMigration(viper.GetString("db.url"), viper.GetString("db.migration_path"))
	logger.Info(ctx, "done initializing migration")

	logger.Info(ctx, "initializing persistence layer ")
	persistenceDB := persistencedb.New(pgxPool, logger)
	persistence := initPersistence(&persistenceDB, logger)
	logger.Info(ctx, "done initializing persistence layer")

	logger.Info(ctx, "initializing waorker pool")
	wp := workerpool.New(viper.GetInt("workerpool.max_workers"), viper.GetInt("workerpool.task_buffer"))
	wp.Start()
	logger.Info(ctx, "done initializing worker pool")

	logger.Info(ctx, "initializing rabbitmq client")
	rabbitMQURL := viper.GetString("rabbitmq.url")
	msgClient, err := messaging.NewRabbitMQClient(rabbitMQURL)
	if err != nil {
		logger.Fatal(ctx, "failed to initialize RabbitMQ client", zap.Error(err))
	}
	logger.Info(ctx, "rabbitmq client initialized")

	logger.Info(ctx, "initializing module layer")
	module := initModule(persistence, msgClient, logger, wp)
	logger.Info(ctx, "done initializing module layer")

	logger.Info(ctx, "initializing handler layer ")
	handler := initHandler(module, logger)
	logger.Info(ctx, "done initializing handler layer")

	logger.Info(ctx, "initializing http server")
	server := echo.New()
	echosrv := server.Group("")
	echosrv.GET("/swagger/*any", echoSwagger.EchoWrapHandler())

	logger.Info(ctx, "initializing route")
	initRoute(echosrv, handler, logger)
	logger.Info(ctx, "done initializing route")

	logger.Info(ctx, "done initializing server")

	srv := &http.Server{
		Addr:              fmt.Sprintf("%s:%d", viper.GetString("app.host"), viper.GetInt("app.port")),
		Handler:           server,
		ReadHeaderTimeout: viper.GetDuration("app.timeout"),
		IdleTimeout:       30 * time.Minute,
	}

	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, syscall.SIGINT, syscall.SIGTERM)
		<-sigint

		log.Info("Shutting down... closing RabbitMQ client")
		if err := msgClient.Close(); err != nil {
			log.Error("failed to close RabbitMQ client", zap.Error(err))
		}

		log.Fatal("HTTP server Shutdown")
	}()

	host := fmt.Sprint(viper.GetString("app.host"), ":", viper.GetInt("app.port"))
	logger.Info(ctx, "server listening at port ", zap.Any("link", host))
	err = srv.ListenAndServe()
	if err != nil {
		logger.Fatal(ctx, fmt.Sprintf("Could not start HTTP server: %s", err))
	}
}
