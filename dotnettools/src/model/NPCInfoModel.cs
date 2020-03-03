using SqlSugar;

namespace dotnettools
{
    [SugarTable("npc")]
    public class NPCInfoModel
    {
        public NPCInfoModel() { }

        [SugarColumn(ColumnName = "id", IsPrimaryKey = true, IsIdentity = true)] //是主键, 还是标识列
        public int Id { get; set; }

        [SugarColumn(ColumnName = "map_id")]
        public int MapIndex { get; set; }

        //[SugarColumn(ColumnName = "npc_index")]
        //public int NPCIndex { get; set; }

        [SugarColumn(ColumnName = "file_name")]
        public string FileName { get; set; }

        [SugarColumn(ColumnName = "name")]
        public string Name { get; set; }

        [SugarColumn(ColumnName = "location_x")]
        // public Point Location;
        public int LocationX { get; set; }

        [SugarColumn(ColumnName = "location_y")]
        public int LocationY { get; set; }

        [SugarColumn(ColumnName = "rate")]
        public ushort Rate { get; set; }

        [SugarColumn(ColumnName = "image")]
        public ushort Image { get; set; }

        [SugarColumn(ColumnName = "time_visible")]
        public bool TimeVisible { get; set; }

        [SugarColumn(ColumnName = "hour_start")]
        public byte HourStart { get; set; }

        [SugarColumn(ColumnName = "minute_start")]
        public byte MinuteStart { get; set; }

        [SugarColumn(ColumnName = "hour_end")]
        public byte HourEnd { get; set; }

        [SugarColumn(ColumnName = "minute_end")]
        public byte MinuteEnd { get; set; }

        [SugarColumn(ColumnName = "min_lev")]
        public short MinLev { get; set; }

        [SugarColumn(ColumnName = "max_lev")]
        public short MaxLev { get; set; }

        [SugarColumn(ColumnName = "day_of_week")]
        public string DayofWeek { get; set; }

        [SugarColumn(ColumnName = "class_required")]
        public string ClassRequired { get; set; }

        [SugarColumn(ColumnName = "flag_needed")]
        public int FlagNeeded { get; set; }

        [SugarColumn(ColumnName = "conquest")]
        public int Conquest { get; set; }
    }
}