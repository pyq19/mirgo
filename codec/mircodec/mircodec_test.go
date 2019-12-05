package mircodec

import (
	"encoding/binary"
	"reflect"
	"testing"
)

type CatShop struct {
	Address string
	Rank    uint64
	Cat     CatInfo
}

type CatInfo struct {
	Name string
	Age  uint16
}

func i(obj interface{}) interface{} {
	return obj
}

func _decode(obj interface{}, bytes []byte) {

}

// TestDecode 把字节数组解码(反序列化)成结构体
func TestDecode(t *testing.T) {

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
			bytes = append(bytes, []byte(vv)...)
		case reflect.Uint16:
			vv := f.Interface().(uint16)
			b := make([]byte, 16)
			binary.LittleEndian.PutUint16(b, uint16(vv))
			bytes = append(bytes)
		case reflect.Uint64:
			vv := f.Interface().(uint64)
			b := make([]byte, 64)
			binary.LittleEndian.PutUint64(b, uint64(vv))
			bytes = append(bytes)
		case reflect.Struct:
			vv := f.Addr().Interface()
			_encode(t, vv)
		default:
			t.Log("error")
		}
	}
	return bytes
}

// TestEncode 把结构体编码(序列化)成字节数组
func TestEncode(t *testing.T) {
	catShop := CatShop{
		"独山路宠物店",
		123,
		CatInfo{
			"花猫",
			2,
		},
	}
	// t.Log(reflect.TypeOf(i(&catShop)))
	// t.Log(reflect.TypeOf(i(&catShop)).Elem())
	// t.Logf("%v", i(catShop))
	// catShopType := reflect.TypeOf(catShop)
	// t.Logf("%v", catShopType)
	// t.Log(catShopType.Name())

	t.Log(_encode(t, &catShop))
}
