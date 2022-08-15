// Read and parse Refractor engine configuration files (.con)
package handler

import (
	"fmt"
	"path/filepath"

	"github.com/cetteup/conman/pkg/config"
	"golang.org/x/sys/windows"
)

type Game string

const (
	GameBf2 Game = "bf2"

	bf2GameDirName     = "Battlefield 2"
	profilesDirName    = "Profiles"
	globalConFileName  = "Global.con"
	profileConFileName = "Profile.con"
)

type FileRepository interface {
	ReadFile(path string) ([]byte, error)
}

type ErrGameNotSupported struct {
	game string
}

func (e *ErrGameNotSupported) Error() string {
	return fmt.Sprintf("game not supported: %s", e.game)
}

type Handler struct {
	repository FileRepository
}

func New(repository FileRepository) *Handler {
	return &Handler{
		repository: repository,
	}
}

func (h *Handler) ReadGlobalConfig(game Game) (*config.Config, error) {
	path, err := buildGlobalConfigPath(game)
	if err != nil {
		return nil, err
	}
	return h.readConfigFile(path)
}

func (h *Handler) ReadProfileConfig(game Game, profile string) (*config.Config, error) {
	path, err := buildProfileConfigPath(game, profile)
	if err != nil {
		return nil, err
	}
	return h.readConfigFile(path)
}

func (h *Handler) readConfigFile(path string) (*config.Config, error) {
	data, err := h.repository.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return config.FromBytes(data), nil
}

func buildGlobalConfigPath(game Game) (string, error) {
	switch game {
	case GameBf2:
		return buildV2GlobalConfigPath(bf2GameDirName)
	default:
		return "", &ErrGameNotSupported{game: string(game)}
	}
}

func buildV2GlobalConfigPath(gameDirName string) (string, error) {
	basePath, err := buildV2BasePath(gameDirName)
	if err != nil {
		return "", err
	}
	return filepath.Join(basePath, globalConFileName), nil
}

func buildProfileConfigPath(game Game, profile string) (string, error) {
	switch game {
	case GameBf2:
		return buildV2ProfileConfigPath(bf2GameDirName, profile)
	default:
		return "", &ErrGameNotSupported{game: string(game)}
	}
}

func buildV2ProfileConfigPath(gameDirName string, profile string) (string, error) {
	basePath, err := buildV2BasePath(gameDirName)
	if err != nil {
		return "", err
	}
	return filepath.Join(basePath, profile, profileConFileName), nil
}

func buildV2BasePath(gameDirName string) (string, error) {
	documentsDirPath, err := windows.KnownFolderPath(windows.FOLDERID_Documents, windows.KF_FLAG_DEFAULT)
	if err != nil {
		return "", err
	}

	return filepath.Join(documentsDirPath, gameDirName, profilesDirName), nil
}
