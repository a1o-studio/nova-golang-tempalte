package crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"io"
)

// EncryptAES 使用 AES 加密明文，返回 Base64 编码的密文
func EncryptAES(plainText, base64Key string) (string, error) {
	key, err := base64.StdEncoding.DecodeString(base64Key)
	if err != nil {
		return "", err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	plainData := pkcs7Pad([]byte(plainText), aes.BlockSize)
	ciphertext := make([]byte, aes.BlockSize+len(plainData))
	iv := ciphertext[:aes.BlockSize]

	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext[aes.BlockSize:], plainData)

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// DecryptAES 使用 AES 解密 Base64 编码的密文，返回明文
func DecryptAES(cipherTextBase64, base64Key string) (string, error) {
	key, err := base64.StdEncoding.DecodeString(base64Key)
	if err != nil {
		return "", err
	}

	cipherData, err := base64.StdEncoding.DecodeString(cipherTextBase64)
	if err != nil {
		return "", err
	}

	if len(cipherData) < aes.BlockSize {
		return "", errors.New("ciphertext too short")
	}
	iv := cipherData[:aes.BlockSize]
	cipherData = cipherData[aes.BlockSize:]

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	if len(cipherData)%aes.BlockSize != 0 {
		return "", errors.New("ciphertext is not a multiple of the block size")
	}

	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(cipherData, cipherData)

	plainData, err := pkcs7Unpad(cipherData)
	if err != nil {
		return "", err
	}

	return string(plainData), nil
}

// padding
func pkcs7Pad(data []byte, blockSize int) []byte {
	padLen := blockSize - len(data)%blockSize
	pad := bytes.Repeat([]byte{byte(padLen)}, padLen)
	return append(data, pad...)
}

func pkcs7Unpad(data []byte) ([]byte, error) {
	length := len(data)
	if length == 0 {
		return nil, errors.New("invalid padding size")
	}
	padLen := int(data[length-1])
	if padLen == 0 || padLen > length {
		return nil, errors.New("invalid padding size")
	}
	// 检查所有 padding 字节
	for i := range padLen {
		if data[length-1-i] != byte(padLen) {
			return nil, errors.New("invalid padding content")
		}
	}
	return data[:length-padLen], nil
}

func Sum256(data []byte) string {
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])
}
