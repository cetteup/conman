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
				repository.EXPECT().WriteFile("C:\\Users\\default\\Documents\\Battlefield 2\\Profiles\\Global.con", []byte("GlobalSettings.setNamePrefix \"=DOG=\""), os.FileMode(0666)).Return(nil)
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
				repository.EXPECT().WriteFile("C:\\Users\\default\\Documents\\Battlefield 2\\Profiles\\Global.con", []byte("GlobalSettings.setNamePrefix \"=DOG=\""), os.FileMode(0666)).Return(fmt.Errorf("some-error"))
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

func TestHandler_BuildBasePath(t *testing.T) {
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
			basePath, err := handler.BuildBasePath(tt.givenGame)

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
