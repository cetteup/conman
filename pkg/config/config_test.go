//go:build unit

package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfigFromBytes(t *testing.T) {
	type test struct {
		name           string
		givenPath      string
		givenData      string
		expectedConfig Config
	}

	tests := []test{
		{
			name:      "parses config with unix line breaks",
			givenPath: "C:\\Users\\default\\Documents\\Battlefield 2\\Profiles\\Global.con",
			givenData: "GlobalSettings.setDefaultUser \"0010\"\nGlobalSettings.setNamePrefix \"=PRE=\"\n",
			expectedConfig: Config{
				Path: "C:\\Users\\default\\Documents\\Battlefield 2\\Profiles\\Global.con",
				content: map[string]Value{
					"GlobalSettings.setDefaultUser": {content: "\"0010\""},
					"GlobalSettings.setNamePrefix":  {content: "\"=PRE=\""},
				},
			},
		},
		{
			name:      "parses config with windows line breaks",
			givenPath: "C:\\Users\\default\\Documents\\Battlefield 2\\Profiles\\Global.con",
			givenData: "GlobalSettings.setDefaultUser \"0010\"\r\nGlobalSettings.setNamePrefix \"=PRE=\"\r\n",
			expectedConfig: Config{
				Path: "C:\\Users\\default\\Documents\\Battlefield 2\\Profiles\\Global.con",
				content: map[string]Value{
					"GlobalSettings.setDefaultUser": {content: "\"0010\""},
					"GlobalSettings.setNamePrefix":  {content: "\"=PRE=\""},
				},
			},
		},
		{
			name:      "parses multiple lines with same key",
			givenPath: "C:\\Users\\default\\Documents\\Battlefield 2\\Profiles\\0001\\General.con",
			givenData: "GeneralSettings.setPlayedVOHelp \"HUD_HELP_A\"\nGeneralSettings.setPlayedVOHelp \"HUD_HELP_B\"\n",
			expectedConfig: Config{
				Path: "C:\\Users\\default\\Documents\\Battlefield 2\\Profiles\\0001\\General.con",
				content: map[string]Value{
					"GeneralSettings.setPlayedVOHelp": {content: "\"HUD_HELP_A\";\"HUD_HELP_B\""},
				},
			},
		},
		{
			name:      "parses empty config",
			givenPath: "C:\\Users\\default\\Documents\\Battlefield 2\\Profiles\\0001\\General.con",
			givenData: "",
			expectedConfig: Config{
				Path:    "C:\\Users\\default\\Documents\\Battlefield 2\\Profiles\\0001\\General.con",
				content: map[string]Value{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// WHEN
			config := FromBytes(tt.givenPath, []byte(tt.givenData))

			// THEN
			assert.Equal(t, &tt.expectedConfig, config)
		})
	}
}

func TestConfig_GetValue(t *testing.T) {
	type test struct {
		name            string
		givenConfig     Config
		givenKey        string
		expectedValue   Value
		wantErrContains string
	}

	tests := []test{
		{
			name: "successfully retrieves value",
			givenConfig: Config{
				content: map[string]Value{
					"some-key": {content: "some-value"},
				},
			},
			givenKey:      "some-key",
			expectedValue: Value{content: "some-value"},
		},
		{
			name: "error for non-existing key",
			givenConfig: Config{
				content: map[string]Value{
					"some-key": {content: "some-value"},
				},
			},
			givenKey:        "some-other-key",
			wantErrContains: "no such key in config",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// WHEN
			value, err := tt.givenConfig.GetValue(tt.givenKey)

			// THEN
			if tt.wantErrContains != "" {
				require.ErrorContains(t, err, tt.wantErrContains)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expectedValue, value)
			}
		})
	}
}

