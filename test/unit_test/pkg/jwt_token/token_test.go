package jwt_token

import (
	"github.com/dgrijalva/jwt-go/v4"
	"mse/pkg"
	"mse/pkg/jwt_token"
	"reflect"
	"testing"
	"time"
)

func TestJwtToken(t *testing.T) {
	should_work_with_hmac_method(t)
	parse_fail_with_different_key(t)
	parse_fail_with_expire_token(t)
	parse_must_fail_with_none_sign_string(t)
}

func should_work_with_hmac_method(t *testing.T) {
	test := func(t *testing.T, key *jwt_token.Key) {
		id := pkg.Identity{Name: "Marry"}

		tokenString, genErr := jwt_token.Gen(id, 86400, key)
		got, parseErr := jwt_token.Parse(tokenString, func(string) *jwt_token.Key {
			return key
		})

		if genErr != nil {
			t.Fatal(genErr)
		}
		if parseErr != nil {
			t.Fatal(parseErr)
		}
		if !reflect.DeepEqual(id, got) {
			t.Fatalf("gen token failed want:%v got:%v", id, got)
		}
	}

	tokenList := map[string]*jwt_token.Key{
		"HS256": jwt_token.NewHS256Key("1", signKey),
		"HS384": jwt_token.NewHS384Key("1", signKey),
		"HS512": jwt_token.NewHS512Key("1", signKey),
	}
	for name, key := range tokenList {
		t.Run(name, func(t *testing.T) {
			test(t, key)
		})
	}
}

func parse_fail_with_different_key(t *testing.T) {
	test := func(t *testing.T, genKey *jwt_token.Key, parseKey *jwt_token.Key) {
		id := pkg.Identity{Name: "Marry"}

		tokenString, genErr := jwt_token.Gen(id, 86400, genKey)
		_, parseErr := jwt_token.Parse(tokenString, func(string) *jwt_token.Key { return parseKey })

		if genErr != nil {
			t.Fatal(genErr)
		}
		if parseErr == nil {
			t.Fatal("parse with an invalid key must fail but passed")
		}
		t.Log(parseErr)
	}

	datas := map[string]struct {
		genKey   *jwt_token.Key
		parseKey *jwt_token.Key
	}{
		"parse_fail_with_different_key(id)": {
			genKey:   jwt_token.NewHS256Key("1", signKey),
			parseKey: jwt_token.NewHS256Key("2", signKey),
		},

		"parse_fail_with_different_key(content)": {
			genKey:   jwt_token.NewHS256Key("1", signKey),
			parseKey: jwt_token.NewHS256Key("1", "error key"),
		},
	}
	for name, data := range datas {
		t.Run(name, func(t *testing.T) {
			test(t, data.genKey, data.parseKey)
		})
	}
}

func parse_fail_with_expire_token(t *testing.T) {
	t.Run("parse_fail_with_expire_token", func(t *testing.T) {
		id := pkg.Identity{Name: "Marry"}
		key := jwt_token.NewHS256Key("1", signKey)

		jwt.TimeFunc = func() time.Time {
			return time.Unix(0, 0)
		}
		tokenString, genErr := jwt_token.Gen(id, 1000, key)
		jwt.TimeFunc = func() time.Time {
			return time.Unix(2000, 0)
		}
		_, parseErr := jwt_token.Parse(tokenString, func(string) *jwt_token.Key { return key })

		if genErr != nil {
			t.Fatal(genErr)
		}
		if parseErr == nil {
			t.Fatal("parse with an expired token must fail but passed")
		}
		t.Log(parseErr)
	})
}

func parse_must_fail_with_none_sign_string(t *testing.T) {
	t.Run("parse_must_fail_with_none_sign_string", func(t *testing.T) {
		tokenString := "eyJhbGciOiJub25lIiwia2lkIjoiMSIsInR5cCI6IkpXVCJ9.eyJleHAiOjAuMDAwMDg2LCJuYW1lIjoiTWFycnkifQ."
		key := jwt_token.NewHS256Key("1", signKey)
		_, err := jwt_token.Parse(tokenString, func(string) *jwt_token.Key { return key })
		if err == nil {
			t.Fatal("parse with an non-signed token must fail but passed")
		}
	})
}

const signKey string = "4c04990f-0654-4121-ab49-7b047c850507"
