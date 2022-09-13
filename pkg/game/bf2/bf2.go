// Methods for working specifically with Battlefield 2 configuration files (.con)
package bf2

import (
	"fmt"
	"path/filepath"
	"strconv"

	"github.com/cetteup/conman/pkg/config"
	"github.com/cetteup/conman/pkg/game"
	"github.com/cetteup/conman/pkg/handler"
)

type ProfileConfigFile string

const (
	ProfileConfigFileAudioCon          ProfileConfigFile = "Audio.con"
	ProfileConfigFileControlsCon       ProfileConfigFile = "Controls.con"
	ProfileConfigFileDemoBookmarksCon  ProfileConfigFile = "DemoBookmarks.con"
	ProfileConfigFileGeneralCon        ProfileConfigFile = "General.con"
	ProfileConfigFileHapticCon         ProfileConfigFile = "Haptic.con"
	ProfileConfigFileMapListCon        ProfileConfigFile = "mapList.con"
	ProfileConfigFileProfileCon        ProfileConfigFile = "Profile.con"
	ProfileConfigFileServerSettingsCon ProfileConfigFile = "ServerSettings.con"
	ProfileConfigFileVideoCon          ProfileConfigFile = "Video.con"

	defaultProfileKey = "Default"
	// profileKeyMaxLength BF2 only uses 4 digit profile keys
	profileKeyMaxLength = 4

	globalConKeyDefaultProfileRef = "GlobalSettings.setDefaultUser"

	profileConKeyName        = "LocalProfile.setName"
	profileConKeyGamespyNick = "LocalProfile.setGamespyNick"
	profileConKeyPassword    = "LocalProfile.setPassword"

	generalConKeyServerHistory       = "GeneralSettings.addServerHistory"
	generalConKeyVoiceOverHelpPlayed = "GeneralSettings.setPlayedVOHelp"
)

// Read a config file from the given Battlefield 2 profile
func ReadProfileConfigFile(h game.Handler, profileKey string, configFile ProfileConfigFile) (*config.Config, error) {
	basePath, err := h.BuildProfilesFolderPath(handler.GameBf2)
	if err != nil {
		return nil, err
	}

	filePath := filepath.Join(basePath, profileKey, string(configFile))
	conFile, err := h.ReadConfigFile(filePath)
	if err != nil {
		return nil, err
	}

	return conFile, nil
}

func GetProfiles(h game.Handler) ([]game.Profile, error) {
	profileKeys, err := h.GetProfileKeys(handler.GameBf2)
	if err != nil {
		return nil, err
	}

	var profiles []game.Profile
	for _, profileKey := range profileKeys {
		// Ignore the default profile
		if profileKey == defaultProfileKey {
			continue
		}

		profileCon, err := h.ReadProfileConfig(handler.GameBf2, profileKey)
		if err != nil {
			return nil, err
		}

		profileName, err := profileCon.GetValue(profileConKeyName)
		if err != nil {
			return nil, err
		}

		profiles = append(profiles, game.Profile{
			Key:  profileKey,
			Name: profileName.String(),
		})
	}

	return profiles, nil
}

// Read and parse the Battlefield 2 Profile.con file for the current default profile
func GetDefaultProfileProfileCon(h game.Handler) (*config.Config, error) {
	profileKey, err := GetDefaultProfileKey(h)
	if err != nil {
		return nil, err
	}

	profileCon, err := h.ReadProfileConfig(handler.GameBf2, profileKey)
	if err != nil {
		return nil, fmt.Errorf("failed to read Profile.con for current default profile (%s): %s", profileKey, err)
	}

	return profileCon, nil
}

// Get the default profile's key by reading and parsing the Battlefield 2 Global.con file
func GetDefaultProfileKey(h game.Handler) (string, error) {
	globalCon, err := h.ReadGlobalConfig(handler.GameBf2)
	if err != nil {
		return "", fmt.Errorf("failed to read Global.con: %s", err)
	}

	defaultUserRef, err := globalCon.GetValue(globalConKeyDefaultProfileRef)
	if err != nil {
		return "", fmt.Errorf("reference to default profile is missing from Global.con")
	}
	// Since BF2 only uses 4 digits for the profile key, 16 bits is plenty to store it
	if _, err := strconv.ParseInt(defaultUserRef.String(), 10, 16); err != nil || len(defaultUserRef.String()) > profileKeyMaxLength {
		return "", fmt.Errorf("reference to default profile in Global.con is not a valid profile key: %s", defaultUserRef.String())
	}

	return defaultUserRef.String(), nil
}

func PurgeShaderCache(h game.Handler) error {
	return h.PurgeShaderCache(handler.GameBf2)
}

func PurgeLogoCache(h game.Handler) error {
	return h.PurgeLogoCache(handler.GameBf2)
}

func SetDefaultProfile(globalCon *config.Config, profileKey string) {
	globalCon.SetValue(globalConKeyDefaultProfileRef, *config.NewQuotedValue(profileKey))
}

// Extract profile name and encrypted password from a parsed Battlefield 2 Profile.con file
func GetEncryptedLogin(profileCon *config.Config) (string, string, error) {
	nickname, err := profileCon.GetValue(profileConKeyGamespyNick)
	// GameSpy nick property is present but empty for local/singleplayer profiles
	if err != nil || nickname.String() == "" {
		return "", "", fmt.Errorf("gamespy nickname is missing/empty")
	}
	encryptedPassword, err := profileCon.GetValue(profileConKeyPassword)
	if err != nil || encryptedPassword.String() == "" {
		return "", "", fmt.Errorf("encrypted password is missing/empty")
	}

	return nickname.String(), encryptedPassword.String(), nil
}

// Remove all server history entries (GeneralSettings.addServerHistory) from given General.con config
func PurgeServerHistory(generalCon *config.Config) {
	generalCon.Delete(generalConKeyServerHistory)
}

// Add all voice over help lines as played (GeneralSettings.setPlayedVOHelp) in given General.con config
func MarkAllVoiceOverHelpAsPlayed(generalCon *config.Config) {
	generalCon.SetValue(generalConKeyVoiceOverHelpPlayed, *config.NewQuotedValueFromSlice(voiceOverHelpLines))
}
