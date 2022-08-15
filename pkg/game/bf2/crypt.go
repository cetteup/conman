//go:build windows

package bf2

import (
	"encoding/hex"
	"strings"
	"unicode"
	"unsafe"

	"golang.org/x/sys/windows"
)

func EncryptProfileConPassword(plain string) (string, error) {
	// Even though the password encrypts and decrypts perfectly fine as is, BF2 needs a NUL character at the end
	enc, err := encrypt([]byte(plain+"\x00"), "This is the description string.")
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(enc), nil
}

func DecryptProfileConPassword(enc string) (string, error) {
	data, err := hex.DecodeString(enc)
	if err != nil {
		return "", err
	}

	dec, _, err := decrypt(data)
	if err != nil {
		return "", err
	}

	clean := strings.Map(func(r rune) rune {
		if unicode.IsPrint(r) {
			return r
		}
		return -1
	}, string(dec))

	return clean, nil
}

func newBlob(data []byte) *windows.DataBlob {
	if len(data) == 0 {
		return &windows.DataBlob{}
	}

	return &windows.DataBlob{
		Size: uint32(len(data)),
		Data: &data[0],
	}
}

func blobToByteArray(blob windows.DataBlob) []byte {
	bytes := make([]byte, blob.Size)
	copy(bytes, unsafe.Slice(blob.Data, blob.Size))
	return bytes
}

func encrypt(data []byte, description string) ([]byte, error) {
	dataIn := newBlob(data)
	var dataOut windows.DataBlob
	name, err := windows.UTF16PtrFromString(description)
	if err != nil {
		return nil, err
	}

	if err = windows.CryptProtectData(dataIn, name, nil, uintptr(0), nil, windows.CRYPTPROTECT_UI_FORBIDDEN, &dataOut); err != nil {
		return nil, err
	}

	defer func() {
		_, _ = windows.LocalFree(windows.Handle(unsafe.Pointer(dataOut.Data)))
	}()

	return blobToByteArray(dataOut), nil
}

func decrypt(data []byte) ([]byte, string, error) {
	dataIn := newBlob(data)
	var dataOut windows.DataBlob
	name, err := windows.UTF16PtrFromString("")
	if err != nil {
		return nil, "", err
	}

	if err = windows.CryptUnprotectData(dataIn, &name, nil, uintptr(0), nil, windows.CRYPTPROTECT_UI_FORBIDDEN, &dataOut); err != nil {
		return nil, "", err
	}

	defer func() {
		_, _ = windows.LocalFree(windows.Handle(unsafe.Pointer(name)))
		_, _ = windows.LocalFree(windows.Handle(unsafe.Pointer(dataOut.Data)))
	}()

	return blobToByteArray(dataOut), windows.UTF16PtrToString(name), nil
}
