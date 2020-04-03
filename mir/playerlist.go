package mir

import "container/list"

// PlayerList 玩家列表
type PlayerList struct {
	List *list.List
}

// NewPlayerList 新建列表
func NewPlayerList() *PlayerList {
	return nil
}

// Get 获取 PlayerList 里索引为 i 的玩家
func (l *PlayerList) Get(i int) *Player {
	return nil
}

// Add 把玩家加到 PlayerList
func (l *PlayerList) Add(p *Player) {

}

// Count 玩家数量
func (l *PlayerList) Count() int {
	return l.List.Len()
}

// Remove 从列表里删除玩家
func (l *PlayerList) Remove(p *Player) {

}
