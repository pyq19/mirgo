using SqlSugar;

namespace dotnettools
{
    [SugarTable("safe_zone_info")]
    public class SafeZoneInfoModel
    {
        public SafeZoneInfoModel() { }

        [SugarColumn(ColumnName = "id", IsPrimaryKey = true, IsIdentity = true)] //是主键, 还是标识列
        public int Id { get; set; }

        [SugarColumn(ColumnName = "map_id")]
        public int MapInfoId { get; set; }

        // public Point Location;
        [SugarColumn(ColumnName = "location_x")]
        public int LocationX { get; set; }
        [SugarColumn(ColumnName = "location_y")]
        public int LocationY { get; set; }

        // public ushort Size;
        [SugarColumn(ColumnName = "size")]
        public ushort Size { get; set; }

        // public bool StartPoint;
        [SugarColumn(ColumnName = "start_point")]
        public int StartPoint { get; set; }
    }
}