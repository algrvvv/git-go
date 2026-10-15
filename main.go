package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/algrvvv/git-go/app"
	"github.com/algrvvv/git-go/git"
	"github.com/algrvvv/git-go/logger"
)

func main() {
	cfg, err := app.ParseConfig(os.Args[1:])
	if err != nil {
		if errors.Is(err, flag.ErrHelp) {
			os.Exit(0)
		}
		fmt.Println(err)
		os.Exit(2)
	}

	log := logger.NewLogger(cfg.Debug)

	cli := git.NewGitCli()
	application := app.NewApp(cli, log, cfg)

	if err = application.Run(); err != nil {
		fmt.Println(err)
	}
}
