package mircodec

import (
	"github.com/davyxu/cellnet"
	"github.com/yenkeia/mirgo/common"
	"github.com/yenkeia/mirgo/proto/server"
	"reflect"
	"testing"
)

func TestConst(t *testing.T) {
	t.Log(server.CONNECTED)
	t.Log(server.CLIENT_VERSION)
	t.Log(reflect.TypeOf(server.CONNECTED))
}

func TestMapInformation(t *testing.T) {
	//[2019/12/8 16:14:38]: --->发送服务端包信息: ServerPackets.MapInformation
	//[2019/12/8 16:14:38]: --->发送服务端包字节信息: 30, 0, 17, 0, 1, 48, 14, 66, 105, 99, 104, 111, 110, 80, 114, 111, 118, 105, 110, 99, 101, 101, 0, 135, 0, 0, 0, 0, 0, 0,

	bytes := []byte{
		1, 48,
		14, 66, 105, 99, 104, 111, 110, 80, 114, 111, 118, 105, 110, 99, 101,
		101, 0,
		135, 0,
		0, 0,
		0,
		0,
		0}
	t.Log(bytes)
	msg := new(server.MapInformation)

	codec := new(MirCodec)
	if err := codec.Decode(bytes, msg); err != nil {
		panic(err)
	}
	t.Log(msg)

	cs := new(cellnet.ContextSet)
	en, err := codec.Encode(msg, *cs)
	if err != nil {
		panic(err)
	}
	t.Log(en)
	//[1 48 14 66 105 99 104 111 110 80 114 111 118 105 110 99 101 101 0 135 0 0 0 0 0 0]
}

// TODO
func TestUserInformation(t *testing.T) {
	// 发送服务端包Index: 18
	// 发送服务端包信息: ServerPackets.UserInformation
	// 发送服务端包字节信息: 29, 2, 18, 0, 128, 3, 1, 0, 1, 0, 0, 0, 10, 104, 101, 108, 108, 111, 119, 111, 114, 108, 100, 0, 0, 255, 255, 255, 255, 0, 0, 1, 0, 28, 1, 0, 0, 97, 2, 0, 0, 6, 5, 17, 0, 14, 0, 0, 0, 0, 0, 0, 0, 0, 0, 100, 0, 0, 0, 0, 0, 0, 0, 0, 1, 46, 0, 0, 0, 1, 3, 0, 0, 0, 0, 0, 0, 0, 146, 2, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 255, 255, 255, 255, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 255, 255, 255, 255, 0, 0, 0, 0, 0, 0, 0, 1, 6, 0, 0, 0, 0, 0, 0, 0, 235, 4, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 255, 255, 255, 255, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 255, 255, 255, 255, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 14, 0, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 221, 0, 0, 0, 160, 15, 160, 15, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 255, 255, 255, 255, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 255, 255, 255, 255, 0, 0, 1, 2, 0, 0, 0, 0, 0, 0, 0, 61, 1, 0, 0, 136, 19, 136, 19, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 255, 255, 255, 255, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 255, 255, 255, 255, 0, 0, 0, 1, 4, 0, 0, 0, 0, 0, 0, 0, 210, 2, 0, 0, 49, 31, 64, 31, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 255, 255, 255, 255, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 255, 255, 255, 255, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 40, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 99, 0,
	bytes := []byte{
		128, 3, 1, 0, // ObjectID
		1, 0, 0, 0, // RealId
		10, 104, 101, 108, 108, 111, 119, 111, 114, 108, 100,
		0,                  // GuildName
		0,                  // GuildRank
		255, 255, 255, 255, // Color	uint32	[4]byte{225, 225, 225, 225}
		0,    // MirClass
		0,    // MirGender
		1, 0, // Level ushort uint16
		28, 1, 0, 0, // Point.X uint32
		97, 2, 0, 0, // Point.Y uint32
		6,	// Direction
		5,	// Hair
		17, 0,	// HP
		14, 0,	// MP
		0, 0, 0, 0, 0, 0, 0, 0, // Experience long int42
		100, 0, 0, 0, 0, 0, 0, 0, // MaxExperience
		0, // LevelEffects
		1, // ReadBoolean !!!!!!!
		46, 0, 0, 0, 1, 3, 0, 0, 0, 0, 0, 0, 0, 146, 2, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 255, 255, 255, 255, 0, 0, 0, 0, 0, 0, 0, 1, 6, 0, 0, 0, 0, 0, 0, 0, 235, 4, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 255, 255, 255, 255, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 255, 255, 255, 255, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 14, 0, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 221, 0, 0, 0, 160, 15, 160, 15, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 255, 255, 255, 255, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 255, 255, 255, 255, 0, 0, 1, 2, 0, 0, 0, 0, 0, 0, 0, 61, 1, 0, 0, 136, 19, 136, 19, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 255, 255, 255, 255, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 255, 255, 255, 255, 0, 0, 0, 1, 4, 0, 0, 0, 0, 0, 0, 0, 210, 2, 0, 0, 49, 31, 64, 31, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 255, 255, 255, 255, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 255, 255, 255, 255, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 40, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 99, 0}
	t.Log(bytes)

	msg := new(server.UserInformation)
	codec := new(MirUserInformationCodec)
	if err := codec.Decode(bytes, msg); err != nil {
		panic(err)
	}
	t.Log(msg)
}

