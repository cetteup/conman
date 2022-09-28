//go:build unit

package bf2

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/cetteup/conman/pkg/config"
	"github.com/cetteup/conman/pkg/game"
	"github.com/cetteup/conman/pkg/handler"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReadProfileConfigFile(t *testing.T) {
	type test struct {
		name            string
		givenProfileKey string
		givenConfigFile ProfileConfigFile
		expect          func(h *MockHandler)
		wantConfig      *config.Config
		wantErrContains string
	}

	tests := []test{
		{
			name:            "successfully reads Profile.con",
			givenProfileKey: "0001",
			givenConfigFile: ProfileConfigFileProfileCon,
			expect: func(h *MockHandler) {
				basePath := "C:\\Users\\default\\Documents\\Battlefield 2\\Profiles"
				profileConPath := filepath.Join(basePath, "0001", "Profile.con")
				h.EXPECT().BuildProfilesFolderPath(handler.GameBf2).Return(basePath, nil)
				h.EXPECT().ReadConfigFile(profileConPath).Return(config.New(
					profileConPath,
					map[string]config.Value{
						"LocalProfile.setName": *config.NewValue("\"mister249\""),
					},
				), nil)
			},
			wantConfig: config.New(
				"C:\\Users\\default\\Documents\\Battlefield 2\\Profiles\\0001\\Profile.con",
				map[string]config.Value{
					"LocalProfile.setName": *config.NewValue("\"mister249\""),
				},
			),
		},
		{
			name:            "errors if base path cannot be determined",
			givenProfileKey: "0001",
			givenConfigFile: ProfileConfigFileProfileCon,
			expect: func(h *MockHandler) {
				h.EXPECT().BuildProfilesFolderPath(handler.GameBf2).Return("", fmt.Errorf("some-error"))
			},
			wantErrContains: "some-error",
		},
		{
			name:            "errors if config file read errors",
			givenProfileKey: "0001",
			givenConfigFile: ProfileConfigFileProfileCon,
			expect: func(h *MockHandler) {
				basePath := "C:\\Users\\default\\Documents\\Battlefield 2\\Profiles"
				h.EXPECT().BuildProfilesFolderPath(handler.GameBf2).Return(basePath, nil)
				h.EXPECT().ReadConfigFile(filepath.Join(basePath, "0001", "Profile.con")).Return(nil, fmt.Errorf("some-error"))
			},
			wantErrContains: "some-error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// GIVEN
			ctrl := gomock.NewController(t)
			h := NewMockHandler(ctrl)

			// EXPECT
			tt.expect(h)

			// WHEN
			readConfig, err := ReadProfileConfigFile(h, tt.givenProfileKey, tt.givenConfigFile)

			// THEN
			if tt.wantErrContains != "" {
				require.ErrorContains(t, err, tt.wantErrContains)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.wantConfig, readConfig)
			}
		})
	}
}

