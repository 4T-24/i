package env

import (
	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

var cfg Config

func Load() {
	godotenv.Load()

	if err := env.Parse(&cfg); err != nil {
		logrus.Fatalf("failed to load env: %v", err)
	}

	switch cfg.LogLevel {
	case "debug":
		logrus.SetLevel(logrus.DebugLevel)
	case "info":
		logrus.SetLevel(logrus.InfoLevel)
	case "warn":
		logrus.SetLevel(logrus.WarnLevel)
	case "error":
		logrus.SetLevel(logrus.ErrorLevel)
	default:
		logrus.SetLevel(logrus.InfoLevel)
	}

	logrus.WithField("level", cfg.LogLevel).Info("Successfully set log level")
}
