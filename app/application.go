package app

import (
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/algrvvv/git-go/compare"
	"github.com/algrvvv/git-go/domain"
	"github.com/algrvvv/git-go/git"
)

type App struct {
	cfg Config

	git git.Git
	log *slog.Logger
}

func NewApp(g git.Git, logger *slog.Logger, cfg Config) *App {
	return &App{
		cfg: cfg,
		git: g,
		log: logger,
	}
}

func (a *App) Run() error {
	switch true {
	case a.cfg.Meta:
		return a.setMetaForBranch()
	case a.cfg.List:
		return a.showBranchList()
	case a.cfg.Checkout:
		return a.createNewBranchAndCheckout()
	case len(a.cfg.Args) != 0:
		return a.checkoutByNameOrMeta(a.cfg.Args)
	}

	// интерактивный режим: нет флагов — показываем выбор ветки
	return a.runInteractive()
}

func (a *App) setMetaForBranch() error {
	var (
		branch string
		err    error
	)

	if len(a.cfg.Args) == 0 {
		branch, err = a.git.GetCurrentBranchName()
		if err != nil {
			return err
		}
	} else {
		branch = a.cfg.Args[0]
	}

	if err = a.git.ChangeBranchMeta(branch); err != nil {
		return err
	}

	return nil
}

func (a *App) showBranchList() error {
	branches, err := a.git.ListBranches()
	if err != nil {
		return err
	}

	for _, branch := range branches {
		fmt.Printf("%s - %s\n", branch.Name, branch.Meta)
	}

	return nil
}

func (a *App) createNewBranchAndCheckout() error {
	if len(a.cfg.Args) != 2 {
		return errors.New("invalid arguments: expected 'branch-name' 'branch-meta'")
	}

	name := a.cfg.Args[0]
	meta := a.cfg.Args[1]

	branch, err := a.git.CreateBranch(name, meta)
	if err != nil {
		return err
	}

	ok(fmt.Sprintf("✔ branch %s (%s) successfully created.\n", branch.Name, branch.Meta))

	return nil
}

func (a *App) handleInteractive(b *domain.GitBranch) error {
	if err := a.git.Checkout(b); err != nil {
		a.log.Error("failed to checkout branch in interactive mode", "branch", b.Name, "err", err)
		return err
	}

	return nil
}

func (a *App) checkoutByNameOrMeta(args []string) error {
	query := strings.Join(args, " ")

	branches, err := a.git.ListBranches()
	if err != nil {
		return err
	}

	a.log.Debug("checking out branches", "branches", branches)

	searcher := compare.NewSearcher(a.log, query, branches)

	var branch *domain.GitBranch
	branch, err = searcher.SearchByQuery()
	if err != nil {
		return err
	}

	a.log.Debug("found branch", "branch", branch)

	fmt.Printf(dim("selected branch: %s - %s\n"), branch.Name, branch.Meta)

	if err = a.git.Checkout(branch); err != nil {
		return err
	}

	return nil
}
