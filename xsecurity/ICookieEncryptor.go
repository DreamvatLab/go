package xsecurity

import (
	"github.com/DreamvatLab/go/xerr"
	"github.com/gorilla/securecookie"
)

type ICookieEncryptor interface {
	Encrypt(name string, value interface{}) (string, error)
	Decrypt(name, value string, dst interface{}) error
}

func NewSecureCookieEncryptor(hashKey, blockKey []byte) ICookieEncryptor {
	s := securecookie.New(hashKey, blockKey)
	return &SecureCookieEncryptor{
		s: s,
	}
}

type SecureCookieEncryptor struct {
	s *securecookie.SecureCookie
}

func (x *SecureCookieEncryptor) Encrypt(name string, value interface{}) (string, error) {
	a, b := x.s.Encode(name, value)
	return a, xerr.WithStack(b)
}

func (x *SecureCookieEncryptor) Decrypt(name, value string, dst interface{}) error {
	err := x.s.Decode(name, value, dst)
	return xerr.WithStack(err)
}
