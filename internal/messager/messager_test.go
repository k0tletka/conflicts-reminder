package messager

import (
	"bytes"
	_ "embed"
	"github.com/stretchr/testify/assert"
	"testing"
)

//go:embed template_test.file
var conflicts2result string

func TestSlackMessageSender_CheckTemplate(t *testing.T) {
	messageData := &SlackConflictsData{
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
	}

	messageBuffer := &bytes.Buffer{}
	if err := conflictsMessageTemplate.ExecuteTemplate(messageBuffer, "conflicts_template", messageData); err != nil {
		t.Error(err)
	}

	assert.Equal(t, conflicts2result, messageBuffer.String())
}
