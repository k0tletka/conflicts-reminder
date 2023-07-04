package reminder

import (
	"context"
	"gitlab.tubecorporate.com/dsp-proxy/conflicts-reminder/internal/config"
	"gitlab.tubecorporate.com/dsp-proxy/conflicts-reminder/internal/messager"
	"time"
)

const (
	defaultRemindCycleMinutes = 60 * 24
)

type ReminderService struct {
	cfg       *config.Config
	ticker    *time.Ticker
	msgSender *messager.SlackMessageSender
}

func NewReminderService(cfg *config.Config) *ReminderService {
	s := &ReminderService{
		cfg:       cfg,
		msgSender: messager.NewSlackMessageSender(cfg),
	}

	remindCycleMinutes := defaultRemindCycleMinutes
	if s.cfg.RemindCycleMinutes != 0 {
		remindCycleMinutes = s.cfg.RemindCycleMinutes
	}

	s.ticker = time.NewTicker(time.Duration(remindCycleMinutes) * time.Minute)

	return s
}

func (s *ReminderService) StartService(ctx context.Context) error {
	err := s.cycle(ctx)
	s.shutdown()
	return err
}

func (s *ReminderService) cycle(ctx context.Context) error {
	for {
		select {
		case <-s.ticker.C:
			// TODO: Use conflict service to notify about new conflicts
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (s *ReminderService) shutdown() {
}
