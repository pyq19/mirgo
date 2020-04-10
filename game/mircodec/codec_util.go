package mircodec

import (
	"bytes"
	"errors"
	"reflect"

	"github.com/yenkeia/mirgo/game/util"
)

// 在指针数组元素前加一个bool值，标示该元素是否为nil
func encodeEmptyFlagArray(v reflect.Value) ([]byte, error) {
	vvv := reflect.ValueOf(v.Interface())
	l := vvv.Len()
	slice := vvv.Slice(0, l)

	buffer := bytes.Buffer{}
	if slice.IsNil() {
		buffer.WriteByte(0)
	} else {
		buffer.WriteByte(1)
	}

	buffer.Write(util.Uint32ToBytes(uint32(l)))

	for i := 0; i < l; i++ {
		if slice.Index(i).IsNil() {
			buffer.WriteByte(0)
			continue
		} else {
			buffer.WriteByte(1)
			b, err := encode(slice.Index(i).Interface())
			if err != nil {
				panic(err)
			}
			buffer.Write(b)
		}
	}
	return buffer.Bytes(), nil
}

func encode(obj interface{}) (bytes []byte, err error) {
	v := reflect.ValueOf(obj)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	for i := 0; i < v.NumField(); i++ {
		tag := v.Type().Field(i).Tag.Get("codec")
		if tag == "-" {
			continue
		}
		if tag == "emptyflag" {
			encodeBytes, err := encodeEmptyFlagArray(v.Field(i))
			if err != nil {
				return bytes, err
			}
			bytes = append(bytes, encodeBytes...)
		} else {
			encodeBytes, err := encodeValue(v.Field(i))
			if err != nil {
				return bytes, err
			}
			bytes = append(bytes, encodeBytes...)
		}
	}
	return
}

// encode 把结构体编码(序列化)成字节数组
func encodeValue(v reflect.Value) (bytes []byte, err error) {
	if v.Type().Kind() == reflect.Ptr {
		v = v.Elem()
	}
	switch v.Kind() {
	case reflect.Struct:
		vv := v.Addr().Interface()
		data, err := encode(vv)
		if err != nil {
			return bytes, err
		}
		bytes = append(bytes, data...)
	case reflect.Bool:
		vv := v.Interface().(bool)
		b := byte(0)
		if vv == true {
			b = 1
		}
		bytes = append(bytes, b)
	case reflect.String:
		vv := v.Interface().(string)
		sb := util.StringToBytes(vv)
		bytes = append(bytes, sb...)
	case reflect.Int:
		vv := uint32(v.Interface().(int))
		bytes = append(bytes, util.Uint32ToBytes(vv)...)
	case reflect.Int8:
		vv := uint8(v.Interface().(int8))
		bytes = append(bytes, vv)
	case reflect.Int16:
		//vv := uint16(v.Interface().(int16))
		//bytes = append(bytes, util.Uint16ToBytes(vv)...)
		if vv, ok := v.Interface().(int16); ok {
			bytes = append(bytes, util.Uint16ToBytes(uint16(vv))...)
			break
		}
		vvv := reflect.ValueOf(v.Interface()).Uint()
		bytes = append(bytes, util.Uint16ToBytes(uint16(vvv))...)
	case reflect.Int32:
		//vv := uint32(v.Interface().(int32))
		//bytes = append(bytes, util.Uint32ToBytes(vv)...)
		if vv, ok := v.Interface().(int32); ok {
			bytes = append(bytes, util.Uint32ToBytes(uint32(vv))...)
			break
		}
		vvv := reflect.ValueOf(v.Interface()).Uint()
		bytes = append(bytes, util.Uint32ToBytes(uint32(vvv))...)
	case reflect.Int64:
		//vv := uint64(v.Interface().(int64))
		//bytes = append(bytes, util.Uint64ToBytes(vv)...)
		if vv, ok := v.Interface().(int64); ok {
			bytes = append(bytes, util.Uint64ToBytes(uint64(vv))...)
			break
		}
		vvv := reflect.ValueOf(v.Interface()).Uint()
		bytes = append(bytes, util.Uint64ToBytes(uint64(vvv))...)
	case reflect.Uint8:
		switch vv := v.Interface().(type) {
		case uint8:
			bytes = append(bytes, vv)
		default:
			vvv := reflect.ValueOf(vv).Uint()
			bytes = append(bytes, uint8(vvv))
		}
	case reflect.Uint16:
		if vv, ok := v.Interface().(uint16); ok {
			bytes = append(bytes, util.Uint16ToBytes(vv)...)
			break
		}
		vvv := reflect.ValueOf(v.Interface()).Uint()
		bytes = append(bytes, util.Uint16ToBytes(uint16(vvv))...)
	case reflect.Uint32:
		if vv, ok := v.Interface().(uint32); ok {
			bytes = append(bytes, util.Uint32ToBytes(vv)...)
			break
		}
		vvv := reflect.ValueOf(v.Interface()).Uint()
		bytes = append(bytes, util.Uint32ToBytes(uint32(vvv))...)
	case reflect.Uint64:
		if vv, ok := v.Interface().(uint64); ok {
			bytes = append(bytes, util.Uint64ToBytes(vv)...)
			break
		}
		vvv := reflect.ValueOf(v.Interface()).Uint()
		bytes = append(bytes, util.Uint64ToBytes(uint64(vvv))...)
	case reflect.Float32:
		if vv, ok := v.Interface().(float32); ok {
			bytes = append(bytes, util.Float32ToBytes(vv)...)
			break
		}
		vvv := reflect.ValueOf(v.Interface()).Float()
		bytes = append(bytes, util.Float32ToBytes(float32(vvv))...)
	case reflect.Float64:
		if vv, ok := v.Interface().(float64); ok {
			bytes = append(bytes, util.Float64ToBytes(vv)...)
			break
		}
		vvv := reflect.ValueOf(v.Interface()).Float()
		bytes = append(bytes, util.Float64ToBytes(vvv)...)
	case reflect.Slice:
		// FIXME 还有别的类型可能会报错
		switch vv := v.Interface().(type) {
		case []uint8:
			bytes = append(bytes, vv...)
		case []string:
			l := len(vv)
			bytes = append(bytes, util.Uint32ToBytes(uint32(l))...)
			for i := 0; i < l; i++ {
				bytes = append(bytes, util.StringToBytes(vv[i])...)
			}
		case []int: // int32
			l := len(vv)
			bytes = append(bytes, util.Uint32ToBytes(uint32(l))...)
			for i := 0; i < l; i++ {
				bytes = append(bytes, util.Uint32ToBytes(uint32(vv[i]))...)
			}
		case []uint: // uint32
			l := len(vv)
			bytes = append(bytes, util.Uint32ToBytes(uint32(l))...)
			for i := 0; i < l; i++ {
				bytes = append(bytes, util.Uint32ToBytes(uint32(vv[i]))...)
			}
		default:
			vvv := reflect.ValueOf(vv)
			l := vvv.Len()
			slice := vvv.Slice(0, l)
			bytes = append(bytes, util.Uint32ToBytes(uint32(l))...)
			for i := 0; i < l; i++ {
				b, err := encodeValue(slice.Index(i))
				if err != nil {
					panic(err)
				}
				bytes = append(bytes, b...)
			}
		}
	default:
		return bytes, errors.New("编码错误: " + v.Kind().String())
	}
	return bytes, nil
}

