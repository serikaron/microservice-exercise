package main

import (
	"context"
	"mse/auth/pkg"
	"mse/auth/proto"
	auth_testing "mse/auth/testing"
	"testing"
)

func TestAuthServer_LoginNew(t *testing.T) {
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
			as := &AuthService{}
			return as.Login(context.Background(), req.(*proto.LoginReq))
		}
		mt := auth_testing.MethodTest{ti, method}

		wrapper := auth_testing.GoTestWrapper{t, &mt, ti}
		return &wrapper
	}

	testCase := auth_testing.TestCase{
		Infos:   tests,
		Builder: builder,
	}
	_ = testCase.Run()
}

//func TestAuthServer_Login(t *testing.T) {
//	type args struct {
//		in0 context.Context
//		in  *proto.LoginReq
//	}
//
//	jwtMap := map[string]string{
//		"Marry": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6Ik1hcnJ5IiwiZXhwIjo4NjQwMH0.h7SvqoYRlXGTh8Qjc-PgZ34iukcveYXMRqGi9eBYec4",
//	}
//
//	rspErr := errors.New("check rsp failed")
//
//	tests := []struct {
//		name    string
//		args    args
//		want    *proto.LoginRsp
//		wantErr error
//	}{
//		{name: "expect success", args: args{
//			in0: context.Background(),
//			in:  &proto.LoginReq{Username: "Marry", Password: "Marry"},
//		}, want: &proto.LoginRsp{Jwt: jwtMap["Marry"]}, wantErr: nil},
//		{name: "login failed", args: args{
//			in0: context.Background(),
//			in:  &proto.LoginReq{Username: "Marry", Password: "error password"},
//		}, want: &proto.LoginRsp{Jwt: jwtMap["Marry"]}, wantErr: pkg.LoginErr},
//		{name: "token failed", args: args{
//			in0: context.Background(),
//			in:  &proto.LoginReq{Username: "Marry", Password: "Marry"},
//		}, want: &proto.LoginRsp{Jwt: ""}, wantErr: rspErr},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			t.Helper()
//			as := &AuthService{}
//			got, err := as.Login(tt.args.in0, tt.args.in)
//			if err == nil &&
//				!reflect.DeepEqual(got, tt.want) {
//				t.Logf("Login() got = %v, want %v", got, tt.want)
//				err = rspErr
//			}
//			if err != tt.wantErr {
//				t.Errorf("Login() error = %v, wantErr %v", err, tt.wantErr)
//				return
//			}
//		})
//	}
//}
