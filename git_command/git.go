package git_command

import (
	"fmt"
	"io"
	"os"
	"os/exec"
)

type GitCommandBuilder struct {
	cmd  *exec.Cmd
	args []string
}

func NewGitCommandBuilder(args ...string) *GitCommandBuilder {
	return &GitCommandBuilder{
		args: args,
	}
}

func (c *GitCommandBuilder) Arg(args ...string) *GitCommandBuilder {
	c.args = append(c.args, args...)
	return c
}

func (c *GitCommandBuilder) Argf(format, arg string) *GitCommandBuilder {
	arg := fmt.Sprintf(format, arg)
	return c.Arg(arg)
}

func (c *GitCommandBuilder) ArgIf(condition bool, arg string) *GitCommandBuilder {
	if condition {
		return c.Arg(arg)
	}

	return c
}

func (c *GitCommandBuilder) ArgElse(condition bool, ifArg, elseArg string) *GitCommandBuilder {
	if condition {
		return c.Arg(ifArg)
	}

	return c.Arg(elseArg)
}

func (c *GitCommandBuilder) Build() *GitCommandBuilder {
	c.cmd = exec.Command("git", c.args...)
	return c
}

func (c *GitCommandBuilder) WithOutput(w ...io.Writer) *GitCommandBuilder {
	c.cmd.Stdout = os.Stdout
	c.cmd.Stderr = os.Stderr
	c.cmd.Stdin = os.Stdin

	if len(w) == 1 {
		c.cmd.Stdout = w[0]
	}

	return c
}

func (c *GitCommandBuilder) Run() error {
	return c.cmd.Run()
}
