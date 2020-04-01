package jwt_token

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go/v4"
	"log"
	"mse/pkg"
	"reflect"
	"time"
)

func getId(token *jwt.Token) (string, error) {
	i, ok := token.Header["kid"]
	if !ok {
		return "", errors.New("kid not found")
	}

	kid, ok := i.(string)
	if !ok {
		return "", fmt.Errorf("kid not string, i.(type):%v", reflect.TypeOf(i))
	}

	return kid, nil
}

type keyFunc func(kid string) *Key

type customClaims struct {
	pkg.Identity
	jwt.StandardClaims
}

func Gen(id pkg.Identity, expireInSecond uint64, key *Key) (string, error) {
	customClaim := customClaims{
		Identity: id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: jwt.At(jwt.Now().Add(time.Duration(expireInSecond) * time.Second)),
		},
	}

	var tk *jwt.Token
	var signKey []byte
	var err error
	var tokenString string

	func() {
		tk = jwt.NewWithClaims(key.method, customClaim)
		tk.Header["kid"] = key.kid

		j, _ := json.Marshal(tk)
		log.Println(string(j))

		signKey, err = key.Read()
		if err != nil {
			return
		}
		tokenString, err = tk.SignedString(signKey)
		if err != nil {
			return
		}
	}()

	log.Printf("JwtToken.Gen id:%v expireIdSecond:%d kid:%s alg:%s signKey.len:%d tokenString:%s err:%v",
		id, expireInSecond,
		key.kid, key.method.Alg(),
		len(signKey), tokenString, err)
	return tokenString, err
}

func Parse(tokenString string, f keyFunc) (pkg.Identity, error) {
	var err error
	var token *jwt.Token
	var id pkg.Identity

	func() {
		token, err = jwt.ParseWithClaims(tokenString, &customClaims{}, func(token *jwt.Token) (interface{}, error) {
			kid := ""
			var key *Key = nil
			var buf []byte = nil
			err := errors.New("jwt.keyFunc error")

			func() {
				kid, err = getId(token)

				key = f(kid)
				if key == nil {
					return
				}

				err = checkAll(token, key)
				if err != nil {
					return
				}

				buf, err = key.Read()
			}()

			log.Printf("jwt_token.Parse jwt.keyFunc kid:%s err:%v",
				kid, err)
			return buf, err
		})

		if err != nil {
			return
		}

		if !token.Valid {
			err = errors.New("token invalid")
			return
		}

		var ok = false
		var claims *customClaims
		claims, ok = token.Claims.(*customClaims)
		if !ok {
			err = errors.New("cast customClaims failed")
		}

		id = claims.Identity
	}()

	log.Printf("JwtToken.Parse tokenString:%s, token.Method:%s token.Valid:%t, id:%v, error:%v",
		tokenString,
		token.Method.Alg(),
		token.Valid,
		id, err)
	return id, err
}