func TestGetProfiles(t *testing.T) {
	type test struct {
		name            string
		expect          func(h *MockHandler)
		wantProfiles    []game.Profile
		wantErrContains string
	}

	tests := []test{
		{
			name: "successfully gets profiles",
			expect: func(h *MockHandler) {
				h.EXPECT().GetProfileKeys(handler.GameBf2).Return([]string{"0001", "0002"}, nil)
				h.EXPECT().ReadProfileConfig(handler.GameBf2, "0001").Return(config.New(
					"C:\\Users\\default\\Documents\\Battlefield 2\\Profiles\\0001\\Profile.con",
					map[string]config.Value{
						ProfileConKeyName:     *config.NewValue("some-multiplayer-profile"),
						ProfileConKeyPassword: *config.NewValue("some-password"),
					},
				), nil)
				h.EXPECT().ReadProfileConfig(handler.GameBf2, "0002").Return(config.New(
					"C:\\Users\\default\\Documents\\Battlefield 2\\Profiles\\0002\\Profile.con",
					map[string]config.Value{
						ProfileConKeyName: *config.NewValue("some-singleplayer-profile"),
					},
				), nil)
			},
			wantProfiles: []game.Profile{
				{
					Key:  "0001",
					Name: "some-multiplayer-profile",
					Type: game.ProfileTypeMultiplayer,
				},
				{
					Key:  "0002",
					Name: "some-singleplayer-profile",
					Type: game.ProfileTypeSingleplayer,
				},
			},
		},
		{
			name: "ignores default profile",
			expect: func(h *MockHandler) {
				h.EXPECT().GetProfileKeys(handler.GameBf2).Return([]string{"0001", DefaultProfileKey}, nil)
				h.EXPECT().ReadProfileConfig(handler.GameBf2, "0001").Return(config.New(
					"C:\\Users\\default\\Documents\\Battlefield 2\\Profiles\\0001\\Profile.con",
					map[string]config.Value{
						ProfileConKeyName: *config.NewValue("some-profile"),
					},
				), nil)
			},
			wantProfiles: []game.Profile{
				{
					Key:  "0001",
					Name: "some-profile",
					Type: game.ProfileTypeSingleplayer,
				},
			},
		},
		{
			name: "error getting profile keys",
			expect: func(h *MockHandler) {
				h.EXPECT().GetProfileKeys(handler.GameBf2).Return([]string{}, fmt.Errorf("some-error"))
			},
			wantErrContains: "some-error",
		},
		{
			name: "error reading profile's Profile.con",
			expect: func(h *MockHandler) {
				h.EXPECT().GetProfileKeys(handler.GameBf2).Return([]string{"0001"}, nil)
				h.EXPECT().ReadProfileConfig(handler.GameBf2, "0001").Return(nil, fmt.Errorf("some-error"))
			},
			wantErrContains: "some-error",
		},
		{
			name: "error for Profile.con not containing profile name",
			expect: func(h *MockHandler) {
				h.EXPECT().GetProfileKeys(handler.GameBf2).Return([]string{"0001"}, nil)
				h.EXPECT().ReadProfileConfig(handler.GameBf2, "0001").Return(config.New(
					"C:\\Users\\default\\Documents\\Battlefield 2\\Profiles\\0001\\Profile.con",
					map[string]config.Value{
						"some-other-key": *config.NewValue("some-other-value"),
					},
				), nil)
			},
			wantErrContains: "no such key in config",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// GIVEN
			ctrl := gomock.NewController(t)
			h := NewMockHandler(ctrl)

			// EXPECT
			tt.expect(h)

			// WHEN
			profiles, err := GetProfiles(h)

			// THEN
			if tt.wantErrContains != "" {
				require.ErrorContains(t, err, tt.wantErrContains)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.wantProfiles, profiles)
			}
		})
	}
}

