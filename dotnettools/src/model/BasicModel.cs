using SqlSugar;

namespace dotnettools
{
    [SugarTable("basic")]
    public class BasicModel
    {
        public BasicModel() {}

        [SugarColumn(ColumnName = "id", IsPrimaryKey = true, IsIdentity = true)] //是主键, 还是标识列
        public int Id { get; set; }

        [SugarColumn(ColumnName = "game_version")]
        public int Version { get; set; }

        [SugarColumn(ColumnName = "custom_version")]
        public int CustomVersion { get; set; }

        [SugarColumn(ColumnName = "map_index")]
        public int MapIndex { get; set; }

        [SugarColumn(ColumnName = "item_index")]
        public int ItemIndex { get; set; }

        [SugarColumn(ColumnName = "monster_index")]
        public int MonsterIndex { get; set; }

        [SugarColumn(ColumnName = "npc_index")]
        public int NPCIndex { get; set; }

        [SugarColumn(ColumnName = "quest_index")]
        public int QuestIndex { get; set; }

        [SugarColumn(ColumnName = "gameshop_index")]
        public int GameshopIndex { get; set; }

        [SugarColumn(ColumnName = "conquest_index")]
        public int ConquestIndex { get; set; }

        [SugarColumn(ColumnName = "respawn_index")]
        public int RespawnIndex { get; set; }
    }
}