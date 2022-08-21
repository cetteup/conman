//go:build unit

package handler

import (
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
			name:            "error for unsupported game",
			givenGame:       "not-a-supported-game",
			expect:          func(repository *MockFileRepository, documentsDirPath string) {},
			wantErrContains: "game not supported",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// GIVEN
			handler, mockRepository := getHandlerWithDependencies(t)
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
				assert.Equal(t, config.New(map[string]config.Value{}), globalConfig)
			}
		})
	}
}

func TestHandler_ReadProfileConfig(t *testing.T) {
	type test struct {
		name            string
		givenGame       Game
		givenProfile    string
		expect          func(repository *MockFileRepository, documentsDirPath string)
		wantErrContains string
	}

	tests := []test{
		{
			name:         "successfully reads config file",
			givenGame:    GameBf2,
			givenProfile: "0001",
			expect: func(repository *MockFileRepository, documentsDirPath string) {
				repository.EXPECT().ReadFile(gomock.Eq(filepath.Join(documentsDirPath, bf2GameDirName, profilesDirName, "0001", profileConFileName))).Return([]byte{}, nil)
			},
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
			handler, mockRepository := getHandlerWithDependencies(t)
			documentsDirPath, err := windows.KnownFolderPath(windows.FOLDERID_Documents, windows.KF_FLAG_DEFAULT)
			require.NoError(t, err)

			// EXPECT
			tt.expect(mockRepository, documentsDirPath)

			// WHEN
			profileConfig, err := handler.ReadProfileConfig(tt.givenGame, tt.givenProfile)

			// THEN
			if tt.wantErrContains != "" {
				require.ErrorContains(t, err, tt.wantErrContains)
			} else {
				require.NoError(t, err)
				assert.Equal(t, config.New(map[string]config.Value{}), profileConfig)
			}
		})
	}
}

func TestHandler_WriteConfigFile(t *testing.T) {
	type test struct {
		name            string
		givenPath       string
		givenConfig     *config.Config
		expect          func(repository *MockFileRepository)
		wantErrContains string
	}

	tests := []test{
		{
			name:      "successfully writes config file",
			givenPath: "C:\\Users\\default\\Documents\\Battlefield 2\\Profiles\\Global.con",
			givenConfig: config.New(map[string]config.Value{
				"GlobalSettings.setNamePrefix": *config.NewValue("\"=DOG=\""),
			}),
			expect: func(repository *MockFileRepository) {
				repository.EXPECT().WriteFile("C:\\Users\\default\\Documents\\Battlefield 2\\Profiles\\Global.con", []byte("GlobalSettings.setNamePrefix \"=DOG=\""), os.FileMode(0666))
			},
		},
	}

	for _, tt := range tests {
		// GIVEN
		handler, mockRepository := getHandlerWithDependencies(t)

		// EXPECT
		tt.expect(mockRepository)

		// WHEN
		err := handler.WriteConfigFile(tt.givenPath, tt.givenConfig)

		// THEN
		if tt.wantErrContains != "" {
			require.ErrorContains(t, err, tt.wantErrContains)
		} else {
			require.NoError(t, err)
		}
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
			handler, _ := getHandlerWithDependencies(t)
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

func getHandlerWithDependencies(t *testing.T) (*Handler, *MockFileRepository) {
	ctrl := gomock.NewController(t)
	mockRepository := NewMockFileRepository(ctrl)
	return New(mockRepository), mockRepository
}
