package config

import (
	"errors"
	"os"

	"github.com/aarzilli/golua/lua"
)

type Config struct {
    Modules *Modules
}

func parseConfig(L *lua.State) (*Config, error) {
    returnCount := L.GetTop()
    if returnCount == 0 {
        return nil, errors.New("config.lua should return a config!")
    }
    if !L.IsTable(1) {
        return nil, errors.New("config.lua should return a table!")
    }
    if returnCount != 1 {
        println("Warning: config.lua does not expect multiple return values")
    }
    L.SetTop(1)
    L.GetField(-1, "modules")
    if !L.IsTable(-1) {
        return nil, errors.New("Config expects field \"modules\" to be a table")
    }
    modules, err := ParseModules(L)
    if err != nil {
        return nil, err
    }

    config := &Config{
        Modules: modules,
    }

    return config, nil
}

func LoadConfig() (*Config, error) {
    L := lua.NewState()
    L.OpenLibs()
    err := L.DoFile(os.Getenv("HOME") + "/.config/modular-music-server/config.lua")
    if err != nil {
        print("there was an error!: ")
        println(err.Error())
    }
    config, err := parseConfig(L)
    if err != nil {
        return nil, err
    }

    L.Close()

    return config, nil
}
