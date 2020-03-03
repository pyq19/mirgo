using SqlSugar;

namespace dotnettools
{
    [SugarTable("item")]
    public class ItemInfoModel
    {
        public ItemInfoModel() { }

        [SugarColumn(ColumnName = "id", IsPrimaryKey = true, IsIdentity = true)] //是主键, 还是标识列
        public int Id { get; set; }

        //[SugarColumn(ColumnName = "item_index")]
        //public int ItemIndex { get; set; }

        [SugarColumn(ColumnName = "name")]
        public string Name { get; set; }

        [SugarColumn(ColumnName = "type")]
        public byte Type { get; set; }       // ItemType

        [SugarColumn(ColumnName = "grade")]
        public byte Grade { get; set; }      // ItemGrade

        [SugarColumn(ColumnName = "required_type")]
        public byte RequiredType { get; set; }       // RequiredType

        [SugarColumn(ColumnName = "required_class")]
        public byte RequiredClass { get; set; }      // RequiredClass

        [SugarColumn(ColumnName = "required_gender")]
        public byte RequiredGender { get; set; }     // RequiredGender

        [SugarColumn(ColumnName = "item_set")]
        public byte Set { get; set; }    // ItemSet

        [SugarColumn(ColumnName = "shape")]
        public short Shape { get; set; }    // int16

        [SugarColumn(ColumnName = "weight")]
        public byte Weight { get; set; }

        [SugarColumn(ColumnName = "light")]
        public byte Light { get; set; }

        [SugarColumn(ColumnName = "required_amount")]
        public byte RequiredAmount { get; set; }

        [SugarColumn(ColumnName = "image")]
        public ushort Image { get; set; }   // uint16

        [SugarColumn(ColumnName = "durability")]
        public ushort Durability { get; set; }

        [SugarColumn(ColumnName = "stack_size")]
        public uint StackSize { get; set; } // uint32

        [SugarColumn(ColumnName = "price")]
        public uint Price { get; set; } // uint32

        [SugarColumn(ColumnName = "min_ac")]
        public byte MinAC { get; set; }

        [SugarColumn(ColumnName = "max_ac")]
        public byte MaxAC { get; set; }

        [SugarColumn(ColumnName = "min_mac")]
        public byte MinMAC { get; set; }

        [SugarColumn(ColumnName = "max_mac")]
        public byte MaxMAC { get; set; }

        [SugarColumn(ColumnName = "min_dc")]
        public byte MinDC { get; set; }

        [SugarColumn(ColumnName = "max_dc")]
        public byte MaxDC { get; set; }

        [SugarColumn(ColumnName = "min_mc")]
        public byte MinMC { get; set; }

        [SugarColumn(ColumnName = "max_mc")]
        public byte MaxMC { get; set; }

        [SugarColumn(ColumnName = "min_sc")]
        public byte MinSC { get; set; }

        [SugarColumn(ColumnName = "max_sc")]
        public byte MaxSC { get; set; }

        [SugarColumn(ColumnName = "hp")]
        public ushort HP { get; set; }

        [SugarColumn(ColumnName = "mp")]
        public ushort MP { get; set; }

        [SugarColumn(ColumnName = "accuracy")]
        public byte Accuracy { get; set; }

        [SugarColumn(ColumnName = "agility")]
        public byte Agility { get; set; }

        [SugarColumn(ColumnName = "luck")]
        public sbyte Luck { get; set; }

        [SugarColumn(ColumnName = "attack_speed")]
        public sbyte AttackSpeed { get; set; }

        [SugarColumn(ColumnName = "start_item")]
        public bool StartItem { get; set; }

        [SugarColumn(ColumnName = "bag_weight")]
        public byte BagWeight { get; set; }

        [SugarColumn(ColumnName = "hand_weight")]
        public byte HandWeight { get; set; }

        [SugarColumn(ColumnName = "wear_weight")]
        public byte WearWeight { get; set; }

        [SugarColumn(ColumnName = "effect")]
        public byte Effect { get; set; }

        [SugarColumn(ColumnName = "strong")]
        public byte Strong { get; set; }

        [SugarColumn(ColumnName = "magic_resist")]
        public byte MagicResist { get; set; }

        [SugarColumn(ColumnName = "poison_resist")]
        public byte PoisonResist { get; set; }

        [SugarColumn(ColumnName = "health_recovery")]
        public byte HealthRecovery { get; set; }

        [SugarColumn(ColumnName = "spell_recovery")]
        public byte SpellRecovery { get; set; }

        [SugarColumn(ColumnName = "poison_recovery")]
        public byte PoisonRecovery { get; set; }

        [SugarColumn(ColumnName = "hp_rate")]
        public byte HPrate { get; set; }

        [SugarColumn(ColumnName = "mp_rate")]
        public byte MPrate { get; set; }

        [SugarColumn(ColumnName = "critical_rate")]
        public byte CriticalRate { get; set; }

        [SugarColumn(ColumnName = "critical_damage")]
        public byte CriticalDamage { get; set; }

        [SugarColumn(ColumnName = "bools")]
        public byte bools { get; set; }

        [SugarColumn(ColumnName = "max_ac_rate")]
        public byte MaxAcRate { get; set; }

        [SugarColumn(ColumnName = "max_mac_rate")]
        public byte MaxMacRate { get; set; }

        [SugarColumn(ColumnName = "holy")]
        public byte Holy { get; set; }

        [SugarColumn(ColumnName = "freezing")]
        public byte Freezing { get; set; }

        [SugarColumn(ColumnName = "poison_attack")]
        public byte PoisonAttack { get; set; }

        [SugarColumn(ColumnName = "bind")]
        public short Bind { get; set; } // BindMode

        [SugarColumn(ColumnName = "reflect")]
        public byte Reflect { get; set; }

        [SugarColumn(ColumnName = "hp_drain_rate")]
        public byte HpDrainRate { get; set; }

        [SugarColumn(ColumnName = "unique_item")]
        public short Unique { get; set; } // SpecialItemMode

        [SugarColumn(ColumnName = "random_stats_id")]
        public byte RandomStatsId { get; set; }

        [SugarColumn(ColumnName = "can_fast_run")]
        public bool CanFastRun { get; set; }

        [SugarColumn(ColumnName = "can_awakening")]
        public bool CanAwakening { get; set; }

        [SugarColumn(ColumnName = "tool_tip")]
        public string ToolTip { get; set; }
    }
}

