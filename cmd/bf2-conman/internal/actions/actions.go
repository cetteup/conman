package actions

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/cetteup/conman/pkg/config"
	"github.com/cetteup/conman/pkg/game/bf2"
	"github.com/cetteup/conman/pkg/handler"
)

const (
	demoBookmarkMaxAge = time.Hour * 24 * 7
)

func SetDefaultProfile(h *handler.Handler, profileKey string) error {
	validKey, err := h.IsValidProfileKey(handler.GameBf2, profileKey)
	if err != nil {
		return err
	}
	if !validKey {
		return fmt.Errorf("given profile key is not valid")
	}

	globalCon, err := h.ReadGlobalConfig(handler.GameBf2)
	if err != nil {
		return err
	}

	bf2.SetDefaultProfile(globalCon, profileKey)

	return h.WriteConfigFile(globalCon)
}

func GetProfilePassword(h *handler.Handler, profileKey string) (string, error) {
	profileCon, err := bf2.ReadProfileConfigFile(h, profileKey, bf2.ProfileConfigFileProfileCon)
	if err != nil {
		return "", err
	}

	encryptedPassword, err := profileCon.GetValue(bf2.ProfileConKeyPassword)
	if err != nil {
		return "", err
	}

	return bf2.DecryptProfileConPassword(encryptedPassword.String())
}

func SetProfilePassword(h *handler.Handler, profileKey string, password string) error {
	profileCon, err := bf2.ReadProfileConfigFile(h, profileKey, bf2.ProfileConfigFileProfileCon)
	if err != nil {
		return err
	}

	encryptedPassword, err := bf2.EncryptProfileConPassword(password)
	if err != nil {
		return err
	}

	profileCon.SetValue(bf2.ProfileConKeyPassword, *config.NewValue(encryptedPassword))

	return h.WriteConfigFile(profileCon)
}

func PurgeServerHistory(h *handler.Handler, profileKey string) error {
	generalCon, err := bf2.ReadProfileConfigFile(h, profileKey, bf2.ProfileConfigFileGeneralCon)
	if err != nil {
		return err
	}

	bf2.PurgeServerHistory(generalCon)

	return h.WriteConfigFile(generalCon)
}

func PurgeServerFavorites(h *handler.Handler, profileKey string) error {
	generalCon, err := bf2.ReadProfileConfigFile(h, profileKey, bf2.ProfileConfigFileGeneralCon)
	if err != nil {
		return err
	}

	bf2.PurgeServerFavorites(generalCon)

	return h.WriteConfigFile(generalCon)
}

func PurgeOldDemoBookmarks(h *handler.Handler, profileKey string) error {
	demoBookmarksCon, err := bf2.ReadProfileConfigFile(h, profileKey, bf2.ProfileConfigFileDemoBookmarksCon)
	if err != nil {
		// We want to clean the demo bookmarks, so we don't consider it an error if the file is missing
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}

	bf2.PurgeOldDemoBookmarks(demoBookmarksCon, time.Now(), demoBookmarkMaxAge)

	return h.WriteConfigFile(demoBookmarksCon)
}

func MarkAllVoiceOverHelpAsPlayed(h *handler.Handler, profileKey string) error {
	generalCon, err := bf2.ReadProfileConfigFile(h, profileKey, bf2.ProfileConfigFileGeneralCon)
	if err != nil {
		return err
	}

	bf2.MarkAllVoiceOverHelpAsPlayed(generalCon)

	return h.WriteConfigFile(generalCon)
}

func PurgeShareCache(h *handler.Handler) error {
	return bf2.PurgeShaderCache(h)
}

func PurgeLogoCache(h *handler.Handler) error {
	return bf2.PurgeLogoCache(h)
}
