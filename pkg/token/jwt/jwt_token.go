package jwt

import (
	"errors"
	"fmt"
	jwt_go "github.com/dgrijalva/jwt-go/v4"
	"io"
	"log"
	"mse/pkg"
	"time"
)

type CustomClaims struct {
	pkg.Identity
	jwt_go.StandardClaims
}

type JwtToken struct {
	method jwt_go.SigningMethod
	buf    []byte
}

func NewHS256Token() *JwtToken {
	return &JwtToken{
		method: jwt_go.SigningMethodHS256,
		buf:    make([]byte, 1024),
	}
}

func (jt *JwtToken) Gen(id pkg.Identity, expireInSecond uint32, signKeyReader io.Reader) (string, error) {
	customClaim := CustomClaims{
		Identity: id,
		StandardClaims: jwt_go.StandardClaims{
			ExpiresAt: jwt_go.At(time.Now().Add(time.Duration(expireInSecond))),
		},
	}

	var token *jwt_go.Token
	var signKey []byte
	var err error
	var tokenString string

	func() {
		token = jwt_go.NewWithClaims(jt.method, customClaim)
		signKey, err = jt.readKey(signKeyReader)
		if err != nil {
			return
		}
		tokenString, err = token.SignedString(signKey)
		if err != nil {
			return
		}
	}()

	log.Printf("JwtToken.Gen id:%v expireIdSecond:%d signKey.len:%d tokenString:%s err:%v",
		id, expireInSecond, len(signKey), tokenString, err)
	return tokenString, err
}

func (jt *JwtToken) Parse(tokenString string, signKeyReader io.Reader) (pkg.Identity, error) {
	var err error
	var token *jwt_go.Token
	var signKey []byte
	var claims *CustomClaims

	func() {
		token, err = jwt_go.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt_go.Token) (interface{}, error) {
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

	log.Printf("JwtToken.Parse tokenString:%s, token.Method:%s jt.method:%s signKey.len:%d token.Valid:%t, id:%s, error:%v",
		tokenString,
		token.Method.Alg(), jt.method.Alg(),
		len(signKey),
		token.Valid,
		claims.Identity, err)
	return claims.Identity, err
}

func (jw *JwtToken) readKey(reader io.Reader) (signKey []byte, err error) {
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
