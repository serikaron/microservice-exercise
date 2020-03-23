package main

import (
	"auth/pkg"
	"auth/proto"
	auth_testing "auth/testing"
	"flag"
	"fmt"
)

func testLogin(ac *pkg.AuthClient) {
	tokenString := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6Ik1hcnJ5IiwiZXhwIjo4NjQwMH0.h7SvqoYRlXGTh8Qjc-PgZ34iukcveYXMRqGi9eBYec4"

	tests := []auth_testing.TestInfo{
		{
			Name:    "expect success",
			Req:     &proto.LoginReq{Username: "Marry", Password: "Marry"},
			Want:    &proto.LoginRsp{Jwt: tokenString},
			WantErr: nil,
		},
		{
			Name:    "login failed",
			Req:     &proto.LoginReq{Username: "Marry", Password: "error password"},
			Want:    &proto.LoginRsp{Jwt: tokenString},
			WantErr: pkg.LoginErr,
		},
		{
			Name:    "token failed",
			Req:     &proto.LoginReq{Username: "Marry", Password: "Marry"},
			Want:    &proto.LoginRsp{Jwt: ""},
			WantErr: auth_testing.RspErr,
		},
	}

	builder := func(ti *auth_testing.TestInfo) auth_testing.Test {
		method := func(req interface{}) (interface{}, error) {
			return ac.Login(req.(*proto.LoginReq))
		}
		mt := auth_testing.MethodTest{
			Info:   ti,
			Method: method,
		}
		return &auth_testing.SimpleWrapper{&mt}
	}

	testCase := auth_testing.TestCase{
		Infos:   tests,
		Builder: builder,
	}
	_ = testCase.Run()
}

func testGetSignKey(client *pkg.AuthInternalClient) {
	tests := []auth_testing.TestInfo{
		{
			Name: "simple test",
			Req:  &proto.GetSignKeyReq{},
			Want: &proto.GetSignKeyRsp{
				Kid: 1,
				Key: "secret-key",
				Alg: "HS256",
			},
			WantErr: nil,
		},
	}

	tc := auth_testing.TestCase{
		Infos: tests,
		Builder: func(info *auth_testing.TestInfo) auth_testing.Test {
			return &auth_testing.SimpleWrapper{Core: &auth_testing.MethodTest{
				Info: info,
				Method: func(req interface{}) (interface{}, error) {
					return client.GetSignKey(req.(*proto.GetSignKeyReq))
				},
			}}
		},
	}
	_ = tc.Run()
}

func main() {
	host := flag.String("auth-service-host", "auth-service", "auth service host")
	port := flag.Uint("auth-service-port", 0, "auth service port")
	internalHost := flag.String("auth-internal-service-host", "auth-service", "auth internal service host")
	internalPort := flag.Uint("auth-internal-service-port", 0, "auth internal service port")
	flag.Parse()
	addr := fmt.Sprintf("%s:%d", *host, *port)
	client := pkg.NewAuthClient(addr)
	testLogin(client)

	internalAddr := fmt.Sprintf("%s:%d", *internalHost, *internalPort)
	internalClient := pkg.NewAuthInternalClient(internalAddr)
	testGetSignKey(internalClient)
}
