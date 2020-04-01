package jwt_token

import (
	"errors"
	"github.com/dgrijalva/jwt-go/v4"
	"io"
)

var (
	ReadKeyErr = errors.New("read key error")
)

type Key struct {
	kid    string
	reader io.Reader
	method jwt.SigningMethod
	buf    []byte
}

func newHSKey(kid string, key string, method jwt.SigningMethod) *Key {
	return &Key{
		kid:    kid,
		reader: nil,
		method: method,
		buf:    []byte(key),
	}
}

func NewHS256Key(kid string, key string) *Key {
	return newHSKey(kid, key, jwt.SigningMethodHS256)
}

func NewHS384Key(kid string, key string) *Key {
	return newHSKey(kid, key, jwt.SigningMethodHS384)
}

func NewHS512Key(kid string, key string) *Key {
	return newHSKey(kid, key, jwt.SigningMethodHS512)
}

func (k *Key) Read() ([]byte, error) {
	if len(k.buf) > 0 {
		return k.buf, nil
	}
	if k.reader == nil {
		return nil, ReadKeyErr
	}

	if k.buf == nil {
		k.buf = make([]byte, 0)
	}
	p := make([]byte, 1024)
	for {
		n, err := k.reader.Read(p)
		if err != nil && err != io.EOF {
			return nil, err
		}
		if n > 0 {
			k.buf = append(k.buf, p[:n]...)
		}
		if err == io.EOF {
			break
		}
	}

	return k.buf, nil
}
