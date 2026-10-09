package git

import "github.com/algrvvv/git-go/domain"

type Git interface {
	ListBranches() ([]*domain.GitBranch, error)
	Checkout(branch *domain.GitBranch) error
	CreateBranch(name string, description string) (*domain.GitBranch, error)
	ChangeBranchMeta(name string) error
}
