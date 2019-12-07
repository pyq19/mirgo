package client

import (
	"github.com/davyxu/cellnet"
	"github.com/yenkeia/mirgo/codec/mircodec"

	// 使用binary协议，因此匿名引用这个包，底层会自动注册
	"reflect"

	_ "github.com/davyxu/cellnet/codec/binary"
)

const (
	CLIENT_VERSION = 1000 + iota
	DISCONNECT
	KEEP_ALIVE
	NEW_ACCOUNT
	CHANGE_PASSWORD
	LOGIN
	NEW_CHARACTER
	DELETE_CHARACTER
	START_GAME
)

type ClientVersion struct {
	VersionHash []uint8
}

type Disconnect struct{}

type KeepAlive struct {
	Time int64
}

// TODO
type NewAccount struct {
	//public string AccountID = string.Empty
	//public string Password = string.Empty
	//public DateTime BirthDate
	//public string UserName = string.Empty
	//public string SecretQuestion = string.Empty
	//public string SecretAnswer = string.Empty
	//public string EMailAddress = string.Empty
}

type ChangePassword struct {
	AccountID       string
	CurrentPassword string
	NewPassword     string
}

type Login struct {
	AccountID string
	Password  string
}

// TODO
type NewCharacter struct {
	//public string Name = string.Empty;
	//public MirGender Gender;
	//public MirClass Class;
}

type DeleteCharacter struct {
	CharacterIndex int16
}

type StartGame struct {
	CharacterIndex int16
}

// 引用消息时，自动注册消息，这个文件可以由代码生成自动生成
func init() {

	mirCodec := new(mircodec.MirCodec)

	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*ClientVersion)(nil)).Elem(),
		ID:    CLIENT_VERSION,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*Disconnect)(nil)).Elem(),
		ID:    DISCONNECT,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*KeepAlive)(nil)).Elem(),
		ID:    KEEP_ALIVE,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*NewAccount)(nil)).Elem(),
		ID:    NEW_ACCOUNT,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*ChangePassword)(nil)).Elem(),
		ID:    CHANGE_PASSWORD,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*Login)(nil)).Elem(),
		ID:    LOGIN,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*NewCharacter)(nil)).Elem(),
		ID:    NEW_CHARACTER,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*DeleteCharacter)(nil)).Elem(),
		ID:    DELETE_CHARACTER,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*StartGame)(nil)).Elem(),
		ID:    START_GAME,
	})
}
