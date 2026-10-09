package git

import (
	"strings"

	"github.com/algrvvv/git-go/domain"
	gitcmd "github.com/algrvvv/git-go/git_command"
)

type GitCli struct{}

func NewGitCli() Git {
	return &GitCli{}
}

func (g *GitCli) Checkout(branch *domain.GitBranch) error {
	return nil
}

// ChangeBranchMeta метод с интерактивным изменением меты ветки
func (g *GitCli) ChangeBranchMeta(name string) error {
	cmd := gitcmd.NewGitCommandBuilder().
		Arg("branch").
		ArgIf(name != "", name).
		Arg("--edit-description").
		Build().
		WithOutput()

	return cmd.Run()
}

func (g *GitCli) CreateBranch(name string, description string) (*domain.GitBranch, error) {
	return nil, nil
}

func (g *GitCli) ListBranches() ([]*domain.GitBranch, error) {
	builder := &strings.Builder{}
	cmd := gitcmd.NewGitCommandBuilder("branch").
		Build().
		WithOutput(builder)

	if err := cmd.Run(); err != nil {
		return nil, err
	}

	// удаляем звездочку у актуальной ветки
	result := strings.ReplaceAll(builder.String(), "*", "")

	var branches []*domain.GitBranch
	for brnch := range strings.FieldsSeq(result) {
		var option domain.GitBranchOptions
		option.Name = strings.TrimSpace(brnch)
		option.Description, _ = g.getDescriptionForBranch(option.Name)

		branches = append(branches, domain.NewGitBranch(option))
	}

	return branches, nil
}

func (g *GitCli) getCurrentBranchName() (string, error) {
	builder := &strings.Builder{}

	cmd := gitcmd.NewGitCommandBuilder("branch", "--show-current").
		Build().
		WithOutput(builder)

	if err := cmd.Run(); err != nil {
		return "", err
	}

	return strings.TrimSpace(builder.String()), nil
}

func (g *GitCli) getDescriptionForCurrent() (string, error) {
	currentBranch, err := g.getCurrentBranchName()
	if err != nil {
		return "", err
	}

	return g.getDescriptionForBranch(currentBranch)
}

func (g *GitCli) getDescriptionForBranch(branch string) (string, error) {
	builder := &strings.Builder{}
	cmd := gitcmd.NewGitCommandBuilder().
		Arg("config").
		Argf("branch.%s.description", branch).
		Build().
		WithOutput(builder)

	if err := cmd.Run(); err != nil {
		return "", err
	}

	return strings.TrimSpace(builder.String()), nil
}
