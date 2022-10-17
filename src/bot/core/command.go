package core

import "fmt"

type CommandHandler func(*Context) error

type CommandCheck func(*Context) (bool, error)

type ArgValidator func(*Context, string) error

type Module struct {
	Name        string
	Description string
	Commands    []*Command
	Check       CommandCheck
}

func (m *Module) Build() ([]string, error) {
	commandNames := []string{}

	for _, command := range m.Commands {
		if err := command.Build(m, nil); err != nil {
			return nil, err
		}

		commandNames = append(commandNames, command.Name)
	}

	return commandNames, nil
}

type Command struct {
	Name        string
	Description string
	Aliases     []string
	Usage       string
	Handler     CommandHandler
	Commands    []*Command
	Check       CommandCheck
	Args        []*Arg
	Flags       []*Flag
	VarArg      *VarArg

	module *Module
	parent *Command
}

func (c *Command) Build(module *Module, parent *Command) error {
	c.module = module
	c.parent = parent

	for _, command := range c.Commands {
		if err := command.Build(module, c); err != nil {
			return err
		}
	}

	return nil
}

type Arg struct {
	Name        string
	Description string
	Required    bool
	Default     string
	Validator   ArgValidator
}

func (a *Arg) Validate(ctx *Context, arg string) error {
	if a.Validator == nil {
		return nil
	}

	return a.Validator(ctx, arg)
}

type Flag struct {
	Name        string
	Description string
	Aliases     []string
	Default     string
	Validator   ArgValidator
}

func (f *Flag) Validate(ctx *Context, arg string) error {
	if f.Validator == nil {
		return nil
	}

	return f.Validator(ctx, arg)
}

type VarArg struct {
	Name        string
	Description string
	Max         int
	Validator   ArgValidator
}

func (v *VarArg) Validate(ctx *Context, args ...string) error {
	if v.Validator == nil {
		return nil
	}

	if len(args) > v.Max {
		return fmt.Errorf("too many arguments, max %d given %d", v.Max, len(args))
	}

	for _, arg := range args {
		if err := v.Validator(ctx, arg); err != nil {
			return err
		}
	}

	return nil
}
