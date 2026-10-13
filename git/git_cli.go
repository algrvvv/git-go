package git

import (
	"errors"
	"strings"

	"github.com/algrvvv/git-go/domain"
	gitcmd "github.com/algrvvv/git-go/gitexec"
)

type GitCli struct{}

func NewGitCli() *GitCli {
	return &GitCli{}
}

func (g *GitCli) Checkout(branch *domain.GitBranch) error {
	if branch == nil {
		return errors.New("empty branch")
	}

	cmd := gitcmd.NewBuilder("checkout", branch.Name).Build()
	return cmd.Run()
}

// ChangeBranchMeta метод с интерактивным изменением меты ветки
func (g *GitCli) ChangeBranchMeta(name string) error {
	cmd := gitcmd.NewBuilder().
		Arg("branch").
		ArgIf(name != "", name).
		Arg("--edit-description").
		Build().
		WithOutput()

	return cmd.Run()
}

func (g *GitCli) CreateBranch(name string, meta string) (*domain.GitBranch, error) {
	branch := domain.NewGitBranch(name, meta)

	checkoutCmd := gitcmd.NewBuilder("checkout", "-b", branch.Name).Build()
	if err := checkoutCmd.Run(); err != nil {
		return nil, err
	}

	metaCmd := gitcmd.NewBuilder().
		Arg("config").
		Argf("branch.%s.description", branch.Name).
		Argf("%q", branch.Meta).
		Build()
	if err := metaCmd.Run(); err != nil {
		return nil, err
	}

	return branch, nil
}

func (g *GitCli) ListBranches() ([]*domain.GitBranch, error) {
	builder := &strings.Builder{}
	cmd := gitcmd.NewBuilder("branch").
		Build().
		WithOutput(builder)

	if err := cmd.Run(); err != nil {
		return nil, err
	}

	// удаляем звездочку у актуальной ветки
	result := strings.ReplaceAll(builder.String(), "*", "")

	var branches []*domain.GitBranch
	for brnch := range strings.FieldsSeq(result) {
		name := strings.TrimSpace(brnch)
		meta, _ := g.GetMetaForBranch(name)

		branches = append(branches, domain.NewGitBranch(name, meta))
	}

	return branches, nil
}

func (g *GitCli) GetCurrentBranchName() (string, error) {
	builder := &strings.Builder{}

	cmd := gitcmd.NewBuilder("branch", "--show-current").
		Build().
		WithOutput(builder)

	if err := cmd.Run(); err != nil {
		return "", err
	}

	return strings.TrimSpace(builder.String()), nil
}

func (g *GitCli) GetMetaForBranch(branch string) (string, error) {
	builder := &strings.Builder{}
	cmd := gitcmd.NewBuilder().
		Arg("config").
		Argf("branch.%s.description", branch).
		Build().
		WithOutput(builder)

	if err := cmd.Run(); err != nil {
		return "", err
	}

	return strings.TrimSpace(builder.String()), nil
}
