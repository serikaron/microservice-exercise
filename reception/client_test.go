package reception

import (
	"bytes"
	"testing"
)

type MockReadWriter struct {
	readBuf   []byte
	writeChan chan []byte
}

func NewMockReaderWriter() *MockReadWriter {
	return &MockReadWriter{
		readBuf:   nil,
		writeChan: make(chan []byte),
	}
}

func (mrw *MockReadWriter) Read(buf []byte) (int, error) {
	n, err := bytes.NewReader(mrw.readBuf).Read(buf)
	mrw.readBuf = nil
	return n, err
}
func (mrw *MockReadWriter) Write(buf []byte) (n int, err error) {
	mrw.writeChan <- buf
	return len(buf), nil
}

func TestClient_Start(t *testing.T) {
	readChan := make(chan []byte)
	mrw := NewMockReaderWriter()
	client := NewClient(0, mrw)
	go client.Start(readChan)

	checkRead := func(t *testing.T, got []byte, want []byte, isEqual bool) {
		t.Helper()
		mrw.readBuf = got

		outBuf := <-readChan
		if isEqual != bytes.Equal(want, outBuf) {
			t.Errorf("checkRead failed, Got:%v, Want:%v, isEqual:%v", outBuf, want, isEqual)
		}
	}

	//readMaxTimes := 0
	t.Run("checkRead, expect same", func(t *testing.T) {
		checkRead(t, []byte("abc"), []byte("abc"), true)
	})
	t.Run("checkRead, expect failed", func(t *testing.T) {
		checkRead(t, []byte("abc"), []byte("cde"), false)
	})

	checkWrite := func(t *testing.T, got []byte, want []byte, isEqual bool) {
		t.Helper()

		client.writeChan <- got
		outBuf := <-mrw.writeChan
		if isEqual != bytes.Equal(outBuf, want) {
			t.Errorf("checkWrite failed, Got:%v, Want:%v, isEqual:%v", outBuf, want, isEqual)
		}
	}

	t.Run("checkWrite, expect same", func(t *testing.T) {
		checkWrite(t, []byte("abc"), []byte("abc"), true)
	})
	t.Run("checkWrite, expect failed", func(t *testing.T) {
		checkWrite(t, []byte("abc"), []byte("cde"), false)
	})

	defer client.Close()
}
