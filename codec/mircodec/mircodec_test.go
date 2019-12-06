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

	newCatShop := new(CatShop)
	err = decode(i(newCatShop), res)
	if err != nil {
		panic(err)
	}

	//t.Log(catShop.Address)
	t.Log(newCatShop.Address, newCatShop.Rank)
	t.Log(newCatShop.Cat.Name, newCatShop.Cat.Age)
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
