package config

import (
	"errors"

	"example.com/modular-music-server/config/modules"
	"github.com/aarzilli/golua/lua"
)

type Modules struct {
    Providers []*modules.Provider
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

