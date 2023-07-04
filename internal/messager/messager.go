package messager

import (
	"bytes"
	"context"
	_ "embed"
	"github.com/slack-go/slack"
	"gitlab.tubecorporate.com/dsp-proxy/conflicts-reminder/internal/config"
	"text/template"
)

const (
	textNoConflict = "Конликтов не обнаружено, все молодцы!🏆"
)

//go:embed template.file
var templateText string

var (
	conflictsMessageTemplate = template.Must(template.New("conflicts_template").Parse(templateText))
)

type SlackConflictsData struct {
	Rows []SlackConflictData
}

type SlackConflictData struct {
	SlackNicknamePing string
	BranchName        string
	MergeRequestLink  string
}

type SlackMessageSender struct {
	cfg *config.Config
	api *slack.Client

	userConfiguration map[string]*slack.User
}

func NewSlackMessageSender(cfg *config.Config) *SlackMessageSender {
	slackApi := slack.New(cfg.Slack.Token)

	return &SlackMessageSender{
		cfg: cfg,
		api: slackApi,
	}
}

func (s *SlackMessageSender) SendMessageWithConflictsData(ctx context.Context, messageData *SlackConflictsData) error {
	messageBuffer := &bytes.Buffer{}
	if err := conflictsMessageTemplate.ExecuteTemplate(messageBuffer, "", messageData); err != nil {
		return err
	}

	_, _, _, err := s.api.SendMessageContext(
		ctx,
		s.cfg.Slack.NotificationChannelID,
		slack.MsgOptionText(messageBuffer.String(), false),
	)

	if err != nil {
		return err
	}

	return nil
}

func (s *SlackMessageSender) SendMessageNoConflicts(ctx context.Context) error {
	_, _, _, err := s.api.SendMessageContext(
		ctx,
		s.cfg.Slack.NotificationChannelID,
		slack.MsgOptionText(textNoConflict, false),
	)

	return err
}

func (s *SlackMessageSender) GetPingNameBySlackId(slackId string) string {
	user, _ := s.userConfiguration[slackId]

	if user != nil {
		return user.Name
	} else {
		return ""
	}
}

func (s *SlackMessageSender) extractUserConfiguration() error {
	usersToCheck := make(map[string]struct{})

	for _, bindCfg := range s.cfg.Binds {
		usersToCheck[bindCfg.SlackID] = struct{}{}
	}

	allUsers, err := s.api.GetUsers()
	if err != nil {
		return err
	}

	for _, user := range allUsers {
		if _, ok := usersToCheck[user.ID]; ok && !user.Deleted {
			s.userConfiguration[user.ID] = &user
		}
	}

	return nil
}
