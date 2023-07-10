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
	SlackID           string
	BranchName        string
	MergeRequestLink  string
	MergeRequestTitle string
}

type SlackMessageSender struct {
	cfg *config.Config
	api *slack.Client
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
