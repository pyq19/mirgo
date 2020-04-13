package cm

// 行会相关的结构体, 因为 server/smsg.go 会用到，所以放在 cm 包下

// Rank 行会职务(头衔)
type Rank struct {
	Options RankOptions
}
