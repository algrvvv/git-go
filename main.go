package main

import (
	"log/slog"
	"os"

	"github.com/algrvvv/git-go/app"
)

var (
	debug    bool
	list     bool
	checkout bool
	meta     bool
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	application := app.NewApp(logger)

	if err := application.Run(); err != nil {
		logger.Error("failed to run application", "error", err)
	}
}
