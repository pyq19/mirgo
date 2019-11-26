using SqlSugar;

namespace dotnettools
{
    [SugarTable("movement_info")]
    public class MovementInfoModel
    {
        [SugarColumn(ColumnName = "id", IsPrimaryKey = true, IsIdentity = true)] //是主键, 还是标识列
        public int Id { get; set; }

        [SugarColumn(ColumnName = "map_index")]
        public int MapIndex { get; set; }

        // Point
        [SugarColumn(ColumnName = "source_x")]
        public int SourceX { get; set; }

        [SugarColumn(ColumnName = "source_y")]
        public int SourceY { get; set; }

        // Point
        [SugarColumn(ColumnName = "destination_x")]
        public int DestinationX { get; set; }

        [SugarColumn(ColumnName = "destination_y")]
        public int DestinationY { get; set; }

        [SugarColumn(ColumnName = "need_hole")]
        public bool NeedHole { get; set; }

        [SugarColumn(ColumnName = "need_move")]
        public bool NeedMove { get; set; }

        [SugarColumn(ColumnName = "conquest_index")]
        public int ConquestIndex { get; set; }
    }
}