package internal_test

import (
	"mse/auth/internal"
	"mse/pkg"
	"testing"

	"google.golang.org/grpc/status"
)

func Test_success_with_correct_password(t *testing.T) {
	t.Parallel()

	id, err := internal.IdentifyWithPassword("Marry", "Marry")

	if err != nil {
		t.Fatal(err)
	}

	if id.Name != "Marry" {
		t.Fatalf("id.Name is %s, wants %s", id.Name, "Marry")
	}
}

func Test_login_failed_with_incorrect_password(t *testing.T) {
	t.Parallel()

	_, err := internal.IdentifyWithPassword("Marry", "InvalidPassword")
	if status.Code(err) != status.Code(pkg.LoginErr) {
		t.Fatalf("errcode not match want:%v got:%v", pkg.LoginErr, err)
	}
}
