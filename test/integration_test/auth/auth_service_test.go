package auth

import (
	"mse/pkg"
	"mse/proto"
	"testing"

	"google.golang.org/grpc/status"
)

func init() {
	pkg.AuthAddr.Attach()
	pkg.CertsPath.Attach()
}

func TestAuthServer_Login(t *testing.T) {
	login_success_with_correct_password(t)
	login_failed_with_incorrect_password(t)
}

func login_success_with_correct_password(t *testing.T) {
	t.Run("login_success_with_correct_password", func(t *testing.T) {
		req := &proto.LoginReq{Username: "Marry", Password: "Marry"}
		sut := pkg.NewAuthClient(pkg.AuthAddr.Addr(), pkg.CertsPath.Pem())

		rsp, err := sut.Login(req)

		if err != nil {
			t.Fatal(err)
		}

		if len(rsp.Jwt) == 0 {
			t.Fatal("invalid jwt")
		}
	})
}

func login_failed_with_incorrect_password(t *testing.T) {
	t.Run("login_failed_with_incorrect_password", func(t *testing.T) {
		req := &proto.LoginReq{Username: "Marry", Password: "IncorrectPassword"}
		sut := pkg.NewAuthClient(pkg.AuthAddr.Addr(), pkg.CertsPath.Pem())

		_, err := sut.Login(req)
		if status.Code(err) != status.Code(pkg.LoginErr) {
			t.Fatalf("errcode not match want:%v got:%v", pkg.LoginErr, err)
		}
	})
}
