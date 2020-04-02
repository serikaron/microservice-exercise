package end_to_end

import (
	"google.golang.org/grpc/status"
	"mse/pkg"
	"mse/proto"
	"testing"
)

func init() {
	pkg.ChatAddr.Attach()
}

func TestAuthentication(t *testing.T) {
	t.Run("chat_without_authentication_is_denied", chat_without_authentication_is_denied)
}

func chat_without_authentication_is_denied(t *testing.T) {
	chat := pkg.NewChatClient(pkg.ChatAddr.Addr())

	err := chat.Say(&proto.SayReq{Msg: "Grettings"})

	if status.Code(err) != status.Code(pkg.MissingToken) {
		t.Fatalf("err not the same want:%v got:%v", pkg.MissingToken, err)
	}
}
