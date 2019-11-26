using System;
using System.IO;
using System.Collections.Generic;

namespace dotnettools
{
    public class Loader
    {
        public int LoadVersion;
        public int LoadCustomVersion;
        public int MapIndex;
        public int ItemIndex;
        public int MonsterIndex;
        public int NPCIndex;
        public int QuestIndex;
        public int GameshopIndex;
        public int ConquestIndex;
        public int RespawnIndex;
        public List<MapInfo> MapInfoList = new List<MapInfo>();
        public List<ItemInfo> ItemInfoList = new List<ItemInfo>();
        public List<MonsterInfo> MonsterInfoList = new List<MonsterInfo>();
        public List<NPCInfo> NPCInfoList = new List<NPCInfo>();
        public List<QuestInfo> QuestInfoList = new List<QuestInfo>();
        public DragonInfo DragonInfo;
        public List<ConquestInfo> ConquestInfos = new List<ConquestInfo>();
        public List<MagicInfo> MagicInfoList = new List<MagicInfo>();
        public List<GameShopItem> GameShopItemList = new List<GameShopItem>();
        public RespawnTimer RespawnTick;

        public long Time { get; private set; }

        public void loadDB(string path)
        {
            using (var stream = File.OpenRead(path))
            using (var reader = new BinaryReader(stream))
            {
                LoadVersion = reader.ReadInt32();
                if (LoadVersion > 57)
                    LoadCustomVersion = reader.ReadInt32();
                MapIndex = reader.ReadInt32();
                ItemIndex = reader.ReadInt32();
                MonsterIndex = reader.ReadInt32();

                if (LoadVersion > 33)
                {
                    NPCIndex = reader.ReadInt32();
                    QuestIndex = reader.ReadInt32();
                }
                if (LoadVersion >= 63)
                {
                    GameshopIndex = reader.ReadInt32();
                }

                if (LoadVersion >= 66)
                {
                    ConquestIndex = reader.ReadInt32();
                }

                if (LoadVersion >= 68)
                    RespawnIndex = reader.ReadInt32();

                var count = reader.ReadInt32();
                MapInfoList.Clear();
                for (var i = 0; i < count; i++)
                    MapInfoList.Add(new MapInfo(reader, this));

                count = reader.ReadInt32();
                ItemInfoList.Clear();
                for (var i = 0; i < count; i++)
                {
                    ItemInfoList.Add(new ItemInfo(reader, LoadVersion, LoadCustomVersion));
                    // TODO 这部分是从 Setting 里面加, 以后改成数据库
                    // if (ItemInfoList[i] != null && ItemInfoList[i].RandomStatsId < Settings.RandomItemStatsList.Count)
                    // {
                    //     ItemInfoList[i].RandomStats = Settings.RandomItemStatsList[ItemInfoList[i].RandomStatsId];
                    // }
                }
                count = reader.ReadInt32();
                MonsterInfoList.Clear();
                for (var i = 0; i < count; i++)
                    MonsterInfoList.Add(new MonsterInfo(reader, this));

                if (LoadVersion > 33)
                {
                    count = reader.ReadInt32();
                    NPCInfoList.Clear();
                    for (var i = 0; i < count; i++)
                        NPCInfoList.Add(new NPCInfo(reader, this));

                    count = reader.ReadInt32();
                    QuestInfoList.Clear();
                    for (var i = 0; i < count; i++)
                        QuestInfoList.Add(new QuestInfo(reader, this));
                }

                DragonInfo = LoadVersion >= 11 ? new DragonInfo(reader) : new DragonInfo();
                if (LoadVersion >= 58)
                {
                    count = reader.ReadInt32();
                    for (var i = 0; i < count; i++)
                    {
                        var m = new MagicInfo(reader, LoadVersion, LoadCustomVersion);
                        if (!MagicExists(m.Spell))
                            MagicInfoList.Add(m);
                    }
                }
                FillMagicInfoList();
                if (LoadVersion <= 70)
                    UpdateMagicInfo();

                if (LoadVersion >= 63)
                {
                    count = reader.ReadInt32();
                    GameShopItemList.Clear();
                    for (var i = 0; i < count; i++)
                    {
                        var item = new GameShopItem(reader, LoadVersion, LoadCustomVersion);
                        if (BindGameShop(item))
                        {
                            GameShopItemList.Add(item);
                        }
                    }
                }

                if (LoadVersion >= 66)
                {
                    ConquestInfos.Clear();
                    count = reader.ReadInt32();
                    for (var i = 0; i < count; i++)
                    {
                        ConquestInfos.Add(new ConquestInfo(reader, this));
                    }
                }

                if (LoadVersion > 67)
                    RespawnTick = new RespawnTimer(reader, this);

            }
        }
        bool MagicExists(Spell spell)
        {

            return true;
        }

        void FillMagicInfoList()
        {

        }

        void UpdateMagicInfo()
        {

        }

        public bool BindGameShop(GameShopItem item, bool editEnvir = true)
        {
            for (var i = 0; i < ItemInfoList.Count; i++)
            {
                var info = ItemInfoList[i];
                if (info.Index != item.ItemIndex) continue;
                item.Info = info;
                return true;
            }
            return false;
        }
    }

}