package testing

import (
	"errors"
	"fmt"
	"log"
	"reflect"
	"testing"
)

type TestInfo struct {
	Name    string
	Req     interface{}
	Want    interface{}
	WantErr error
}

var RspErr = errors.New("check rsp failed")

type Test interface {
	Run() error
}

type TestCase struct {
	Infos   []TestInfo
	Builder func(info *TestInfo) Test
}

func (tc *TestCase) Run() error {
	for _, ti := range tc.Infos {
		t := tc.Builder(&ti)
		_ = t.Run()
	}
	return nil
}

type MethodTest struct {
	Info   *TestInfo
	Method func(interface{}) (interface{}, error)
}

func (mt *MethodTest) Run() error {
	got, err := mt.Method(mt.Info.Req)
	gotMatch := err == nil && reflect.DeepEqual(got, mt.Info.Want)
	switch mt.Info.WantErr {
	case RspErr:
		if gotMatch {
			return errors.New(fmt.Sprintf("[%s] test failed, GotErr %v, WantErr %v", mt.Info.Name, err, mt.Info.WantErr))
		}
	case nil:
		if !gotMatch {
			return errors.New(fmt.Sprintf("[%s] test failed, Got %v, Want %v", mt.Info.Name, got, mt.Info.Want))
		}
	default:
		if err == nil || mt.Info.WantErr.Error() != err.Error() {
			return errors.New(fmt.Sprintf("[%s] GotErr %v, WantErr %v", mt.Info.Name, err, mt.Info.WantErr))
		}
	}
	return nil
}

type GoTestWrapper struct {
	T    *testing.T
	Core Test
	Info *TestInfo
}

func (gtw *GoTestWrapper) Run() error {
	gtw.T.Run(gtw.Info.Name, func(t *testing.T) {
		//t.Helper()
		err := gtw.Core.Run()
		if err != nil {
			t.Error(err)
		}
	})
	return nil
}

type SimpleWrapper struct {
	Core Test
}

func (sw *SimpleWrapper) Run() error {
	if err := sw.Core.Run(); err != nil {
		log.Fatalln(err)
	}
	return nil
}