func TestStartGame(t *testing.T) {
	codec := new(MirCodec)
	b := []byte{4, 0, 4, 0, 0}
	s1 := new(server.StartGame)
	if err := codec.Decode(b, s1); err != nil {
		panic(err)
	}
	t.Log("decode:", s1)

	s2 := new(server.StartGame)
	s2.Result = 4
	s2.Resolution = 1024
	ctx := new(cellnet.ContextSet)
	bytes, _ := codec.Encode(s2, *ctx)
	t.Log("encode:", bytes)
}

func TestSetConcentration(t *testing.T) {
	codec := new(MirCodec)
	bytes := []byte{128, 3, 1, 0, 0, 0}
	obj := new(server.SetConcentration)
	codec.Decode(bytes, obj)
	t.Log(obj)

	obj2 := new(server.SetConcentration)
	obj2.ObjectID = 66432
	obj2.Enabled = false
	obj2.Interrupted = false
	ctx := new(cellnet.ContextSet)
	bytes2, _ := codec.Encode(obj2, *ctx)
	t.Log(bytes2)
}

func TestEncodeLoginSuccess(t *testing.T) {
	codec := new(MirCodec)
	ctx := new(cellnet.ContextSet)

	res := new(server.LoginSuccess)

	c1 := new(common.SelectInfo)
	c1.Name = "测试登陆1"
	c1.Index = 1
	c1.Gender = common.MirGenderFemale
	c1.Class = common.MirClassArcher
	res.Characters = append(res.Characters, *c1)

	c2 := new(common.SelectInfo)
	c2.Name = "测试登陆2"
	c2.Index = 2
	c2.Gender = common.MirGenderFemale
	c2.Class = common.MirClassAssassin
	res.Characters = append(res.Characters, *c2)

	bytes, _ := codec.Encode(res, *ctx)
	t.Log(bytes)
}

func TestDecodeLoginSuccess(t *testing.T) {
	codec := new(MirCodec)
	bytes1 := []byte{2, 0, 0, 0, 1, 0, 0, 0, 13, 230, 181, 139, 232, 175, 149, 231, 153, 187, 233, 153, 134, 49, 0, 0, 4, 1, 0, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 13, 230, 181, 139, 232, 175, 149, 231, 153, 187, 233, 153, 134, 50, 0, 0, 3, 1, 0, 0, 0, 0, 0, 0, 0, 0}
	ls1 := new(server.LoginSuccess)
	codec.Decode(bytes1, ls1)
	t.Log(ls1)
}