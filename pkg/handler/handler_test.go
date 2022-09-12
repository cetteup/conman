//go:build unit

package handler

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/cetteup/conman/pkg/config"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/sys/windows"
)

func TestHandler_ReadGlobalConfig(t *testing.T) {
	type test struct {
		name            string
		givenGame       Game
		expect          func(repository *MockFileRepository, documentsDirPath string)
		wantErrContains string
	}

	tests := []test{
		{
			name:      "successfully reads config file",
			givenGame: GameBf2,
			expect: func(repository *MockFileRepository, documentsDirPath string) {
				repository.EXPECT().ReadFile(gomock.Eq(filepath.Join(documentsDirPath, bf2GameDirName, profilesDirName, globalConFileName))).Return([]byte{}, nil)
			},
		},
		{
			name:      "error reading config file",
			givenGame: GameBf2,
			expect: func(repository *MockFileRepository, documentsDirPath string) {
				repository.EXPECT().ReadFile(gomock.Eq(filepath.Join(documentsDirPath, bf2GameDirName, profilesDirName, globalConFileName))).Return([]byte{}, fmt.Errorf("some-error"))
			},
			wantErrContains: "some-error",
		},
		{
			name:            "error for unsupported game",
			givenGame:       "not-a-supported-game",
			expect:          func(repository *MockFileRepository, documentsDirPath string) {},
			wantErrContains: "game not supported",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// GIVEN
			_, handler, mockRepository := getHandlerWithDependencies(t)
			documentsDirPath, err := windows.KnownFolderPath(windows.FOLDERID_Documents, windows.KF_FLAG_DEFAULT)
			require.NoError(t, err)

			// EXPECT
			tt.expect(mockRepository, documentsDirPath)

			// WHEN
			globalConfig, err := handler.ReadGlobalConfig(tt.givenGame)

			// THEN
			if tt.wantErrContains != "" {
				require.ErrorContains(t, err, tt.wantErrContains)
			} else {
				require.NoError(t, err)
				assert.Equal(t, config.New(filepath.Join(documentsDirPath, bf2GameDirName, profilesDirName, globalConFileName), map[string]config.Value{}), globalConfig)
			}
		})
	}
}

