//go:build windows && unit

package bf2

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestEncryptProfileConPassword Since the dpapi.h ensures that encrypted data can only be read on the same machine by the same user,
// we can only test encryption round trips (any hard-coded encrypted string would fail to decrypt on any other system)
func TestEncryptProfileConPassword(t *testing.T) {
	t.Run("successfully encrypts password", func(t *testing.T) {
		// GIVEN
		plain := "som3-p@assw0rd"

		// WHEN
		encrypted, err := EncryptProfileConPassword(plain)

		// THEN
		require.NoError(t, err)
		// Encrypted string should have an even length
		assert.True(t, len(encrypted)%2 == 0)
		decrypted, err := DecryptProfileConPassword(encrypted)
		require.NoError(t, err)
		assert.Equal(t, plain, decrypted)
	})
}

func TestDecryptProfileConPassword(t *testing.T) {
	t.Run("successfully decrypts password", func(t *testing.T) {
		// GIVEN
		plain := "som3-p@assw0rd"
		encrypted, err := EncryptProfileConPassword(plain)
		require.NoError(t, err)

		// WHEN
		decrypted, err := DecryptProfileConPassword(encrypted)

		// THEN
		require.NoError(t, err)
		assert.Equal(t, plain, decrypted)
	})

	t.Run("error decrypting password from another machine", func(t *testing.T) {
		// GIVEN
		// "som3-p@assw0rd" encrypted on a throwaway VM
		encrypted := "01000000d08c9ddf0115d1118c7a00c04fc297eb01000000a3ab32c19e0b704d83adb2ae27d8f7b100000000020000000000106600000001000020000000736460651ebfb936c41250dff9fb10961856a509908567118679adcda3843e24000000000e800000000200002000000011bd9dd1b15cc690e27ac0688f2f9eeea0d451e691768ca08ec06f700caa9415100000001dd8a15adc4b79fa2e93c14863fe404b40000000404b685423f21bce7db488fb6bda5ec9772379996ea4749019185e7b8e18fac16054e877c719da94782537881d1aec5ba560763f0d3bb9bad03b2bca2a6d7b12"

		// WHEN
		_, err := DecryptProfileConPassword(encrypted)

		// THEN
		require.ErrorContains(t, err, "Key not valid for use in specified state")
	})

	t.Run("error decrypting odd-length encrypted string", func(t *testing.T) {
		// GIVEN
		encrypted := "0af"

		// WHEN
		_, err := DecryptProfileConPassword(encrypted)

		// THEN
		require.ErrorContains(t, err, "odd length hex string")
	})

	t.Run("error decrypting non-hex encrypted string", func(t *testing.T) {
		// GIVEN
		encrypted := "not-a-hex-string"

		// WHEN
		_, err := DecryptProfileConPassword(encrypted)

		// THEN
		require.ErrorContains(t, err, "invalid byte")
	})
}
