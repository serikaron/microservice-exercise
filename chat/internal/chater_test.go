package internal_test

import (
	"mse/chat/internal"
	"testing"
)

func TestBuilder(t *testing.T) {
	t.Run("message_will_prefix_with_who_said", message_will_prefix_with_who_said)
}

func message_will_prefix_with_who_said(t *testing.T) {
	msg := internal.Chat("Marry", "Greetings")

	if msg != "Marry: Greetings" {
		t.Fatalf("Chater.Chat return %s want %s", msg, "Marry: Greetings")
	}
}
