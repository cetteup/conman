// Read and parse Refractor engine configuration files (.con)
package handler

import (
	"fmt"
	"os"
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
	FileExists(path string) (bool, error)
	WriteFile(path string, data []byte, perm os.FileMode) error
	ReadFile(path string) ([]byte, error)
	ReadDir(path string) ([]os.DirEntry, error)
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
	return h.ReadConfigFile(path)
}

// Retrieve a list of profile keys (valid profile directories in the game's profile folder)
func (h *Handler) GetProfileKeys(game Game) ([]string, error) {
	path, err := h.BuildBasePath(game)
	if err != nil {
		return nil, err
	}

	entries, err := h.repository.ReadDir(path)
	if err != nil {
		return nil, err
	}

	var profiles []string
	for _, entry := range entries {
		if entry.IsDir() {
			valid, err := h.isValidProfileDir(game, path, entry.Name())
			if err != nil {
				return nil, err
			}
			if valid {
				profiles = append(profiles, entry.Name())
			}
		}
	}

	return profiles, nil
}

func (h *Handler) isValidProfileDir(game Game, basePath string, profile string) (bool, error) {
	var conFileName string
	switch game {
	case GameBf2:
		conFileName = profileConFileName
	default:
		return false, &ErrGameNotSupported{game: string(game)}
	}

	conFilePath := filepath.Join(basePath, profile, conFileName)

	return h.repository.FileExists(conFilePath)
}

func (h *Handler) ReadProfileConfig(game Game, profile string) (*config.Config, error) {
	path, err := buildProfileConfigPath(game, profile)
	if err != nil {
		return nil, err
	}
	return h.ReadConfigFile(path)
}

func (h *Handler) ReadConfigFile(path string) (*config.Config, error) {
	data, err := h.repository.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return config.FromBytes(data), nil
}

func (h *Handler) WriteConfigFile(path string, c *config.Config) error {
	return h.repository.WriteFile(path, c.ToBytes(), 0666)
}

func (h *Handler) BuildBasePath(game Game) (string, error) {
	switch game {
	case GameBf2:
		return buildV2BasePath(bf2GameDirName)
	default:
		return "", &ErrGameNotSupported{game: string(game)}
	}
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
