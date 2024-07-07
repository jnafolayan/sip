package cli

import (
	"flag"
	"fmt"
)

type Command struct {
	Name        string
	Desc        string
	Run         func(cmd *Command, args []string) error
	Init        func(cmd *Command)
	subcommands map[string]*Command
	FlagSet     *flag.FlagSet
}

type CommandArgs map[string]any
type CommandFlags map[string]any

// RegisterCmd registers a commands to be executed if there is a match.
func (c *Command) RegisterCmd(cmd *Command) {
	if c.subcommands == nil {
		c.subcommands = make(map[string]*Command)
	}
	c.subcommands[cmd.Name] = cmd
	if cmd.Init != nil {
		cmd.Init(cmd)
	}
}

// hasSubCommands return true if this command has any subcommands.
func (c *Command) hasSubCommands() bool {
	return c.subcommands != nil && len(c.subcommands) > 0
}

// Execute runs this command or one of its subcommands if required.
func (c *Command) Execute(args []string) error {
	if c.hasSubCommands() && len(args) > 0 {
		if cmd, ok := c.subcommands[args[0]]; ok {
			return cmd.Execute(args[1:])
		}
	}

	if c.Run == nil {
		return fmt.Errorf("no runner defined for %q", c.Name)
	}

	return c.Run(c, args)
}