func TestHandler_GetProfileKeys(t *testing.T) {
	type test struct {
		name            string
		givenGame       Game
		expect          func(ctrl *gomock.Controller, repository *MockFileRepository, documentsDirPath string)
		wantProfileKeys []string
		wantErrContains string
	}

	tests := []test{
		{
			name:      "successfully gets profile keys",
			givenGame: GameBf2,
			expect: func(ctrl *gomock.Controller, repository *MockFileRepository, documentsDirPath string) {
				dirEntry := NewMockDirEntry(ctrl)
				dirEntry.EXPECT().IsDir().Return(true)
				dirEntry.EXPECT().Name().Return("0001").Times(2)
				fileEntry := NewMockDirEntry(ctrl)
				fileEntry.EXPECT().IsDir().Return(false)
				repository.EXPECT().ReadDir(gomock.Eq(filepath.Join(documentsDirPath, bf2GameDirName, profilesDirName))).Return([]os.DirEntry{
					dirEntry,
					fileEntry,
				}, nil)
				repository.EXPECT().FileExists(gomock.Eq(filepath.Join(documentsDirPath, bf2GameDirName, profilesDirName, "0001", profileConFileName))).Return(true, nil)
			},
			wantProfileKeys: []string{
				"0001",
			},
		},
		{
			name:      "ignores profile dir which does not contain Profile.con",
			givenGame: GameBf2,
			expect: func(ctrl *gomock.Controller, repository *MockFileRepository, documentsDirPath string) {
				dirEntry := NewMockDirEntry(ctrl)
				dirEntry.EXPECT().IsDir().Return(true)
				dirEntry.EXPECT().Name().Return("0001").Times(2)
				invalidDirEntry := NewMockDirEntry(ctrl)
				invalidDirEntry.EXPECT().IsDir().Return(true)
				invalidDirEntry.EXPECT().Name().Return("0002")
				repository.EXPECT().ReadDir(gomock.Eq(filepath.Join(documentsDirPath, bf2GameDirName, profilesDirName))).Return([]os.DirEntry{
					dirEntry,
					invalidDirEntry,
				}, nil)
				repository.EXPECT().FileExists(gomock.Eq(filepath.Join(documentsDirPath, bf2GameDirName, profilesDirName, "0001", profileConFileName))).Return(true, nil)
				repository.EXPECT().FileExists(gomock.Eq(filepath.Join(documentsDirPath, bf2GameDirName, profilesDirName, "0002", profileConFileName))).Return(false, nil)
			},
			wantProfileKeys: []string{
				"0001",
			},
		},
		{
			name:      "error listing profile dir entries",
			givenGame: GameBf2,
			expect: func(ctrl *gomock.Controller, repository *MockFileRepository, documentsDirPath string) {
				repository.EXPECT().ReadDir(gomock.Eq(filepath.Join(documentsDirPath, bf2GameDirName, profilesDirName))).Return([]os.DirEntry{}, fmt.Errorf("some-error"))
			},
			wantErrContains: "some-error",
		},
		{
			name:      "error checking if profile dir is valid",
			givenGame: GameBf2,
			expect: func(ctrl *gomock.Controller, repository *MockFileRepository, documentsDirPath string) {
				dirEntry := NewMockDirEntry(ctrl)
				dirEntry.EXPECT().IsDir().Return(true)
				dirEntry.EXPECT().Name().Return("0001").Times(1)
				repository.EXPECT().ReadDir(gomock.Eq(filepath.Join(documentsDirPath, bf2GameDirName, profilesDirName))).Return([]os.DirEntry{
					dirEntry,
				}, nil)
				repository.EXPECT().FileExists(gomock.Eq(filepath.Join(documentsDirPath, bf2GameDirName, profilesDirName, "0001", profileConFileName))).Return(false, fmt.Errorf("some-error"))
			},
			wantErrContains: "some-error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// GIVEN
			ctrl, handler, mockRepository := getHandlerWithDependencies(t)
			documentsDirPath, err := windows.KnownFolderPath(windows.FOLDERID_Documents, windows.KF_FLAG_DEFAULT)
			require.NoError(t, err)

			// EXPECT
			tt.expect(ctrl, mockRepository, documentsDirPath)

			// WHEN
			profileKeys, err := handler.GetProfileKeys(tt.givenGame)

			// THEN
			if tt.wantErrContains != "" {
				require.ErrorContains(t, err, tt.wantErrContains)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.wantProfileKeys, profileKeys)
			}
		})
	}
}

