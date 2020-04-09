package mircodec

import (
	"reflect"

	"github.com/yenkeia/mirgo/game/util"
)

type BytesWrapper struct {
	Bytes *[]byte
}

func (r *BytesWrapper) Last() (res []byte) {
	res = *r.Bytes
	*r.Bytes = make([]byte, 0)
	return
}

func (r *BytesWrapper) ReadByte() (n int8) {
	return r.ReadInt8()
}

func (r *BytesWrapper) ReadInt8() (n int8) {
	b := (*r.Bytes)[0]
	*r.Bytes = (*r.Bytes)[1:]
	return int8(b)
}

func (r *BytesWrapper) ReadInt16() (n int16) {
	b := (*r.Bytes)[:2]
	*r.Bytes = (*r.Bytes)[2:]
	return int16(util.BytesToUint16(b))
}

func (r *BytesWrapper) ReadInt32() (n int32) {
	b := (*r.Bytes)[:4]
	*r.Bytes = (*r.Bytes)[4:]
	return int32(util.BytesToUint32(b))
}

func (r *BytesWrapper) ReadInt64() (n int64) {
	b := (*r.Bytes)[:8]
	*r.Bytes = (*r.Bytes)[8:]
	return int64(util.BytesToUint64(b))
}

func (r *BytesWrapper) ReadSByte() (n uint8) {
	return r.ReadUInt8()
}

func (r *BytesWrapper) ReadUInt8() (n uint8) {
	b := (*r.Bytes)[0]
	*r.Bytes = (*r.Bytes)[1:]
	return b
}

func (r *BytesWrapper) ReadUInt16() (n uint16) {
	b := (*r.Bytes)[:2]
	*r.Bytes = (*r.Bytes)[2:]
	return util.BytesToUint16(b)
}

func (r *BytesWrapper) ReadUInt32() (n uint32) {
	b := (*r.Bytes)[:4]
	*r.Bytes = (*r.Bytes)[4:]
	return util.BytesToUint32(b)
}

func (r *BytesWrapper) ReadUInt64() (n uint64) {
	b := (*r.Bytes)[:8]
	*r.Bytes = (*r.Bytes)[8:]
	return util.BytesToUint64(b)
}

func (r *BytesWrapper) ReadBoolean() bool {
	b := (*r.Bytes)[0]
	*r.Bytes = (*r.Bytes)[1:]
	if int(b) == 0 {
		return false
	}
	return true
}

func (r *BytesWrapper) ReadString() string {
	b := int((*r.Bytes)[0])
	s := string((*r.Bytes)[1 : b+1])
	*r.Bytes = (*r.Bytes)[b+1:]
	return s
}

func (r *BytesWrapper) Write(obj interface{}) {
	if res, err := encodeValue(reflect.ValueOf(obj)); err != nil {
		log.Errorln(err)
	} else {
		*r.Bytes = append(*r.Bytes, res...)
	}
}