func TestConfig_SetValue(t *testing.T) {
	type test struct {
		name           string
		givenConfig    Config
		givenKey       string
		givenValue     Value
		expectedConfig Config
	}

	tests := []test{
		{
			name: "adds value under new key",
			givenConfig: Config{
				content: map[string]Value{
					"some-key": {content: "some-value"},
				},
			},
			givenKey:   "other-key",
			givenValue: Value{content: "other-value"},
			expectedConfig: Config{
				content: map[string]Value{
					"some-key":  {content: "some-value"},
					"other-key": {content: "other-value"},
				},
			},
		},
		{
			name: "overwrites value at existing key",
			givenConfig: Config{
				content: map[string]Value{
					"some-key": {content: "old-value"},
				},
			},
			givenKey:   "some-key",
			givenValue: Value{content: "new-value"},
			expectedConfig: Config{
				content: map[string]Value{
					"some-key": {content: "new-value"},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// GIVEN
			config := tt.givenConfig

			// WHEN
			config.SetValue(tt.givenKey, tt.givenValue)

			// THEN
			assert.Equal(t, tt.expectedConfig, config)
		})
	}
}

func TestConfig_Delete(t *testing.T) {
	type test struct {
		name           string
		givenConfig    Config
		givenKey       string
		expectedConfig Config
	}

	tests := []test{
		{
			name: "removes existing key",
			givenConfig: Config{
				content: map[string]Value{
					"some-key": {content: "some-value"},
				},
			},
			givenKey: "some-key",
			expectedConfig: Config{
				content: map[string]Value{},
			},
		},
		{
			name: "noop for non-existing key",
			givenConfig: Config{
				content: map[string]Value{
					"some-key": {content: "some-value"},
				},
			},
			givenKey: "other-key",
			expectedConfig: Config{
				content: map[string]Value{
					"some-key": {content: "some-value"},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// GIVEN
			config := tt.givenConfig

			// WHEN
			config.Delete(tt.givenKey)

			// THEN
			assert.Equal(t, tt.expectedConfig, config)
		})
	}
}

func TestConfig_ToBytes(t *testing.T) {
	type test struct {
		name         string
		givenConfig  Config
		expectedData string
	}

	tests := []test{
		{
			name: "serializes config with unquoted single value",
			givenConfig: Config{
				content: map[string]Value{
					"LocalProfile.setNumTimesLoggedIn": {content: "8"},
				},
			},
			expectedData: "LocalProfile.setNumTimesLoggedIn 8\r\n",
		},
		{
			name: "serializes config with quoted single value",
			givenConfig: Config{
				content: map[string]Value{
					"LocalProfile.setName": {content: "\"mister249\""},
				},
			},
			expectedData: "LocalProfile.setName \"mister249\"\r\n",
		},
		{
			name: "serializes config with unquoted multi value",
			givenConfig: Config{
				content: map[string]Value{
					"LocalProfile.setNumTimesLoggedIn": {content: "8;9;10"},
				},
			},
			expectedData: "LocalProfile.setNumTimesLoggedIn 10\r\nLocalProfile.setNumTimesLoggedIn 8\r\nLocalProfile.setNumTimesLoggedIn 9\r\n",
		},
		{
			name: "serializes config with quoted multi value",
			givenConfig: Config{
				content: map[string]Value{
					"LocalProfile.setName": {content: "\"mister249\";\"mister250\";\"mister251\""},
				},
			},
			expectedData: "LocalProfile.setName \"mister249\"\r\nLocalProfile.setName \"mister250\"\r\nLocalProfile.setName \"mister251\"\r\n",
		},
		{
			name: "serializes config with server history entries",
			givenConfig: Config{
				content: map[string]Value{
					"GeneralSettings.addServerHistory": {content: "\"135.125.56.26\" 29940 \"=DOG= No Explosives (Infantry)\" 1025;\"37.230.210.130\" 29900 \"PlayBF2! T~GAMER #1 Allmaps\" 360"},
				},
			},
			expectedData: "GeneralSettings.addServerHistory \"135.125.56.26\" 29940 \"=DOG= No Explosives (Infantry)\" 1025\r\nGeneralSettings.addServerHistory \"37.230.210.130\" 29900 \"PlayBF2! T~GAMER #1 Allmaps\" 360\r\n",
		},
		{
			name: "serializes config in correct sort order",
			givenConfig: Config{
				content: map[string]Value{
					"LocalProfile.setName":        {content: "\"mister249\""},
					"LocalProfile.setNick":        {content: "\"mister249\""},
					"LocalProfile.setGamespyNick": {content: "\"mister249\""},
				},
			},
			expectedData: "LocalProfile.setGamespyNick \"mister249\"\r\nLocalProfile.setName \"mister249\"\r\nLocalProfile.setNick \"mister249\"\r\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// WHEN
			bytes := tt.givenConfig.ToBytes()

			// THEN
			assert.Equal(t, tt.expectedData, string(bytes))
		})
	}
}

func TestValue_String(t *testing.T) {
	type test struct {
		name           string
		givenValue     Value
		expectedString string
	}

	tests := []test{
		{
			name:           "returns non-quoted string as is",
			givenValue:     Value{content: "some-unquoted-value"},
			expectedString: "some-unquoted-value",
		},
		{
			name:           "returns string containing quotes as is",
			givenValue:     Value{content: "\"some-quoted-sub-value\" some-unquoted-sub-value"},
			expectedString: "\"some-quoted-sub-value\" some-unquoted-sub-value",
		},
		{
			name:           "returns quoted string without quotes",
			givenValue:     Value{content: "\"some-quoted-value\""},
			expectedString: "some-quoted-value",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// WHEN
			asString := tt.givenValue.String()

			// THEN
			assert.Equal(t, tt.expectedString, asString)
		})
	}
}