func TestHandler_IsValidProfileKey(t *testing.T) {
	type test struct {
		name                  string
		givenGame             Game
		givenProfileKey       string
		expect                func(ctrl *gomock.Controller, repository *MockFileRepository, documentsDirPath string)
		wantIsValidProfileKey bool
		wantErrContains       string
	}

	tests := []test{
		{
			name:            "true for key of existing profile",
			givenGame:       GameBf2,
			givenProfileKey: "0001",
			expect: func(ctrl *gomock.Controller, repository *MockFileRepository, documentsDirPath string) {
				dirEntry := NewMockDirEntry(ctrl)
				dirEntry.EXPECT().IsDir().Return(true)
				dirEntry.EXPECT().Name().Return("0001").Times(2)
				repository.EXPECT().ReadDir(gomock.Eq(filepath.Join(documentsDirPath, bf2GameDirName, profilesDirName))).Return([]os.DirEntry{dirEntry}, nil)
				repository.EXPECT().FileExists(gomock.Eq(filepath.Join(documentsDirPath, bf2GameDirName, profilesDirName, "0001", profileConFileName))).Return(true, nil)
			},
			wantIsValidProfileKey: true,
		},
		{
			name:      "false for profile dir which does not contain a Profile.con",
			givenGame: GameBf2,
			expect: func(ctrl *gomock.Controller, repository *MockFileRepository, documentsDirPath string) {
				invalidDirEntry := NewMockDirEntry(ctrl)
				invalidDirEntry.EXPECT().IsDir().Return(true)
				invalidDirEntry.EXPECT().Name().Return("0001")
				repository.EXPECT().ReadDir(gomock.Eq(filepath.Join(documentsDirPath, bf2GameDirName, profilesDirName))).Return([]os.DirEntry{invalidDirEntry}, nil)
				repository.EXPECT().FileExists(gomock.Eq(filepath.Join(documentsDirPath, bf2GameDirName, profilesDirName, "0001", profileConFileName))).Return(false, nil)
			},
			wantIsValidProfileKey: false,
		},
		{
			name:      "error listing profile dir entries",
			givenGame: GameBf2,
			expect: func(ctrl *gomock.Controller, repository *MockFileRepository, documentsDirPath string) {
				repository.EXPECT().ReadDir(gomock.Eq(filepath.Join(documentsDirPath, bf2GameDirName, profilesDirName))).Return([]os.DirEntry{}, fmt.Errorf("some-error"))
			},
			wantErrContains: "some-error",
		},
		{
			name:      "error checking if profile dir is valid",
			givenGame: GameBf2,
			expect: func(ctrl *gomock.Controller, repository *MockFileRepository, documentsDirPath string) {
				dirEntry := NewMockDirEntry(ctrl)
				dirEntry.EXPECT().IsDir().Return(true)
				dirEntry.EXPECT().Name().Return("0001").Times(1)
				repository.EXPECT().ReadDir(gomock.Eq(filepath.Join(documentsDirPath, bf2GameDirName, profilesDirName))).Return([]os.DirEntry{
					dirEntry,
				}, nil)
				repository.EXPECT().FileExists(gomock.Eq(filepath.Join(documentsDirPath, bf2GameDirName, profilesDirName, "0001", profileConFileName))).Return(false, fmt.Errorf("some-error"))
			},
			wantErrContains: "some-error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// GIVEN
			ctrl, handler, mockRepository := getHandlerWithDependencies(t)
			documentsDirPath, err := windows.KnownFolderPath(windows.FOLDERID_Documents, windows.KF_FLAG_DEFAULT)
			require.NoError(t, err)

			// EXPECT
			tt.expect(ctrl, mockRepository, documentsDirPath)

			// WHEN
			isValidProfileKey, err := handler.IsValidProfileKey(tt.givenGame, tt.givenProfileKey)

			// THEN
			if tt.wantErrContains != "" {
				require.ErrorContains(t, err, tt.wantErrContains)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.wantIsValidProfileKey, isValidProfileKey)
			}
		})
	}
}

