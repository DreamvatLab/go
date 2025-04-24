package xsecurity

import (
	"crypto/rand"
	"encoding/base64"
	"io"

	"github.com/DreamvatLab/go/xbytes"
	"github.com/DreamvatLab/go/xlog"
	"golang.org/x/crypto/scrypt"
)

func GeneratePasswordSalt(keyLen int) string {
	salt := make([]byte, keyLen)
	_, err := io.ReadFull(rand.Reader, salt)
	if err != nil {
		xlog.Error(err)
		return err.Error()
	}
	saltString := base64.StdEncoding.EncodeToString(salt)
	return saltString
}

func HashPassword(salt, pass string, keyLen int) string {
	saltBytes, err := base64.StdEncoding.DecodeString(salt)
	if err != nil {
		xlog.Error(err)
		return err.Error()
	}

	hash, err := scrypt.Key(xbytes.StrToBytes(pass), saltBytes, 1<<14, 8, 1, keyLen)
	if err != nil {
		xlog.Error(err)
		return err.Error()
	}

	hashString := base64.StdEncoding.EncodeToString(hash)
	return hashString
}