func TestGetDefaultProfileProfileCon(t *testing.T) {
	type test struct {
		name               string
		expect             func(h *MockHandler)
		expectedProfileCon *config.Config
		wantErrContains    string
	}

	tests := []test{
		{
			name: "successfully retrieves default user's Profile.con",
			expect: func(h *MockHandler) {
				profileKey := "0001"
				h.EXPECT().ReadGlobalConfig(handler.GameBf2).Return(config.New(
					"C:\\Users\\default\\Documents\\Battlefield 2\\Profiles\\Global.con",
					map[string]config.Value{
						GlobalConKeyDefaultProfileRef: *config.NewValue(profileKey),
					},
				), nil)
				h.EXPECT().ReadProfileConfig(handler.GameBf2, profileKey).Return(config.New(
					"C:\\Users\\default\\Documents\\Battlefield 2\\Profiles\\0001\\Profile.con",
					map[string]config.Value{
						ProfileConKeyGamespyNick: *config.NewValue("some-nick"),
						ProfileConKeyPassword:    *config.NewValue("some-encrypted-password"),
					},
				), nil)
			},
			expectedProfileCon: config.New(
				"C:\\Users\\default\\Documents\\Battlefield 2\\Profiles\\0001\\Profile.con",
				map[string]config.Value{
					ProfileConKeyGamespyNick: *config.NewValue("some-nick"),
					ProfileConKeyPassword:    *config.NewValue("some-encrypted-password"),
				},
			),
		},
		{
			name: "error if default profile detection errors",
			expect: func(h *MockHandler) {
				h.EXPECT().ReadGlobalConfig(handler.GameBf2).Return(nil, fmt.Errorf("some-default-profile-detection-error"))
			},
			wantErrContains: "some-default-profile-detection-error",
		},
		{
			name: "error if Profile.con read errors",
			expect: func(h *MockHandler) {
				profileKey := "0001"
				h.EXPECT().ReadGlobalConfig(handler.GameBf2).Return(config.New(
					"C:\\Users\\default\\Documents\\Battlefield 2\\Profiles\\Global.con",
					map[string]config.Value{
						GlobalConKeyDefaultProfileRef: *config.NewValue(profileKey),
					},
				), nil)
				h.EXPECT().ReadProfileConfig(handler.GameBf2, profileKey).Return(nil, fmt.Errorf("some-profile-con-read-error"))
			},
			wantErrContains: "some-profile-con-read-error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// GIVEN
			ctrl := gomock.NewController(t)
			h := NewMockHandler(ctrl)

			// EXPECT
			tt.expect(h)

			// WHEN
			profileCon, err := GetDefaultProfileProfileCon(h)

			// THEN
			if tt.wantErrContains != "" {
				require.ErrorContains(t, err, tt.wantErrContains)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expectedProfileCon, profileCon)
			}
		})
	}
}

func TestGetDefaultProfileKey(t *testing.T) {
	type test struct {
		name               string
		expect             func(h *MockHandler)
		expectedProfileKey string
		wantErrContains    string
	}

	tests := []test{
		{
			name: "successfully retrieves default user profile key",
			expect: func(h *MockHandler) {
				h.EXPECT().ReadGlobalConfig(handler.GameBf2).Return(config.New(
					"C:\\Users\\default\\Documents\\Battlefield 2\\Profiles\\Global.con",
					map[string]config.Value{
						GlobalConKeyDefaultProfileRef: *config.NewValue("0001"),
					},
				), nil)
			},
			expectedProfileKey: "0001",
		},
		{
			name: "error if Global.con read errors",
			expect: func(h *MockHandler) {
				h.EXPECT().ReadGlobalConfig(handler.GameBf2).Return(nil, fmt.Errorf("some-global-con-read-error"))
			},
			wantErrContains: "some-global-con-read-error",
		},
		{
			name: "error if default user reference is missing from Global.con",
			expect: func(h *MockHandler) {
				h.EXPECT().ReadGlobalConfig(handler.GameBf2).Return(config.New(
					"C:\\Users\\default\\Documents\\Battlefield 2\\Profiles\\Global.con",
					map[string]config.Value{},
				), nil)
			},
			wantErrContains: "reference to default profile is missing from Global.con",
		},
		{
			name: "error if default user reference is non-numeric",
			expect: func(h *MockHandler) {
				h.EXPECT().ReadGlobalConfig(handler.GameBf2).Return(config.New(
					"C:\\Users\\default\\Documents\\Battlefield 2\\Profiles\\Global.con",
					map[string]config.Value{
						GlobalConKeyDefaultProfileRef: *config.NewValue("abcd"),
					},
				), nil)
			},
			wantErrContains: "reference to default profile in Global.con is not a valid profile key",
		},
		{
			name: "error if default user reference exceeds max length",
			expect: func(h *MockHandler) {
				h.EXPECT().ReadGlobalConfig(handler.GameBf2).Return(config.New(
					"C:\\Users\\default\\Documents\\Battlefield 2\\Profiles\\Global.con",
					map[string]config.Value{
						GlobalConKeyDefaultProfileRef: *config.NewValue("00001"),
					},
				), nil)
			},
			wantErrContains: "reference to default profile in Global.con is not a valid profile key",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// GIVEN
			ctrl := gomock.NewController(t)
			h := NewMockHandler(ctrl)

			// EXPECT
			tt.expect(h)

			// WHEN
			profileKey, err := GetDefaultProfileKey(h)

			// THEN
			if tt.wantErrContains != "" {
				require.ErrorContains(t, err, tt.wantErrContains)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expectedProfileKey, profileKey)
			}
		})
	}
}

