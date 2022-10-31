package config

import (
	"errors"
	"os"

	"github.com/eleven-sh/hetzner-cloud-provider/userconfig"
	toml "github.com/pelletier/go-toml"
)

type TOMLConfig struct {
	ActiveContext string              `toml:"active_context,omitempty"`
	Contexts      []TOMLConfigContext `toml:"contexts"`
}

type TOMLConfigContext struct {
	Name  string `toml:"name"`
	Token string `toml:"token"`
}

type ContextLoader struct{}

func NewContextLoader() ContextLoader {
	return ContextLoader{}
}

func (ContextLoader) Load(
	context string,
	configFilePath string,
) (*userconfig.Config, error) {

	configContent, err := os.ReadFile(configFilePath)

	if err != nil && errors.Is(err, os.ErrNotExist) {
		if len(context) > 0 {
			return nil, ErrContextNotFound{
				Context:        context,
				ConfigFilePath: configFilePath,
			}
		}

		return userconfig.NewConfig("", ""), nil
	}

	if err != nil {
		return nil, err
	}

	return unmarshalConfig(
		configFilePath,
		configContent,
		context,
	)
}

func unmarshalConfig(
	configFilePath string,
	configContent []byte,
	context string,
) (*userconfig.Config, error) {

	var tomlConfig TOMLConfig
	if err := toml.Unmarshal(configContent, &tomlConfig); err != nil {
		return nil, err
	}

	userConfig := &userconfig.Config{}

	contextToSearchFor := context
	if len(contextToSearchFor) == 0 {
		contextToSearchFor = tomlConfig.ActiveContext
	}

	for _, tomlConfigContext := range tomlConfig.Contexts {
		if tomlConfigContext.Name == contextToSearchFor {
			return userconfig.NewConfig(
				tomlConfigContext.Token,
				"",
			), nil
		}
	}

	if len(context) > 0 {
		return nil, ErrContextNotFound{
			Context:        context,
			ConfigFilePath: configFilePath,
		}
	}

	return userConfig, nil
}
