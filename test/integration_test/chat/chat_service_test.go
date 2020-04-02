package chat

import (
	"mse/pkg"
	"mse/proto"
	"testing"
)

func init() {
	pkg.ChatAddr.Attach()
	pkg.CertsPath.Attach()
}

func TestChatService(t *testing.T) {
	t.Run("monitor_can_know_what_sender_say", monitor_can_know_what_sender_say)
}

func monitor_can_know_what_sender_say(t *testing.T) {
	monitor := pkg.NewChatClient(pkg.ChatAddr.Addr(), pkg.CertsPath.Pem())
	listenDone := make(chan bool)
	defer close(listenDone)
	c, err := monitor.Listen(listenDone)
	if err != nil {
		t.Fatal(err)
	}

	sender := pkg.NewChatClient(pkg.ChatAddr.Addr(), pkg.CertsPath.Pem())

	err = sender.Say(&proto.SayReq{Msg: "Greetings"})
	if err != nil {
		t.Fatal(err)
	}

	got := <-c
	if got != "Greetings" {
		t.Fatalf("monitor receive wrong msg, want:%s got:%s", "Greetings", got)
	}
}
