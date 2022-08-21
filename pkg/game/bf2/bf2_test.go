//go:build unit

package bf2

import (
	"fmt"
	"testing"

	"github.com/cetteup/conman/pkg/config"
	"github.com/cetteup/conman/pkg/handler"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
