package end_to_end

import (
	"mse/pkg"
	"mse/proto"
	"testing"

	"google.golang.org/grpc/status"
)

func init() {
	pkg.AuthAddr.Attach()
	pkg.ChatAddr.Attach()
	pkg.CertsPath.Attach()
}

func TestAuthentication(t *testing.T) {
	t.Run("chat_without_authentication_is_denied", chat_without_authentication_is_denied)
	t.Run("expired_token_should_be_denied", expired_token_should_be_denied)
	t.Run("run_monitor_chater_serially", run_monitor_chater_serially)
}

func chat_without_authentication_is_denied(t *testing.T) {
	sut := pkg.NewChatClient(pkg.ChatAddr.Addr(), pkg.CertsPath.Pem())

	err := sut.Say(&proto.SayReq{Msg: "Grettings"})

	if status.Code(err) != status.Code(pkg.MissingToken) {
		t.Fatalf("err not the same want:%v got:%v", pkg.MissingToken, err)
	}
}

func expired_token_should_be_denied(t *testing.T) {
	sut := pkg.NewChatClient(pkg.ChatAddr.Addr(), pkg.CertsPath.Pem())
	sut.UpdateToken("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6Ik1hcnJ5IiwiZXhwIjo4NjQwMH0.h7SvqoYRlXGTh8Qjc-PgZ34iukcveYXMRqGi9eBYec4")

	err := sut.Say(&proto.SayReq{Msg: "Grettings"})

	if status.Code(err) != status.Code(pkg.InvalidToken) {
		t.Fatalf("err not the same want:%v got:%v", pkg.InvalidToken, err)
	}
}

func run_monitor_chater_serially(t *testing.T) {
	monitorStep := func(t *testing.T) chan string {
		monitor := pkg.NewChatClient(pkg.ChatAddr.Addr(), pkg.CertsPath.Pem())

		c := make(chan string)
		var err error
		go func() {
			err = monitor.Listen(c)
		}()
		<-c
		if status.Code(err) != status.Code(pkg.MissingToken) {
			t.Fatalf("before login monitor.Listen return err:%v, want:%v", err, pkg.MissingToken)
		}

		auth := pkg.NewAuthClient(pkg.AuthAddr.Addr(), pkg.CertsPath.Pem())
		loginRsp, err := auth.Login(&proto.LoginReq{Username: "Marry", Password: "Marry"})
		if err != nil {
			t.Fatalf("auth.Login return err:%v, want:%v", err, nil)
		}

		monitor.UpdateToken(loginRsp.Jwt)

		c1 := make(chan string)
		go func() {
			err = monitor.Listen(c1)
			if err != nil {
				t.Fatalf("after login monitor.Listen return err:%v, want:%v", err, nil)
			}
		}()

		return c1

	}

	chaterStep := func(t *testing.T) {
		chater := pkg.NewChatClient(pkg.ChatAddr.Addr(), pkg.CertsPath.Pem())

		err := chater.Say(&proto.SayReq{Msg: "Greeting"})
		if status.Code(err) != status.Code(pkg.MissingToken) {
			t.Fatalf("before login client.Say return err:%v want:%v", err, pkg.MissingToken)
		}

		auth := pkg.NewAuthClient(pkg.AuthAddr.Addr(), pkg.CertsPath.Pem())
		rsp, err := auth.Login(&proto.LoginReq{Username: "John", Password: "John"})
		if err != nil {
			t.Fatalf("auth.Login return err:%v want:%v", err, nil)
		}

		chater.UpdateToken(rsp.Jwt)

		err = chater.Say(&proto.SayReq{Msg: "Greetings"})
		if err != nil {
			t.Fatalf("after login client.Say return err:%v want:%v", err, nil)
		}
	}

	c1 := monitorStep(t)
	chaterStep(t)

	msg := <-c1

	if msg != "John: Greetings" {
		t.Fatalf("monitor receive message:%s want:%s", msg, "Greetings")
	}
}
