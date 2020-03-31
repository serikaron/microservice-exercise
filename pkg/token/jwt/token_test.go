package jwt

import (
	"errors"
	"mse/pkg"
	"mse/pkg/token"
	"testing"
	"time"
)

`
func Test_read_key_from_string_buf(t *testing.T) {
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
			jw := &JwtToken{
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

func Test_read_key_from_file(t *testing.T) {
	tests := []struct {
		name     string
		filepath string
		want     []byte
	}{
		{
			name:     "read-key-from-file",
			filepath: "../res/certs/cert.pem",
			want:     []byte(certPem),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jw := &JwtToken{
				method: nil,
				buf:    make([]byte, 1024),
			}
			file, err := os.Open(tt.filepath)
			if err != nil {
				t.Fatal(err)
			}

			reader := bufio.NewReader(file)

			got, err := jw.readKey(reader)
			if err != nil {
				t.Fatal(err)
			}

			if string(got) != string(tt.want) {
				t.Fatalf("get wrong key, want:%s got:%s", string(tt.want), string(got))
			}
		})
	}
}
`

func Test_gen_jwt_token_string(t *testing.T) {
	var readKeyErr = errors.New("readKeyErr")

	type args struct {
	}
	tests := map[string]struct {
		sut       *JwtToken
		key       token.Key
		identity  pkg.Identity
		wantToken string
		wantErr   error
	}{
		"HS256 token should work properly": {
			sut:       NewHS256Token(),
			key:       TestingSignKey{86400, "1", []byte(certPem), nil},
			identity:  pkg.Identity{Name: "Marry"},
			wantToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhbGciOiJIUzI1NiIsImV4cCI6ODY0MDAsImtpZCI6MSwibmFtZSI6Ik1hcnJ5IiwidHlwIjoiSldUIn0.bXNdIPoVRenyCiCU6N1Bz6njXwsTSYwHQwQ4sWKmimU",
			wantErr:   nil,
		},
		"invalid key should return error": {
			sut:       NewHS256Token(),
			key:       TestingSignKey{86400, "1", []byte(certPem), readKeyErr},
			identity:  pkg.Identity{Name: "Marry"},
			wantToken: "",
			wantErr:   readKeyErr,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			tokenString, err := tt.sut.Gen(tt.identity, tt.key)
			if tt.wantErr != err {
				t.Fatalf("unexpect error want:%v got:%v", tt.wantErr, err)
			}
			if tokenString != tt.wantToken {
				t.Fatalf("tokenString want:%s got:%s", tt.wantToken, tokenString)
			}
		})
	}
}

type TestingSignKey struct {
	expireAt uint64
	kid      string
	key      []byte
	readErr  error
}

func (tsk *TestingSignKey) Kid() string {
	return tsk.kid
}

func (tsk *TestingSignKey) Read() ([]byte, error) {
	return tsk.key, tsk.readErr
}

func (tsk *TestingSignKey) ExpireAt() *time.Time {
	t := new(time.Time)
	t.Add(86400 * time.Second)
	return t
}

const certPem string = `-----BEGIN RSA PRIVATE KEY-----
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
-----END RSA PRIVATE KEY-----
`
