package main

import (
	"context"
	"echo/proto"
	"testing"
)

func TestEchoService_Echo(t *testing.T) {
	type testParam struct {
		sendMsg string
		recvMsg string
		same    bool
	}

	check := func(server *EchoService, t *testing.T, param testParam) {
		t.Helper()

		req := &proto.EchoReq{Msg: param.sendMsg}
		rsp, err := server.Echo(context.Background(), req)
		if err != nil {
			t.Errorf("call server.Echo failed: %v", err)
			return
		}

		actual := rsp.Msg == param.recvMsg
		//t.Logf("test failed, param: %v, rsp: %v, actual:%v", param, rsp, actual)
		if actual != param.same {
			t.Errorf("test failed, param: %v, rsp: %v", param, rsp)
			//} else {
			//	t.Log("test pass")
		}
	}

	server := &EchoService{}

	t.Run("Expect same", func(t *testing.T) {
		check(server, t, testParam{
			sendMsg: "expected message",
			recvMsg: "expected message",
			same:    true,
		})
	})

	t.Run("Expect fail", func(t *testing.T) {
		check(server, t, testParam{
			sendMsg: "expected message",
			recvMsg: "another message",
			same:    false,
		})
	})
}
