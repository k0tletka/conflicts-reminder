package main

import (
	"context"
	"gitlab.tubecorporate.com/dsp-proxy/conflicts-reminder/internal/config"
	"gitlab.tubecorporate.com/dsp-proxy/conflicts-reminder/internal/reminder"
	"go.uber.org/zap"
	"os/signal"
	"syscall"
)

func main() {
	logger, _ := zap.NewProduction()
	zap.ReplaceGlobals(logger)

	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGTERM,
		syscall.SIGKILL,
		syscall.SIGINT,
		syscall.SIGHUP,
	)
	defer cancel()

	cfg := config.ReadConfig()
	slackService := reminder.NewReminderService(cfg)

	if err := slackService.StartService(ctx); err != nil && err != context.Canceled {
		zap.S().With("err", err).Error("Reminder exited with errors")
	}
}
