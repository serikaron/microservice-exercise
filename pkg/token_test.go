package pkg

import (
	"bytes"
	"io"
	"reflect"
	"testing"
)

func Test_jwtToken_readKey(t *testing.T) {
	type fields struct {
		buf []byte
	}
	type args struct {
		reader io.Reader
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		wantSignKey []byte
		wantErr     error
	}{
		{
			"empty buf",
			fields{buf: make([]byte, 0)},
			args{bytes.NewReader([]byte("some string"))},
			nil,
			nil,
		},
		{
			"nothing to read",
			fields{buf: make([]byte, 10)},
			args{bytes.NewReader(make([]byte, 0))},
			nil,
			nil,
		},
		{
			"short buf not aligned",
			fields{buf: make([]byte, 3)},
			args{bytes.NewReader([]byte("1234567890"))},
			[]byte("1234567890"),
			nil,
		},
		{
			"short buf aligned",
			fields{buf: make([]byte, 5)},
			args{bytes.NewReader([]byte("1234567890"))},
			[]byte("1234567890"),
			nil,
		},
		{
			"length match",
			fields{buf: make([]byte, 10)},
			args{bytes.NewReader([]byte("1234567890"))},
			[]byte("1234567890"),
			nil,
		},
		{
			"long buf",
			fields{buf: make([]byte, 15)},
			args{bytes.NewReader([]byte("1234567890"))},
			[]byte("1234567890"),
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jw := &jwtToken{
				method: nil,
				buf:    tt.fields.buf,
			}
			gotSignKey, err := jw.readKey(tt.args.reader)
			if err != tt.wantErr {
				t.Errorf("readKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotSignKey, tt.wantSignKey) {
				t.Errorf("readKey() gotSignKey = %v, want %v", gotSignKey, tt.wantSignKey)
			}
		})
	}
}
