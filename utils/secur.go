package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/json"
	"os"
	"strings"

	"github.com/ASV-Aachen/Seereisenplan-backend/modules/gocloak"
)

var secret string = os.Getenv("secret")
var Bes = []byte{35, 46, 57, 24, 85, 35, 24, 74, 87, 35, 88, 98, 66, 32, 14, 05}

func EncodeUser(data gocloak.User) (string, error) {
	var buf bytes.Buffer
	encoder := base64.NewEncoder(base64.StdEncoding, &buf)
	err := json.NewEncoder(encoder).Encode(data)
	if err != nil {
		return "", err
	}
	encoder.Close()
	return buf.String(), nil
}

func DecodeUser(data string) (gocloak.User, error) {
	var resultJson gocloak.User
	error := json.NewDecoder(base64.NewDecoder(base64.StdEncoding, strings.NewReader(data))).Decode(&resultJson)

	return resultJson, error
}

func Encode_ByteToString(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

func Decode_stringToByte(s string) []byte {
	data, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		panic(err)
	}
	return data
}

func Encrypt(data string) (string, error) {
	block, err := aes.NewCipher([]byte(secret))
	if err != nil {
		return "", err
	}
	plainText := []byte(data)
	cfb := cipher.NewCFBEncrypter(block, Bes)
	cipherText := make([]byte, len(plainText))
	cfb.XORKeyStream(cipherText, plainText)

	return Encode_ByteToString(cipherText), nil
}

func Decrypt(data string) (string, error) {
	block, err := aes.NewCipher([]byte(secret))
	if err != nil {
		return "", err
	}
	cipherText := Decode_stringToByte(data)
	cfb := cipher.NewCFBDecrypter(block, Bes)
	plainText := make([]byte, len(cipherText))
	cfb.XORKeyStream(plainText, cipherText)
	return string(plainText), nil
}