func TestValue_Slice(t *testing.T) {
	type test struct {
		name          string
		givenValue    Value
		expectedSlice []string
	}

	tests := []test{
		{
			name:          "returns unquoted single value string as is in slice with one element",
			givenValue:    Value{content: "some-unquoted-single-value"},
			expectedSlice: []string{"some-unquoted-single-value"},
		},
		{
			name:          "returns single value string containing quotes as is in slice with one element",
			givenValue:    Value{content: "\"some-quoted-single-sub-value\" some-unquoted-single-sub-value"},
			expectedSlice: []string{"\"some-quoted-single-sub-value\" some-unquoted-single-sub-value"},
		},
		{
			name:          "returns quoted single value string without quotes in slice with one element",
			givenValue:    Value{content: "\"some-quoted-single-value\""},
			expectedSlice: []string{"some-quoted-single-value"},
		},
		{
			name:          "returns unquoted multi value string as is in slice with multiple elements",
			givenValue:    Value{content: "some-unquoted-value;some-other-unquoted-value"},
			expectedSlice: []string{"some-unquoted-value", "some-other-unquoted-value"},
		},
		{
			name:          "returns multi value string containing quotes as is in slice with multiple elements",
			givenValue:    Value{content: "\"some-quoted-sub-value\" some-unquoted-sub-value;\"some-other-quoted-sub-value\" some-other-unquoted-sub-value"},
			expectedSlice: []string{"\"some-quoted-sub-value\" some-unquoted-sub-value", "\"some-other-quoted-sub-value\" some-other-unquoted-sub-value"},
		},
		{
			name:          "returns quoted multi value string without quotes in slice with multiple elements",
			givenValue:    Value{content: "\"some-quoted-value\";\"some-other-quoted-value\""},
			expectedSlice: []string{"some-quoted-value", "some-other-quoted-value"},
		},
		{
			name:          "returns mixed quoted multi value string without quotes and as is in slice with multiple elements",
			givenValue:    Value{content: "\"some-quoted-value\";some-unquoted-value"},
			expectedSlice: []string{"some-quoted-value", "some-unquoted-value"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// WHEN
			asSlice := tt.givenValue.Slice()

			// THEN
			assert.Equal(t, tt.expectedSlice, asSlice)
		})
	}
}
