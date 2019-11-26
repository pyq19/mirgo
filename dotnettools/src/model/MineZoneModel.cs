using SqlSugar;

namespace dotnettools
{
    [SugarTable("mine_zone")]
    public class MineZoneModel
    {
        [SugarColumn(ColumnName = "id", IsPrimaryKey = true, IsIdentity = true)] //是主键, 还是标识列
        public int Id { get; set; }

        [SugarColumn(ColumnName = "map_index")]
        public int MapIndex { get; set; }

        [SugarColumn(ColumnName = "mine")]
        public byte Mine { get; set; }

        // Point
        [SugarColumn(ColumnName = "location_x")]
        public int LocationX { get; set; }
        [SugarColumn(ColumnName = "location_y")]
        public int LocationY { get; set; }

        [SugarColumn(ColumnName = "size")]
        public ushort Size { get; set; }
    }
}