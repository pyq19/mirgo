package mircodec

import (
	"github.com/davyxu/cellnet"
	"github.com/yenkeia/mirgo/common"
	"github.com/yenkeia/mirgo/proto/server"
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

func TestEncodeDecode(t *testing.T) {
	rawCatShop := CatShop{
		"独山路宠物店111",
		123,
		CatInfo{
			"花猫2222",
			555,
		},
	}
	res, err := encode(&rawCatShop)
	t.Log("before:", res)

	newCatShop := new(CatShop)
	err = decode(i(newCatShop), res)
	if err != nil {
		panic(err)
	}

	//t.Log(catShop.Address)
	t.Log(newCatShop.Address, newCatShop.Rank)
	t.Log(newCatShop.Cat.Name, newCatShop.Cat.Age)

	t.Log("after:", res)
}

// TestDecode 把字节数组解码(反序列化)成结构体
func TestDecode(t *testing.T) {
	bytes := []byte{18, 231, 139, 172, 229, 177, 177, 232, 183, 175, 229, 174, 160, 231, 137, 169, 229, 186, 151,
		111,
		6, 232, 138, 177, 231, 140, 171,
		172, 8, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	catShop := new(CatShop)
	decode(i(catShop), bytes)
	//t.Log(catShop.Address)
	t.Log(catShop.Cat.Name, catShop.Cat.Age)

	t.Log(bytes)
}

func TestByteSlice(t *testing.T) {
	bytes := []byte{5, 11, 1, 66, 3, 63, 98}
	t.Log(bytes[3])
	t.Log(bytes[:2]) // [5 11]
	t.Log(bytes[2:]) // [1 66 3 63 98]
	t.Log(bytes[3:]) // [66 3 63 98]
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
	res, err := encode(&catShop)
	if err != nil {
		panic(err)
	}
	t.Log("->", res)
}

type Book struct {
	Name string
	Page uint8
}

type BookStore struct {
	Books []Book
}

func TestEncodeSlice(t *testing.T) {
	b1 := new(Book)
	b1.Name = "庆余年"
	b1.Page = 123
	b2 := new(Book)
	b2.Name = "红楼梦"
	b2.Page = 32
	bs := new(BookStore)
	bs.Books = append(bs.Books, *b1)
	bs.Books = append(bs.Books, *b2)
	t.Log(bs) // &{[{庆余年 123} {红楼梦 32}]}

	bytes, err := encode(i(bs))
	if err != nil {
		panic(err)
	}
	t.Log(bytes)
	//[2 0 0 0 9 229 186 134 228 189 153 229 185 180 123 9 231 186 162 230 165 188 230 162 166 32]
}

func TestDecodeSlice(t *testing.T) {
	// &{[{庆余年 123} {红楼梦 32}]}
	bytes := []byte{2, 0, 0, 0, 9, 229, 186, 134, 228, 189, 153, 229, 185, 180, 123, 9, 231, 186, 162, 230, 165, 188, 230, 162, 166, 32}
	t.Log("before:", bytes)
	bs := i(new(BookStore))

	if err := decode(bs, bytes); err != nil {
		panic(err)
	}
	t.Log("after:", bytes)
	t.Log(bs)
}

type ClientVersion struct {
	Hash []uint8
}

func TestDecodeSlice2(t *testing.T) {
	cv := new(ClientVersion)
	cv.Hash = []byte{5, 0, 0, 0, 22, 33, 11, 55, 100}
	bytes, err := encode(i(cv))
	if err != nil {
		panic(err)
	}
	t.Log("before:", bytes)

	cv2 := new(ClientVersion)
	if err := decode(cv2, bytes); err != nil {
		panic(err)
	}
	t.Log("after:", bytes)
	t.Log(cv2)
}

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

func newLoginSuccessStruct() *server.LoginSuccess {
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
	return res
}

func TestEncodeLoginSuccess(t *testing.T) {
	codec := new(MirCodec)
	ctx := new(cellnet.ContextSet)

	res := newLoginSuccessStruct()

	bytes, _ := codec.Encode(res, *ctx)
	t.Log(bytes.([]byte))

	t.Log(len(bytes.([]byte)))
}

func TestDecodeLoginSuccess(t *testing.T) {
	codec := new(MirCodec)
	bytes1 := []byte{2, 0, 0, 0, 1, 0, 0, 0, 13, 230, 181, 139, 232, 175, 149, 231, 153, 187, 233, 153, 134, 49, 0, 0, 4, 1, 0, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 13, 230, 181, 139, 232, 175, 149, 231, 153, 187, 233, 153, 134, 50, 0, 0, 3, 1, 0, 0, 0, 0, 0, 0, 0, 0}
	t.Log(len(bytes1))
	ls1 := new(server.LoginSuccess)
	codec.Decode(bytes1, ls1)
	t.Log(ls1)
	//&{[{1 测试登陆1 0 4 1 0} {2 测试登陆2 0 3 1 0}]}
}

func TestDecodeEncodeLoginSuccess(t *testing.T) {
	codec := new(MirCodec)

	ls2 := newLoginSuccessStruct()
	bytes, err := codec.Encode(ls2, *new(cellnet.ContextSet))
	if err != nil {
		panic(err)
	}
	t.Log(bytes)

	ls1 := new(server.LoginSuccess)
	codec.Decode(bytes, ls1)
	t.Log(ls1)
	//&{[{1 测试登陆1 0 4 1 0} {2 测试登陆2 0 3 1 0}]}
}

func TestUserInformation(t *testing.T) {
	// --->发送服务端包字节信息: 157, 1, 18, 0, 128, 3, 1, 0, 1, 0, 0, 0, 6, 99, 99, 99, 99, 99, 99, 0, 0, 255, 255, 255, 255, 1, 1, 1, 0, 28, 1, 0, 0, 96, 2, 0, 0, 1, 8, 15, 0, 17, 0, 0, 0, 0, 0, 0, 0, 0, 0, 100, 0, 0, 0, 0, 0, 0, 0, 0, 1, 46, 0, 0, 0, 1, 3, 0, 0, 0, 0, 0, 0, 0, 146, 2, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 255, 255, 255, 255, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 5, 0, 0, 0, 0, 0, 0, 0, 235, 4, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 255, 255, 255, 255, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 4, 0, 0, 0, 0, 0, 0, 0, 210, 2, 0, 0, 54, 31, 64, 31, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 255, 255, 255, 255, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 14, 0, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 221, 0, 0, 0, 160, 15, 160, 15, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 255, 255, 255, 255, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 2, 0, 0, 0, 0, 0, 0, 0, 62, 1, 0, 0, 136, 19, 136, 19, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 255, 255, 255, 255, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 40, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	bytes := []byte{
		128, 3, 1, 0, // ObjectID
		1, 0, 0, 0, // RealId
		6, 99, 99, 99, 99, 99, 99, // Name
		0,
		0,
		255, 255, 255, 255,
		1, 1, // class, gender
		1, 0, // level
		28, 1, 0, 0, // location x
		96, 2, 0, 0, // location y
		1,     // direction
		8,     // hair
		15, 0, // hp
		17, 0, // mp
		0, 0, 0, 0, 0, 0, 0, 0, // Experience
		100, 0, 0, 0, 0, 0, 0, 0, // MaxExperience
		0, // level effect

		///// Inventory 开始

		1,           // !!!ReadBoolean
		46, 0, 0, 0, // Inventory count
		1, // ReadBoolean = true 第一个 UserItem
		// UserItem 开始
		3, 0, 0, 0, 0, 0, 0, 0, // UniqueId
		146, 2, 0, 0, // ItemIndex
		0, 0,
		0, 0,
		1, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, // AC ~ Luck
		255, 255, 255, 255, // SoulBoundId
		0,                            // Bools
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, // strong ~ poisonAttack
		// UserItem 结束

		// ReadBoolean = false
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,

		1, // ReadBoolean = true
		// 第二个 UserItem 开始
		5, 0, 0, 0, 0, 0, 0, 0,
		235, 4, 0, 0,
		0, 0,
		0, 0,
		1, 0, 0, 0, // Count
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, // AC ~ Luck
		255, 255, 255, 255,
		0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, // strong ~ poisonAttack

		// ReadBoolean = false
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		1, // ReadBoolean = true
		// 第三个 UserItem 开始
		4, 0, 0, 0, 0, 0, 0, 0,
		210, 2, 0, 0,
		54, 31,
		64, 31,
		1, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		255, 255, 255, 255,
		0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0,

		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		////// Inventory 结束

		////// Equipment 开始
		1,           // ReadBoolean = true
		14, 0, 0, 0, // Equipment count
		1, // ReadBoolean = true
		// 第一个 Equipment UserItem
		1, 0, 0, 0, 0, 0, 0, 0,
		221, 0, 0, 0,
		160, 15,
		160, 15,
		1, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		255, 255, 255, 255,
		0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0,

		1,
		// 第二个 Equipment UserItem
		2, 0, 0, 0, 0, 0, 0, 0,
		62, 1, 0, 0,
		136, 19, 136, 19,
		1, 0,
		0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		255, 255, 255, 255,
		0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0,

		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		////// Equipment 结束

		////// QuestInventory 开始
		1,
		40, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		////// QuestInventory 结束

		0, 0, 0, 0, // Gold
		0, 0, 0, 0, // Credit
	}
	t.Log(bytes)

	msg := new(server.UserInformation)
	codec := new(MirUserInformationCodec)
	if err := codec.Decode(bytes, msg); err != nil {
		panic(err)
	}
	t.Log(msg)

	obj, err := codec.Encode(msg, *new(cellnet.ContextSet))
	if err != nil {
		t.Error(err)
	}
	t.Log(obj)
}

func TestEmptySlice(t *testing.T) {
	slice := make([]common.UserItem, 5)
	slice[0].ItemId = 1
	t.Log(len(slice))
	t.Log(slice[0])
	for i := 0; i < len(slice); i++ {
		t.Log(IsNull(slice[i]))
	}
}
