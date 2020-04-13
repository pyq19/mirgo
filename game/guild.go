package game

import (
	"container/list"

	"github.com/yenkeia/mirgo/game/cm"
)

// Guild 对应 C# 里的 GuildObject, 工会
type Guild struct {
	ID    int     // GuildIndex
	Name  string  // 工会名称
	Ranks []*Rank // 工会所有职务(头衔)
}

func NewGuild(player *Player, name string) *Guild {
	return nil
}

// SendGuildStatus ...
func (g *Guild) SendGuildStatus(member *Player) {
	/*
		string gName = Name;
		string conquest = "";
		if (Conquest != null)
		{
		    conquest = "[" + Conquest.Info.Name + "]";
		    gName = gName + conquest;
		}

		member.Enqueue(new ServerPackets.GuildStatus()
		    {
		        GuildName = gName,
		        GuildRankName = member.MyGuildRank != null? member.MyGuildRank.Name: "",
		        Experience = Experience,
		        MaxExperience = MaxExperience,
		        MemberCount = Membercount,
		        MaxMembers = MemberCap,
		        Gold = Gold,
		        Level = Level,
		        Voting = Voting,
		        SparePoints = SparePoints,
		        ItemCount = (byte)StoredItems.Length,
		        BuffCount = (byte)0,//(byte)BuffList.Count,
		        MyOptions = member.MyGuildRank != null? member.MyGuildRank.Options: (RankOptions)0,
		        MyRankId = member.MyGuildRank != null? member.MyGuildRank.Index: 256
		    });
	*/
}

// IsAtWar 是否在行会战期间
func (g *Guild) IsAtWar() bool {
	return false
}

// DeleteMember 行会删除成员
func (g *Guild) DeleteMember(kicker *Player, membername string) bool {
	return true
}

// ChangeRank 改变行会成员头衔
func (g *Guild) ChangeRank(self *Player, membername string, rankIndex byte, rankName string) bool {
	return true
}

func (g *Guild) NewRank(self *Player) bool {
	return true
}

func (g *Guild) ChangeRankOption(self *Player, rankIndex byte, option int, enabled string) bool {
	return true
}

func (g *Guild) ChangeRankName(self *Player, rankName string, rankIndex byte) bool {
	return true
}

type GuildList struct {
	List *list.List
}

func NewGuildList() *GuildList {
	return nil
}

func (l *GuildList) Add(guild *Guild) {

}

// Rank 行会职务(头衔)
type Rank struct {
	Options cm.RankOptions
}
