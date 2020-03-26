package main

import (
	"context"
	"mse/auth/proto"
	"reflect"
	"testing"
)

func TestAuthInternalService_GetSignKey(t *testing.T) {
	type args struct {
		in0 context.Context
		in  *proto.GetSignKeyReq
	}
	tests := []struct {
		name    string
		args    args
		want    *proto.GetSignKeyRsp
		wantErr bool
	}{
		{name: "simple check", args: args{
			in0: context.Background(),
			in:  &proto.GetSignKeyReq{},
		}, want: &proto.GetSignKeyRsp{
			Kid: 1,
			Key: "secret-key",
			Alg: "HS256",
		}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ais := &AuthInternalService{}
			got, err := ais.GetSignKey(tt.args.in0, tt.args.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetSign() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetSign() got = %v, want %v", got, tt.want)
			}
		})
	}
}
