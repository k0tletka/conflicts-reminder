package conflict

import (
	"context"
	"github.com/AlekSi/pointer"
	"github.com/xanzy/go-gitlab"
	"gitlab.tubecorporate.com/dsp-proxy/conflicts-reminder/internal/config"
	"path"
)

type ConflictsData struct {
	Conflicts []ConflictData
}

type ConflictData struct {
	AuthorID          int
	MergeRequestURL   string
	MergeRequestTitle string
	BranchName        string
}

type GitlabConflictDetector struct {
	cfg          *config.Config
	gitlabClient *gitlab.Client
}

func NewGitlabConflictDetector(cfg *config.Config) (*GitlabConflictDetector, error) {
	gitlabClient, err := gitlab.NewClient(
		cfg.Gitlab.Token,
		gitlab.WithBaseURL(path.Join(cfg.Gitlab.GitlabAddress, "api/v4")),
	)

	if err != nil {
		return nil, err
	}

	d := &GitlabConflictDetector{
		cfg:          cfg,
		gitlabClient: gitlabClient,
	}

	return d, nil
}

func (d *GitlabConflictDetector) DetectConflicts(ctx context.Context) (*ConflictsData, error) {
	listMRsCfg := &gitlab.ListProjectMergeRequestsOptions{
		State: pointer.To("opened"),
		WIP:   pointer.To("no"),
	}

	mrs, _, err := d.gitlabClient.MergeRequests.ListProjectMergeRequests(
		d.cfg.Gitlab.ProjectName,
		listMRsCfg,
		gitlab.WithContext(ctx),
	)
	if err != nil {
		return nil, err
	}

	conflicts := &ConflictsData{Conflicts: []ConflictData{}}

	for _, mr := range mrs {
		if !mr.HasConflicts {
			continue
		}

		if !d.cfg.CheckGitlabId(mr.Author.ID) {
			continue
		}

		conflicts.Conflicts = append(conflicts.Conflicts, ConflictData{
			AuthorID:          mr.Author.ID,
			MergeRequestURL:   mr.WebURL,
			MergeRequestTitle: mr.Title,
			BranchName:        mr.SourceBranch,
		})
	}

	return conflicts, nil
}
