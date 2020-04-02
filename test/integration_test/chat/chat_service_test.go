package chat

import (
	"mse/pkg"
	"mse/proto"
	"testing"

	"google.golang.org/grpc/status"
)

func init() {
	pkg.ChatAddr.Attach()
	pkg.CertsPath.Attach()
	pkg.IntegrationKey.Attach()
}

func TestChatService(t *testing.T) {
	t.Run("monitor_can_know_what_sender_say", monitor_can_know_what_sender_say)
	t.Run("reject_if_token_invalid", reject_if_token_invalid)
}

func monitor_can_know_what_sender_say(t *testing.T) {
	monitor := pkg.NewChatClient(pkg.ChatAddr.Addr(), pkg.CertsPath.Pem())
	monitor.UpdateToken(pkg.IntegrationKey.Val)
	listenDone := make(chan bool)
	defer close(listenDone)
	c, err := monitor.Listen(listenDone)
	if err != nil {
		t.Fatal(err)
	}
	sender := pkg.NewChatClient(pkg.ChatAddr.Addr(), pkg.CertsPath.Pem())
	sender.UpdateToken(pkg.IntegrationKey.Val)

	err = sender.Say(&proto.SayReq{Msg: "Greetings"})
	if err != nil {
		t.Fatal(err)
	}

	got := <-c
	if got != "Greetings" {
		t.Fatalf("monitor receive wrong msg, want:%s got:%s", "Greetings", got)
	}
}

func reject_if_token_invalid(t *testing.T) {
	test := func(t *testing.T, token string, want error) {
		sut := pkg.NewChatClient(pkg.ChatAddr.Addr(), pkg.CertsPath.Pem())
		sut.UpdateToken(token)

		err := sut.Say(&proto.SayReq{Msg: "Grettings"})

		if status.Code(err) != status.Code(want) {
			t.Fatalf("sut.Say(_) return [%v], wants:%v", err, want)
		}
	}

	t.Run("empty token", func(t *testing.T) {
		test(t, "", pkg.InvalidToken)
	})
	t.Run("invalid token", func(t *testing.T) {
		test(t, "invalid-token", pkg.InvalidToken)
	})
}
