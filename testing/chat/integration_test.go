package chat

import (
	"errors"
	"flag"
	"fmt"
	"mse/pkg"
	"mse/proto"
	"reflect"
	"testing"
	"time"
)

type MonitorClient struct {
	gotMsgList []string
	data       *TestData
}

func (mc *MonitorClient) Attach(caller Caller) error {
	r, err := caller.Call(nil)
	if err != nil {
		return err
	}

	c, ok := r.(chan string)
	if !ok {
		return fmt.Errorf("MonitorClient.Attach failed, got.(type):%v", reflect.TypeOf(r))
	}

	for msg := range c {
		mc.gotMsgList = append(mc.gotMsgList, msg)
	}
	return nil
}

func (mc *MonitorClient) Check() error {
	for retry := 3; retry > 0 && len(mc.gotMsgList) == 0; retry-- {
		time.Sleep(100 * time.Millisecond)
	}
	if len(mc.gotMsgList) == 0 {
		return errors.New("MonitorClient got nothing")
	}

	if mc.data.wantMsg != mc.gotMsgList[0] {
		return fmt.Errorf("MonitorClient got a wrong msg, want:%s got:%s", mc.data.wantMsg, mc.gotMsgList[0])
	}

	return nil
}

func (mc *MonitorClient) Setup(data *TestData) {
	mc.data = data
}

type SenderClient struct {
	data   *TestData
	gotMsg string
}

func (sc *SenderClient) Send(caller Caller) error {
	_, err := caller.Call(sc.data.msgToSend)
	if err != nil {
		return err
	}

	//rsp, ok := r.(*proto.SayRsp)
	//if !ok {
	//	return fmt.Errorf("SenderClient.Send failed, got.(type):%v", reflect.TypeOf(r))
	//}
	//
	//sc.gotMsg = rsp.Msg
	return nil
}

func (sc *SenderClient) Check() error {
	//if sc.gotMsg != sc.data.wantMsg {
	//	return fmt.Errorf("SenderClient got wrong msg, want:%s got:%s", sc.data.wantMsg, sc.gotMsg)
	//}
	return nil
}

func (sc *SenderClient) Setup(data *TestData) {
	sc.data = data
}

type CallerClient struct {
	data   *TestData
	client *pkg.ChatClient
	done   chan bool
}

func (cc *CallerClient) Setup(data *TestData) {
	cc.data = data
}

func (cc *CallerClient) Call(req interface{}) (interface{}, error) {
	if req == nil {
		return cc.Listen()
	}

	msg, ok := req.(string)
	if !ok {
		return nil, fmt.Errorf("CallerClient.Call sender failed, req.(type):%v", reflect.TypeOf(req))
	}
	return cc.Say(&proto.SayReq{Msg: msg})
}

func (cc *CallerClient) Listen() (chan string, error) {
	return cc.client.Listen(cc.done)
}

func (cc *CallerClient) Say(req *proto.SayReq) (*proto.SayRsp, error) {
	return nil, cc.client.Say(req)
}

var (
	host string
	port uint
	addr string
)

func init() {
	flag.StringVar(&host, "chat-service-host", "chat-service", "chat service host")
	flag.UintVar(&port, "chat-service-port", 0, "chat service port")
	flag.Parse()
	addr = fmt.Sprintf("%s:%d", host, port)
}

func TestChatService_Integration(t *testing.T) {
	client := pkg.NewChatClient(addr)
	done := make(chan bool)
	t.Run("chat service integration testing", func(t *testing.T) {
		TestChatService(t,
			&CallerClient{data: nil, client: client, done: done},
			&MonitorClient{gotMsgList: make([]string, 0), data: nil},
			&SenderClient{data: nil, gotMsg: ""},
		)
	})
}

func TestChatService_Unauthentic(t *testing.T) {
	client := pkg.NewChatClient(addr)
	err := client.Say(&proto.SayReq{Msg: "should failed"})
	if err != pkg.AuthErr {
		t.Errorf("want err:[%v] got:[%v]", pkg.AuthErr, err)
	}
}
