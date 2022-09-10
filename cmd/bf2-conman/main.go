//go:generate goversioninfo

package main

import (
	"flag"
	"os"

	"github.com/cetteup/conman/cmd/bf2-conman/internal/actions"
	"github.com/cetteup/conman/cmd/bf2-conman/internal/gui"
	"github.com/cetteup/conman/pkg/game/bf2"
	"github.com/cetteup/conman/pkg/handler"
	filerepo "github.com/cetteup/filerepo/pkg"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	logKeyProfile string = "profile"
)

func init() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
}

func main() {
	var noGUI bool
	var doPurgeServerHistory bool
	var doMarkAllVoiceOverHelpAsPlayed bool
	var setDefaultProfileKey string
	flag.BoolVar(&noGUI, "no-gui", false, "do not open/use the graphical user interface")
	flag.BoolVar(&doPurgeServerHistory, "purge-server-history", false, "purge all server history entries from the current default profile's General.con")
	flag.BoolVar(&doMarkAllVoiceOverHelpAsPlayed, "disable-help-voice-overs", false, "mark all help voice over lines as played for the current default profile")
	flag.StringVar(&setDefaultProfileKey, "default-profile", "", "set the given profile as the current default profile")
	flag.Parse()

	fileRepository := filerepo.New()
	h := handler.New(fileRepository)

	profiles, err := bf2.GetProfiles(h)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to get list of available profiles")
		os.Exit(1)
	}

	defaultProfileKey, err := bf2.GetDefaultProfileKey(h)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to determine current default user profile key")
		os.Exit(1)
	}

	if !noGUI {
		mw, err := gui.CreateMainWindow(h, profiles, defaultProfileKey)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to create main window")
			os.Exit(1)
		}

		mw.Run()
	} else {
		if setDefaultProfileKey != "" {
			err = actions.SetDefaultProfile(h, setDefaultProfileKey)
			if err != nil {
				log.Error().Err(err).Str(logKeyProfile, setDefaultProfileKey).Msg("Failed to update current default profile")
			} else {
				log.Info().Str(logKeyProfile, setDefaultProfileKey).Msg("Updated current default profile")
				// Use new default profile from here on
				defaultProfileKey = setDefaultProfileKey
			}
		}

		if doPurgeServerHistory {
			err = actions.PurgeServerHistory(h, defaultProfileKey)
			if err != nil {
				log.Error().Err(err).Str(logKeyProfile, defaultProfileKey).Msg("Failed to purge server history from current default profile")
			} else {
				log.Info().Str(logKeyProfile, defaultProfileKey).Msg("Purged server history from current default profile")
			}
		}

		if doMarkAllVoiceOverHelpAsPlayed {
			err = actions.MarkAllVoiceOverHelpAsPlayed(h, defaultProfileKey)
			if err != nil {
				log.Error().Err(err).Str(logKeyProfile, defaultProfileKey).Msg("Failed to mark all voice over help lines as played for current default profile")
			} else {
				log.Info().Str(logKeyProfile, defaultProfileKey).Msg("Marked all voice over help lines as played for current default profile")
			}
		}
	}
}
