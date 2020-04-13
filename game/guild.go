package game

import (
	"container/list"

	"github.com/yenkeia/mirgo/game/cm"
)

// Guild 对应 C# 里的 GuildObject, 行会
type Guild struct {
	GuildIndex  int
	Name        string     // 行会名称
	Ranks       []*cm.Rank // 行会所有职务(头衔)
	Notice      []string   // 行会公告
	Membercount int        // 行会人数
	NeedSave    bool
	FlagColour  cm.Color
}

func NewGuild(player *Player, name string) *Guild {
	g := &Guild{}
	leader := &cm.GuildMember{
		Name:      player.Name,
		PlayerID:  player.ID,
		LastLogin: 0,
		Online:    true,
	}
	rank := &cm.Rank{
		Name:    "会长", // Leader
		Options: cm.RankOptions(255),
		Index:   0,
	}
	rank.Members = []*cm.GuildMember{leader}
	g.Ranks = []*cm.Rank{rank}
	g.Membercount++
	g.NeedSave = true
	/* FIXME
	if (Level < Settings.Guild_ExperienceList.Count)
		MaxExperience = Settings.Guild_ExperienceList[Level];
	if (Level < Settings.Guild_MembercapList.Count)
		MemberCap = Settings.Guild_MembercapList[Level];
	*/
	g.FlagColour = cm.ColorWhite
	return g
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
	var member *cm.GuildMember
	var memberRank *cm.Rank
	if (kicker.MyGuild.GuildIndex != g.GuildIndex) || (kicker.MyGuildRank == nil) {
		return false
	}
	for i := 0; i < len(g.Ranks); i++ {
		for j := 0; j < len(g.Ranks[i].Members); j++ {
			if g.Ranks[i].Members[j].Name == membername {
				member = g.Ranks[i].Members[j]
				memberRank = g.Ranks[i]
				goto FOUND
			}
		}
	}
FOUND:
	if member == nil {
		return false
	}
	if ((kicker.MyGuildRank.Index >= memberRank.Index) && (kicker.MyGuildRank.Index != 0)) && (kicker.Name != membername) {
		kicker.ReceiveChat("你的职位权限不够。", cm.ChatTypeSystem)
		return false
	}
	if memberRank.Index == 0 {
		if len(memberRank.Members) < 2 {
			kicker.ReceiveChat("你不能离开行会，因为你是领导者。", cm.ChatTypeSystem)
			return false
		}
		for i := 0; i < len(memberRank.Members); i++ {
			if (memberRank.Members[i].Online) && (memberRank.Members[i].PlayerID != member.PlayerID) {
				goto AllOk
			}
		}
		kicker.ReceiveChat("需要至少1个行会领导在线。", cm.ChatTypeSystem)
		return false
	}
AllOk:
	// FIXME
	// g.MemberDeleted(membername, member.PlayerID, member.Name == kicker.Name)
	if member.PlayerID != 0 {
		leavingMember := env.Players.GetPlayerByID(member.PlayerID)
		leavingMember.RefreshStats()
	}
	// FIXME
	// memberRank.Members.Remove(Member)
	g.NeedSave = true
	g.Membercount--
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

func (g *Guild) NewNotice(notes []string) {

}

func (g *Guild) HasRoom() bool {
	return true
}

func (g *Guild) NewMember(newmember *Player) {

}

func (g *Guild) FindRank(name string) *cm.Rank {
	return nil
}

type GuildList struct {
	List *list.List
}

func NewGuildList() *GuildList {
	return nil
}

func (l *GuildList) Add(guild *Guild) {

}
