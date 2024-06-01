//go:generate goversioninfo

package main

import (
	"flag"
	"os"

	filerepo "github.com/cetteup/filerepo/pkg"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/cetteup/conman/cmd/bf2-conman/internal/actions"
	"github.com/cetteup/conman/cmd/bf2-conman/internal/gui"
	"github.com/cetteup/conman/pkg/game/bf2"
	"github.com/cetteup/conman/pkg/handler"
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
	var doPurgeServerFavorites bool
	var doPurgeOldDemoBookmarks bool
	var doMarkAllVoiceOverHelpAsPlayed bool
	var doPurgeShaderCache bool
	var doPurgeLogoCache bool
	var setDefaultProfileKey string
	flag.BoolVar(&noGUI, "no-gui", false, "do not open/use the graphical user interface")
	flag.BoolVar(&doPurgeServerHistory, "purge-server-history", false, "purge all server history entries from the current default profile")
	flag.BoolVar(&doPurgeServerFavorites, "purge-server-favorites", false, "purge all server favorites from the current default profile")
	flag.BoolVar(&doPurgeOldDemoBookmarks, "purge-old-demo-bookmarks", false, "purge all old demo bookmarks (older than 1 week) from the current default profile")
	flag.BoolVar(&doMarkAllVoiceOverHelpAsPlayed, "disable-help-voice-overs", false, "mark all help voice over lines as played for the current default profile")
	flag.BoolVar(&doPurgeShaderCache, "purge-shader-cache", false, "purge all shader cache files and folders")
	flag.BoolVar(&doPurgeLogoCache, "purge-logo-cache", false, "purge cached server banner images")
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

		if doPurgeServerFavorites {
			err = actions.PurgeServerFavorites(h, defaultProfileKey)
			if err != nil {
				log.Error().Err(err).Str(logKeyProfile, defaultProfileKey).Msg("Failed to purge server favorites from current default profile")
			} else {
				log.Info().Str(logKeyProfile, defaultProfileKey).Msg("Purged server favorites from current default profile")
			}
		}

		if doPurgeOldDemoBookmarks {
			err = actions.PurgeOldDemoBookmarks(h, defaultProfileKey)
			if err != nil {
				log.Error().Err(err).Str(logKeyProfile, defaultProfileKey).Msg("Failed to purge old demo bookmarks from current default profile")
			} else {
				log.Info().Str(logKeyProfile, defaultProfileKey).Msg("Purged old demo bookmarks from current default profile")
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

		if doPurgeShaderCache {
			err = actions.PurgeShareCache(h)
			if err != nil {
				log.Error().Err(err).Msg("Failed to purge shader cache")
			} else {
				log.Info().Msg("Purged shader cache")
			}
		}

		if doPurgeLogoCache {
			err = actions.PurgeLogoCache(h)
			if err != nil {
				log.Error().Err(err).Msg("Failed to purge logo cache")
			} else {
				log.Info().Msg("Purged logo cache")
			}
		}
	}
}
