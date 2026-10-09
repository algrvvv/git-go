package domain

import "time"

type GitBranch struct {
	Name string
	Desc string
}

type GitBranchOptions struct {
	Name        string
	Description string

	LastCommitAt time.Time
	LastCheckout time.Time

	IsRemote  bool
	IsMerged  bool
	IsCurrent bool
}

func NewGitBranch(options ...GitBranchOptions) *GitBranch {
	branch := &GitBranch{}

	if len(options) == 1 {
		option := options[0]
		branch.Name = option.Name
		branch.Desc = option.Description
	}

	return branch
}
