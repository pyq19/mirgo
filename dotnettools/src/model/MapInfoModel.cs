using SqlSugar;

namespace dotnettools
{
    [SugarTable("map_info")]
    public class MapInfoModel
    {
        public MapInfoModel() { }

        [SugarColumn(ColumnName = "id", IsPrimaryKey = true, IsIdentity = true)] //是主键, 还是标识列
        public int Id { get; set; }

        [SugarColumn(ColumnName = "map_index")]
        public int MapIndex { get; set; }

        [SugarColumn(ColumnName = "file_name")]
        public string FileName { get; set; }

        [SugarColumn(ColumnName = "title")]
        public string Title { get; set; }

        [SugarColumn(ColumnName = "mini_map")]
        public ushort MiniMap { get; set; }

        [SugarColumn(ColumnName = "big_map")]
        public ushort BigMap { get; set; }

        [SugarColumn(ColumnName = "music")]
        public ushort Music { get; set; }

        [SugarColumn(ColumnName = "light")]
        public byte Light { get; set; } // LightSetting

        [SugarColumn(ColumnName = "map_dark_light")]
        public byte MapDarkLight { get; set; }

        [SugarColumn(ColumnName = "mine_index")]
        public byte MineIndex { get; set; }

        [SugarColumn(ColumnName = "no_teleport")]
        public bool NoTeleport { get; set; }

        [SugarColumn(ColumnName = "no_reconnect")]
        public bool NoReconnect { get; set; }

        [SugarColumn(ColumnName = "no_random")]
        public bool NoRandom { get; set; }

        [SugarColumn(ColumnName = "no_escape")]
        public bool NoEscape { get; set; }

        [SugarColumn(ColumnName = "no_recall")]
        public bool NoRecall { get; set; }

        [SugarColumn(ColumnName = "no_drug")]
        public bool NoDrug { get; set; }

        [SugarColumn(ColumnName = "no_position")]
        public bool NoPosition { get; set; }

        [SugarColumn(ColumnName = "no_fight")]
        public bool NoFight { get; set; }

        [SugarColumn(ColumnName = "no_throw_item")]
        public bool NoThrowItem { get; set; }

        [SugarColumn(ColumnName = "no_drop_player")]
        public bool NoDropPlayer { get; set; }

        [SugarColumn(ColumnName = "no_drop_monster")]
        public bool NoDropMonster { get; set; }

        [SugarColumn(ColumnName = "no_names")]
        public bool NoNames { get; set; }

        [SugarColumn(ColumnName = "no_mount")]
        public bool NoMount { get; set; }

        [SugarColumn(ColumnName = "need_bridle")]
        public bool NeedBridle { get; set; }

        [SugarColumn(ColumnName = "fight")]
        public bool Fight { get; set; }

        [SugarColumn(ColumnName = "fire")]
        public bool Fire { get; set; }

        [SugarColumn(ColumnName = "lightning")]
        public bool Lightning { get; set; }

        [SugarColumn(ColumnName = "no_town_teleport")]
        public bool NoTownTeleport { get; set; }

        [SugarColumn(ColumnName = "no_reincarnation")]
        public bool NoReincarnation { get; set; }

        [SugarColumn(ColumnName = "no_reconnect_map")]
        public string NoReconnectMap { get; set; }

        [SugarColumn(ColumnName = "fire_damage")]
        public int FireDamage { get; set; }

        [SugarColumn(ColumnName = "lightning_damage")]
        public int LightningDamage { get; set; }
    }
}