func TestHandler_ReadProfileConfig(t *testing.T) {
	type test struct {
		name            string
		givenGame       Game
		givenProfileKey string
		expect          func(repository *MockFileRepository, documentsDirPath string)
		wantErrContains string
	}

	tests := []test{
		{
			name:            "successfully reads config file",
			givenGame:       GameBf2,
			givenProfileKey: "0001",
			expect: func(repository *MockFileRepository, documentsDirPath string) {
				repository.EXPECT().ReadFile(gomock.Eq(filepath.Join(documentsDirPath, bf2GameDirName, profilesDirName, "0001", profileConFileName))).Return([]byte{}, nil)
			},
		},
		{
			name:            "error reading config file",
			givenGame:       GameBf2,
			givenProfileKey: "0001",
			expect: func(repository *MockFileRepository, documentsDirPath string) {
				repository.EXPECT().ReadFile(gomock.Eq(filepath.Join(documentsDirPath, bf2GameDirName, profilesDirName, "0001", profileConFileName))).Return([]byte{}, fmt.Errorf("some-error"))
			},
			wantErrContains: "some-error",
		},
		{
			name:            "error for unsupported game",
			givenGame:       "not-a-supported-game",
			expect:          func(repository *MockFileRepository, documentsDirPath string) {},
			wantErrContains: "game not supported",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// GIVEN
			_, handler, mockRepository := getHandlerWithDependencies(t)
			documentsDirPath, err := windows.KnownFolderPath(windows.FOLDERID_Documents, windows.KF_FLAG_DEFAULT)
			require.NoError(t, err)

			// EXPECT
			tt.expect(mockRepository, documentsDirPath)

			// WHEN
			profileConfig, err := handler.ReadProfileConfig(tt.givenGame, tt.givenProfileKey)

			// THEN
			if tt.wantErrContains != "" {
				require.ErrorContains(t, err, tt.wantErrContains)
			} else {
				require.NoError(t, err)
				assert.Equal(t, config.New(filepath.Join(documentsDirPath, bf2GameDirName, profilesDirName, "0001", profileConFileName), map[string]config.Value{}), profileConfig)
			}
		})
	}
}

