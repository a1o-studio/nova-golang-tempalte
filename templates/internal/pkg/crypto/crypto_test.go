package crypto

import (
	"encoding/base64"
	"testing"

	"github.com/a1ostudio/nova/internal/pkg/util"

	"github.com/stretchr/testify/require"
)

func TestEncryptDecryptAES(t *testing.T) {
	plainText := util.RandomString(20)
	key := make([]byte, 32) // 256-bit key
	for i := range key {
		key[i] = byte(i)
	}
	base64Key := base64.StdEncoding.EncodeToString(key)

	cipherText, err := EncryptAES(plainText, base64Key)
	require.NoError(t, err)
	require.NotEmpty(t, cipherText)

	decrypted, err := DecryptAES(cipherText, base64Key)
	require.NoError(t, err)
	require.Equal(t, plainText, decrypted)
}

func TestEncryptAESInvalidKey(t *testing.T) {
	plainText := util.RandomString(20)
	invalidKey := base64.StdEncoding.EncodeToString([]byte(util.RandomString(8)))
	_, err := EncryptAES(plainText, invalidKey)
	require.Error(t, err)
}

func TestDecryptAESInvalidKey(t *testing.T) {
	key := make([]byte, 32)
	base64Key := base64.StdEncoding.EncodeToString(key)
	invalidKey := base64.StdEncoding.EncodeToString([]byte(util.RandomString(8)))
	cipherText, _ := EncryptAES(util.RandomString(4), base64Key)
	_, err := DecryptAES(cipherText, invalidKey)
	require.Error(t, err)
}

func TestDecryptAESInvalidCipherText(t *testing.T) {
	key := make([]byte, 32)
	base64Key := base64.StdEncoding.EncodeToString(key)
	_, err := DecryptAES(util.RandomString(20), base64Key)
	require.Error(t, err)
}

func TestPKCS7PadUnpad(t *testing.T) {
	blockSize := 16
	data := []byte(util.RandomString(10))
	padded := pkcs7Pad(data, blockSize)
	require.Equal(t, 16, len(padded))
	unpadded, err := pkcs7Unpad(padded)
	require.NoError(t, err)
	require.Equal(t, data, unpadded)
}

func TestPKCS7UnpadInvalid(t *testing.T) {
	// Padding too large
	data := []byte(util.RandomString(10) + string([]byte{20, 20, 20, 20}))
	_, err := pkcs7Unpad(data)
	require.Error(t, err)

	// Padding zero
	data = []byte(util.RandomString(10) + string([]byte{0, 0, 0, 0}))
	_, err = pkcs7Unpad(data)
	require.Error(t, err)
}
