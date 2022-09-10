//go:build unit

package bf2

import (
	"fmt"
	"path/filepath"
	"testing"

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
		givenProfile    string
		givenConfigFile ProfileConfigFile
		expect          func(h *MockHandler)
		wantConfig      *config.Config
		wantConfigPath  string
		wantErrContains string
	}

	tests := []test{
		{
			name:            "successfully reads Profile.con",
			givenProfile:    "0001",
			givenConfigFile: ProfileConfigFileProfileCon,
			expect: func(h *MockHandler) {
				basePath := "C:\\Users\\default\\Documents\\Battlefield 2\\Profiles"
				h.EXPECT().BuildBasePath(handler.GameBf2).Return(basePath, nil)
				h.EXPECT().ReadConfigFile(filepath.Join(basePath, "0001", "Profile.con")).Return(config.New(map[string]config.Value{
					"LocalProfile.setName": *config.NewValue("\"mister249\""),
				}), nil)
			},
			wantConfig: config.New(map[string]config.Value{
				"LocalProfile.setName": *config.NewValue("\"mister249\""),
			}),
			wantConfigPath: "C:\\Users\\default\\Documents\\Battlefield 2\\Profiles\\0001\\Profile.con",
		},
		{
			name:            "errors if base path cannot be determined",
			givenProfile:    "0001",
			givenConfigFile: ProfileConfigFileProfileCon,
			expect: func(h *MockHandler) {
				h.EXPECT().BuildBasePath(handler.GameBf2).Return("", fmt.Errorf("some-error"))
			},
			wantErrContains: "some-error",
		},
		{
			name:            "errors if config file read errors",
			givenProfile:    "0001",
			givenConfigFile: ProfileConfigFileProfileCon,
			expect: func(h *MockHandler) {
				basePath := "C:\\Users\\default\\Documents\\Battlefield 2\\Profiles"
				h.EXPECT().BuildBasePath(handler.GameBf2).Return(basePath, nil)
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
			readConfig, readConfigPath, err := ReadProfileConfigFile(h, tt.givenProfile, tt.givenConfigFile)

			// THEN
			if tt.wantErrContains != "" {
				require.ErrorContains(t, err, tt.wantErrContains)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.wantConfig, readConfig)
				assert.Equal(t, tt.wantConfigPath, readConfigPath)
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
				h.EXPECT().GetProfileKeys(handler.GameBf2).Return([]string{"0001"}, nil)
				h.EXPECT().ReadProfileConfig(handler.GameBf2, "0001").Return(config.New(map[string]config.Value{
					profileConKeyName: *config.NewValue("some-profile"),
				}), nil)
			},
			wantProfiles: []game.Profile{
				{
					Key:  "0001",
					Name: "some-profile",
				},
			},
		},
		{
			name: "ignores default profile",
			expect: func(h *MockHandler) {
				h.EXPECT().GetProfileKeys(handler.GameBf2).Return([]string{"0001", defaultProfileKey}, nil)
				h.EXPECT().ReadProfileConfig(handler.GameBf2, "0001").Return(config.New(map[string]config.Value{
					profileConKeyName: *config.NewValue("some-profile"),
				}), nil)
			},
			wantProfiles: []game.Profile{
				{
					Key:  "0001",
					Name: "some-profile",
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
				h.EXPECT().ReadProfileConfig(handler.GameBf2, "0001").Return(config.New(map[string]config.Value{
					"some-other-key": *config.NewValue("some-other-value"),
				}), nil)
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

func TestGetDefaultUserProfileCon(t *testing.T) {
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
				profileNumber := "0001"
				h.EXPECT().ReadGlobalConfig(handler.GameBf2).Return(config.New(map[string]config.Value{
					globalConKeyDefaultUserRef: *config.NewValue(profileNumber),
				}), nil)
				h.EXPECT().ReadProfileConfig(handler.GameBf2, profileNumber).Return(config.New(map[string]config.Value{
					profileConKeyGamespyNick: *config.NewValue("some-nick"),
					profileConKeyPassword:    *config.NewValue("some-encrypted-password"),
				}), nil)
			},
			expectedProfileCon: config.New(map[string]config.Value{
				profileConKeyGamespyNick: *config.NewValue("some-nick"),
				profileConKeyPassword:    *config.NewValue("some-encrypted-password"),
			}),
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
				profileNumber := "0001"
				h.EXPECT().ReadGlobalConfig(handler.GameBf2).Return(config.New(map[string]config.Value{
					globalConKeyDefaultUserRef: *config.NewValue(profileNumber),
				}), nil)
				h.EXPECT().ReadProfileConfig(handler.GameBf2, profileNumber).Return(nil, fmt.Errorf("some-profile-con-read-error"))
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
			profileCon, err := GetDefaultUserProfileCon(h)

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

func TestGetDefaultUserProfileNumber(t *testing.T) {
	type test struct {
		name                  string
		expect                func(h *MockHandler)
		expectedProfileNumber string
		wantErrContains       string
	}

	tests := []test{
		{
			name: "successfully retrieves default user profile number",
			expect: func(h *MockHandler) {
				h.EXPECT().ReadGlobalConfig(handler.GameBf2).Return(config.New(map[string]config.Value{
					globalConKeyDefaultUserRef: *config.NewValue("0001"),
				}), nil)
			},
			expectedProfileNumber: "0001",
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
				h.EXPECT().ReadGlobalConfig(handler.GameBf2).Return(config.New(map[string]config.Value{}), nil)
			},
			wantErrContains: "reference to default profile is missing from Global.con",
		},
		{
			name: "error if default user reference is non-numeric",
			expect: func(h *MockHandler) {
				h.EXPECT().ReadGlobalConfig(handler.GameBf2).Return(config.New(map[string]config.Value{
					globalConKeyDefaultUserRef: *config.NewValue("abcd"),
				}), nil)
			},
			wantErrContains: "reference to default profile in Global.con is not a valid profile number",
		},
		{
			name: "error if default user reference exceeds max length",
			expect: func(h *MockHandler) {
				h.EXPECT().ReadGlobalConfig(handler.GameBf2).Return(config.New(map[string]config.Value{
					globalConKeyDefaultUserRef: *config.NewValue("00001"),
				}), nil)
			},
			wantErrContains: "reference to default profile in Global.con is not a valid profile number",
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
			profileNumber, err := GetDefaultUserProfileNumber(h)

			// THEN
			if tt.wantErrContains != "" {
				require.ErrorContains(t, err, tt.wantErrContains)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expectedProfileNumber, profileNumber)
			}
		})
	}
}

func TestGetEncryptedProfileConLogin(t *testing.T) {
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
				profileCon.Delete(profileConKeyGamespyNick)
			},
			wantErrContains: "gamespy nickname is missing/empty",
		},
		{
			name: "fails if nickname is empty",
			prepareProfileConMap: func(profileCon *config.Config) {
				profileCon.SetValue(profileConKeyGamespyNick, *config.NewValue(""))
			},
			wantErrContains: "gamespy nickname is missing/empty",
		},
		{
			name: "fails if password is missing",
			prepareProfileConMap: func(profileCon *config.Config) {
				profileCon.Delete(profileConKeyPassword)
			},
			wantErrContains: "encrypted password is missing/empty",
		},
		{
			name: "fails if password is empty",
			prepareProfileConMap: func(profileCon *config.Config) {
				profileCon.SetValue(profileConKeyPassword, *config.NewValue(""))
			},
			wantErrContains: "encrypted password is missing/empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// GIVEN
			bytes := []byte(fmt.Sprintf("%s \"mister249\"\r\n%s \"some-encrypted-password\"\r\n", profileConKeyGamespyNick, profileConKeyPassword))
			profileCon := config.FromBytes(bytes)
			tt.prepareProfileConMap(profileCon)

			// WHEN
			nickname, encryptedPassword, err := GetEncryptedProfileConLogin(profileCon)

			// THEN
			if tt.wantErrContains != "" {
				require.ErrorContains(t, err, tt.wantErrContains)
			} else {
				require.NoError(t, err)
				expectedNickname, err := profileCon.GetValue(profileConKeyGamespyNick)
				require.NoError(t, err)
				assert.Equal(t, expectedNickname.String(), nickname)
				expectedPassword, err := profileCon.GetValue(profileConKeyPassword)
				require.NoError(t, err)
				assert.Equal(t, expectedPassword.String(), encryptedPassword)
			}
		})
	}
}

