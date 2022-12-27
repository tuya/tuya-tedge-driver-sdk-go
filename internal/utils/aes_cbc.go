package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"errors"
)

func AESCBCEncrypt(content, key, iv []byte) ([]byte, error) {
	var crypted []byte
	if len(content) <= 0 {
		return crypted, errors.New("plain content empty")
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return crypted, err
	}
	ecb := cipher.NewCBCEncrypter(block, iv)
	content = PKCS5Padding(content, block.BlockSize())
	crypted = make([]byte, len(content))
	ecb.CryptBlocks(crypted, content)
	return crypted, nil
}

func AESCBCDecrypt(crypted, key, iv []byte) ([]byte, error) {
	var decrypted []byte
	if len(crypted) == 0 {
		return decrypted, errors.New("plain content empty")
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return decrypted, err
	}
	ecb := cipher.NewCBCDecrypter(block, iv)
	decrypted = make([]byte, len(crypted))
	ecb.CryptBlocks(decrypted, crypted)
	return PKCS5Trimming(decrypted), nil
}

func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS5Trimming(encrypt []byte) []byte {
	padding := encrypt[len(encrypt)-1]
	return encrypt[:len(encrypt)-int(padding)]
}
