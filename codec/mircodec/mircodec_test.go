package mircodec

import (
	"github.com/davyxu/cellnet"
	"github.com/jinzhu/gorm"
	"github.com/yenkeia/mirgo/common"
	"github.com/yenkeia/mirgo/proto/server"
	"os"
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

	//saveToDB(msg)

	obj, err := codec.Encode(msg, *new(cellnet.ContextSet))
	if err != nil {
		t.Error(err)
	}
	t.Log(obj)
}

func saveToDB(msg *server.UserInformation) {
	gopath := os.Getenv("GOPATH")
	db, err := gorm.Open("sqlite3", gopath+mirDB)
	if err != nil {
		panic("failed to connect database")
	}
	//defer db.Close()

	i := msg.Inventory
	for _, v := range i {
		if v.ID == 0 {
			continue
		}
		db.Table("user_item").Create(v)
	}
	e := msg.Equipment
	for _, v := range e {
		if v.ID == 0 {
			continue
		}
		db.Table("user_item").Create(v)
	}
	q := msg.QuestInventory
	for _, v := range q {
		if v.ID == 0 {
			continue
		}
		db.Table("user_item").Create(v)
	}
}

func TestEmptySlice(t *testing.T) {
	slice := make([]common.UserItem, 5)
	slice[0].ItemID = 1
	t.Log(len(slice))
	t.Log(slice[0])
	for i := 0; i < len(slice); i++ {
		t.Log(IsNull(slice[i]))
	}
}

func TestDecodeUserInformation(t *testing.T) {
	bytes := []byte{
		1, 0, 0, 0,
		1, 0, 0, 0,
		6, 232, 140, 131, 233, 151, 178,
		0,
		0,
		255, 255, 255, 255,
		2, 0,
		20, 0,
		28, 1, 0, 0,
		96, 2, 0, 0,
		1,
		1,
		15, 0,
		17, 0,
		55, 0, 0, 0, 0, 0, 0, 0,
		100, 0, 0, 0, 0, 0, 0, 0,
		1,

		1,
		46, 0, 0, 0,
		1,
		3, 0, 0, 0, 0, 0, 0, 0, 146, 2, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 255, 255, 255, 255, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		1,
		4, 0, 0, 0, 0, 0, 0, 0, 210, 2, 0, 0, 54, 31, 64, 31, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 255, 255, 255, 255, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		1,
		5, 0, 0, 0, 0, 0, 0, 0, 235, 4, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 255, 255, 255, 255, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,

		1,
		14, 0, 0, 0,
		1,
		1, 0, 0, 0, 0, 0, 0, 0, 221, 0, 0, 0, 160, 15, 160, 15, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 255, 255, 255, 255, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		1,
		2, 0, 0, 0, 0, 0, 0, 0, 62, 1, 0, 0, 136, 19, 136, 19, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 255, 255, 255, 255, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,

		1,
		40, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,

		100, 0, 0, 0,
		100, 0, 0, 0}
	msg := new(server.UserInformation)
	codec := new(MirUserInformationCodec)
	if err := codec.Decode(bytes, msg); err != nil {
		panic(err)
	}
	t.Log(msg)
}

