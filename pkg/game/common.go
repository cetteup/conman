package game

import (
	"github.com/cetteup/conman/pkg/config"
	"github.com/cetteup/conman/pkg/handler"
)

type Profile struct {
	Key  string
	Name string
}

type Handler interface {
	ReadConfigFile(path string) (*config.Config, error)
	ReadGlobalConfig(game handler.Game) (*config.Config, error)
	GetProfileKeys(game handler.Game) ([]string, error)
	ReadProfileConfig(game handler.Game, profileKey string) (*config.Config, error)
	PurgeShaderCache(game handler.Game) error
	PurgeLogoCache(game handler.Game) error
	BuildBasePath(game handler.Game) (string, error)
}
