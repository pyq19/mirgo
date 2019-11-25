using System;
using System.IO;

namespace dotnettools
{
    class Loader
    {
        public Database loadDB(string path)
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
                    MapInfoList.Add(new MapInfo(reader));

                count = reader.ReadInt32();
                ItemInfoList.Clear();
                for (var i = 0; i < count; i++)
                {
                    ItemInfoList.Add(new ItemInfo(reader, LoadVersion, LoadCustomVersion));
                    if (ItemInfoList[i] != null && ItemInfoList[i].RandomStatsId < Settings.RandomItemStatsList.Count)
                    {
                        ItemInfoList[i].RandomStats = Settings.RandomItemStatsList[ItemInfoList[i].RandomStatsId];
                    }
                }
                count = reader.ReadInt32();
                MonsterInfoList.Clear();
                for (var i = 0; i < count; i++)
                    MonsterInfoList.Add(new MonsterInfo(reader));

                if (LoadVersion > 33)
                {
                    count = reader.ReadInt32();
                    NPCInfoList.Clear();
                    for (var i = 0; i < count; i++)
                        NPCInfoList.Add(new NPCInfo(reader));

                    count = reader.ReadInt32();
                    QuestInfoList.Clear();
                    for (var i = 0; i < count; i++)
                        QuestInfoList.Add(new QuestInfo(reader));
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
                    GameShopList.Clear();
                    for (var i = 0; i < count; i++)
                    {
                        var item = new GameShopItem(reader, LoadVersion, LoadCustomVersion);
                        if (Main.BindGameShop(item))
                        {
                            GameShopList.Add(item);
                        }
                    }
                }

                if (LoadVersion >= 66)
                {
                    ConquestInfos.Clear();
                    count = reader.ReadInt32();
                    for (var i = 0; i < count; i++)
                    {
                        ConquestInfos.Add(new ConquestInfo(reader));
                    }
                }

                if (LoadVersion > 67)
                    RespawnTick = new RespawnTimer(reader);

            }
            return null;
        }
    }

}