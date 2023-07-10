package reminder

import (
	"context"
	"gitlab.tubecorporate.com/dsp-proxy/conflicts-reminder/internal/config"
	"gitlab.tubecorporate.com/dsp-proxy/conflicts-reminder/internal/conflict"
	"gitlab.tubecorporate.com/dsp-proxy/conflicts-reminder/internal/messager"
	"time"
)

const (
	defaultRemindCycleMinutes = 60 * 24
)

type ReminderService struct {
	cfg               *config.Config
	ticker            *time.Ticker
	msgSender         *messager.SlackMessageSender
	conflictsDetector *conflict.GitlabConflictDetector
}

func NewReminderService(cfg *config.Config) (*ReminderService, error) {
	conflictsDetector, err := conflict.NewGitlabConflictDetector(cfg)
	if err != nil {
		return nil, err
	}

	s := &ReminderService{
		cfg:               cfg,
		msgSender:         messager.NewSlackMessageSender(cfg),
		conflictsDetector: conflictsDetector,
	}

	remindCycleMinutes := defaultRemindCycleMinutes
	if s.cfg.RemindCycleMinutes != 0 {
		remindCycleMinutes = s.cfg.RemindCycleMinutes
	}

	s.ticker = time.NewTicker(time.Duration(remindCycleMinutes) * time.Minute)
	return s, nil
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
			if err := s.remindAboutConflicts(ctx); err != nil {
				return err
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (s *ReminderService) remindAboutConflicts(ctx context.Context) error {
	conflicts, err := s.conflictsDetector.DetectConflicts(ctx)
	if err != nil {
		return err
	}

	if len(conflicts.Conflicts) == 0 {
		// There is no conflicts, we need to send successfull message
		return s.msgSender.SendMessageNoConflicts(ctx)
	}

	slackConflictsData := s.convertConflictDataForSlack(conflicts)
	return s.msgSender.SendMessageWithConflictsData(ctx, slackConflictsData)
}

func (s *ReminderService) convertConflictDataForSlack(conflictsData *conflict.ConflictsData) *messager.SlackConflictsData {
	slackReminderMessage := &messager.SlackConflictsData{
		make([]messager.SlackConflictData, 0, len(conflictsData.Conflicts)),
	}

	for _, conflictData := range conflictsData.Conflicts {
		slackReminderMessage.Rows = append(slackReminderMessage.Rows, messager.SlackConflictData{
			SlackID:           s.cfg.GetSlackIDByGitlabID(conflictData.AuthorID),
			BranchName:        conflictData.BranchName,
			MergeRequestLink:  conflictData.MergeRequestURL,
			MergeRequestTitle: conflictData.MergeRequestTitle,
		})
	}

	return slackReminderMessage
}

func (s *ReminderService) shutdown() {
}
