package main

import (
	"context"
	"echo/proto"
	"flag"
	"google.golang.org/grpc"
	"testing"
)

var host string

func init() {
	flag.StringVar(&host, "host", "", "echo-service host")
}

type testParam struct {
	sendMsg string
	recvMsg string
	same    bool
}

type I interface {
	Echo(req *proto.EchoReq) (*proto.EchoRsp, error)
}

func runTest(server I, t *testing.T) {
	t.Helper()

	check := func(server I, t *testing.T, param testParam) {
		req := &proto.EchoReq{Msg: param.sendMsg}
		var rsp *proto.EchoRsp
		var err error
		rsp, err = server.Echo(req)
		if err != nil {
			t.Errorf("call service.Echo failed: %v", err)
			return
		}

		actual := rsp.Msg == param.recvMsg
		//t.Logf("test failed, param: %v, rsp: %v, actual:%v", param, rsp, actual)
		if actual != param.same {
			t.Errorf("test failed, param: %v, rsp: %v", param, rsp)
		}
	}

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

type functionTestingServer struct {
	server *EchoService
}

func (s *functionTestingServer) Echo(req *proto.EchoReq) (*proto.EchoRsp, error) {
	return s.server.Echo(context.Background(), req)
}

func TestFunctional(t *testing.T) {
	server := &functionTestingServer{server: &EchoService{}}
	runTest(server, t)
}

type integrationTestingServer struct {
	conn *grpc.ClientConn
}

func (s *integrationTestingServer) Echo(req *proto.EchoReq) (*proto.EchoRsp, error) {
	client := proto.NewEchoClient(s.conn)
	return client.Echo(context.Background(), req)
}

func TestIntegration(t *testing.T) {
	t.Logf("Host :%v", host)
	conn, err := grpc.Dial(host, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		t.Fatal("grpc.Dial failed: ", err)
	}

	runTest(&integrationTestingServer{conn: conn}, t)
}
