using SqlSugar;

namespace dotnettools
{
    [SugarTable("monster")]
    public class MonsterInfoModel
    {
        [SugarColumn(ColumnName = "id", IsPrimaryKey = true, IsIdentity = true)] //是主键, 还是标识列
        public int Id { get; set; }

        //[SugarColumn(ColumnName = "monster_index")]
        //public int MonsterIndex { get; set; }

        [SugarColumn(ColumnName = "name")]
        public string Name { get; set; }

        [SugarColumn(ColumnName = "image")]
        public ushort Image { get; set; }

        [SugarColumn(ColumnName = "ai")]
        public byte AI { get; set; }

        [SugarColumn(ColumnName = "effect")]
        public byte Effect { get; set; }

        [SugarColumn(ColumnName = "level")]
        public ushort Level { get; set; }

        [SugarColumn(ColumnName = "view_range")]
        public byte ViewRange { get; set; }

        [SugarColumn(ColumnName = "cool_eye")]
        public byte CoolEye { get; set; }

        [SugarColumn(ColumnName = "hp")]
        public uint HP { get; set; }

        [SugarColumn(ColumnName = "min_ac")]
        public ushort MinAC { get; set; }

        [SugarColumn(ColumnName = "max_ac")]
        public ushort MaxAC { get; set; }

        [SugarColumn(ColumnName = "min_mac")]
        public ushort MinMAC { get; set; }

        [SugarColumn(ColumnName = "max_mac")]
        public ushort MaxMAC { get; set; }

        [SugarColumn(ColumnName = "min_dc")]
        public ushort MinDC { get; set; }

        [SugarColumn(ColumnName = "max_dc")]
        public ushort MaxDC { get; set; }

        [SugarColumn(ColumnName = "min_mc")]
        public ushort MinMC { get; set; }

        [SugarColumn(ColumnName = "max_mc")]
        public ushort MaxMC { get; set; }

        [SugarColumn(ColumnName = "min_sc")]
        public ushort MinSC { get; set; }

        [SugarColumn(ColumnName = "max_sc")]
        public ushort MaxSC { get; set; }

        [SugarColumn(ColumnName = "accuracy")]
        public byte Accuracy { get; set; }

        [SugarColumn(ColumnName = "agility")]
        public byte Agility { get; set; }

        [SugarColumn(ColumnName = "light")]
        public byte Light { get; set; }

        [SugarColumn(ColumnName = "attack_speed")]
        public ushort AttackSpeed { get; set; }

        [SugarColumn(ColumnName = "move_speed")]
        public ushort MoveSpeed { get; set; }

        [SugarColumn(ColumnName = "experience")]
        public uint Experience { get; set; }

        [SugarColumn(ColumnName = "can_push")]
        public bool CanPush { get; set; }

        [SugarColumn(ColumnName = "can_tame")]
        public bool CanTame { get; set; }

        [SugarColumn(ColumnName = "auto_rev")]
        public bool AutoRev { get; set; }

        [SugarColumn(ColumnName = "undead")]
        public bool Undead { get; set; }
    }
}