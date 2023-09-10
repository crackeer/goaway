package util

import (
	"bytes"
	"crypto/cipher"
	"crypto/des"
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"
)

const salt = "AKJsahs_&xshja"

// PlainObject ...
type PlainObject struct {
	AppID    string
	ExpireAt int64
}

// GenerateToken
//
//	@param ak
//	@param sk
//	@return string
//	@return int64
//	@return error
func GenerateToken(sk string, object *PlainObject) (string, error) {
	encText, err := desEncrypt([]byte(fmt.Sprintf("%d", object.ExpireAt)), []byte(desKey(sk)))
	if err != nil {
		return "", err
	}

	encodeText := strings.Join([]string{object.AppID, string(encText)}, ":")

	return base64.StdEncoding.EncodeToString([]byte(encodeText)), nil
}

// ParseToken
//
//	@param token
//	@param funcSK
//	@return int64
//	@return error
func ParseToken(token string, funcSK func(string) string) (*PlainObject, error) {
	bytes, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		return nil, fmt.Errorf("base64 decode error: %v", err)
	}
	parts := strings.Split(string(bytes), ":")
	if len(parts) != 2 {
		return nil, fmt.Errorf("token not the right format")
	}

	sk := funcSK(parts[0])
	if len(sk) < 1 {
		return nil, fmt.Errorf("app not found")
	}

	result, err := desDecrypt([]byte(parts[1]), []byte(desKey(sk)))
	if err != nil {
		return nil, fmt.Errorf("decrypt failed: %v", err)
	}

	expire, _ := strconv.Atoi(string(result))
	return &PlainObject{
		AppID:    parts[0],
		ExpireAt: int64(expire),
	}, nil
}

func desKey(sk string) string {
	return strings.Join([]string{salt, sk}, ":")
}

func desDecrypt(crypted, key []byte) ([]byte, error) {
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockMode := cipher.NewCBCDecrypter(block, key)
	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)
	origData = pkcs5UnPadding(origData)
	return origData, nil
}

func desEncrypt(origData, key []byte) ([]byte, error) {
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}
	origData = pkcs5Padding(origData, block.BlockSize())
	blockMode := cipher.NewCBCEncrypter(block, key)
	crypted := make([]byte, len(origData))
	blockMode.CryptBlocks(crypted, origData)
	return crypted, nil
}

func pkcs5UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

func pkcs5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}
