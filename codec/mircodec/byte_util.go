package mircodec

import (
	"encoding/binary"
	"reflect"
)

// ReadString 根据传入的索引 返回读完string后所在bytes 的下一个索引 及string
func ReadString(bytes []byte, index int) (int, string) {
	//[0, 1, 2, 5, 9, 10, 11, 12, 50, 23, 77, 99]
	if len(bytes) == 0 {
		return -1, ""
	}
	strLen := int(bytes[index])
	lastIndex := index + strLen + 1
	msg := string(bytes[index+1 : lastIndex])
	return lastIndex, msg
}

// StringToBytes ...
func StringToBytes(s string) []byte {
	// 根据传入的string 返回字节数组，索引为0（即数组开头）的值为数组的长度
	if s == "" {
		return []byte{0}
	}
	stringBytes := []byte(s)
	stringBytesLenBytes := []byte{byte(len(stringBytes))}
	return append(stringBytesLenBytes, stringBytes...)
}

// Uint16ToBytes ...
func Uint16ToBytes(num uint16) []byte {
	res := make([]byte, 2)
	binary.LittleEndian.PutUint16(res, uint16(num))
	return res
}

// Uint32ToBytes ...
func Uint32ToBytes(num uint32) []byte {
	res := make([]byte, 4)
	binary.LittleEndian.PutUint32(res, uint32(num))
	return res
}

// Uint64ToBytes ...
func Uint64ToBytes(num uint64) []byte {
	res := make([]byte, 8)
	binary.LittleEndian.PutUint64(res, uint64(num))
	return res
}

// BoolToBytes ...
func BoolToBytes(b bool) []byte {
	if b == true {
		return []byte{1}
	}
	return []byte{0}
}

// BytesToUint16 ...
func BytesToUint16(bytes []byte) uint16 {
	return binary.LittleEndian.Uint16(bytes)
}

// BytesToUint32 ...
func BytesToUint32(bytes []byte) uint32 {
	return binary.LittleEndian.Uint32(bytes)
}

// BytesToUint64 ...
func BytesToUint64(bytes []byte) uint64 {
	return binary.LittleEndian.Uint64(bytes)
}

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
	return int16(BytesToUint16(b))
}

func (r *BytesWrapper) ReadInt32() (n int32) {
	b := (*r.Bytes)[:4]
	*r.Bytes = (*r.Bytes)[4:]
	return int32(BytesToUint32(b))
}

func (r *BytesWrapper) ReadInt64() (n int64) {
	b := (*r.Bytes)[:8]
	*r.Bytes = (*r.Bytes)[8:]
	return int64(BytesToUint64(b))
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
	return BytesToUint16(b)
}

func (r *BytesWrapper) ReadUInt32() (n uint32) {
	b := (*r.Bytes)[:4]
	*r.Bytes = (*r.Bytes)[4:]
	return BytesToUint32(b)
}

func (r *BytesWrapper) ReadUInt64() (n uint64) {
	b := (*r.Bytes)[:8]
	*r.Bytes = (*r.Bytes)[8:]
	return BytesToUint64(b)
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
