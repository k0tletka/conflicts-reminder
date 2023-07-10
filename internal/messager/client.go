package messager

import (
	"context"
	"github.com/slack-go/slack"
)

//go:generate mockery --name=SlackClient --case=snake
type SlackClient interface {
	SendMessageContext(context.Context, string, ...slack.MsgOption) (string, string, string, error)
}
