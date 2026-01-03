package initiator

import (
	"fmt"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type Config struct {
	Name   string
	Path   string
	Type   string
	Logger *zap.Logger
}

func InitConfig(config Config) error {
	if config.Logger == nil {
		return fmt.Errorf("logger cannot be nil")
	}

	if config.Name == "" || config.Path == "" {
		return fmt.Errorf("config name and path cannot be empty")
	}

	if config.Type == "" {
		config.Type = "yaml"
	}

	viper.SetConfigName(config.Name)
	viper.AddConfigPath(config.Path)
	viper.SetConfigType(config.Type)

	viper.AutomaticEnv()
	viper.SetEnvPrefix("APPLICATION")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	const maxRetries = 3
	var err error
	for attempt := 1; attempt <= maxRetries; attempt++ {
		err = viper.ReadInConfig()
		if err == nil {
			break
		}
		config.Logger.Warn("Failed to read config",
			zap.Int("attempt", attempt),
			zap.Error(err),
		)
		if attempt < maxRetries {
			time.Sleep(time.Second * time.Duration(attempt))
			continue
		}
		config.Logger.Fatal("Failed to read config after retries",
			zap.Int("maxRetries", maxRetries),
			zap.Error(err),
		)
	}

	if err := validateConfig(); err != nil {
		return fmt.Errorf("configuration validation failed: %v", err)
	}

	if err := setupConfigWatcher(config.Logger); err != nil {
		config.Logger.Warn("Failed to setup config watcher",
			zap.Error(err),
		)
	}

	config.Logger.Info("Configuration initialized successfully",
		zap.String("config_name", config.Name),
		zap.String("config_path", config.Path),
	)

	return nil
}

func validateConfig() error {

	requiredKeys := []string{}

	for _, key := range requiredKeys {
		if !viper.IsSet(key) || viper.GetString(key) == "" {
			return fmt.Errorf("missing required configuration key: %s", key)
		}
	}
	return nil
}

func setupConfigWatcher(log *zap.Logger) error {
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Info("Configuration file changed",
			zap.String("file", e.Name),
			zap.String("operation", e.Op.String()),
		)

		if err := validateConfig(); err != nil {
			log.Error("Invalid configuration after change",
				zap.Error(err),
			)
		}
	})
	return nil
}
