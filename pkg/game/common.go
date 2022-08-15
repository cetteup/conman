package game

import (
	"github.com/cetteup/conman/pkg/config"
	"github.com/cetteup/conman/pkg/handler"
)

type Handler interface {
	ReadGlobalConfig(game handler.Game) (*config.Config, error)
	ReadProfileConfig(game handler.Game, profile string) (*config.Config, error)
}
