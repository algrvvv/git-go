package gitexec

import (
	"fmt"
	"io"
	"os"
	"os/exec"
)

type Builder struct {
	args []string
}

func NewBuilder(args ...string) *Builder {
	return &Builder{args: args}
}

func (b *Builder) Arg(args ...string) *Builder {
	b.args = append(b.args, args...)
	return b
}

func (b *Builder) Argf(format, arg string) *Builder {
	return b.Arg(fmt.Sprintf(format, arg))
}

func (b *Builder) ArgIf(condition bool, arg string) *Builder {
	if condition {
		return b.Arg(arg)
	}

	return b
}

func (b *Builder) ArgElse(condition bool, ifArg, elseArg string) *Builder {
	if condition {
		return b.Arg(ifArg)
	}

	return b.Arg(elseArg)
}

func (b *Builder) Build() *Command {
	return NewCommand("git", b.args...)
}

// Command — собранная git-команда, готовая к запуску.
type Command struct {
	cmd *exec.Cmd
}

func NewCommand(command string, args ...string) *Command {
	return &Command{cmd: exec.Command(command, args...)}
}

func (c *Command) WithOutput(w ...io.Writer) *Command {
	c.cmd.Stdout = os.Stdout
	c.cmd.Stderr = os.Stderr
	c.cmd.Stdin = os.Stdin

	if len(w) == 1 {
		c.cmd.Stdout = w[0]
	}

	return c
}

func (c *Command) Run() error {
	return c.cmd.Run()
}
