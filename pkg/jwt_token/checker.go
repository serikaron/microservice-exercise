package jwt_token

import (
	"fmt"
	"github.com/dgrijalva/jwt-go/v4"
)

type checkFunc func(token *jwt.Token, key *Key) error

func checkAll(token *jwt.Token, key *Key) error {
	checkList := []checkFunc{
		checkKid,
		checkAlg,
	}

	for _, check := range checkList {
		err := check(token, key)
		if err != nil {
			return err
		}
	}
	return nil
}

var checkKid = func(token *jwt.Token, key *Key) error {
	kid, err := getId(token)
	if err != nil {
		return err
	}

	if key.kid != kid {
		return fmt.Errorf("check kid failed, token.kid:%s key.kid:%s", kid, key.kid)
	}

	return nil
}

var checkAlg = func(token *jwt.Token, key *Key) error {
	if token.Method.Alg() != key.method.Alg() {
		return fmt.Errorf("check alg failed, token.alg:%s key.alg:%s", token.Method.Alg(), key.method.Alg())
	}
	return nil
}
