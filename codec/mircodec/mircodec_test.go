package mircodec

import (
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