func TestHandler_WriteConfigFile(t *testing.T) {
	type test struct {
		name            string
		givenConfig     *config.Config
		expect          func(repository *MockFileRepository)
		wantErrContains string
	}

	tests := []test{
		{
			name: "successfully writes config file",
			givenConfig: config.New(
				"C:\\Users\\default\\Documents\\Battlefield 2\\Profiles\\Global.con",
				map[string]config.Value{
					"GlobalSettings.setNamePrefix": *config.NewValue("\"=DOG=\""),
				},
			),
			expect: func(repository *MockFileRepository) {
				repository.EXPECT().WriteFile("C:\\Users\\default\\Documents\\Battlefield 2\\Profiles\\Global.con", []byte("GlobalSettings.setNamePrefix \"=DOG=\"\r\n"), os.FileMode(0666)).Return(nil)
			},
		},
		{
			name: "error writing config file",
			givenConfig: config.New(
				"C:\\Users\\default\\Documents\\Battlefield 2\\Profiles\\Global.con",
				map[string]config.Value{
					"GlobalSettings.setNamePrefix": *config.NewValue("\"=DOG=\""),
				},
			),
			expect: func(repository *MockFileRepository) {
				repository.EXPECT().WriteFile("C:\\Users\\default\\Documents\\Battlefield 2\\Profiles\\Global.con", []byte("GlobalSettings.setNamePrefix \"=DOG=\"\r\n"), os.FileMode(0666)).Return(fmt.Errorf("some-error"))
			},
			wantErrContains: "some-error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// GIVEN
			_, handler, mockRepository := getHandlerWithDependencies(t)

			// EXPECT
			tt.expect(mockRepository)

			// WHEN
			err := handler.WriteConfigFile(tt.givenConfig)

			// THEN
			if tt.wantErrContains != "" {
				require.ErrorContains(t, err, tt.wantErrContains)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestHandler_PurgeShaderCache(t *testing.T) {
	type test struct {
		name            string
		givenGame       Game
		expect          func(repository *MockFileRepository, documentsDirPath string)
		wantErrContains string
	}

	tests := []test{
		{
			name:      "successfully purges shader cache",
			givenGame: GameBf2,
			expect: func(repository *MockFileRepository, documentsDirPath string) {
				basePath := filepath.Join(documentsDirPath, bf2GameDirName)
				pattern := filepath.Join(basePath, modsDirName, "*", cacheDirName, "*")
				bf2CachePath := filepath.Join(basePath, modsDirName, "bf2", cacheDirName, "{D7B71E3E-5F43-11CF-726F-0D3C0EC2D335}_112_1")
				xpackCachePath := filepath.Join(basePath, modsDirName, "xpack", cacheDirName, "{D7B71E3E-5F43-11CF-726F-0D3C0EC2D335}_112_2")
				repository.EXPECT().Glob(pattern).Return([]string{
					bf2CachePath,
					xpackCachePath,
				}, nil)
				repository.EXPECT().RemoveAll(bf2CachePath)
				repository.EXPECT().RemoveAll(xpackCachePath)
			},
		},
		{
			name:      "does nothing if no cache folders exist",
			givenGame: GameBf2,
			expect: func(repository *MockFileRepository, documentsDirPath string) {
				basePath := filepath.Join(documentsDirPath, bf2GameDirName)
				pattern := filepath.Join(basePath, modsDirName, "*", cacheDirName, "*")
				repository.EXPECT().Glob(pattern).Return([]string{}, nil)
			},
		},
		{
			name:      "error in glob pattern",
			givenGame: GameBf2,
			expect: func(repository *MockFileRepository, documentsDirPath string) {
				basePath := filepath.Join(documentsDirPath, bf2GameDirName)
				pattern := filepath.Join(basePath, modsDirName, "*", cacheDirName, "*")
				repository.EXPECT().Glob(pattern).Return([]string{}, fmt.Errorf("some-error"))
			},
			wantErrContains: "some-error",
		},
		{
			name:      "error removing cache folder",
			givenGame: GameBf2,
			expect: func(repository *MockFileRepository, documentsDirPath string) {
				basePath := filepath.Join(documentsDirPath, bf2GameDirName)
				pattern := filepath.Join(basePath, modsDirName, "*", cacheDirName, "*")
				bf2CachePath := filepath.Join(basePath, modsDirName, "bf2", cacheDirName, "{D7B71E3E-5F43-11CF-726F-0D3C0EC2D335}_112_1")
				repository.EXPECT().Glob(pattern).Return([]string{
					bf2CachePath,
				}, nil)
				repository.EXPECT().RemoveAll(bf2CachePath).Return(fmt.Errorf("some-error"))
			},
			wantErrContains: "some-error",
		},
		{
			name:            "error for unsupported game",
			givenGame:       "not-a-supported-game",
			expect:          func(repository *MockFileRepository, documentsDirPath string) {},
			wantErrContains: "game not supported",
		},
		// TODO Add test for supported game not suppored by this action once more games are implemented
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// GIVEN
			_, handler, mockRepository := getHandlerWithDependencies(t)
			documentsDirPath, err := windows.KnownFolderPath(windows.FOLDERID_Documents, windows.KF_FLAG_DEFAULT)
			require.NoError(t, err)

			// EXPECT
			tt.expect(mockRepository, documentsDirPath)

			// WHEN
			err = handler.PurgeShaderCache(tt.givenGame)

			// THEN
			if tt.wantErrContains != "" {
				require.ErrorContains(t, err, tt.wantErrContains)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestHandler_PurgeLogoCache(t *testing.T) {
	type test struct {
		name            string
		givenGame       Game
		expect          func(repository *MockFileRepository, documentsDirPath string)
		wantErrContains string
	}

	tests := []test{
		{
			name:      "successfully purges logo cache",
			givenGame: GameBf2,
			expect: func(repository *MockFileRepository, documentsDirPath string) {
				basePath := filepath.Join(documentsDirPath, bf2GameDirName)
				pattern := filepath.Join(basePath, logoCacheDirName, "*")
				dogclanDotNetCachePath := filepath.Join(basePath, logoCacheDirName, "www.dogclan.net")
				superinfantryclanDotComCachePath := filepath.Join(basePath, logoCacheDirName, "www.superinfantryclan.com")
				repository.EXPECT().Glob(pattern).Return([]string{
					dogclanDotNetCachePath,
					superinfantryclanDotComCachePath,
				}, nil)
				repository.EXPECT().RemoveAll(dogclanDotNetCachePath)
				repository.EXPECT().RemoveAll(superinfantryclanDotComCachePath)
			},
		},
		{
			name:      "does nothing if no cache folders exist",
			givenGame: GameBf2,
			expect: func(repository *MockFileRepository, documentsDirPath string) {
				basePath := filepath.Join(documentsDirPath, bf2GameDirName)
				pattern := filepath.Join(basePath, logoCacheDirName, "*")
				repository.EXPECT().Glob(pattern).Return([]string{}, nil)
			},
		},
		{
			name:      "error in glob pattern",
			givenGame: GameBf2,
			expect: func(repository *MockFileRepository, documentsDirPath string) {
				basePath := filepath.Join(documentsDirPath, bf2GameDirName)
				pattern := filepath.Join(basePath, logoCacheDirName, "*")
				repository.EXPECT().Glob(pattern).Return([]string{}, fmt.Errorf("some-error"))
			},
			wantErrContains: "some-error",
		},
		{
			name:      "error removing cache folder",
			givenGame: GameBf2,
			expect: func(repository *MockFileRepository, documentsDirPath string) {
				basePath := filepath.Join(documentsDirPath, bf2GameDirName)
				pattern := filepath.Join(basePath, logoCacheDirName, "*")
				dogclanDotNetCachePath := filepath.Join(basePath, logoCacheDirName, "www.dogclan.net")
				repository.EXPECT().Glob(pattern).Return([]string{
					dogclanDotNetCachePath,
				}, nil)
				repository.EXPECT().RemoveAll(dogclanDotNetCachePath).Return(fmt.Errorf("some-error"))
			},
			wantErrContains: "some-error",
		},
		{
			name:            "error for unsupported game",
			givenGame:       "not-a-supported-game",
			expect:          func(repository *MockFileRepository, documentsDirPath string) {},
			wantErrContains: "game not supported",
		},
		// TODO Add test for supported game not suppored by this action once more games are implemented
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// GIVEN
			_, handler, mockRepository := getHandlerWithDependencies(t)
			documentsDirPath, err := windows.KnownFolderPath(windows.FOLDERID_Documents, windows.KF_FLAG_DEFAULT)
			require.NoError(t, err)

			// EXPECT
			tt.expect(mockRepository, documentsDirPath)

			// WHEN
			err = handler.PurgeLogoCache(tt.givenGame)

			// THEN
			if tt.wantErrContains != "" {
				require.ErrorContains(t, err, tt.wantErrContains)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestHandler_BuildProfilesFolderPath(t *testing.T) {
	type test struct {
		name                     string
		givenGame                Game
		expectedPathFromDocument string
		wantErrContains          string
	}

	tests := []test{
		{
			name:                     "builds base path for Battlefield 2",
			givenGame:                GameBf2,
			expectedPathFromDocument: "Battlefield 2\\Profiles",
		},
		{
			name:            "error for unsupported game",
			givenGame:       "not-a-supported-game",
			wantErrContains: "game not supported",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// GIVEN
			_, handler, _ := getHandlerWithDependencies(t)
			documentsDirPath, err := windows.KnownFolderPath(windows.FOLDERID_Documents, windows.KF_FLAG_DEFAULT)
			require.NoError(t, err)

			// WHEN
			basePath, err := handler.BuildProfilesFolderPath(tt.givenGame)

			// THEN
			if tt.wantErrContains != "" {
				require.ErrorContains(t, err, tt.wantErrContains)
			} else {
				require.NoError(t, err)
				assert.Equal(t, filepath.Join(documentsDirPath, tt.expectedPathFromDocument), basePath)
			}
		})
	}
}

func getHandlerWithDependencies(t *testing.T) (*gomock.Controller, *Handler, *MockFileRepository) {
	ctrl := gomock.NewController(t)
	mockRepository := NewMockFileRepository(ctrl)
	return ctrl, New(mockRepository), mockRepository
}
