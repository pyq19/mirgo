package mircodec

import (
	"encoding/binary"
	"reflect"
	"testing"
)

func i(obj interface{}) interface{} {
	return obj
}

type CatShop struct {
	Address string
	Rank    uint8
	Cat     CatInfo
}

type CatInfo struct {
	Name string
	Age  uint16
}

func _decode(t *testing.T, obj interface{}, bytes []byte) {
	if len(bytes) == 0 {
		return
	}
	t.Log("->", reflect.TypeOf(obj))
	v := reflect.ValueOf(obj)
	if v.Type().Kind() == reflect.Ptr {
		v = v.Elem()
	}
	t.Log("-->", v.NumField(), v.Type(), v.Type().Name())
	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		t.Log("--->", f.Type().Name())
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
			// t.Log("!!!", f.NumField())
			o := f.Interface()
			_decode(t, o, bytes)
		}
	}
}

// TestDecode 把字节数组解码(反序列化)成结构体
func TestDecode(t *testing.T) {
	bytes := []byte{18, 231, 139, 172, 229, 177, 177, 232, 183, 175, 229, 174, 160, 231, 137, 169, 229, 186, 151,
		111,
		6, 232, 138, 177, 231, 140, 171,
		172, 8, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	catShop := new(CatShop)
	_decode(t, i(catShop), bytes)
	t.Log(catShop.Address)
}

func TestByteSlice(t *testing.T) {
	bytes := []byte{5, 11, 1, 66, 3, 63, 98}
	t.Log(bytes[3])
	t.Log(bytes[:2]) // [5 11]
	t.Log(bytes[2:]) // [1 66 3 63 98]
	t.Log(bytes[3:]) // [66 3 63 98]
}

func _encode(t *testing.T, obj interface{}) (bytes []byte) {
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
			bytes = append(bytes, _encode(t, vv)...)
		default:
			t.Log("!!!!!error")
		}
	}
	return bytes
}

// TestEncode 把结构体编码(序列化)成字节数组
func TestEncode(t *testing.T) {
	catShop := CatShop{
		"独山路宠物店",
		111,
		CatInfo{
			"花猫",
			2220,
		},
	}
	res := _encode(t, &catShop)
	t.Log("->", res)
}
