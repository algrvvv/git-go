package domain

type GitBranch struct {
	Name string
	Meta string
}

func NewGitBranch(name, meta string) *GitBranch {
	return &GitBranch{Name: name, Meta: meta}
}
