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
	modsDirName        = "mods"
	cacheDirName       = "cache"
	logoCacheDirName   = "LogoCache"
	profilesDirName    = "Profiles"
	globalConFileName  = "Global.con"
	profileConFileName = "Profile.con"
)

type FileRepository interface {
	FileExists(path string) (bool, error)
	DirExists(path string) (bool, error)
	WriteFile(path string, data []byte, perm os.FileMode) error
	ReadFile(path string) ([]byte, error)
	ReadDir(path string) ([]os.DirEntry, error)
	Glob(pattern string) ([]string, error)
	RemoveAll(path string) error
}

type ErrGameNotSupported struct {
	game string
}

func (e *ErrGameNotSupported) Error() string {
	return fmt.Sprintf("game not supported: %s", e.game)
}

type ErrActionNotSupportedForGame struct {
	game   string
	action string
}

func (e *ErrActionNotSupportedForGame) Error() string {
	return fmt.Sprintf("action not supported for game: %s, %s", e.action, e.game)
}

type Handler struct {
	repository FileRepository
}

func New(repository FileRepository) *Handler {
	return &Handler{
		repository: repository,
	}
}

// Read central profile configuration file (primarily contains reference to current default profile)
func (h *Handler) ReadGlobalConfig(game Game) (*config.Config, error) {
	path, err := h.buildGlobalConfigPath(game)
	if err != nil {
		return nil, err
	}
	return h.ReadConfigFile(path)
}

// Retrieve a list of profile keys (valid profile directories in the game's profile folder)
func (h *Handler) GetProfileKeys(game Game) ([]string, error) {
	path, err := h.BuildProfilesFolderPath(game)
	if err != nil {
		return nil, err
	}

	entries, err := h.repository.ReadDir(path)
	if err != nil {
		return nil, err
	}

	var profileKeys []string
	for _, entry := range entries {
		if entry.IsDir() {
			valid, err := h.isValidProfileDir(game, path, entry.Name())
			if err != nil {
				return nil, err
			}
			if valid {
				profileKeys = append(profileKeys, entry.Name())
			}
		}
	}

	return profileKeys, nil
}

func (h *Handler) isValidProfileDir(game Game, basePath string, profileKey string) (bool, error) {
	var conFileName string
	switch game {
	case GameBf2:
		conFileName = profileConFileName
	default:
		return false, &ErrGameNotSupported{game: string(game)}
	}

	conFilePath := filepath.Join(basePath, profileKey, conFileName)

	return h.repository.FileExists(conFilePath)
}

// Checks whether a given profile key is valid (a profile with the given key exists)
func (h *Handler) IsValidProfileKey(game Game, profileKey string) (bool, error) {
	profileKeys, err := h.GetProfileKeys(game)
	if err != nil {
		return false, err
	}

	for _, pk := range profileKeys {
		if profileKey == pk {
			return true, nil
		}
	}

	return false, nil
}

// Read the profile configuration file for the specified profile
func (h *Handler) ReadProfileConfig(game Game, profileKey string) (*config.Config, error) {
	path, err := h.buildProfileConfigPath(game, profileKey)
	if err != nil {
		return nil, err
	}
	return h.ReadConfigFile(path)
}

// Read the config file at given path
func (h *Handler) ReadConfigFile(path string) (*config.Config, error) {
	data, err := h.repository.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return config.FromBytes(path, data), nil
}

// Write the given config file to disk
func (h *Handler) WriteConfigFile(c *config.Config) error {
	return h.repository.WriteFile(c.Path, c.ToBytes(), 0666)
}

// Delete all shader cache (.cfx) files [Refractor v2 games only]
func (h *Handler) PurgeShaderCache(game Game) error {
	if !isSupportedGame(game) {
		return &ErrGameNotSupported{game: string(game)}
	}
	if !isRefractorV2Game(game) {
		return &ErrActionNotSupportedForGame{
			action: "PurgeShaderCache",
			game:   string(game),
		}
	}

	/*
		Looking from the base path, the shader cache files are stored in:
		mods/
		├──[mod]/
		   ├──cache/
			  ├──[cache dir with uuid-looking name]/
				 ├──[cache file].cfx
		Best practise is to delete all folder inside each mod's /cache directory
	*/
	basePath, err := h.BuildBasePath(game)
	if err != nil {
		return err
	}
	return h.globRemoveAll(filepath.Join(basePath, modsDirName, "*", cacheDirName, "*"))
}

func (h *Handler) globRemoveAll(pattern string) error {
	matches, err := h.repository.Glob(pattern)
	if err != nil {
		return err
	}

	for _, match := range matches {
		err = h.repository.RemoveAll(match)
		if err != nil {
			return err
		}
	}

	return nil
}

// Delete all cached server banner images [Refractor v2 games only]
func (h *Handler) PurgeLogoCache(game Game) error {
	if !isSupportedGame(game) {
		return &ErrGameNotSupported{game: string(game)}
	}
	if !isRefractorV2Game(game) {
		return &ErrActionNotSupportedForGame{
			action: "PurgeShaderCache",
			game:   string(game),
		}
	}

	/*
		Looking from the base path, the logo cache files are stored in:
		LogoCache/
		├──[server hosting banner image]/
		   ├──[...path to file on server]/
		Simply delete all folders in LogoCache/
	*/
	basePath, err := h.BuildBasePath(game)
	if err != nil {
		return err
	}
	return h.globRemoveAll(filepath.Join(basePath, logoCacheDirName, "*"))
}

// Build path to the root folder for given game's configuration
func (h *Handler) BuildBasePath(game Game) (string, error) {
	switch game {
	case GameBf2:
		return buildV2BasePath(bf2GameDirName)
	default:
		return "", &ErrGameNotSupported{game: string(game)}
	}
}

// Build path to the folder containing given game's profile configuration
func (h *Handler) BuildProfilesFolderPath(game Game) (string, error) {
	basePath, err := h.BuildBasePath(game)
	if err != nil {
		return "", err
	}
	return filepath.Join(basePath, profilesDirName), nil
}

func (h *Handler) buildGlobalConfigPath(game Game) (string, error) {
	switch game {
	case GameBf2:
		return h.buildV2GlobalConfigPath(game)
	default:
		return "", &ErrGameNotSupported{game: string(game)}
	}
}

func (h *Handler) buildV2GlobalConfigPath(game Game) (string, error) {
	profilesPath, err := h.BuildProfilesFolderPath(game)
	if err != nil {
		return "", err
	}
	return filepath.Join(profilesPath, globalConFileName), nil
}

func (h *Handler) buildProfileConfigPath(game Game, profileKey string) (string, error) {
	switch game {
	case GameBf2:
		return h.buildV2ProfileConfigPath(game, profileKey)
	default:
		return "", &ErrGameNotSupported{game: string(game)}
	}
}

func (h *Handler) buildV2ProfileConfigPath(game Game, profileKey string) (string, error) {
	basePath, err := h.BuildProfilesFolderPath(game)
	if err != nil {
		return "", err
	}
	return filepath.Join(basePath, profileKey, profileConFileName), nil
}

func buildV2BasePath(gameDirName string) (string, error) {
	documentsDirPath, err := windows.KnownFolderPath(windows.FOLDERID_Documents, windows.KF_FLAG_DEFAULT)
	if err != nil {
		return "", err
	}

	return filepath.Join(documentsDirPath, gameDirName), nil
}

func isSupportedGame(game Game) bool {
	switch game {
	case GameBf2:
		return true
	default:
		return false
	}
}

func isRefractorV2Game(game Game) bool {
	switch game {
	case GameBf2:
		return true
	default:
		return false
	}
}
