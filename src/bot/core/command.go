package core

import (
	"fmt"
	"strings"

	"github.com/vcokltfre/volcan/src/utils"
)

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
		commandNames = append(commandNames, command.Aliases...)
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

	module   *Module
	parent   *Command
	required int
}

func (c *Command) Build(module *Module, parent *Command) error {
	c.module = module
	c.parent = parent

	for _, command := range c.Commands {
		if err := command.Build(module, c); err != nil {
			return err
		}
	}

	for _, arg := range c.Args {
		if arg.Required {
			c.required++
		}
	}

	return nil
}

func (c *Command) Find(parts []string, index int) (*Command, int) {
	if len(parts) == 0 {
		return c, index
	}

	for _, command := range c.Commands {
		if command.Name == parts[0] || utils.Contains(command.Aliases, parts[0]) {
			return command.Find(parts[1:], index+1)
		}
	}

	return c, index
}

func (c *Command) Run(ctx *Context, args []string) error {
	if c.module.Check != nil {
		_, err := c.module.Check(ctx)
		if err != nil {
			return err
		}
	}

	if c.Check != nil {
		_, err := c.Check(ctx)
		if err != nil {
			return err
		}
	}

	parent := c.parent
	for parent != nil {
		if parent.Check != nil {
			_, err := parent.Check(ctx)
			if err != nil {
				return err
			}
		}

		parent = parent.parent
	}

	if err := c.match(ctx, args); err != nil {
		return err
	}

	return c.Handler(ctx)
}

func (c *Command) match(ctx *Context, args []string) error {
	cleanArgs := []string{}
	index := 0

	flagsDone := []string{}

	for index < len(args) {
		arg := args[index]

		if strings.HasPrefix(arg, "--") {
			flagName := strings.TrimPrefix(arg, "--")
			flag := c.findFlag(flagName)

			if flag == nil {
				return fmt.Errorf("unknown flag %s", flagName)
			}

			if flag.Boolean {
				ctx.bools[flag.Name] = true
				flagsDone = append(flagsDone, flag.Name)
				index++
				continue
			}

			if index+1 >= len(args) {
				return fmt.Errorf("expected value for flag %s", flagName)
			}

			if err := flag.Validate(ctx, args[index+1]); err != nil {
				return err
			}

			ctx.flags[flag.Name] = args[index+1]
			flagsDone = append(flagsDone, flag.Name)
			index += 2
		}

		if strings.HasPrefix(arg, "-") {
			flagNames := strings.TrimPrefix(arg, "-")
			for _, flagName := range flagNames {
				flag := c.findFlag(string(flagName))

				if flag == nil {
					return fmt.Errorf("unknown flag %s", string(flagName))
				}

				if !flag.Boolean && len(flagNames) > 1 {
					return fmt.Errorf("flag %s is not a boolean flag but is being used in a multi-flag", string(flagName))
				}

				if flag.Boolean {
					ctx.bools[string(flagName)] = true
					flagsDone = append(flagsDone, string(flagName))
					continue
				}

				if index+1 >= len(args) {
					return fmt.Errorf("expected value for flag %s", string(flagName))
				}

				if err := flag.Validate(ctx, args[index+1]); err != nil {
					return err
				}

				ctx.flags[string(flagName)] = args[index+1]
				flagsDone = append(flagsDone, string(flagName))
				index += 1
			}

			index++
			continue
		}

		cleanArgs = append(cleanArgs, arg)
		index++
	}

	for _, flag := range c.Flags {
		if !utils.Contains(flagsDone, flag.Name) {
			if flag.Boolean {
				ctx.bools[flag.Name] = false
			} else {
				ctx.flags[flag.Name] = flag.Default
			}
		}
	}

	if c.VarArg == nil && len(cleanArgs) > len(c.Args) {
		return fmt.Errorf("too many arguments, max %d given %d", len(c.Args), len(cleanArgs))
	}

	if len(cleanArgs) < c.required {
		return fmt.Errorf("not enough arguments, required %d given %d", c.required, len(cleanArgs))
	}

	usedArgs := 0

	for index, arg := range c.Args {
		if index >= len(cleanArgs) {
			break
		}

		if err := arg.Validate(ctx, cleanArgs[index]); err != nil {
			return err
		}

		ctx.args[arg.Name] = cleanArgs[index]
		usedArgs++
	}

	if c.VarArg != nil {
		err := c.VarArg.Validate(ctx, cleanArgs[usedArgs:]...)
		if err != nil {
			return err
		}
		ctx.varArgs = cleanArgs[usedArgs:]
	}

	return nil
}

func (c *Command) findFlag(name string) *Flag {
	for _, flag := range c.Flags {
		if flag.Name == name || utils.Contains(flag.Aliases, name) {
			return flag
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
	Boolean     bool
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
