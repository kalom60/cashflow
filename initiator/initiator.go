package initiator

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/spf13/viper"
	echoSwagger "github.com/swaggo/echo-swagger"
	"go.uber.org/zap"
)

func Initiate() {

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

	// initializing migration
	logger.Info(ctx, "initializing migration")
	InitMigration(viper.GetString("db.url"), viper.GetString("db.migration_path"))
	logger.Info(ctx, "done initializing migration")

	logger.Info(ctx, "initializing http server")
	server := echo.New()
	echosrv := server.Group("")
	echosrv.GET("/swagger/*any", echoSwagger.EchoWrapHandler())

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

		log.Fatal("HTTP server Shutdown")
	}()

	host := fmt.Sprint(viper.GetString("app.host"), ":", viper.GetInt("app.port"))
	logger.Info(ctx, "server listening at port ", zap.Any("link", host))
	err = srv.ListenAndServe()
	if err != nil {
		logger.Fatal(ctx, fmt.Sprintf("Could not start HTTP server: %s", err))
	}
}
