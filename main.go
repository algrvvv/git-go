package main

import (
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/algrvvv/git-go/app"
	"github.com/algrvvv/git-go/git"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	cfg, err := app.ParseConfig(os.Args[1:])
	if err != nil {
		if errors.Is(err, flag.ErrHelp) {
			os.Exit(0)
		}
		fmt.Println(err)
		os.Exit(2)
	}

	cli := git.NewGitCli()
	application := app.NewApp(cli, logger, cfg)

	if err = application.Run(); err != nil {
		fmt.Println(err)
	}
}
