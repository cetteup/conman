package game

import (
	"github.com/cetteup/conman/pkg/config"
	"github.com/cetteup/conman/pkg/handler"
)

type ProfileType int

const (
	ProfileTypeMultiplayer ProfileType = iota
	ProfileTypeSingleplayer
)

type Profile struct {
	Key  string
	Name string
	Type ProfileType
}

type Handler interface {
	ReadConfigFile(path string) (*config.Config, error)
	ReadGlobalConfig(game handler.Game) (*config.Config, error)
	GetProfileKeys(game handler.Game) ([]string, error)
	ReadProfileConfig(game handler.Game, profileKey string) (*config.Config, error)
	PurgeShaderCache(game handler.Game) error
	PurgeLogoCache(game handler.Game) error
	BuildProfilesFolderPath(game handler.Game) (string, error)
}
