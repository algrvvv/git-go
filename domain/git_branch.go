package domain

type GitBranch struct {
	Name string
	Meta string
	Date string
}

func NewGitBranch(name, meta string) *GitBranch {
	return &GitBranch{Name: name, Meta: meta}
}

func (b *GitBranch) SetDate(date string) *GitBranch {
	b.Date = date
	return b
}

func (b *GitBranch) String() string {
	return b.Name
}
