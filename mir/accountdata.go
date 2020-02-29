package mir

import "github.com/yenkeia/mirgo/common"

// AccountData 游戏运行时实时保存的数据
type AccountData struct {
	Account           []*common.Account
	AccountCharacter  []*common.AccountCharacter
	Character         []*common.Character
	CharacterUserItem []*common.CharacterUserItem
	UserItem          []*common.UserItem
	UserMagics        []*common.UserMagic
}