func TestPurgeShaderCache(t *testing.T) {
	type test struct {
		name            string
		expect          func(h *MockHandler)
		wantErrContains string
	}

	tests := []test{
		{
			name: "successfully purges shader cache",
			expect: func(h *MockHandler) {
				h.EXPECT().PurgeShaderCache(handler.GameBf2)
			},
		},
		{
			name: "error purging shader cache",
			expect: func(h *MockHandler) {
				h.EXPECT().PurgeShaderCache(handler.GameBf2).Return(fmt.Errorf("some-error"))
			},
			wantErrContains: "some-error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// GIVEN
			ctrl := gomock.NewController(t)
			h := NewMockHandler(ctrl)

			// EXPECT
			tt.expect(h)

			// WHEN
			err := PurgeShaderCache(h)

			// THEN
			if tt.wantErrContains != "" {
				require.ErrorContains(t, err, tt.wantErrContains)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestPurgeLogoCache(t *testing.T) {
	type test struct {
		name            string
		expect          func(h *MockHandler)
		wantErrContains string
	}

	tests := []test{
		{
			name: "successfully purges logo cache",
			expect: func(h *MockHandler) {
				h.EXPECT().PurgeLogoCache(handler.GameBf2)
			},
		},
		{
			name: "error purging logo cache",
			expect: func(h *MockHandler) {
				h.EXPECT().PurgeLogoCache(handler.GameBf2).Return(fmt.Errorf("some-error"))
			},
			wantErrContains: "some-error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// GIVEN
			ctrl := gomock.NewController(t)
			h := NewMockHandler(ctrl)

			// EXPECT
			tt.expect(h)

			// WHEN
			err := PurgeLogoCache(h)

			// THEN
			if tt.wantErrContains != "" {
				require.ErrorContains(t, err, tt.wantErrContains)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestSetDefaultProfile(t *testing.T) {
	type test struct {
		name            string
		givenGlobalCon  *config.Config
		givenProfileKey string
		wantGlobalCon   *config.Config
	}

	tests := []test{
		{
			name: "sets default profile in Global.con",
			givenGlobalCon: config.New(
				"C:\\Users\\default\\Documents\\Battlefield 2\\Profiles\\0001\\General.con",
				map[string]config.Value{
					"GlobalSettings.setNamePrefix": *config.NewQuotedValue("=DOG="),
				},
			),
			givenProfileKey: "0001",
			wantGlobalCon: config.New(
				"C:\\Users\\default\\Documents\\Battlefield 2\\Profiles\\0001\\General.con",
				map[string]config.Value{
					GlobalConKeyDefaultProfileRef:  *config.NewQuotedValue("0001"),
					"GlobalSettings.setNamePrefix": *config.NewQuotedValue("=DOG="),
				},
			),
		},
		{
			name: "overwrites existing default profile in Global.con",
			givenGlobalCon: config.New(
				"C:\\Users\\default\\Documents\\Battlefield 2\\Profiles\\0001\\General.con",
				map[string]config.Value{
					GlobalConKeyDefaultProfileRef:  *config.NewQuotedValue("0001"),
					"GlobalSettings.setNamePrefix": *config.NewQuotedValue("=DOG="),
				},
			),
			givenProfileKey: "0002",
			wantGlobalCon: config.New(
				"C:\\Users\\default\\Documents\\Battlefield 2\\Profiles\\0001\\General.con",
				map[string]config.Value{
					GlobalConKeyDefaultProfileRef:  *config.NewQuotedValue("0002"),
					"GlobalSettings.setNamePrefix": *config.NewQuotedValue("=DOG="),
				},
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// GIVEN
			globalCon := tt.givenGlobalCon

			// WHEN
			SetDefaultProfile(globalCon, tt.givenProfileKey)

			// THEN
			assert.Equal(t, tt.wantGlobalCon, globalCon)
		})
	}
}

func TestGetEncryptedLogin(t *testing.T) {
	type test struct {
		name                 string
		prepareProfileConMap func(profileCon *config.Config)
		wantErrContains      string
	}

	tests := []test{
		{
			name:                 "successfully extracts encrypted login details",
			prepareProfileConMap: func(profileCon *config.Config) {},
		},
		{
			name: "fails if nickname is missing",
			prepareProfileConMap: func(profileCon *config.Config) {
				profileCon.Delete(ProfileConKeyGamespyNick)
			},
			wantErrContains: "gamespy nickname is missing/empty",
		},
		{
			name: "fails if nickname is empty",
			prepareProfileConMap: func(profileCon *config.Config) {
				profileCon.SetValue(ProfileConKeyGamespyNick, *config.NewValue(""))
			},
			wantErrContains: "gamespy nickname is missing/empty",
		},
		{
			name: "fails if password is missing",
			prepareProfileConMap: func(profileCon *config.Config) {
				profileCon.Delete(ProfileConKeyPassword)
			},
			wantErrContains: "encrypted password is missing/empty",
		},
		{
			name: "fails if password is empty",
			prepareProfileConMap: func(profileCon *config.Config) {
				profileCon.SetValue(ProfileConKeyPassword, *config.NewValue(""))
			},
			wantErrContains: "encrypted password is missing/empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// GIVEN
			bytes := []byte(fmt.Sprintf("%s \"mister249\"\r\n%s \"some-encrypted-password\"\r\n", ProfileConKeyGamespyNick, ProfileConKeyPassword))
			profileCon := config.FromBytes("C:\\Users\\default\\Documents\\Battlefield 2\\Profiles\\0001\\Profile.con", bytes)
			tt.prepareProfileConMap(profileCon)

			// WHEN
			nickname, encryptedPassword, err := GetEncryptedLogin(profileCon)

			// THEN
			if tt.wantErrContains != "" {
				require.ErrorContains(t, err, tt.wantErrContains)
			} else {
				require.NoError(t, err)
				expectedNickname, err := profileCon.GetValue(ProfileConKeyGamespyNick)
				require.NoError(t, err)
				assert.Equal(t, expectedNickname.String(), nickname)
				expectedPassword, err := profileCon.GetValue(ProfileConKeyPassword)
				require.NoError(t, err)
				assert.Equal(t, expectedPassword.String(), encryptedPassword)
			}
		})
	}
}

func TestPurgeServerHistory(t *testing.T) {
	type test struct {
		name               string
		givenGeneralCon    *config.Config
		expectedGeneralCon *config.Config
	}

	tests := []test{
		{
			name: "removes single server history from General.con",
			givenGeneralCon: config.New(
				"C:\\Users\\default\\Documents\\Battlefield 2\\Profiles\\0001\\General.con",
				map[string]config.Value{
					"GeneralSettings.setHUDTransparency": *config.NewValue("67.7346"),
					"GeneralSettings.addServerHistory":   *config.NewValue("\"135.125.56.26\" 29940 \"=DOG= No Explosives (Infantry)\" 1025"),
				},
			),
			expectedGeneralCon: config.New(
				"C:\\Users\\default\\Documents\\Battlefield 2\\Profiles\\0001\\General.con",
				map[string]config.Value{
					"GeneralSettings.setHUDTransparency": *config.NewValue("67.7346"),
				},
			),
		},
		{
			name: "removes multiple server history items from General.con",
			givenGeneralCon: config.New(
				"C:\\Users\\default\\Documents\\Battlefield 2\\Profiles\\0001\\General.con",
				map[string]config.Value{
					"GeneralSettings.setHUDTransparency": *config.NewValue("67.7346"),
					"GeneralSettings.addServerHistory":   *config.NewValue("\"135.125.56.26\" 29940 \"=DOG= No Explosives (Infantry)\" 1025;\"37.230.210.130\" 29900 \"PlayBF2! T~GAMER #1 Allmaps\" 360"),
				},
			),
			expectedGeneralCon: config.New(
				"C:\\Users\\default\\Documents\\Battlefield 2\\Profiles\\0001\\General.con",
				map[string]config.Value{
					"GeneralSettings.setHUDTransparency": *config.NewValue("67.7346"),
				},
			),
		},
		{
			name: "does nothing if General.con does not contain any server history items",
			givenGeneralCon: config.New(
				"C:\\Users\\default\\Documents\\Battlefield 2\\Profiles\\0001\\General.con",
				map[string]config.Value{
					"GeneralSettings.setHUDTransparency": *config.NewValue("67.7346"),
				},
			),
			expectedGeneralCon: config.New(
				"C:\\Users\\default\\Documents\\Battlefield 2\\Profiles\\0001\\General.con",
				map[string]config.Value{
					"GeneralSettings.setHUDTransparency": *config.NewValue("67.7346"),
				},
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// GIVEN
			generalCon := tt.givenGeneralCon

			// WHEN
			PurgeServerHistory(generalCon)

			// THEN
			assert.Equal(t, tt.expectedGeneralCon, generalCon)
		})
	}
}

func TestPurgeOldDemoBookmarks(t *testing.T) {
	type test struct {
		name                     string
		givenDemoBookmarksCon    *config.Config
		givenReference           time.Time
		givenMaxAge              time.Duration
		expectedDemoBookmarksCon *config.Config
	}

	reference := time.Date(2022, 9, 27, 12, 36, 0, 0, time.UTC)
	tests := []test{
		{
			name: "only removes bookmarks older than max age",
			givenDemoBookmarksCon: config.New(
				"C:\\Users\\default\\Documents\\Battlefield 2\\Profiles\\0001\\General.con",
				map[string]config.Value{
					DemoBookmarksConKeyDemoBookmark: *config.NewValueFromSlice([]string{
						buildDemoBookmark("some-server", "strike_at_karkand", reference.Add(-time.Hour*24*7-1)),
						buildDemoBookmark("some-other-server", "dragon_valley", reference.Add(-time.Hour*24*7+1)),
						buildDemoBookmark("some-other-server", "road_to_jalalabad", reference.Add(-time.Hour*24*7)),
					}),
				},
			),
			givenReference: reference,
			givenMaxAge:    time.Hour * 24 * 7,
			expectedDemoBookmarksCon: config.New(
				"C:\\Users\\default\\Documents\\Battlefield 2\\Profiles\\0001\\General.con",
				map[string]config.Value{
					DemoBookmarksConKeyDemoBookmark: *config.NewValueFromSlice([]string{
						buildDemoBookmark("some-other-server", "dragon_valley", reference.Add(-time.Hour*24*7+1)),
						buildDemoBookmark("some-other-server", "road_to_jalalabad", reference.Add(-time.Hour*24*7)),
					}),
				},
			),
		},
		{
			name: "deletes bookmarks key if no entries are kept",
			givenDemoBookmarksCon: config.New(
				"C:\\Users\\default\\Documents\\Battlefield 2\\Profiles\\0001\\General.con",
				map[string]config.Value{
					DemoBookmarksConKeyDemoBookmark: *config.NewValueFromSlice([]string{
						buildDemoBookmark("some-server", "strike_at_karkand", reference.Add(-time.Hour*24*7-1)),
					}),
				},
			),
			givenReference: reference,
			givenMaxAge:    time.Hour * 24 * 7,
			expectedDemoBookmarksCon: config.New(
				"C:\\Users\\default\\Documents\\Battlefield 2\\Profiles\\0001\\General.con",
				map[string]config.Value{},
			),
		},
		{
			name: "handles unquoted single-word values",
			givenDemoBookmarksCon: config.New(
				"C:\\Users\\default\\Documents\\Battlefield 2\\Profiles\\0001\\General.con",
				map[string]config.Value{
					DemoBookmarksConKeyDemoBookmark: *config.NewValueFromSlice([]string{
						fmt.Sprintf("first-server \"strike_at_karkand\" \"http://unused/first-server/strike_at_karkand/some-demo.bf2demo\" \"%s\"", formatDemoBookmarkTimestamp(reference)),
						fmt.Sprintf("\"second-server\" dragon_valley \"http://unused/second-server/dragon_valley/some-demo.bf2demo\" \"%s\"", formatDemoBookmarkTimestamp(reference)),
						fmt.Sprintf("\"third-server\" \"dragon_valley\" http://unused/second-server/dragon_valley/some-demo.bf2demo \"%s\"", formatDemoBookmarkTimestamp(reference)),
					}),
				},
			),
			givenReference: reference,
			givenMaxAge:    time.Hour * 24 * 7,
			expectedDemoBookmarksCon: config.New(
				"C:\\Users\\default\\Documents\\Battlefield 2\\Profiles\\0001\\General.con",
				map[string]config.Value{
					DemoBookmarksConKeyDemoBookmark: *config.NewValueFromSlice([]string{
						fmt.Sprintf("first-server \"strike_at_karkand\" \"http://unused/first-server/strike_at_karkand/some-demo.bf2demo\" \"%s\"", formatDemoBookmarkTimestamp(reference)),
						fmt.Sprintf("\"second-server\" dragon_valley \"http://unused/second-server/dragon_valley/some-demo.bf2demo\" \"%s\"", formatDemoBookmarkTimestamp(reference)),
						fmt.Sprintf("\"third-server\" \"dragon_valley\" http://unused/second-server/dragon_valley/some-demo.bf2demo \"%s\"", formatDemoBookmarkTimestamp(reference)),
					}),
				},
			),
		},
		{
			name: "removes bookmarks with invalid number of elements",
			givenDemoBookmarksCon: config.New(
				"C:\\Users\\default\\Documents\\Battlefield 2\\Profiles\\0001\\General.con",
				map[string]config.Value{
					DemoBookmarksConKeyDemoBookmark: *config.NewValueFromSlice([]string{
						fmt.Sprintf("\"strike_at_karkand\" \"http://unused/some-server/strike_at_karkand/some-demo.bf2demo\" \"%s\"", formatDemoBookmarkTimestamp(reference)),
						fmt.Sprintf("\"some-other-server\" \"second-server\" dragon_valley \"http://unused/some-other-server/dragon_valley/some-demo.bf2demo\" \"%s\" \"extra-element\"", formatDemoBookmarkTimestamp(reference)),
					}),
				},
			),
			givenReference: reference,
			givenMaxAge:    time.Hour * 24 * 7,
			expectedDemoBookmarksCon: config.New(
				"C:\\Users\\default\\Documents\\Battlefield 2\\Profiles\\0001\\General.con",
				map[string]config.Value{},
			),
		},
		{
			name: "removes bookmarks with invalid timestamps",
			givenDemoBookmarksCon: config.New(
				"C:\\Users\\default\\Documents\\Battlefield 2\\Profiles\\0001\\General.con",
				map[string]config.Value{
					DemoBookmarksConKeyDemoBookmark: *config.NewValueFromSlice([]string{
						"\"some-server\" \"strike_at_karkand\" \"http://unused/some-server/strike_at_karkand/some-demo.bf2demo\" \"not-a-valid-bookmark-timestamp\"",
					}),
				},
			),
			givenReference: reference,
			givenMaxAge:    time.Hour * 24 * 7,
			expectedDemoBookmarksCon: config.New(
				"C:\\Users\\default\\Documents\\Battlefield 2\\Profiles\\0001\\General.con",
				map[string]config.Value{},
			),
		},
		{
			name: "does nothing if bookmark key is missing",
			givenDemoBookmarksCon: config.New(
				"C:\\Users\\default\\Documents\\Battlefield 2\\Profiles\\0001\\General.con",
				map[string]config.Value{
					"some-other-key": *config.NewValue("some-value"),
				},
			),
			givenReference: reference,
			givenMaxAge:    time.Hour * 24 * 7,
			expectedDemoBookmarksCon: config.New(
				"C:\\Users\\default\\Documents\\Battlefield 2\\Profiles\\0001\\General.con",
				map[string]config.Value{
					"some-other-key": *config.NewValue("some-value"),
				},
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// GIVEN
			demoBookmarksCon := tt.givenDemoBookmarksCon

			// WHEN
			PurgeOldDemoBookmarks(demoBookmarksCon, tt.givenReference, tt.givenMaxAge)

			// THEN
			assert.Equal(t, tt.expectedDemoBookmarksCon, demoBookmarksCon)
		})
	}
}

func TestMarkAllVoiceOverHelpAsPlayed(t *testing.T) {
	type test struct {
		name               string
		givenGeneralCon    *config.Config
		expectedGeneralCon *config.Config
	}

	quoted := make([]string, 0, len(voiceOverHelpLines))
	for _, l := range voiceOverHelpLines {
		quoted = append(quoted, fmt.Sprintf("\"%s\"", l))
	}

	tests := []test{
		{
			name: "marks all voice over help lines as played in General.con",
			givenGeneralCon: config.New(
				"C:\\Users\\default\\Documents\\Battlefield 2\\Profiles\\0001\\General.con",
				map[string]config.Value{
					"GeneralSettings.setHUDTransparency": *config.NewValue("67.7346"),
				},
			),
			expectedGeneralCon: config.New(
				"C:\\Users\\default\\Documents\\Battlefield 2\\Profiles\\0001\\General.con",
				map[string]config.Value{
					"GeneralSettings.setHUDTransparency": *config.NewValue("67.7346"),
					GeneralConKeyVoiceOverHelpPlayed:     *config.NewValue(strings.Join(quoted, ";")),
				},
			),
		},
		{
			name: "overwrites existing voice over help lines which are marked as played in General.con",
			givenGeneralCon: config.New(
				"C:\\Users\\default\\Documents\\Battlefield 2\\Profiles\\0001\\General.con",
				map[string]config.Value{
					"GeneralSettings.setHUDTransparency": *config.NewValue("67.7346"),
					GeneralConKeyVoiceOverHelpPlayed:     *config.NewValue("GeneralSettings.setPlayedVOHelp \"HUD_HELP_COMMANDER_commanderApply\";GeneralSettings.setPlayedVOHelp \"HUD_HELP_KIT_SUPPORT_inVehicle\""),
				},
			),
			expectedGeneralCon: config.New(
				"C:\\Users\\default\\Documents\\Battlefield 2\\Profiles\\0001\\General.con",
				map[string]config.Value{
					"GeneralSettings.setHUDTransparency": *config.NewValue("67.7346"),
					GeneralConKeyVoiceOverHelpPlayed:     *config.NewValue(strings.Join(quoted, ";")),
				},
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// GIVEN
			generalCon := tt.givenGeneralCon

			// WHEN
			MarkAllVoiceOverHelpAsPlayed(generalCon)

			// THEN
			assert.Equal(t, tt.expectedGeneralCon, generalCon)
		})
	}
}

func buildDemoBookmark(serverName string, mapName string, timestamp time.Time) string {
	ts := formatDemoBookmarkTimestamp(timestamp)
	return fmt.Sprintf(
		"\"%[1]s\" \"%[2]s\" \"http://unused/%[1]s/%[2]s/%[3]s.bf2demo\" \"%[3]s\"",
		serverName,
		mapName,
		ts,
	)
}

func formatDemoBookmarkTimestamp(timestamp time.Time) string {
	return timestamp.Format(demoBookmarkTimestampLayout)
}