func TestPurgeGeneralConServerHistory(t *testing.T) {
	type test struct {
		name               string
		givenGeneralCon    *config.Config
		expectedGeneralCon *config.Config
	}

	tests := []test{
		{
			name: "removes single server history from General.con",
			givenGeneralCon: config.New(map[string]config.Value{
				"GeneralSettings.setHUDTransparency": *config.NewValue("67.7346"),
				"GeneralSettings.addServerHistory":   *config.NewValue("\"135.125.56.26\" 29940 \"=DOG= No Explosives (Infantry)\" 1025"),
			}),
			expectedGeneralCon: config.New(map[string]config.Value{
				"GeneralSettings.setHUDTransparency": *config.NewValue("67.7346"),
			}),
		},
		{
			name: "removes multiple server history items from General.con",
			givenGeneralCon: config.New(map[string]config.Value{
				"GeneralSettings.setHUDTransparency": *config.NewValue("67.7346"),
				"GeneralSettings.addServerHistory":   *config.NewValue("\"135.125.56.26\" 29940 \"=DOG= No Explosives (Infantry)\" 1025;\"37.230.210.130\" 29900 \"PlayBF2! T~GAMER #1 Allmaps\" 360"),
			}),
			expectedGeneralCon: config.New(map[string]config.Value{
				"GeneralSettings.setHUDTransparency": *config.NewValue("67.7346"),
			}),
		},
		{
			name: "does nothing if General.con does not contain any server history items",
			givenGeneralCon: config.New(map[string]config.Value{
				"GeneralSettings.setHUDTransparency": *config.NewValue("67.7346"),
			}),
			expectedGeneralCon: config.New(map[string]config.Value{
				"GeneralSettings.setHUDTransparency": *config.NewValue("67.7346"),
			}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// GIVEN
			generalCon := tt.givenGeneralCon

			// WHEN
			PurgeGeneralConServerHistory(generalCon)

			// THEN
			assert.Equal(t, tt.expectedGeneralCon, generalCon)
		})
	}
}
