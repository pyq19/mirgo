package mircodec

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/codec"
	"github.com/davyxu/golog"
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
	// TODO 测试用
	//v := reflect.ValueOf(msgObj)
	//if v.Kind() == reflect.Ptr {
	//	v = v.Elem()
	//}
	//log.Debugf("Encode %v\n", v.Type().Name())

	return encode(msgObj)
}

// Decode 将字节数组转换为数据
func (*MirCodec) Decode(data interface{}, msgObj interface{}) error {
	bytes := data.([]byte)
	return decode(msgObj, bytes)
}

func init() {
	codec.RegisterCodec(new(MirCodec))
}
