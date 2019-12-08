package mircodec

import (
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

// Encode 将数据转换为字节数组
func (*MirCodec) Encode(msgObj interface{}, ctx cellnet.ContextSet) (data interface{}, err error) {
	data, err = encode(msgObj)

	// TODO 测试用
	v := reflect.ValueOf(msgObj)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	log.Debugf("---> %v 编码后得到: %v\n", v.Type().Name(), data)

	return data, err
}

// Decode 将字节数组转换为数据
func (*MirCodec) Decode(data interface{}, msgObj interface{}) error {
	bytes := data.([]byte)
	log.Debugln("<--- 收到:", bytes)
	return decode(msgObj, bytes)
}

// encode 把结构体编码(序列化)成字节数组
func encode(obj interface{}) (bytes []byte, err error) {
	v := reflect.ValueOf(obj)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		// t.Log(fieldType.Name, fieldType.Type.Kind())
		switch f.Kind() {
		case reflect.Struct:
			vv := f.Addr().Interface()
			data, err := encode(vv)
			if err != nil {
				return bytes, err
			}
			bytes = append(bytes, data...)
		case reflect.String:
			vv := f.Interface().(string)
			sb := StringToBytes(vv)
			bytes = append(bytes, sb...)
		case reflect.Int8:
			vv := uint8(f.Interface().(int8))
			bytes = append(bytes, vv)
		case reflect.Int16:
			vv := uint16(f.Interface().(int16))
			bytes = append(bytes, Uint16ToBytes(vv)...)
		case reflect.Int32:
			vv := uint32(f.Interface().(int32))
			bytes = append(bytes, Uint32ToBytes(vv)...)
		case reflect.Int64:
			vv := uint64(f.Interface().(int64))
			bytes = append(bytes, Uint64ToBytes(vv)...)
		case reflect.Uint8:
			switch vv := f.Interface().(type) {
			case uint8:
				bytes = append(bytes, vv)
			default:
				vvv := reflect.ValueOf(vv).Uint()
				bytes = append(bytes, uint8(vvv))
			}
		case reflect.Uint16:
			vv := f.Interface().(uint16)
			bytes = append(bytes, Uint16ToBytes(vv)...)
		case reflect.Uint32:
			vv := f.Interface().(uint32)
			bytes = append(bytes, Uint32ToBytes(vv)...)
		case reflect.Uint64:
			vv := f.Interface().(uint64)
			bytes = append(bytes, Uint64ToBytes(vv)...)
		case reflect.Slice:
			// FIXME 还有别的类型可能会报错
			switch vv := f.Interface().(type) {
			case []uint8:
				bytes = append(bytes, vv...)
			default:
				vvv := reflect.ValueOf(vv)
				l := vvv.Len()
				slice := vvv.Slice(0, l)
				bytes = append(bytes, Uint32ToBytes(uint32(l))...)
				for i := 0; i < l; i++ {
					b, err := encode(slice.Index(i).Interface())
					if err != nil {
						panic(err)
					}
					bytes = append(bytes, b...)
				}
			}
		default:
			log.Errorln("编码错误")
			return bytes, errors.New("编码错误")
		}
	}
	if bytes == nil {
		bytes = []byte{}
	}
	return bytes, nil
}

func decodeValue(f reflect.Value, bytes []byte) []byte {
	switch f.Type().Kind() {
	case reflect.Struct:
		l := f.NumField()
		for i := 0; i < l; i++ {
			bytes = decodeValue(f.Field(i), bytes)
		}
	case reflect.String:
		i, s := ReadString(bytes, 0)
		f.SetString(s)
		bytes = bytes[i:]
	case reflect.Int8:
		f.SetInt(int64(bytes[0]))
		bytes = bytes[1:]
	case reflect.Int16:
		f.SetInt(int64(BytesToUint16(bytes[:2])))
		bytes = bytes[2:]
	case reflect.Int32:
		f.SetInt(int64(BytesToUint32(bytes[:4])))
		bytes = bytes[4:]
	case reflect.Int64:
		f.SetInt(int64(BytesToUint64(bytes[:8])))
		bytes = bytes[8:]
	case reflect.Uint8:
		f.SetUint(uint64(bytes[0]))
		bytes = bytes[1:]
	case reflect.Uint16:
		f.SetUint(uint64(BytesToUint16(bytes[:2])))
		bytes = bytes[2:]
	case reflect.Uint32:
		f.SetUint(uint64(BytesToUint32(bytes[:4])))
		bytes = bytes[4:]
	case reflect.Uint64:
		f.SetUint(BytesToUint64(bytes[:8]))
		bytes = bytes[8:]
	case reflect.Slice:
		e := f.Type().Elem()
		if e.Kind() == reflect.Uint8 {
			l := BytesToUint32(bytes[:4])
			f.SetBytes(bytes[:l+4])
			bytes = bytes[l+4:]
		} else if e.Kind() == reflect.Struct {
			l := e.NumField()
			bytes = bytes[4:]
			slice := reflect.MakeSlice(f.Type(), l, l)
			for i := 0; i < l; i++ {
				sliceValue := reflect.New(slice.Type().Elem()).Elem()
				bytes = decodeValue(sliceValue, bytes)
				slice.Index(i).Set(sliceValue)
			}
			f.Set(slice)
		} else {
			// FIXME 还有别的类型可能会报错
			log.Errorln("!!!暂不支持的类型解码，待完善")
		}
	default:
		log.Errorln(f.Type(), "解码错误")
	}
	//kk := f.Type().Kind()
	//log.Debugln(kk)
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

func init() {
	codec.RegisterCodec(new(MirCodec))
}
