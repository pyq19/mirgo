using SqlSugar;

namespace dotnettools
{
    [SugarTable("magic_info")]
    public class MagicInfoModel
    {
        [SugarColumn(ColumnName = "id", IsPrimaryKey = true, IsIdentity = true)] //是主键, 还是标识列
        public int Id { get; set; }

        [SugarColumn(ColumnName = "name")]
        public string Name { get; set; }

        // public Spell Spell
        [SugarColumn(ColumnName = "spell")]
        public byte Spell { get; set; }

        [SugarColumn(ColumnName = "base_cost")]
        public byte BaseCost { get; set; }

        [SugarColumn(ColumnName = "level_cost")]
        public byte LevelCost { get; set; }

        [SugarColumn(ColumnName = "icon")]
        public byte Icon { get; set; }

        [SugarColumn(ColumnName = "level_1")]
        public byte Level1 { get; set; }

        [SugarColumn(ColumnName = "level_2")]
        public byte Level2 { get; set; }

        [SugarColumn(ColumnName = "level_3")]
        public byte Level3 { get; set; }

        [SugarColumn(ColumnName = "need_1")]
        public ushort Need1 { get; set; }

        [SugarColumn(ColumnName = "need_2")]
        public ushort Need2 { get; set; }

        [SugarColumn(ColumnName = "need_3")]
        public ushort Need3 { get; set; }

        [SugarColumn(ColumnName = "delay_base")]
        public uint DelayBase { get; set; }

        [SugarColumn(ColumnName = "delay_reduction")]
        public uint DelayReduction { get; set; }

        [SugarColumn(ColumnName = "power_base")]
        public ushort PowerBase { get; set; }

        [SugarColumn(ColumnName = "power_bonus")]
        public ushort PowerBonus { get; set; }

        [SugarColumn(ColumnName = "m_power_base")]
        public ushort MPowerBase { get; set; }

        [SugarColumn(ColumnName = "m_power_bonus")]
        public ushort MPowerBonus { get; set; }

        [SugarColumn(ColumnName = "magic_range")]
        public byte Range { get; set; }

        [SugarColumn(ColumnName = "multiplier_base")]
        public float MultiplierBase { get; set; }

        [SugarColumn(ColumnName = "multiplier_bonus")]
        public float MultiplierBonus { get; set; }
    }
}