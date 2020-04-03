package mir

// PlayerList 玩家列表
type PlayerList struct {
	List []*Player
}

// NewPlayerList 新建列表
func NewPlayerList() *PlayerList {
	return &PlayerList{List: make([]*Player, 0)}
}

// Get 获取 PlayerList 里索引为 i 的玩家
func (l *PlayerList) Get(i int) *Player {
	return l.List[i]
}

// Add 把玩家加到 PlayerList
func (l *PlayerList) Add(p *Player) {
	l.List = append(l.List, p)
}

// Count 玩家数量
func (l *PlayerList) Count() int {
	return len(l.List)
}

// Remove 从列表里删除玩家
func (l *PlayerList) Remove(p *Player) {
	for k, v := range l.List {
		if v.GetID() == p.GetID() {
			l.List = append(l.List[:k], l.List[k+1:]...)
		}
	}
}
