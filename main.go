package main

import (
	"context"
	"os"
	"os/signal"
	"time"

	"github.com/Systemnick/user-service/config"
	"github.com/rs/zerolog"
)

const serviceStoppingTimeout = 60*time.Second

func main() {
	c := config.GetConfigFromEnv()

	logger := initLogger(c.LogLevel)

	application, err := NewApplication(c, logger)
	if err != nil {
		logger.Fatal().Err(err).Msg("Couldn't run application")
		return
	}

	go func() {
		if err := application.Run(); err != nil {
			logger.Fatal().Err(err).Msg("Service run error")
			return
		}
	}()
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	logger.Debug().Msg("Get signal for shutdown")

	ctx, cancel := context.WithTimeout(context.Background(), serviceStoppingTimeout)
	defer cancel()
	if err := application.Stop(ctx); err != nil {
		logger.Fatal().Err(err).Msg("Service shutdown error")
		return
	}
}

// Possible values of level: debug, info, warn, error, fatal, panic
// Pass empty string for debug level by default
func initLogger(level string) *zerolog.Logger {
	zerolog.TimeFieldFormat = time.RFC3339Nano
	zerolog.DurationFieldInteger = true
	zerolog.TimestampFieldName = "timestamp"

	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()

	if level == "" {
		return &logger
	}

	l, err := zerolog.ParseLevel(level)
	if err != nil {
		zerolog.SetGlobalLevel(l)
	} else {
		logger.Debug().Err(err).Msg("Default log level is debug")
	}

	return &logger
}
