package config

import (
	"errors"

	"github.com/aarzilli/golua/lua"
	"github.com/modular-music-server/server/config/modules"
)

type Modules struct {
	Providers map[string]*modules.Provider
}

func ParseModules(L *lua.State) (*Modules, error) {
	L.GetField(-1, "providers")
	if !L.IsTable(-1) {
		return nil, errors.New("Config expects field \"modules.providers\" to be a table")
	}
	providerEntries, err := modules.ParseProviderEntries(L)
	if err != nil {
		return nil, err
	}
	providers, err := modules.HandleProviderEntries(providerEntries)
	if err != nil {
		return nil, err
	}

	modules := &Modules{
		Providers: providers,
	}

	return modules, nil
}
