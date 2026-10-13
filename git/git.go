package git

import "github.com/algrvvv/git-go/domain"

type Git interface {
	Brancher
	MetaBrancher
}

type Brancher interface {
	ListBranches() ([]*domain.GitBranch, error)
	GetCurrentBranchName() (string, error)
	Checkout(branch *domain.GitBranch) error
}

type MetaBrancher interface {
	CreateBranch(name string, meta string) (*domain.GitBranch, error)
	ChangeBranchMeta(name string) error
}
