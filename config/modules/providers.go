package modules

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/aarzilli/golua/lua"
	"github.com/go-git/go-git/v5"
)

type ProviderEntry struct {
	Type     string
	Location string
	Id       string
}

type Provider struct {
	Location    string
	Name        string
	Author      string
	Description string
}

func ParseProviderEntries(L *lua.State) ([]*ProviderEntry, error) {
	var providerEntries []*ProviderEntry
	L.PushNil()
	for L.Next(-2) != 0 {
		providerEntry, err := parseProviderEntry(L)
		if err != nil {
			return nil, err
		}
		err = ensureProviderEntryClean(providerEntry, providerEntries)
		if err != nil {
			return nil, err
		}
		providerEntries = append(providerEntries, providerEntry)
	}
	return providerEntries, nil
}

func parseProviderEntry(L *lua.State) (*ProviderEntry, error) {
	if !L.IsTable(-1) {
		return nil, errors.New("Config expects values in \"modules.providers\" to be tables")
	}
	L.GetField(-1, "type")
	if !L.IsString(-1) {
		return nil, errors.New("Config expects \"provider.type\" to be a string")
	}
	providerType := L.ToString(-1)
	L.Pop(1)
	L.GetField(-1, "location")
	if !L.IsString(-1) {
		return nil, errors.New("Config expects \"provider.location\" to be a string")
	}
	providerLocation := L.ToString(-1)
	L.Pop(1)
	L.GetField(-1, "id")
	if !L.IsString(-1) {
		return nil, errors.New("Config expects \"provider.id\" to be a string")
	}
	providerId := L.ToString(-1)
	L.Pop(1)

	L.Pop(1)

	return &ProviderEntry{
		Type:     providerType,
		Location: providerLocation,
		Id:       providerId,
	}, nil
}

func ensureProviderEntryClean(providerEntry *ProviderEntry, entries []*ProviderEntry) error {
	if providerEntry.Type != "git" {
		return errors.New("Provider type must be one of \"git\"")
	}
	if providerEntry.Id == "." || providerEntry.Id == ".." {
		return errors.New("Provider id cannot be \".\" or \"..\"")
	}
	if strings.Contains(providerEntry.Id, "/") {
		return errors.New("Provider id cannot contain a forward slash (/)")
	}
	for _, entry := range entries {
		if providerEntry.Id == entry.Id {
			return errors.New("Cannot have duplicate provider id (\"" + providerEntry.Id + "\")")
		}
	}
	return nil
}

func HandleProviderEntries(providerEntries []*ProviderEntry) (map[string]*Provider, error) {
	providers := make(map[string]*Provider)
	for _, providerEntry := range providerEntries {
		provider, id, err := handleProviderEntry(providerEntry)
		if err != nil {
			return nil, err
		}
		providers[id] = provider
	}
	return providers, nil
}

func handleProviderEntry(providerEntry *ProviderEntry) (*Provider, string, error) {
	switch providerEntry.Type {
	case "git":
		_, err := git.PlainOpen(os.Getenv("HOME") + "/.local/share/modular-music-server/" + providerEntry.Id)
		if err != nil {
			if err == git.ErrRepositoryNotExists {
				_, err = git.PlainClone(os.Getenv("HOME")+"/.local/share/modular-music-server/"+providerEntry.Id, false, &git.CloneOptions{
					URL:      providerEntry.Location,
					Progress: os.Stdout,
				})
				if err != nil {
					log.Fatalf("Error cloning repository: %v\n", err)
				}
			} else {
				log.Fatalf("Error opening repository: %v\n", err)
			}
		}
	default:
		return nil, "", errors.New(fmt.Sprintf("Unexpected provider entry type %q\n", providerEntry.Type))
	}

	L := lua.NewState()
	L.OpenLibs()
	err := L.DoFile(os.Getenv("HOME") + "/.local/share/modular-music-server/" + providerEntry.Id + "/init.lua")
	if err != nil {
		return nil, "", err
	}
	returnCount := L.GetTop()
	if returnCount == 0 {
		return nil, "", errors.New("init.lua should return a config!")
	}
	if !L.IsTable(1) {
		return nil, "", errors.New("init.lua should return a table!")
	}
	if returnCount != 1 {
		println("Warning: init.lua does not expect multiple return values")
	}
	L.SetTop(1)

	L.GetField(-1, "name")
	if !L.IsString(-1) {
		return nil, "", errors.New("init.lua expects field \"name\" to be a string")
	}
	providerName := L.ToString(-1)
	L.Pop(1)

	L.GetField(-1, "author")
	if !L.IsString(-1) {
		return nil, "", errors.New("init.lua expects field \"author\" to be a string")
	}
	providerAuthor := L.ToString(-1)
	L.Pop(1)

	L.GetField(-1, "description")
	if !L.IsString(-1) {
		return nil, "", errors.New("init.lua expects field \"description\" to be a string")
	}
	providerDescription := L.ToString(-1)
	L.Pop(1)

	L.Close()

	provider := &Provider{
		Location:    providerEntry.Location,
		Name:        providerName,
		Author:      providerAuthor,
		Description: providerDescription,
	}
	return provider, providerEntry.Id, nil
}
