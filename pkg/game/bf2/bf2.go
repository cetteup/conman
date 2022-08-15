// Methods for working specifically with Battlefield 2 configuration files (.con)
package bf2

import (
	"fmt"
	"strconv"

	"github.com/cetteup/conman/pkg/config"
	"github.com/cetteup/conman/pkg/game"
	"github.com/cetteup/conman/pkg/handler"
)

const (
	globalConKeyDefaultUserRef = "GlobalSettings.setDefaultUser"
	profileConKeyGamespyNick   = "LocalProfile.setGamespyNick"
	profileConKeyPassword      = "LocalProfile.setPassword"
	// profileNumberMaxLength BF2 only uses 4 digit profile numbers
	profileNumberMaxLength = 4
)

// Read and parse the Battlefield 2 Profile.con file for the current default profile/user
func GetDefaultUserProfileCon(h game.Handler) (*config.Config, error) {
	profileNumber, err := GetDefaultUserProfileNumber(h)
	if err != nil {
		return nil, err
	}

	profileCon, err := h.ReadProfileConfig(handler.GameBf2, profileNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to read Profile.con for current default profile (%s): %s", profileNumber, err)
	}

	return profileCon, nil
}

// Get the default user profile number by reading and parsing the Battlefield 2 Global.con file
func GetDefaultUserProfileNumber(h game.Handler) (string, error) {
	globalCon, err := h.ReadGlobalConfig(handler.GameBf2)
	if err != nil {
		return "", fmt.Errorf("failed to read Global.con: %s", err)
	}

	defaultUserRef, err := globalCon.GetValue(globalConKeyDefaultUserRef)
	if err != nil {
		return "", fmt.Errorf("reference to default profile is missing from Global.con")
	}
	// Since BF2 only uses 4 digits for the profile number, 16 bits is plenty to store it
	if _, err := strconv.ParseInt(defaultUserRef.String(), 10, 16); err != nil || len(defaultUserRef.String()) > profileNumberMaxLength {
		return "", fmt.Errorf("reference to default profile in Global.con is not a valid profile number: %s", defaultUserRef.String())
	}

	return defaultUserRef.String(), nil
}

// Extract profile name and encrypted password from a parsed Battlefield 2 Profile.con file
func GetEncryptedProfileConLogin(profileCon *config.Config) (string, string, error) {
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
