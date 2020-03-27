package pkg

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go/v4"
	"io"
	"log"
	"time"
)

type CustomClaims struct {
	Identity
	jwt.StandardClaims
}

type JwtToken interface {
	Gen(identity Identity, expireInSecond uint32, signKeyReader io.Reader) (string, error)
	Parse(tokenString string, signKeyReader io.Reader) (Identity, error)
}

func NewJwtToken() JwtToken {
	return &jwtToken{
		method: jwt.SigningMethodHS256,
		buf:    make([]byte, 1024),
	}
}

type jwtToken struct {
	method jwt.SigningMethod
	buf    []byte
}

func (jt *jwtToken) Gen(id Identity, expireInSecond uint32, signKeyReader io.Reader) (string, error) {
	customClaim := CustomClaims{
		Identity: id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: jwt.At(time.Now().Add(time.Duration(expireInSecond))),
		},
	}

	var token *jwt.Token
	var signKey []byte
	var err error
	var tokenString string

	func() {
		token = jwt.NewWithClaims(jt.method, customClaim)
		signKey, err = jt.readKey(signKeyReader)
		if err != nil {
			return
		}
		tokenString, err = token.SignedString(signKey)
		if err != nil {
			return
		}
	}()

	log.Printf("jwtToken.Gen id:%v expireIdSecond:%d signKey.len:%d tokenString:%s err:%v",
		id, expireInSecond, len(signKey), tokenString, err)
	return tokenString, err
}

func (jt *jwtToken) Parse(tokenString string, signKeyReader io.Reader) (Identity, error) {
	var err error
	var token *jwt.Token
	var signKey []byte
	var claims *CustomClaims

	func() {
		token, err = jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
			if token.Method != jt.method {
				return nil, fmt.Errorf("method invalid")
			}
			signKey, err = jt.readKey(signKeyReader)
			if err != nil {
				return nil, err
			}
			return signKey, nil
		})

		if err != nil {
			return
		}

		if !token.Valid {
			err = errors.New("token invalid")
			return
		}

		var ok = false
		claims, ok = token.Claims.(*CustomClaims)
		if !ok {
			err = errors.New("cast CustomClaims failed")
		}
	}()

	log.Printf("jwtToken.Parse tokenString:%s, token.Method:%s jt.method:%s signKey.len:%d token.Valid:%t, id:%s, error:%v",
		tokenString,
		token.Method.Alg(), jt.method.Alg(),
		len(signKey),
		token.Valid,
		claims.Identity, err)
	return claims.Identity, err
}

func (jw *jwtToken) readKey(reader io.Reader) (signKey []byte, err error) {
	var n = 0
	for err == nil {
		n, err = io.ReadFull(reader, jw.buf)
		if err == nil && n > 0 {
			signKey = append(signKey, jw.buf...)
		} else if err == nil && n == 0 {
			break
		}
	}
	if err == io.ErrUnexpectedEOF {
		if n > 0 {
			signKey = append(signKey, jw.buf[:n]...)
		}
		err = nil
	} else if n == 0 && err == io.EOF {
		err = nil
	}
	return
}
