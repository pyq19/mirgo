using System;
using SqlSugar;

namespace dotnettools
{
    [SugarTable("game_shop_item")]
    public class GameShopItemModel
    {
        public GameShopItemModel() { }

        [SugarColumn(ColumnName = "id", IsPrimaryKey = true, IsIdentity = true)] //是主键, 还是标识列
        public int Id { get; set; }

        //[SugarColumn(ColumnName = "game_shop_item_index")]
        //public int GameShopItemIndex { get; set; }

        // [SugarColumn(ColumnName = "info")]
        // public ItemInfo Info { get; set; }
        [SugarColumn(ColumnName = "item_id")]
        public int ItemIndex { get; set; }

        [SugarColumn(ColumnName = "gold_price")]
        public uint GoldPrice { get; set; }

        [SugarColumn(ColumnName = "credit_price")]
        public uint CreditPrice { get; set; }

        [SugarColumn(ColumnName = "count")]
        public uint Count { get; set; }

        [SugarColumn(ColumnName = "class")]
        public string Class { get; set; }

        [SugarColumn(ColumnName = "category")]
        public string Category { get; set; }

        [SugarColumn(ColumnName = "stock")]
        public int Stock { get; set; }

        [SugarColumn(ColumnName = "i_stock")]
        public bool iStock { get; set; }

        [SugarColumn(ColumnName = "deal")]
        public bool Deal { get; set; }

        [SugarColumn(ColumnName = "top_item")]
        public bool TopItem { get; set; }

        [SugarColumn(ColumnName = "create_date")]
        public DateTime Date { get; set; }
    }
}