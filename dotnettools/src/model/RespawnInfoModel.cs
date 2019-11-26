using SqlSugar;

namespace dotnettools
{
    [SugarTable("respawn_info")]
    public class RespawnInfoModel
    {
        [SugarColumn(ColumnName = "id", IsPrimaryKey = true, IsIdentity = true)] //是主键, 还是标识列
        public int Id { get; set; }

        [SugarColumn(ColumnName = "map_index")]
        public int MapIndex { get; set; }

        [SugarColumn(ColumnName = "monster_index")]
        public int MonsterIndex { get; set; }

        [SugarColumn(ColumnName = "location_x")]
        // public Point Location;
        public int LocationX { get; set; }

        [SugarColumn(ColumnName = "location_y")]
        public int LocationY { get; set; }

        [SugarColumn(ColumnName = "count")]
        public ushort Count { get; set; }

        [SugarColumn(ColumnName = "spread")]
        public ushort Spread { get; set; }

        [SugarColumn(ColumnName = "delay")]
        public ushort Delay { get; set; }

        [SugarColumn(ColumnName = "random_delay")]
        public ushort RandomDelay { get; set; }

        [SugarColumn(ColumnName = "direction")]
        public byte Direction { get; set; }

        [SugarColumn(ColumnName = "route_path")]
        public string RoutePath { get; set; }

        [SugarColumn(ColumnName = "respawn_index")]
        public int RespawnIndex { get; set; }

        [SugarColumn(ColumnName = "save_respawn_time")]
        public bool SaveRespawnTime { get; set; }

        [SugarColumn(ColumnName = "respawn_ticks")]
        public ushort RespawnTicks { get; set; }
    }
}