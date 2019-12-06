package mircodec

import (
	"encoding/binary"
	"errors"
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/codec"
	"github.com/davyxu/golog"
	"reflect"
)

var log = golog.New("codec.mircodec")

// MirCodec 编码解码
type MirCodec struct {
}

// Name 返回名字
func (m *MirCodec) Name() string {
	return "MirCodec"
}

// MimeType 我也不知道是干嘛的
func (m *MirCodec) MimeType() string {
	return "application/binary"
}

// encode 把结构体编码(序列化)成字节数组
func encode(obj interface{}) (bytes []byte, err error) {
	v := reflect.ValueOf(obj).Elem()
	// 1. 遍历结构体每个字段
	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		// t.Log(fieldType.Name, fieldType.Type.Kind())
		switch f.Kind() {
		case reflect.String:
			vv := f.Interface().(string)
			sb := StringToBytes(vv)
			bytes = append(bytes, sb...)
		case reflect.Uint8:
			vv := f.Interface().(uint8)
			bytes = append(bytes, vv)
		case reflect.Uint16:
			vv := f.Interface().(uint16)
			b := make([]byte, 16)
			binary.LittleEndian.PutUint16(b, uint16(vv))
			bytes = append(bytes, b...)
		case reflect.Uint64:
			vv := f.Interface().(uint64)
			b := make([]byte, 64)
			binary.LittleEndian.PutUint64(b, uint64(vv))
			bytes = append(bytes, b...)
		case reflect.Struct:
			vv := f.Addr().Interface()
			data, err := encode(vv)
			if err != nil {
				return bytes, err
			}
			bytes = append(bytes, data...)
		default:
			errors.New("error")
		}
	}
	return bytes, nil
}

// Encode 将数据转换为字节数组
func (*MirCodec) Encode(msgObj interface{}, ctx cellnet.ContextSet) (data interface{}, err error) {
	return encode(msgObj)
}

func decodeValue(f reflect.Value, bytes []byte) []byte {
	switch f.Type().Kind() {
	case reflect.String:
		i, s := ReadString(bytes, 0)
		f.SetString(s)
		bytes = bytes[i:]
	case reflect.Uint8:
		f.SetUint(uint64(bytes[0]))
		bytes = bytes[1:]
	case reflect.Uint16:
		f.SetUint(uint64(BytesToUint16(bytes[:2])))
		bytes = bytes[2:]
	case reflect.Uint64:
		f.SetUint(uint64(BytesToUint64(bytes[:8])))
		bytes = bytes[8:]
	case reflect.Struct:
		l := f.NumField()
		for i := 0; i < l; i++ {
			bytes = decodeValue(f.Field(i), bytes)
		}
	}
	return bytes
}

// decode 把字节数组解码(反序列化)成结构
func decode(obj interface{}, bytes []byte) (err error) {
	if len(bytes) == 0 {
		return
	}
	v := reflect.ValueOf(obj)
	if v.Type().Kind() == reflect.Ptr {
		v = v.Elem()
	}
	for i := 0; i < v.NumField(); i++ {
		bytes = decodeValue(v.Field(i), bytes)
	}
	return
}

// Decode 将字节数组转换为数据
func (*MirCodec) Decode(data interface{}, msgObj interface{}) error {
	return decode(msgObj, data.([]uint8))
}

func init() {
	codec.RegisterCodec(new(MirCodec))
}
