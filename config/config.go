package config

import (
	"fmt"
	"os"

	"github.com/festeh/doubt/engine"
	"github.com/yosuke-furukawa/json5/encoding/json5"
)

type CommandConfig struct {
	Type string `json:"type"`
	Sleep *struct {
		Duration int `json:"duration"`
	} `json:"sleep,omitempty"`
	Keypress *struct {
		Key string `json:"key"`
	} `json:"keypress,omitempty"`
}

type Config struct {
	ExecutablePath string          `json:"executable_path"`
	Commands       []CommandConfig `json:"commands"`
}

func LoadConfig(configPath string) (*Config, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := json5.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return &config, nil
}

func (c *Config) ToEngineCommands() ([]engine.Command, error) {
	var commands []engine.Command

	for i, cmdConfig := range c.Commands {
		switch cmdConfig.Type {
		case "sleep":
			if cmdConfig.Sleep == nil {
				return nil, fmt.Errorf("command %d: sleep command missing duration", i)
			}
			commands = append(commands, engine.NewSleepCommand(cmdConfig.Sleep.Duration))
		case "keypress":
			if cmdConfig.Keypress == nil {
				return nil, fmt.Errorf("command %d: keypress command missing key", i)
			}
			commands = append(commands, engine.NewKeypressCommand(cmdConfig.Keypress.Key))
		default:
			return nil, fmt.Errorf("command %d: unknown command type: %s", i, cmdConfig.Type)
		}
	}

	return commands, nil
}
