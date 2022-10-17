package core

import (
	"fmt"

	"github.com/vcokltfre/volcan/src/utils"
)

type Bot struct {
	Modules []*Module

	modules map[string]*Module
}

func (b *Bot) Build() error {
	commandNames := []string{}

	for _, module := range b.modules {
		if _, ok := b.modules[module.Name]; ok {
			return fmt.Errorf("module %s already exists", module.Name)
		}

		cmds, err := module.Build()
		if err != nil {
			return err
		}

		for _, cmd := range cmds {
			if utils.Contains(commandNames, cmd) {
				return fmt.Errorf("command %s already exists", cmd)
			}

			commandNames = append(commandNames, cmd)
		}

		b.modules[module.Name] = module
	}

	return nil
}
