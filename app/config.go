package app

import "flag"

// Config — распарсенные флаги и позиционные аргументы CLI.
type Config struct {
	Debug    bool
	List     bool
	Checkout bool
	Meta     bool
	Args     []string
}

func ParseConfig(argv []string) (Config, error) {
	var cfg Config

	fs := flag.NewFlagSet("git-go", flag.ContinueOnError)
	fs.BoolVar(&cfg.Debug, "debug", false, "enable debug logging")
	fs.BoolVar(&cfg.List, "list", false, "list all branches")
	fs.BoolVar(&cfg.Checkout, "b", false, "checkout and create branch")
	fs.BoolVar(&cfg.Meta, "meta", false, "set meta for current branch")

	if err := fs.Parse(argv); err != nil {
		return Config{}, err
	}
	cfg.Args = fs.Args()

	return cfg, nil
}
