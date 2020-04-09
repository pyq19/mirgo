package util

import (
	"encoding/binary"
	"math"
)

// ReadString 根据传入的索引 返回读完string后所在bytes 的下一个索引 及string
func ReadString(bytes []byte, index int) (int, string) {
	//[0, 1, 2, 5, 9, 10, 11, 12, 50, 23, 77, 99]
	if len(bytes) == 0 {
		return -1, ""
	}
	strLen := int(bytes[index])
	lastIndex := index + strLen + 1
	msg := string(bytes[index+1 : lastIndex])
	return lastIndex, msg
}

// StringToBytes ...
func StringToBytes(s string) []byte {
	// 根据传入的string 返回字节数组，索引为0（即数组开头）的值为数组的长度
	if s == "" {
		return []byte{0}
	}
	stringBytes := []byte(s)
	stringBytesLenBytes := []byte{byte(len(stringBytes))}
	return append(stringBytesLenBytes, stringBytes...)
}

// Uint16ToBytes ...
func Uint16ToBytes(num uint16) []byte {
	res := make([]byte, 2)
	binary.LittleEndian.PutUint16(res, uint16(num))
	return res
}

// Uint32ToBytes ...
func Uint32ToBytes(num uint32) []byte {
	res := make([]byte, 4)
	binary.LittleEndian.PutUint32(res, uint32(num))
	return res
}

// Uint64ToBytes ...
func Uint64ToBytes(num uint64) []byte {
	res := make([]byte, 8)
	binary.LittleEndian.PutUint64(res, uint64(num))
	return res
}

// BoolToBytes ...
func BoolToBytes(b bool) []byte {
	if b == true {
		return []byte{1}
	}
	return []byte{0}
}

// BytesToUint16 ...
func BytesToUint16(bytes []byte) uint16 {
	return binary.LittleEndian.Uint16(bytes)
}

// BytesToUint32 ...
func BytesToUint32(bytes []byte) uint32 {
	return binary.LittleEndian.Uint32(bytes)
}

// BytesToUint64 ...
func BytesToUint64(bytes []byte) uint64 {
	return binary.LittleEndian.Uint64(bytes)
}

func Float32ToBytes(float float32) []byte {
	bits := math.Float32bits(float)
	bytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(bytes, bits)
	return bytes
}

func BytesToFloat32(bytes []byte) float32 {
	bits := binary.LittleEndian.Uint32(bytes)
	return math.Float32frombits(bits)
}

func Float64ToBytes(float float64) []byte {
	bits := math.Float64bits(float)
	bytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(bytes, bits)
	return bytes
}

func BytesToFloat64(bytes []byte) float64 {
	bits := binary.LittleEndian.Uint64(bytes)
	return math.Float64frombits(bits)
}