func decodeValue(f reflect.Value, bytes []byte) []byte {
	defer func() {
		err := recover()
		if err != nil {
			log.Warnln(err)
		}
	}()
	if f.Type().Kind() == reflect.Ptr {
		f = f.Elem()
	}
	switch f.Type().Kind() {
	case reflect.Struct:
		l := f.NumField()
		for i := 0; i < l; i++ {
			bytes = decodeValue(f.Field(i), bytes)
		}
	case reflect.String:
		i, s := util.ReadString(bytes, 0)
		f.SetString(s)
		bytes = bytes[i:]
	case reflect.Bool:
		tmp := bytes[0]
		b := true
		if tmp == 0 {
			b = false
		}
		f.SetBool(b)
		bytes = bytes[1:]
	case reflect.Int8:
		f.SetInt(int64(bytes[0]))
		bytes = bytes[1:]
	case reflect.Int16:
		f.SetInt(int64(util.BytesToUint16(bytes[:2])))
		bytes = bytes[2:]
	case reflect.Int, reflect.Int32:
		f.SetInt(int64(util.BytesToUint32(bytes[:4])))
		bytes = bytes[4:]
	case reflect.Int64:
		f.SetInt(int64(util.BytesToUint64(bytes[:8])))
		bytes = bytes[8:]
	case reflect.Uint8:
		f.SetUint(uint64(bytes[0]))
		bytes = bytes[1:]
	case reflect.Uint16:
		f.SetUint(uint64(util.BytesToUint16(bytes[:2])))
		bytes = bytes[2:]
	case reflect.Uint32:
		f.SetUint(uint64(util.BytesToUint32(bytes[:4])))
		bytes = bytes[4:]
	case reflect.Uint64:
		f.SetUint(util.BytesToUint64(bytes[:8]))
		bytes = bytes[8:]
	case reflect.Float32:
		f.SetFloat(float64(util.BytesToFloat32(bytes[:4])))
		bytes = bytes[4:]
	case reflect.Float64:
		f.SetFloat(util.BytesToFloat64(bytes[:8]))
		bytes = bytes[8:]
	case reflect.Slice:
		l := int(util.BytesToUint32(bytes[:4]))
		e := f.Type().Elem()
		switch e.Kind() {
		case reflect.Int8, reflect.Uint8:
			f.SetBytes(bytes[:l+4])
			bytes = bytes[l+4:]
		case reflect.Struct:
			bytes = bytes[4:]
			slice := reflect.MakeSlice(f.Type(), l, l)
			for i := 0; i < l; i++ {
				sliceValue := reflect.New(slice.Type().Elem()).Elem()
				bytes = decodeValue(sliceValue, bytes)
				slice.Index(i).Set(sliceValue)
			}
			f.Set(slice)
		case reflect.String:
			sl := int(util.BytesToUint32(bytes[:4]))
			bytes = bytes[4:]
			slice := reflect.MakeSlice(f.Type(), l, l)
			for si := 0; si < sl; si++ {
				sliceValue := reflect.New(slice.Type().Elem()).Elem()
				i, s := util.ReadString(bytes, 0)
				bytes = bytes[i:]
				sliceValue.SetString(s)
				slice.Index(si).Set(sliceValue)
			}
			f.Set(slice)
		case reflect.Int, reflect.Uint, reflect.Int32, reflect.Uint32:
			il := int(util.BytesToUint32(bytes[:4]))
			bytes = bytes[4:]
			slice := reflect.MakeSlice(f.Type(), l, l)
			for ii := 0; ii < il; ii++ {
				sliceValue := reflect.New(slice.Type().Elem()).Elem()
				i := util.BytesToUint32(bytes)
				bytes = bytes[4:]
				sliceValue.SetInt(int64(i))
				slice.Index(ii).Set(sliceValue)
			}
			f.Set(slice)
		default:
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
		tag := v.Type().Field(i).Tag.Get("codec")
		if tag == "-" {
			continue
		}
		bytes = decodeValue(v.Field(i), bytes)
	}
	return
}
