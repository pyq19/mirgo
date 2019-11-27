using SqlSugar;

namespace dotnettools
{
    [SugarTable("quest_info")]
    public class QuestInfoModel
    {
        public QuestInfoModel() { }

        [SugarColumn(ColumnName = "id", IsPrimaryKey = true, IsIdentity = true)] //是主键, 还是标识列
        public int Id { get; set; }

        [SugarColumn(ColumnName = "quest_index")]
        public int QuestIndex { get; set; }

        [SugarColumn(ColumnName = "name")]
        public string Name { get; set; }

        [SugarColumn(ColumnName = "quest_group")]
        public string Group { get; set; }

        [SugarColumn(ColumnName = "file_name")]
        public string FileName { get; set; }

        [SugarColumn(ColumnName = "required_min_level")]
        public int RequiredMinLevel { get; set; }

        [SugarColumn(ColumnName = "required_max_level")]
        public int RequiredMaxLevel { get; set; }

        [SugarColumn(ColumnName = "required_quest")]
        public int RequiredQuest { get; set; }

        // public RequiredClass RequiredClass
        [SugarColumn(ColumnName = "required_class")]
        public byte RequiredClass { get; set; }

        // public QuestType Type
        [SugarColumn(ColumnName = "quest_type")]
        public byte Type { get; set; }

        [SugarColumn(ColumnName = "goto_message")]
        public string GotoMessage { get; set; }

        [SugarColumn(ColumnName = "kill_message")]
        public string KillMessage { get; set; }

        [SugarColumn(ColumnName = "item_message")]
        public string ItemMessage { get; set; }

        [SugarColumn(ColumnName = "flag_message")]
        public string FlagMessage { get; set; }
    }
}