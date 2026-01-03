package initiator

import (
	"os"
	"path/filepath"
	"time"

	"github.com/kalom60/cashflow/platform/logger"
	"github.com/spf13/viper"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func getLogFilePath(level string) string {
	now := time.Now()
	baseDir := "logs"
	year := now.Format("2006")
	month := now.Format("01")
	day := now.Format("02")

	dir := filepath.Join(baseDir, year, month, day, level)
	os.MkdirAll(dir, os.ModePerm) // create directories if not exist

	return filepath.Join(dir, level+".log")
}

func getZapCore(filePath string, level zapcore.Level) zapcore.Core {
	cfg := zap.NewProductionEncoderConfig()
	cfg.TimeKey = "time"
	cfg.EncodeTime = zapcore.ISO8601TimeEncoder
	encoder := zapcore.NewJSONEncoder(cfg)

	file, _ := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	writeSyncer := zapcore.AddSync(file)

	return zapcore.NewCore(encoder, writeSyncer, level)
}

func NewFileLoggerWithLevel() *zap.Logger {
	lvl := zapcore.Level(viper.GetInt("logger.level"))

	// File paths
	infoPath := getLogFilePath("info")
	errorPath := getLogFilePath("error")

	// JSON encoder for files
	fileEncoderConfig := zap.NewProductionEncoderConfig()
	fileEncoderConfig.TimeKey = "time"
	fileEncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	fileEncoder := zapcore.NewJSONEncoder(fileEncoderConfig)

	// Console encoder (human-readable)
	consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())

	// File writers
	infoFile, _ := os.OpenFile(infoPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	errorFile, _ := os.OpenFile(errorPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)

	infoSync := zapcore.AddSync(infoFile)
	errorSync := zapcore.AddSync(errorFile)
	consoleSync := zapcore.AddSync(os.Stdout)

	// Info core → only logs >= lvl and < Error
	infoCore := zapcore.NewCore(fileEncoder, infoSync,
		zap.LevelEnablerFunc(func(l zapcore.Level) bool {
			return l >= lvl && l < zapcore.ErrorLevel
		}),
	)

	// Error core → logs >= Error
	errorCore := zapcore.NewCore(fileEncoder, errorSync,
		zap.LevelEnablerFunc(func(l zapcore.Level) bool {
			return l >= zapcore.ErrorLevel
		}),
	)

	// Console core → logs everything >= lvl
	consoleCore := zapcore.NewCore(consoleEncoder, consoleSync,
		zap.LevelEnablerFunc(func(l zapcore.Level) bool {
			return l >= lvl
		}),
	)

	// Combine all outputs
	combinedCore := zapcore.NewTee(infoCore, errorCore, consoleCore)

	return zap.New(combinedCore, zap.AddCaller(), zap.AddCallerSkip(1))
}

func InitLogger() logger.Logger {
	lg := NewFileLoggerWithLevel()
	return logger.New(lg)
}
