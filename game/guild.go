package game

type Guild struct{}

type GuildMember struct {
	Name      string
	LastLogin int64 // DateTime
	Hasvoted  bool
	Online    bool
	PlayerID  uint32 // ID
}

// Rank 行会职务(头衔)
type Rank struct {
	Name    string
	Index   int
	Options RankOptions
	Members []*GuildMember
}
