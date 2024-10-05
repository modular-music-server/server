package config

import (
	"log"
	"os"

	"github.com/aarzilli/golua/lua"
	"github.com/go-git/go-git/v5"
)

func luaConfig(L *lua.State) int {
    argumentCount := L.GetTop()
    if argumentCount == 0 {
        L.RaiseError("Config requires arguments!");
        return 0
    }
    if !L.IsTable(1) {
        L.RaiseError("Config expects a table as its first argument!")
        return 0
    }
    if argumentCount != 1 {
        println("Warning: Config does not expect multiple arguments")
    }
    L.SetTop(1)
    L.GetField(-1, "modules")
    if !L.IsTable(-1) {
        L.RaiseError("Config expects field \"modules\" to be a table")
        return 0
    }
    L.GetField(-1, "providers")
    if !L.IsTable(-1) {
        L.RaiseError("Config expects field \"modules.providers\" to be a table")
        return 0
    }
    L.PushNil()
    for L.Next(-2) != 0 {
        if !L.IsTable(-1) {
            L.RaiseError("Config expects values in \"modules.providers\" to be tables")
            return 0
        }
        L.GetField(-1, "type")
        if !L.IsString(-1) {
            L.RaiseError("Config expects \"provider.type\" to be a string")
            return 0
        }
        providerType := L.ToString(-1)
        L.Pop(1)
        L.GetField(-1, "location")
        if !L.IsString(-1) {
            L.RaiseError("Config expects \"provider.location\" to be a string")
            return 0
        }
        providerLocation := L.ToString(-1)
        L.Pop(1)
        L.GetField(-1, "name")
        if !L.IsString(-1) {
            L.RaiseError("Config expects \"provider.name\" to be a string")
            return 0
        }
        providerName := L.ToString(-1)
        L.Pop(1)
        switch providerType {
        case "git":
            _, err := git.PlainClone(os.Getenv("HOME") + "/.local/share/modular-music-server/" + providerName, false, &git.CloneOptions{
                URL: providerLocation,
                Progress: os.Stdout,
            })
            if err != nil {
                println("There was an error!!!!!!!!!!!!! ahhhhhhh")
                log.Print(err)
            }
        default:
            L.RaiseError("Unexpected provider type \"" + providerType + "\"")
            return 0
        }
        L.Pop(1)
    }
    println("Lookin good!")
    return 0
}

func initModule(L *lua.State) int {
    // println("Loading test module...")
    L.PushGoFunction(luaConfig)
    return 1
}

func Config() {
    L := lua.NewState()
    L.OpenLibs()
    L.GetGlobal("package")
    L.GetField(-1, "preload")
    L.PushGoClosure(initModule)
    L.SetField(-2, "modular-music-server")
    L.SetTop(0)
    err := L.DoFile(os.Getenv("HOME") + "/.config/modular-music-server/config.lua")
    if err != nil {
        print("there was an error!: ")
        println(err.Error())
    }
    defer L.Close()
}
