package app

import (
	"flag"
	"log/slog"

	"github.com/algrvvv/git-go/git"
)

type App struct {
	debug    bool
	list     bool
	checkout bool
	meta     bool

	git git.Git
	log *slog.Logger
}

func NewApp(git git.Git, logger *slog.Logger) *App {
	return &App{
		log: logger,
		git: git,
	}
}

func (a *App) Run() error {
	a.parse()

	switch true {
	case a.meta:
		return a.setMetaForBranch()
	case a.list:
		return a.showBranchList()
	}

	// todo: interactive mode
	return nil
}

func (a *App) parse() {
	flag.BoolVar(&a.debug, "debug", false, "enable debug logging")
	flag.BoolVar(&a.list, "list", false, "list all branches")
	flag.BoolVar(&a.checkout, "b", false, "checkout and create branch")
	flag.BoolVar(&a.meta, "meta", false, "set meta for current branch")
	flag.Parse()
}

func (a *App) setMetaForBranch() error {
	return nil
}

func (a *App) showBranchList() error {
	return nil
}
