package mircodec

import "testing"

func TestBytesWrapper_ReadString(t *testing.T) {
	bytes := []byte{10, 104, 101, 108, 108, 111, 119, 111, 114, 108, 100}
	reader := &BytesWrapper{Bytes: &bytes}
	s := reader.ReadString()
	t.Log(s)
}