func TestDecodeUserInformation2(t *testing.T) {
	bytes := []byte{128, 3, 1, 0, 2, 0, 0, 0, 5, 116, 101, 115, 116, 49, 0, 0, 255, 255, 255, 255, 1, 0, 1, 0, 29, 1, 0, 0, 99, 2, 0, 0, 7, 5, 15, 0, 17, 0, 0, 0, 0, 0, 0, 0, 0, 0, 100, 0, 0, 0, 0, 0, 0, 0, 0, 1, 46, 0, 0, 0, 1, 8, 0, 0, 0, 0, 0, 0, 0, 146, 2, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 255, 255, 255, 255, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 6, 0, 0, 0, 0, 0, 0, 0, 221, 0, 0, 0, 160, 15, 160, 15, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 255, 255, 255, 255, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 7, 0, 0, 0, 0, 0, 0, 0, 61, 1, 0, 0, 136, 19, 136, 19, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 255, 255, 255, 255, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 9, 0, 0, 0, 0, 0, 0, 0, 210, 2, 0, 0, 64, 31, 64, 31, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 255, 255, 255, 255, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 10, 0, 0, 0, 0, 0, 0, 0, 235, 4, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 255, 255, 255, 255, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 14, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 40, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	msg := new(server.UserInformation)
	codec := new(MirUserInformationCodec)
	if err := codec.Decode(bytes, msg); err != nil {
		panic(err)
	}
	t.Log(msg)
}

func TestDecodeEncodeObjectPlayer(t *testing.T) {
	bytes := []byte{1, 0, 0, 0,
		6, 232, 140, 131, 233, 151, 178,
		0, 0,
		255, 255, 255, 255,
		0, 0,
		20, 0,
		28, 1, 0, 0, 96, 2, 0, 0,
		1, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0} // TODO
	msg := new(server.ObjectPlayer)
	codec := new(MirObjectPlayerCodec)
	if err := codec.Decode(bytes, msg); err != nil {
		panic(err)
	}
	t.Log(msg)

	obj, err := codec.Encode(msg, *new(cellnet.ContextSet))
	if err != nil {
		t.Error(err)
	}
	t.Log(obj)
	//[1 0 0 0
	//6 232 140 131 233 151 178
	//0 0
	//255 255 255 255
	//0 0
	//20 0
	//28 1 0 0 96 2 0 0
	//1 0 0 0 0 0 0 1 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0]
}

func TestEncodeDecodeNPCResponse(t *testing.T) {
	bytes := []byte{2, 0, 0, 0, 6, 232, 140, 131, 233, 151, 178, 6, 232, 140, 131, 233, 151, 178}
	msg := new(server.NPCResponse)
	codec := new(MirCodec)
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

type SliceTest struct {
	Nums []int
}

func TestEncodeDecodeSlice(t *testing.T) {
	// []int
	//bytes := []byte{2, 0, 0, 0, 3, 0, 0, 0, 4, 0, 0, 0}
	bytes := []byte{}
	msg := new(SliceTest)
	codec := new(MirCodec)
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

func TestEncodeDecodePlayerInspect(t *testing.T) {
	codec := new(MirPlayerInspectCodec)
	msg := &server.PlayerInspect{
		Name:      "testName",
		GuildName: "testGuildName",
		GuildRank: "testGuildRank",
		Equipment: make([]common.UserItem, 14),
		Class:     common.MirClassTaoist,
		Gender:    common.MirGenderFemale,
		Hair:      1,
		Level:     10,
		LoverName: "testLoverName",
	}
	obj, err := codec.Encode(msg, *new(cellnet.ContextSet))
	if err != nil {
		t.Error(err)
	}
	bytes := obj.([]byte)
	t.Log(bytes)

	if err := codec.Decode(bytes, msg); err != nil {
		panic(err)
	}
	t.Log(msg)
}

func TestDecodeEncodeObjectMonster(t *testing.T) {
	bytes := []byte{
		34, 6, 0, 0,
		4, 68, 101, 101, 114,
		255, 255, 255, 255,
		24, 1, 0, 0, 86, 2, 0, 0,
		4, 0, 1, 0, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	codec := new(MirCodec)
	msg := &server.ObjectMonster{}
	if err := codec.Decode(bytes, msg); err != nil {
		panic(err)
	}
	t.Log(msg)

	obj, err := codec.Encode(msg, *new(cellnet.ContextSet))
	if err != nil {
		t.Error(err)
	}
	t.Log(obj.([]byte))
	//34 6 0 0
	//4 68 101 101 114
	//255 255 255 255
	//24 1 0 0 86 2 0 0
	//4 0 1 0 2 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0]
}

func TestDecodeEncodePlayerUpdate(t *testing.T) {
	bytes := []byte{129, 3, 1, 0, 0, 0, 0, 0, 0, 1, 0, 0}
	codec := new(MirCodec)
	msg := &server.PlayerUpdate{}
	if err := codec.Decode(bytes, msg); err != nil {
		panic(err)
	}
	t.Log(msg)

	obj, err := codec.Encode(msg, *new(cellnet.ContextSet))
	if err != nil {
		t.Error(err)
	}
	t.Log(obj.([]byte))
}

func TestDecodeEncodeObjectNPC(t *testing.T) {
	//bytes := []byte{1, 0, 0, 0, 16, 84, 101, 108, 101, 112, 111, 114, 116, 95, 71, 105, 108, 98, 101, 114, 116, 0, 255, 0, 255, 15, 0, 0, 0, 0, 0, 31, 1, 0, 0, 103, 2, 0, 0, 1, 0, 0, 0, 0}
	bytes := []byte{5, 0, 0, 0, 14, 77, 101, 114, 99, 104, 97, 110, 116, 95, 83, 109, 105, 116, 104, 0, 255, 0, 255, 0, 0, 0, 0, 0, 0, 41, 1, 0, 0, 100, 2, 0, 0, 2, 2, 0, 0, 0, 5, 0, 0, 0, 6, 0, 0, 0}
	codec := new(MirObjectNPCCodec)
	msg := &server.ObjectNPC{}
	if err := codec.Decode(bytes, msg); err != nil {
		panic(err)
	}
	t.Log(msg)
	obj, err := codec.Encode(msg, *new(cellnet.ContextSet))
	if err != nil {
		t.Error(err)
	}
	t.Log(obj.([]byte))
	//[1 0 0 0 16 84 101 108 101 112 111 114 116 95 71 105 108 98 101 114 116 0 255 0 255 15 0 0 0 0 0 31 1 0 0 103 2 0 0 1 0 0 0 0]
	//[5 0 0 0 14 77 101 114 99 104 97 110 116 95 83 109 105 116 104 0 255 0 255 0 0 0 0 0 0 41 1 0 0 100 2 0 0 2 2 0 0 0 5 0 0 0 6 0 0 0]
}

func TestDecodeEncodeNPCResponse(t *testing.T) {
	bytes := []byte{0, 0, 0, 0}
	codec := new(MirCodec)
	msg := &server.NPCResponse{}
	if err := codec.Decode(bytes, msg); err != nil {
		panic(err)
	}
	t.Log(msg)
	obj, err := codec.Encode(msg, *new(cellnet.ContextSet))
	if err != nil {
		t.Error(err)
	}
	t.Log(obj.([]byte))
}

// mirgo
//(server.OBJECT_NPC) [48, 0, 77, 0, 25, 142, 1, 0, 16, 84, 101, 108, 101, 112, 111, 114, 116, 95, 71, 105, 108, 98, 101, 114, 116, 255, 255, 255, 255, 15, 0, 255, 255, 255, 255, 31, 1, 0, 0, 103, 2, 0, 0, 4, 0, 0, 0, 0]
//(server.OBJECT_NPC) [47, 0, 77, 0, 28, 142, 1, 0, 15, 67, 114, 97, 102, 116, 115, 76, 97, 100, 121, 95, 74, 117, 100, 101, 255, 255, 255, 255, 7, 0, 255, 255, 255, 255, 38, 1, 0, 0, 107, 2, 0, 0, 4, 0, 0, 0, 0]
//(server.OBJECT_NPC) [45, 0, 77, 0, 30, 142, 1, 0, 13, 77, 101, 114, 99, 104, 97, 110, 116, 95, 74, 111, 104, 110, 255, 255, 255, 255, 11, 0, 255, 255, 255, 255, 37, 1, 0, 0, 91, 2, 0, 0, 4, 0, 0, 0, 0]
//(server.OBJECT_NPC) [51, 0, 77, 0, 26, 142, 1, 0, 19, 66, 111, 114, 100, 101, 114, 86, 105, 108, 108, 97, 103, 101, 95, 66, 111, 97, 114, 100, 255, 255, 255, 255, 45, 0, 255, 255, 255, 255, 28, 1, 0, 0, 103, 2, 0, 0, 4, 0, 0, 0, 0]
//(server.OBJECT_NPC) [46, 0, 77, 0, 27, 142, 1, 0, 14, 65, 115, 115, 105, 115, 116, 97, 110, 116, 95, 74, 97, 110, 101, 255, 255, 255, 255, 5, 0, 255, 255, 255, 255, 28, 1, 0, 0, 94, 2, 0, 0, 4, 0, 0, 0, 0]
//(server.OBJECT_NPC) [46, 0, 77, 0, 29, 142, 1, 0, 14, 77, 101, 114, 99, 104, 97, 110, 116, 95, 83, 109, 105, 116, 104, 255, 255, 255, 255, 0, 0, 255, 255, 255, 255, 41, 1, 0, 0, 100, 2, 0, 0, 4, 0, 0, 0, 0]
//(server.OBJECT_NPC) [46, 0, 77, 0, 32, 142, 1, 0, 14, 77, 101, 114, 99, 104, 97, 110, 116, 95, 82, 117, 98, 101, 110, 255, 255, 255, 255, 1, 0, 255, 255, 255, 255, 35, 1, 0, 0, 98, 2, 0, 0, 4, 0, 0, 0, 0]
//(server.OBJECT_NPC) [54, 0, 77, 0, 65, 142, 1, 0, 22, 67, 114, 97, 102, 116, 105, 110, 103, 86, 105, 108, 108, 97, 103, 101, 95, 80, 111, 114, 116, 97, 108, 255, 255, 255, 255, 0, 0, 255, 255, 255, 255, 40, 1, 0, 0, 107, 2, 0, 0, 4, 0, 0, 0, 0]
//(server.OBJECT_NPC) [57, 0, 77, 0, 66, 142, 1, 0, 25, 84, 114, 97, 118, 101, 108, 108, 105, 110, 103, 77, 101, 114, 99, 104, 97, 110, 116, 95, 68, 97, 109, 105, 97, 110, 255, 255, 255, 255, 94, 0, 255, 255, 255, 255, 31, 1, 0, 0, 94, 2, 0, 0, 4, 0, 0, 0, 0]
// C#
//[2020/1/15 23:49:52]: --->发送服务端包信息: ServerPackets.ObjectNPC
//[2020/1/15 23:49:52]: --->发送服务端包字节信息: 48, 0, 77, 0, 7, 0, 0, 0, 16, 77, 101, 114, 99, 104, 97, 110, 116, 95, 87, 104, 105, 116, 110, 101, 121, 0, 255, 0, 255, 7, 0, 0, 0, 0, 0, 49, 1, 0, 0, 95, 2, 0, 0, 0, 0, 0, 0, 0,
//[2020/1/15 23:49:52]: --->发送服务端包信息: ServerPackets.ObjectNPC
//[2020/1/15 23:49:52]: --->发送服务端包字节信息: 50, 0, 77, 0, 8, 0, 0, 0, 14, 77, 101, 114, 99, 104, 97, 110, 116, 95, 82, 117, 98, 101, 110, 0, 255, 0, 255, 1, 0, 0, 0, 0, 0, 35, 1, 0, 0, 98, 2, 0, 0, 0, 1, 0, 0, 0, 48, 0, 0, 0,

func TestDecodeEncodeObjectNPC_2(t *testing.T) {
	//[2020/1/15 23:49:52]: --->发送服务端包字节信息: 50, 0, 77, 0, 8, 0, 0, 0, 14, 77, 101, 114, 99, 104, 97, 110, 116, 95, 82, 117, 98, 101, 110, 0, 255, 0, 255, 1, 0, 0, 0, 0, 0, 35, 1, 0, 0, 98, 2, 0, 0, 0, 1, 0, 0, 0, 48, 0, 0, 0,
	bytes := []byte{8, 0, 0, 0, 14, 77, 101, 114, 99, 104, 97, 110, 116, 95, 82, 117, 98, 101, 110, 0, 255, 0, 255, 1, 0, 0, 0, 0, 0, 35, 1, 0, 0, 98, 2, 0, 0, 0, 1, 0, 0, 0, 48, 0, 0, 0}
	codec := new(MirObjectNPCCodec)
	msg := &server.ObjectNPC{}
	if err := codec.Decode(bytes, msg); err != nil {
		panic(err)
	}
	//ObjectNPC: ID(8) Name(Merchant_Ruben) NameColor(-16711936) Image(1) Color(0) Location(291,610) Direction(0) QuestIDs([48])
	t.Log(msg)
	obj, err := codec.Encode(msg, *new(cellnet.ContextSet))
	if err != nil {
		t.Error(err)
	}
	t.Log(obj.([]byte))
}

func TestDecodeEncodeObjectNPC_3(t *testing.T) {
	//(server.OBJECT_NPC) [46, 0, 77, 0, 32, 142, 1, 0, 14, 77, 101, 114, 99, 104, 97, 110, 116, 95, 82, 117, 98, 101, 110, 255, 255, 255, 255, 1, 0, 255, 255, 255, 255, 35, 1, 0, 0, 98, 2, 0, 0, 4, 0, 0, 0, 0]
	bytes := []byte{32, 142, 1, 0, 14, 77, 101, 114, 99, 104, 97, 110, 116, 95, 82, 117, 98, 101, 110, 255, 255, 255, 255, 1, 0, 255, 255, 255, 255, 35, 1, 0, 0, 98, 2, 0, 0, 4, 0, 0, 0, 0}
	codec := new(MirObjectNPCCodec)
	msg := &server.ObjectNPC{}
	if err := codec.Decode(bytes, msg); err != nil {
		panic(err)
	}
	t.Log(msg)
	//ObjectNPC: ID(101920) Name(Merchant_Ruben) NameColor(-1) Image(1) Color(-1) Location(291,610) Direction(4) QuestIDs([])
	obj, err := codec.Encode(msg, *new(cellnet.ContextSet))
	if err != nil {
		t.Error(err)
	}
	t.Log(obj.([]byte))
}

func TestEncodeCharacter(t *testing.T) {
	path := os.Getenv("GOPATH") + "/src/github.com/yenkeia/mirgo/dotnettools/mir.sqlite"
	db, err := gorm.Open("sqlite3", path)
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	c := new(common.Character)
	db.Table("character").Where("id = ?", 1).Find(&c)
	t.Log(c)

	codec := new(MirCodec)
	obj, err := codec.Encode(c, *new(cellnet.ContextSet))
	if err != nil {
		t.Error(err)
	}
	t.Log(obj.([]byte))
}
