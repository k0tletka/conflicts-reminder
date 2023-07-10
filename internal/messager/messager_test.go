package messager

import (
	"context"
	"github.com/slack-go/slack"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gitlab.tubecorporate.com/dsp-proxy/conflicts-reminder/internal/config"
	"gitlab.tubecorporate.com/dsp-proxy/conflicts-reminder/internal/messager/mocks"
	"testing"
)

const (
	conflicts2result = `При проверке были обнаружены следующие конфликты💩:
    <@UG29420>: конфликт в ветке master: <https://mrlink.1|Sample merge request>
    <@UG29550>: конфликт в ветке developer-1: <https://mrlink.2|Sample merge request to developer branch>

Кто до завтра их не исправит, тот крокодил🐊
`
)

func TestSlackMessageSender_SendMessageWithConflictsData(t *testing.T) {
	tests := []struct {
		name        string
		messageData *SlackConflictsData
		cfg         *config.Config
		setupMock   func(*mocks.SlackClient)
	}{
		{
			name: "2 conflicts",
			messageData: &SlackConflictsData{
				Rows: []SlackConflictData{
					{
						SlackID:           "UG29490",
						BranchName:        "master",
						MergeRequestTitle: "Sample merge request",
						MergeRequestLink:  "https://mrlink.1",
					},
					{
						SlackID:           "UG29550",
						BranchName:        "developer-1",
						MergeRequestTitle: "Sample merge request to developer branch",
						MergeRequestLink:  "https://mrlink.2",
					},
				},
			},
			cfg: &config.Config{
				Slack: config.SlackConfig{
					NotificationChannelID: "1234",
				},
			},
			setupMock: func(m *mocks.SlackClient) {
				m.On("SendMessageContext",
					mock.Anything,
					"1234",
					slack.MsgOptionText(conflicts2result, false),
				).Return("", "", "", nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clientMock := mocks.NewSlackClient(t)
			slackClient := &SlackMessageSender{
				cfg: tt.cfg,
				api: clientMock,
			}

			tt.setupMock(clientMock)
			err := slackClient.SendMessageWithConflictsData(context.Background(), tt.messageData)
			assert.NoError(t, err)
		})
	}
}
