package mircodec

import (
	"fmt"
	"reflect"

	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/codec"
	"github.com/davyxu/golog"
	"github.com/davyxu/goobjfmt"
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
func encode(obj interface{}) (bytes []byte) {
	return bytes
}

// Encode 将数据转换为字节数组
func (*MirCodec) Encode(msgObj interface{}, ctx cellnet.ContextSet) (data interface{}, err error) {
	return goobjfmt.BinaryWrite(msgObj)
}

// decode 把字节数组解码(反序列化)成结构
func decode(obj interface{}, bytes []byte) {

}

// Decode 将字节数组转换为数据
func (*MirCodec) Decode(data interface{}, msgObj interface{}) error {
	var bytes []uint8 = data.([]uint8)
	log.Debugln(bytes)
	objType := reflect.TypeOf(msgObj).Elem()
	//log.Debugln(":::" + objType.Elem().Name())
	// 遍历结构体所有成员
	for i := 0; i < objType.NumField(); i++ {
		// 获取每个成员的结构体字段类型
		fieldType := objType.Field(i)
		// 输出成员名和tag
		fmt.Printf("name: %v  tag: '%v'\n", fieldType.Name, fieldType.Tag)
	}
	return goobjfmt.BinaryRead(data.([]byte), msgObj)
}

func init() {
	codec.RegisterCodec(new(MirCodec))
}
