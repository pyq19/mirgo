using System;
using System.IO;
using System.Collections.Generic;
using SqlSugar;

namespace dotnettools
{
    public class Manager
    {
        public static SqlSugarClient DB;

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
        public RespawnTimer RespawnTick = new RespawnTimer();
        public long Time { get; private set; }

        // 从原始的文件数据库中读取数据
        public void loadFromFile(string path)
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
                    // TODO 
                    // 这部分是从 Setting 里面加, 以后改成数据库
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

        public bool BindGameShop(GameShopItem item, bool editEnvir = true)
        {
            for (var i = 0; i < ItemInfoList.Count; i++)
            {
                var info = ItemInfoList[i];
                if (info.ItemIndex != item.ItemIndex) continue;
                item.Info = info;
                return true;
            }
            return false;
        }

        #region 魔法相关
        private bool MagicExists(Spell spell)
        {
            for (var i = 0; i < MagicInfoList.Count; i++)
            {
                if (MagicInfoList[i].Spell == spell) return true;
            }
            return false;
        }

        private void FillMagicInfoList()
        {
            //Warrior
            if (!MagicExists(Spell.Fencing)) MagicInfoList.Add(new MagicInfo { Name = "Fencing", Spell = Spell.Fencing, Icon = 2, Level1 = 7, Level2 = 9, Level3 = 12, Need1 = 270, Need2 = 600, Need3 = 1300, Range = 0 });
            if (!MagicExists(Spell.Slaying)) MagicInfoList.Add(new MagicInfo { Name = "Slaying", Spell = Spell.Slaying, Icon = 6, Level1 = 15, Level2 = 17, Level3 = 20, Need1 = 500, Need2 = 1100, Need3 = 1800, Range = 0 });
            if (!MagicExists(Spell.Thrusting)) MagicInfoList.Add(new MagicInfo { Name = "Thrusting", Spell = Spell.Thrusting, Icon = 11, Level1 = 22, Level2 = 24, Level3 = 27, Need1 = 2000, Need2 = 3500, Need3 = 6000, Range = 0, MultiplierBase = 0.25f, MultiplierBonus = 0.25f });
            if (!MagicExists(Spell.HalfMoon)) MagicInfoList.Add(new MagicInfo { Name = "HalfMoon", Spell = Spell.HalfMoon, Icon = 24, Level1 = 26, Level2 = 28, Level3 = 31, Need1 = 5000, Need2 = 8000, Need3 = 14000, BaseCost = 3, Range = 0, MultiplierBase = 0.3f, MultiplierBonus = 0.1f });
            if (!MagicExists(Spell.ShoulderDash)) MagicInfoList.Add(new MagicInfo { Name = "ShoulderDash", Spell = Spell.ShoulderDash, Icon = 26, Level1 = 30, Level2 = 32, Level3 = 34, Need1 = 3000, Need2 = 4000, Need3 = 6000, BaseCost = 4, LevelCost = 4, DelayBase = 2500, Range = 0, MPowerBase = 4 });
            if (!MagicExists(Spell.TwinDrakeBlade)) MagicInfoList.Add(new MagicInfo { Name = "TwinDrakeBlade", Spell = Spell.TwinDrakeBlade, Icon = 37, Level1 = 32, Level2 = 34, Level3 = 37, Need1 = 4000, Need2 = 6000, Need3 = 10000, BaseCost = 10, Range = 0, MultiplierBase = 0.8f, MultiplierBonus = 0.1f });
            if (!MagicExists(Spell.Entrapment)) MagicInfoList.Add(new MagicInfo { Name = "Entrapment", Spell = Spell.Entrapment, Icon = 46, Level1 = 32, Level2 = 35, Level3 = 37, Need1 = 2000, Need2 = 3500, Need3 = 5500, BaseCost = 15, LevelCost = 3, Range = 9 });
            if (!MagicExists(Spell.FlamingSword)) MagicInfoList.Add(new MagicInfo { Name = "FlamingSword", Spell = Spell.FlamingSword, Icon = 25, Level1 = 35, Level2 = 37, Level3 = 40, Need1 = 2000, Need2 = 4000, Need3 = 6000, BaseCost = 7, Range = 0, MultiplierBase = 1.4f, MultiplierBonus = 0.4f });
            if (!MagicExists(Spell.LionRoar)) MagicInfoList.Add(new MagicInfo { Name = "LionRoar", Spell = Spell.LionRoar, Icon = 42, Level1 = 36, Level2 = 39, Level3 = 41, Need1 = 5000, Need2 = 8000, Need3 = 12000, BaseCost = 14, LevelCost = 4, Range = 0 });
            if (!MagicExists(Spell.CrossHalfMoon)) MagicInfoList.Add(new MagicInfo { Name = "CrossHalfMoon", Spell = Spell.CrossHalfMoon, Icon = 33, Level1 = 38, Level2 = 40, Level3 = 42, Need1 = 7000, Need2 = 11000, Need3 = 16000, BaseCost = 6, Range = 0, MultiplierBase = 0.4f, MultiplierBonus = 0.1f });
            if (!MagicExists(Spell.BladeAvalanche)) MagicInfoList.Add(new MagicInfo { Name = "BladeAvalanche", Spell = Spell.BladeAvalanche, Icon = 43, Level1 = 38, Level2 = 41, Level3 = 43, Need1 = 5000, Need2 = 8000, Need3 = 12000, BaseCost = 14, LevelCost = 4, Range = 0, MultiplierBonus = 0.3f });
            if (!MagicExists(Spell.ProtectionField)) MagicInfoList.Add(new MagicInfo { Name = "ProtectionField", Spell = Spell.ProtectionField, Icon = 50, Level1 = 39, Level2 = 42, Level3 = 45, Need1 = 6000, Need2 = 12000, Need3 = 18000, BaseCost = 23, LevelCost = 6, Range = 0 });
            if (!MagicExists(Spell.Rage)) MagicInfoList.Add(new MagicInfo { Name = "Rage", Spell = Spell.Rage, Icon = 49, Level1 = 44, Level2 = 47, Level3 = 50, Need1 = 8000, Need2 = 14000, Need3 = 20000, BaseCost = 20, LevelCost = 5, Range = 0 });
            if (!MagicExists(Spell.CounterAttack)) MagicInfoList.Add(new MagicInfo { Name = "CounterAttack", Spell = Spell.CounterAttack, Icon = 72, Level1 = 47, Level2 = 51, Level3 = 55, Need1 = 7000, Need2 = 11000, Need3 = 15000, BaseCost = 12, LevelCost = 4, DelayBase = 24000, Range = 0, MultiplierBonus = 0.4f });
            if (!MagicExists(Spell.SlashingBurst)) MagicInfoList.Add(new MagicInfo { Name = "SlashingBurst", Spell = Spell.SlashingBurst, Icon = 55, Level1 = 50, Level2 = 53, Level3 = 56, Need1 = 10000, Need2 = 16000, Need3 = 24000, BaseCost = 25, LevelCost = 4, MPowerBase = 1, PowerBase = 3, DelayBase = 14000, DelayReduction = 4000, Range = 0, MultiplierBase = 3.25f, MultiplierBonus = 0.25f });
            if (!MagicExists(Spell.Fury)) MagicInfoList.Add(new MagicInfo { Name = "Fury", Spell = Spell.Fury, Icon = 76, Level1 = 45, Level2 = 48, Level3 = 51, Need1 = 8000, Need2 = 14000, Need3 = 20000, BaseCost = 10, LevelCost = 4, DelayBase = 600000, DelayReduction = 120000, Range = 0 });
            if (!MagicExists(Spell.ImmortalSkin)) MagicInfoList.Add(new MagicInfo { Name = "ImmortalSkin", Spell = Spell.ImmortalSkin, Icon = 80, Level1 = 60, Level2 = 61, Level3 = 62, Need1 = 1560, Need2 = 2200, Need3 = 3000, BaseCost = 10, LevelCost = 4, DelayBase = 600000, DelayReduction = 120000, Range = 0 });

            //Wizard
            if (!MagicExists(Spell.FireBall)) MagicInfoList.Add(new MagicInfo { Name = "FireBall", Spell = Spell.FireBall, Icon = 0, Level1 = 7, Level2 = 9, Level3 = 11, Need1 = 200, Need2 = 350, Need3 = 700, BaseCost = 3, LevelCost = 2, MPowerBase = 8, PowerBase = 2, Range = 9 });
            if (!MagicExists(Spell.Repulsion)) MagicInfoList.Add(new MagicInfo { Name = "Repulsion", Spell = Spell.Repulsion, Icon = 7, Level1 = 12, Level2 = 15, Level3 = 19, Need1 = 500, Need2 = 1300, Need3 = 2200, BaseCost = 2, LevelCost = 2, Range = 0, MPowerBase = 4 });
            if (!MagicExists(Spell.ElectricShock)) MagicInfoList.Add(new MagicInfo { Name = "ElectricShock", Spell = Spell.ElectricShock, Icon = 19, Level1 = 13, Level2 = 18, Level3 = 24, Need1 = 530, Need2 = 1100, Need3 = 2200, BaseCost = 3, LevelCost = 1, Range = 9 });
            if (!MagicExists(Spell.GreatFireBall)) MagicInfoList.Add(new MagicInfo { Name = "GreatFireBall", Spell = Spell.GreatFireBall, Icon = 4, Level1 = 15, Level2 = 18, Level3 = 21, Need1 = 2000, Need2 = 2700, Need3 = 3500, BaseCost = 5, LevelCost = 1, MPowerBase = 6, PowerBase = 10, Range = 9 });
            if (!MagicExists(Spell.HellFire)) MagicInfoList.Add(new MagicInfo { Name = "HellFire", Spell = Spell.HellFire, Icon = 8, Level1 = 16, Level2 = 20, Level3 = 24, Need1 = 700, Need2 = 2700, Need3 = 3500, BaseCost = 10, LevelCost = 3, MPowerBase = 14, PowerBase = 6, Range = 0 });
            if (!MagicExists(Spell.ThunderBolt)) MagicInfoList.Add(new MagicInfo { Name = "ThunderBolt", Spell = Spell.ThunderBolt, Icon = 10, Level1 = 17, Level2 = 20, Level3 = 23, Need1 = 500, Need2 = 2000, Need3 = 3500, BaseCost = 9, LevelCost = 2, MPowerBase = 8, MPowerBonus = 20, PowerBase = 9, Range = 9 });
            if (!MagicExists(Spell.Teleport)) MagicInfoList.Add(new MagicInfo { Name = "Teleport", Spell = Spell.Teleport, Icon = 20, Level1 = 19, Level2 = 22, Level3 = 25, Need1 = 350, Need2 = 1000, Need3 = 2000, BaseCost = 10, LevelCost = 3, Range = 0 });
            if (!MagicExists(Spell.FireBang)) MagicInfoList.Add(new MagicInfo { Name = "FireBang", Spell = Spell.FireBang, Icon = 22, Level1 = 22, Level2 = 25, Level3 = 28, Need1 = 3000, Need2 = 5000, Need3 = 10000, BaseCost = 14, LevelCost = 4, MPowerBase = 8, PowerBase = 8, Range = 9 });
            if (!MagicExists(Spell.FireWall)) MagicInfoList.Add(new MagicInfo { Name = "FireWall", Spell = Spell.FireWall, Icon = 21, Level1 = 24, Level2 = 28, Level3 = 33, Need1 = 4000, Need2 = 10000, Need3 = 20000, BaseCost = 30, LevelCost = 5, MPowerBase = 3, PowerBase = 3, Range = 9 });
            if (!MagicExists(Spell.Lightning)) MagicInfoList.Add(new MagicInfo { Name = "Lightning", Spell = Spell.Lightning, Icon = 9, Level1 = 26, Level2 = 29, Level3 = 32, Need1 = 3000, Need2 = 6000, Need3 = 12000, BaseCost = 38, LevelCost = 7, MPowerBase = 12, PowerBase = 12, Range = 0 });
            if (!MagicExists(Spell.FrostCrunch)) MagicInfoList.Add(new MagicInfo { Name = "FrostCrunch", Spell = Spell.FrostCrunch, Icon = 38, Level1 = 28, Level2 = 30, Level3 = 33, Need1 = 3000, Need2 = 5000, Need3 = 8000, BaseCost = 15, LevelCost = 3, MPowerBase = 12, PowerBase = 12, Range = 9 });
            if (!MagicExists(Spell.ThunderStorm)) MagicInfoList.Add(new MagicInfo { Name = "ThunderStorm", Spell = Spell.ThunderStorm, Icon = 23, Level1 = 30, Level2 = 32, Level3 = 34, Need1 = 4000, Need2 = 8000, Need3 = 12000, BaseCost = 29, LevelCost = 9, MPowerBase = 10, MPowerBonus = 20, PowerBase = 10, PowerBonus = 20, Range = 0 });
            if (!MagicExists(Spell.MagicShield)) MagicInfoList.Add(new MagicInfo { Name = "MagicShield", Spell = Spell.MagicShield, Icon = 30, Level1 = 31, Level2 = 34, Level3 = 38, Need1 = 3000, Need2 = 7000, Need3 = 10000, BaseCost = 35, LevelCost = 5, Range = 0 });
            if (!MagicExists(Spell.TurnUndead)) MagicInfoList.Add(new MagicInfo { Name = "TurnUndead", Spell = Spell.TurnUndead, Icon = 31, Level1 = 32, Level2 = 35, Level3 = 39, Need1 = 3000, Need2 = 7000, Need3 = 10000, BaseCost = 52, LevelCost = 13, Range = 9 });
            if (!MagicExists(Spell.Vampirism)) MagicInfoList.Add(new MagicInfo { Name = "Vampirism", Spell = Spell.Vampirism, Icon = 47, Level1 = 33, Level2 = 36, Level3 = 40, Need1 = 3000, Need2 = 5000, Need3 = 8000, BaseCost = 26, LevelCost = 13, MPowerBase = 12, PowerBase = 12, Range = 9 });
            if (!MagicExists(Spell.IceStorm)) MagicInfoList.Add(new MagicInfo { Name = "IceStorm", Spell = Spell.IceStorm, Icon = 32, Level1 = 35, Level2 = 37, Level3 = 40, Need1 = 4000, Need2 = 8000, Need3 = 12000, BaseCost = 33, LevelCost = 3, MPowerBase = 12, PowerBase = 14, Range = 9 });
            if (!MagicExists(Spell.FlameDisruptor)) MagicInfoList.Add(new MagicInfo { Name = "FlameDisruptor", Spell = Spell.FlameDisruptor, Icon = 34, Level1 = 38, Level2 = 40, Level3 = 42, Need1 = 5000, Need2 = 9000, Need3 = 14000, BaseCost = 28, LevelCost = 3, MPowerBase = 15, MPowerBonus = 20, PowerBase = 9, Range = 9 });
            if (!MagicExists(Spell.Mirroring)) MagicInfoList.Add(new MagicInfo { Name = "Mirroring", Spell = Spell.Mirroring, Icon = 41, Level1 = 41, Level2 = 43, Level3 = 45, Need1 = 6000, Need2 = 11000, Need3 = 16000, BaseCost = 21, Range = 0 });
            if (!MagicExists(Spell.FlameField)) MagicInfoList.Add(new MagicInfo { Name = "FlameField", Spell = Spell.FlameField, Icon = 44, Level1 = 42, Level2 = 43, Level3 = 45, Need1 = 6000, Need2 = 11000, Need3 = 16000, BaseCost = 45, LevelCost = 8, MPowerBase = 100, PowerBase = 25, Range = 9 });
            if (!MagicExists(Spell.Blizzard)) MagicInfoList.Add(new MagicInfo { Name = "Blizzard", Spell = Spell.Blizzard, Icon = 51, Level1 = 44, Level2 = 47, Level3 = 50, Need1 = 8000, Need2 = 16000, Need3 = 24000, BaseCost = 65, LevelCost = 10, MPowerBase = 30, MPowerBonus = 10, PowerBase = 20, PowerBonus = 5, Range = 9 });
            if (!MagicExists(Spell.MagicBooster)) MagicInfoList.Add(new MagicInfo { Name = "MagicBooster", Spell = Spell.MagicBooster, Icon = 73, Level1 = 47, Level2 = 49, Level3 = 52, Need1 = 12000, Need2 = 18000, Need3 = 24000, BaseCost = 150, LevelCost = 15, Range = 0 });
            if (!MagicExists(Spell.MeteorStrike)) MagicInfoList.Add(new MagicInfo { Name = "MeteorStrike", Spell = Spell.MeteorStrike, Icon = 52, Level1 = 49, Level2 = 52, Level3 = 55, Need1 = 15000, Need2 = 20000, Need3 = 25000, BaseCost = 115, LevelCost = 17, MPowerBase = 40, MPowerBonus = 10, PowerBase = 20, PowerBonus = 15, Range = 9 });
            if (!MagicExists(Spell.IceThrust)) MagicInfoList.Add(new MagicInfo { Name = "IceThrust", Spell = Spell.IceThrust, Icon = 56, Level1 = 53, Level2 = 56, Level3 = 59, Need1 = 17000, Need2 = 22000, Need3 = 27000, BaseCost = 100, LevelCost = 20, MPowerBase = 100, PowerBase = 50, Range = 0 });
            if (!MagicExists(Spell.Blink)) MagicInfoList.Add(new MagicInfo { Name = "Blink", Spell = Spell.Blink, Icon = 20, Level1 = 19, Level2 = 22, Level3 = 25, Need1 = 350, Need2 = 1000, Need3 = 2000, BaseCost = 10, LevelCost = 3, Range = 9 });
            //if (!MagicExists(Spell.FastMove)) MagicInfoList.Add(new MagicInfo { Name = "FastMove", Spell = Spell.ImmortalSkin, Icon = ?, Level1 = ?, Level2 = ?, Level3 = ?, Need1 = ?, Need2 = ?, Need3 = ?, BaseCost = ?, LevelCost = ?, DelayBase = ?, DelayReduction = ? });
            if (!MagicExists(Spell.StormEscape)) MagicInfoList.Add(new MagicInfo { Name = "StormEscape", Spell = Spell.StormEscape, Icon = 23, Level1 = 60, Level2 = 61, Level3 = 62, Need1 = 2200, Need2 = 3300, Need3 = 4400, BaseCost = 65, LevelCost = 8, MPowerBase = 12, PowerBase = 4, Range = 9 });


            //Taoist
            if (!MagicExists(Spell.Healing)) MagicInfoList.Add(new MagicInfo { Name = "Healing", Spell = Spell.Healing, Icon = 1, Level1 = 7, Level2 = 11, Level3 = 14, Need1 = 150, Need2 = 350, Need3 = 700, BaseCost = 3, LevelCost = 2, MPowerBase = 14, Range = 9 });
            if (!MagicExists(Spell.SpiritSword)) MagicInfoList.Add(new MagicInfo { Name = "SpiritSword", Spell = Spell.SpiritSword, Icon = 3, Level1 = 9, Level2 = 12, Level3 = 15, Need1 = 350, Need2 = 1300, Need3 = 2700, Range = 0 });
            if (!MagicExists(Spell.Poisoning)) MagicInfoList.Add(new MagicInfo { Name = "Poisoning", Spell = Spell.Poisoning, Icon = 5, Level1 = 14, Level2 = 17, Level3 = 20, Need1 = 700, Need2 = 1300, Need3 = 2700, BaseCost = 2, LevelCost = 1, Range = 9 });
            if (!MagicExists(Spell.SoulFireBall)) MagicInfoList.Add(new MagicInfo { Name = "SoulFireBall", Spell = Spell.SoulFireBall, Icon = 12, Level1 = 18, Level2 = 21, Level3 = 24, Need1 = 1300, Need2 = 2700, Need3 = 4000, BaseCost = 3, LevelCost = 1, MPowerBase = 8, PowerBase = 3, Range = 9 });
            if (!MagicExists(Spell.SummonSkeleton)) MagicInfoList.Add(new MagicInfo { Name = "SummonSkeleton", Spell = Spell.SummonSkeleton, Icon = 16, Level1 = 19, Level2 = 22, Level3 = 26, Need1 = 1000, Need2 = 2000, Need3 = 3500, BaseCost = 12, LevelCost = 4, Range = 0 });
            if (!MagicExists(Spell.Hiding)) MagicInfoList.Add(new MagicInfo { Name = "Hiding", Spell = Spell.Hiding, Icon = 17, Level1 = 20, Level2 = 23, Level3 = 26, Need1 = 1300, Need2 = 2700, Need3 = 5300, BaseCost = 1, LevelCost = 1, Range = 0 });
            if (!MagicExists(Spell.MassHiding)) MagicInfoList.Add(new MagicInfo { Name = "MassHiding", Spell = Spell.MassHiding, Icon = 18, Level1 = 21, Level2 = 25, Level3 = 29, Need1 = 1300, Need2 = 2700, Need3 = 5300, BaseCost = 2, LevelCost = 2, Range = 9 });
            if (!MagicExists(Spell.SoulShield)) MagicInfoList.Add(new MagicInfo { Name = "SoulShield", Spell = Spell.SoulShield, Icon = 13, Level1 = 22, Level2 = 24, Level3 = 26, Need1 = 2000, Need2 = 3500, Need3 = 7000, BaseCost = 2, LevelCost = 2, Range = 9 });
            if (!MagicExists(Spell.Revelation)) MagicInfoList.Add(new MagicInfo { Name = "Revelation", Spell = Spell.Revelation, Icon = 27, Level1 = 23, Level2 = 25, Level3 = 28, Need1 = 1500, Need2 = 2500, Need3 = 4000, BaseCost = 4, LevelCost = 4, Range = 9 });
            if (!MagicExists(Spell.BlessedArmour)) MagicInfoList.Add(new MagicInfo { Name = "BlessedArmour", Spell = Spell.BlessedArmour, Icon = 14, Level1 = 25, Level2 = 27, Level3 = 29, Need1 = 4000, Need2 = 6000, Need3 = 10000, BaseCost = 2, LevelCost = 2, Range = 9 });
            if (!MagicExists(Spell.EnergyRepulsor)) MagicInfoList.Add(new MagicInfo { Name = "EnergyRepulsor", Spell = Spell.EnergyRepulsor, Icon = 36, Level1 = 27, Level2 = 29, Level3 = 31, Need1 = 1800, Need2 = 2400, Need3 = 3200, BaseCost = 2, LevelCost = 2, Range = 0, MPowerBase = 4 });
            if (!MagicExists(Spell.TrapHexagon)) MagicInfoList.Add(new MagicInfo { Name = "TrapHexagon", Spell = Spell.TrapHexagon, Icon = 15, Level1 = 28, Level2 = 30, Level3 = 32, Need1 = 2500, Need2 = 5000, Need3 = 10000, BaseCost = 7, LevelCost = 3, Range = 9 });
            if (!MagicExists(Spell.Purification)) MagicInfoList.Add(new MagicInfo { Name = "Purification", Spell = Spell.Purification, Icon = 39, Level1 = 30, Level2 = 32, Level3 = 35, Need1 = 3000, Need2 = 5000, Need3 = 8000, BaseCost = 14, LevelCost = 2, Range = 9 });
            if (!MagicExists(Spell.MassHealing)) MagicInfoList.Add(new MagicInfo { Name = "MassHealing", Spell = Spell.MassHealing, Icon = 28, Level1 = 31, Level2 = 33, Level3 = 36, Need1 = 2000, Need2 = 4000, Need3 = 8000, BaseCost = 28, LevelCost = 3, MPowerBase = 10, PowerBase = 4, Range = 9 });
            if (!MagicExists(Spell.Hallucination)) MagicInfoList.Add(new MagicInfo { Name = "Hallucination", Spell = Spell.Hallucination, Icon = 48, Level1 = 31, Level2 = 34, Level3 = 36, Need1 = 4000, Need2 = 6000, Need3 = 9000, BaseCost = 22, LevelCost = 10, Range = 9 });
            if (!MagicExists(Spell.UltimateEnhancer)) MagicInfoList.Add(new MagicInfo { Name = "UltimateEnchancer", Spell = Spell.UltimateEnhancer, Icon = 35, Level1 = 33, Level2 = 35, Level3 = 38, Need1 = 5000, Need2 = 7000, Need3 = 10000, BaseCost = 28, LevelCost = 4, Range = 9 });
            if (!MagicExists(Spell.SummonShinsu)) MagicInfoList.Add(new MagicInfo { Name = "SummonShinsu", Spell = Spell.SummonShinsu, Icon = 29, Level1 = 35, Level2 = 37, Level3 = 40, Need1 = 2000, Need2 = 4000, Need3 = 6000, BaseCost = 28, LevelCost = 4, Range = 0 });
            if (!MagicExists(Spell.Reincarnation)) MagicInfoList.Add(new MagicInfo { Name = "Reincarnation", Spell = Spell.Reincarnation, Icon = 53, Level1 = 37, Level2 = 39, Level3 = 41, Need1 = 2000, Need2 = 6000, Need3 = 10000, BaseCost = 125, LevelCost = 17, Range = 9 });
            if (!MagicExists(Spell.SummonHolyDeva)) MagicInfoList.Add(new MagicInfo { Name = "SummonHolyDeva", Spell = Spell.SummonHolyDeva, Icon = 40, Level1 = 38, Level2 = 41, Level3 = 43, Need1 = 4000, Need2 = 6000, Need3 = 9000, BaseCost = 28, LevelCost = 4, Range = 0 });
            if (!MagicExists(Spell.Curse)) MagicInfoList.Add(new MagicInfo { Name = "Curse", Spell = Spell.Curse, Icon = 45, Level1 = 40, Level2 = 42, Level3 = 44, Need1 = 4000, Need2 = 6000, Need3 = 9000, BaseCost = 17, LevelCost = 3, Range = 9, MPowerBase = 20 });
            if (!MagicExists(Spell.Plague)) MagicInfoList.Add(new MagicInfo { Name = "Plague", Spell = Spell.Plague, Icon = 74, Level1 = 42, Level2 = 44, Level3 = 47, Need1 = 5000, Need2 = 9000, Need3 = 13000, BaseCost = 20, LevelCost = 5, Range = 9 });
            if (!MagicExists(Spell.PoisonCloud)) MagicInfoList.Add(new MagicInfo { Name = "PoisonCloud", Spell = Spell.PoisonCloud, Icon = 54, Level1 = 43, Level2 = 45, Level3 = 48, Need1 = 4000, Need2 = 8000, Need3 = 12000, BaseCost = 30, LevelCost = 5, MPowerBase = 40, PowerBase = 20, DelayBase = 18000, DelayReduction = 2000, Range = 9 });
            if (!MagicExists(Spell.EnergyShield)) MagicInfoList.Add(new MagicInfo { Name = "EnergyShield", Spell = Spell.EnergyShield, Icon = 57, Level1 = 48, Level2 = 51, Level3 = 54, Need1 = 5000, Need2 = 9000, Need3 = 13000, BaseCost = 50, LevelCost = 20, Range = 9 });
            if (!MagicExists(Spell.PetEnhancer)) MagicInfoList.Add(new MagicInfo { Name = "PetEnhancer", Spell = Spell.PetEnhancer, Icon = 78, Level1 = 45, Level2 = 48, Level3 = 51, Need1 = 4000, Need2 = 8000, Need3 = 12000, BaseCost = 30, LevelCost = 40, Range = 0 });
            //if (!MagicExists(Spell.HealingCircle)) MagicInfoList.Add(new MagicInfo { Name = "HealingCircle", Spell = Spell.ImmortalSkin, Icon = ?, Level1 = ?, Level2 = ?, Level3 = ?, Need1 = ?, Need2 = ?, Need3 = ?, BaseCost = ?, LevelCost = ?, DelayBase = ?, DelayReduction = ? });

            //Assassin
            if (!MagicExists(Spell.FatalSword)) MagicInfoList.Add(new MagicInfo { Name = "FatalSword", Spell = Spell.FatalSword, Icon = 58, Level1 = 7, Level2 = 9, Level3 = 12, Need1 = 500, Need2 = 1000, Need3 = 2300, Range = 0 });
            if (!MagicExists(Spell.DoubleSlash)) MagicInfoList.Add(new MagicInfo { Name = "DoubleSlash", Spell = Spell.DoubleSlash, Icon = 59, Level1 = 15, Level2 = 17, Level3 = 19, Need1 = 700, Need2 = 1500, Need3 = 2200, BaseCost = 2, LevelCost = 1 });
            if (!MagicExists(Spell.Haste)) MagicInfoList.Add(new MagicInfo { Name = "Haste", Spell = Spell.Haste, Icon = 60, Level1 = 20, Level2 = 22, Level3 = 25, Need1 = 2000, Need2 = 3000, Need3 = 6000, BaseCost = 3, LevelCost = 2, Range = 0 });
            if (!MagicExists(Spell.FlashDash)) MagicInfoList.Add(new MagicInfo { Name = "FlashDash", Spell = Spell.FlashDash, Icon = 61, Level1 = 25, Level2 = 27, Level3 = 30, Need1 = 4000, Need2 = 7000, Need3 = 9000, BaseCost = 12, LevelCost = 2, DelayBase = 200, Range = 0 });
            if (!MagicExists(Spell.LightBody)) MagicInfoList.Add(new MagicInfo { Name = "LightBody", Spell = Spell.LightBody, Icon = 68, Level1 = 27, Level2 = 29, Level3 = 32, Need1 = 5000, Need2 = 7000, Need3 = 10000, BaseCost = 11, LevelCost = 2, Range = 0 });
            if (!MagicExists(Spell.HeavenlySword)) MagicInfoList.Add(new MagicInfo { Name = "HeavenlySword", Spell = Spell.HeavenlySword, Icon = 62, Level1 = 30, Level2 = 32, Level3 = 35, Need1 = 4000, Need2 = 8000, Need3 = 10000, BaseCost = 13, LevelCost = 2, MPowerBase = 8, Range = 0 });
            if (!MagicExists(Spell.FireBurst)) MagicInfoList.Add(new MagicInfo { Name = "FireBurst", Spell = Spell.FireBurst, Icon = 63, Level1 = 33, Level2 = 35, Level3 = 38, Need1 = 4000, Need2 = 6000, Need3 = 8000, BaseCost = 10, LevelCost = 1, Range = 0 });
            if (!MagicExists(Spell.Trap)) MagicInfoList.Add(new MagicInfo { Name = "Trap", Spell = Spell.Trap, Icon = 64, Level1 = 33, Level2 = 35, Level3 = 38, Need1 = 2000, Need2 = 4000, Need3 = 6000, BaseCost = 14, LevelCost = 2, DelayBase = 60000, DelayReduction = 15000, Range = 9 });
            if (!MagicExists(Spell.PoisonSword)) MagicInfoList.Add(new MagicInfo { Name = "PoisonSword", Spell = Spell.PoisonSword, Icon = 69, Level1 = 34, Level2 = 36, Level3 = 39, Need1 = 5000, Need2 = 8000, Need3 = 11000, BaseCost = 14, LevelCost = 3, Range = 0 });
            if (!MagicExists(Spell.MoonLight)) MagicInfoList.Add(new MagicInfo { Name = "MoonLight", Spell = Spell.MoonLight, Icon = 65, Level1 = 36, Level2 = 39, Level3 = 42, Need1 = 3000, Need2 = 5000, Need3 = 8000, BaseCost = 36, LevelCost = 3, Range = 0 });
            if (!MagicExists(Spell.MPEater)) MagicInfoList.Add(new MagicInfo { Name = "MPEater", Spell = Spell.MPEater, Icon = 66, Level1 = 38, Level2 = 41, Level3 = 44, Need1 = 5000, Need2 = 8000, Need3 = 11000, Range = 0 });
            if (!MagicExists(Spell.SwiftFeet)) MagicInfoList.Add(new MagicInfo { Name = "SwiftFeet", Spell = Spell.SwiftFeet, Icon = 67, Level1 = 40, Level2 = 43, Level3 = 46, Need1 = 4000, Need2 = 6000, Need3 = 9000, BaseCost = 17, LevelCost = 5, DelayBase = 210000, DelayReduction = 40000, Range = 0 });
            if (!MagicExists(Spell.DarkBody)) MagicInfoList.Add(new MagicInfo { Name = "DarkBody", Spell = Spell.DarkBody, Icon = 70, Level1 = 46, Level2 = 49, Level3 = 52, Need1 = 6000, Need2 = 10000, Need3 = 14000, BaseCost = 40, LevelCost = 7, Range = 0 });
            if (!MagicExists(Spell.Hemorrhage)) MagicInfoList.Add(new MagicInfo { Name = "Hemorrhage", Spell = Spell.Hemorrhage, Icon = 75, Level1 = 47, Level2 = 51, Level3 = 55, Need1 = 9000, Need2 = 15000, Need3 = 21000, Range = 0 });
            if (!MagicExists(Spell.CrescentSlash)) MagicInfoList.Add(new MagicInfo { Name = "CresentSlash", Spell = Spell.CrescentSlash, Icon = 71, Level1 = 50, Level2 = 53, Level3 = 56, Need1 = 12000, Need2 = 16000, Need3 = 24000, BaseCost = 19, LevelCost = 5, Range = 0 });
            //if (!MagicExists(Spell.MoonMist)) MagicInfoList.Add(new MagicInfo { Name = "MoonMist", Spell = Spell.ImmortalSkin, Icon = ?, Level1 = ?, Level2 = ?, Level3 = ?, Need1 = ?, Need2 = ?, Need3 = ?, BaseCost = ?, LevelCost = ?, DelayBase = ?, DelayReduction = ? });

            //Archer
            if (!MagicExists(Spell.Focus)) MagicInfoList.Add(new MagicInfo { Name = "Focus", Spell = Spell.Focus, Icon = 88, Level1 = 7, Level2 = 13, Level3 = 17, Need1 = 270, Need2 = 600, Need3 = 1300, Range = 0 });
            if (!MagicExists(Spell.StraightShot)) MagicInfoList.Add(new MagicInfo { Name = "StraightShot", Spell = Spell.StraightShot, Icon = 89, Level1 = 9, Level2 = 12, Level3 = 16, Need1 = 350, Need2 = 750, Need3 = 1400, BaseCost = 3, LevelCost = 2, MPowerBase = 8, PowerBase = 3, Range = 9 });
            if (!MagicExists(Spell.DoubleShot)) MagicInfoList.Add(new MagicInfo { Name = "DoubleShot", Spell = Spell.DoubleShot, Icon = 90, Level1 = 14, Level2 = 18, Level3 = 21, Need1 = 700, Need2 = 1500, Need3 = 2100, BaseCost = 3, LevelCost = 2, MPowerBase = 6, PowerBase = 2, Range = 9 });
            if (!MagicExists(Spell.ExplosiveTrap)) MagicInfoList.Add(new MagicInfo { Name = "ExplosiveTrap", Spell = Spell.ExplosiveTrap, Icon = 91, Level1 = 22, Level2 = 25, Level3 = 30, Need1 = 2000, Need2 = 3500, Need3 = 5000, BaseCost = 10, LevelCost = 3, MPowerBase = 15, PowerBase = 15, Range = 0 });
            if (!MagicExists(Spell.DelayedExplosion)) MagicInfoList.Add(new MagicInfo { Name = "DelayedExplosion", Spell = Spell.DelayedExplosion, Icon = 92, Level1 = 31, Level2 = 34, Level3 = 39, Need1 = 3000, Need2 = 7000, Need3 = 10000, BaseCost = 8, LevelCost = 2, MPowerBase = 30, PowerBase = 15, Range = 9 });
            if (!MagicExists(Spell.Meditation)) MagicInfoList.Add(new MagicInfo { Name = "Meditation", Spell = Spell.Meditation, Icon = 93, Level1 = 19, Level2 = 24, Level3 = 29, Need1 = 1800, Need2 = 2600, Need3 = 5600, BaseCost = 8, LevelCost = 2, Range = 0 });
            if (!MagicExists(Spell.ElementalShot)) MagicInfoList.Add(new MagicInfo { Name = "ElementalShot", Spell = Spell.ElementalShot, Icon = 94, Level1 = 20, Level2 = 25, Level3 = 31, Need1 = 1800, Need2 = 2700, Need3 = 6000, BaseCost = 8, LevelCost = 2, MPowerBase = 6, PowerBase = 3, Range = 9 });
            if (!MagicExists(Spell.Concentration)) MagicInfoList.Add(new MagicInfo { Name = "Concentration", Spell = Spell.Concentration, Icon = 96, Level1 = 23, Level2 = 27, Level3 = 32, Need1 = 2100, Need2 = 3800, Need3 = 6500, BaseCost = 8, LevelCost = 2, Range = 0 });
            if (!MagicExists(Spell.ElementalBarrier)) MagicInfoList.Add(new MagicInfo { Name = "ElementalBarrier", Spell = Spell.ElementalBarrier, Icon = 98, Level1 = 33, Level2 = 38, Level3 = 44, Need1 = 3000, Need2 = 7000, Need3 = 10000, BaseCost = 10, LevelCost = 2, MPowerBase = 15, PowerBase = 5, Range = 0 });
            if (!MagicExists(Spell.BackStep)) MagicInfoList.Add(new MagicInfo { Name = "BackStep", Spell = Spell.BackStep, Icon = 95, Level1 = 30, Level2 = 34, Level3 = 38, Need1 = 2400, Need2 = 3000, Need3 = 6000, BaseCost = 12, LevelCost = 2, DelayBase = 2500, Range = 0 });
            if (!MagicExists(Spell.BindingShot)) MagicInfoList.Add(new MagicInfo { Name = "BindingShot", Spell = Spell.BindingShot, Icon = 97, Level1 = 35, Level2 = 39, Level3 = 42, Need1 = 400, Need2 = 7000, Need3 = 9500, BaseCost = 7, LevelCost = 3, Range = 9 });
            if (!MagicExists(Spell.SummonVampire)) MagicInfoList.Add(new MagicInfo { Name = "SummonVampire", Spell = Spell.SummonVampire, Icon = 99, Level1 = 28, Level2 = 33, Level3 = 41, Need1 = 2000, Need2 = 2700, Need3 = 7500, BaseCost = 10, LevelCost = 5, Range = 9 });
            if (!MagicExists(Spell.VampireShot)) MagicInfoList.Add(new MagicInfo { Name = "VampireShot", Spell = Spell.VampireShot, Icon = 100, Level1 = 26, Level2 = 32, Level3 = 36, Need1 = 3000, Need2 = 6000, Need3 = 12000, BaseCost = 12, LevelCost = 3, MPowerBase = 10, PowerBase = 7, Range = 9 });
            if (!MagicExists(Spell.SummonToad)) MagicInfoList.Add(new MagicInfo { Name = "SummonToad", Spell = Spell.SummonToad, Icon = 101, Level1 = 37, Level2 = 43, Level3 = 47, Need1 = 5800, Need2 = 10000, Need3 = 13000, BaseCost = 10, LevelCost = 5, Range = 9 });
            if (!MagicExists(Spell.PoisonShot)) MagicInfoList.Add(new MagicInfo { Name = "PoisonShot", Spell = Spell.PoisonShot, Icon = 102, Level1 = 40, Level2 = 45, Level3 = 49, Need1 = 6000, Need2 = 14000, Need3 = 16000, BaseCost = 10, LevelCost = 4, MPowerBase = 10, PowerBase = 10, Range = 9 });
            if (!MagicExists(Spell.CrippleShot)) MagicInfoList.Add(new MagicInfo { Name = "CrippleShot", Spell = Spell.CrippleShot, Icon = 103, Level1 = 43, Level2 = 47, Level3 = 50, Need1 = 12000, Need2 = 15000, Need3 = 18000, BaseCost = 15, LevelCost = 3, MPowerBase = 10, MPowerBonus = 20, PowerBase = 10, Range = 9 });
            if (!MagicExists(Spell.SummonSnakes)) MagicInfoList.Add(new MagicInfo { Name = "SummonSnakes", Spell = Spell.SummonSnakes, Icon = 104, Level1 = 46, Level2 = 51, Level3 = 54, Need1 = 14000, Need2 = 17000, Need3 = 20000, BaseCost = 10, LevelCost = 5, Range = 9 });
            if (!MagicExists(Spell.NapalmShot)) MagicInfoList.Add(new MagicInfo { Name = "NapalmShot", Spell = Spell.NapalmShot, Icon = 105, Level1 = 48, Level2 = 52, Level3 = 55, Need1 = 15000, Need2 = 18000, Need3 = 21000, BaseCost = 40, LevelCost = 10, MPowerBase = 25, MPowerBonus = 25, PowerBase = 25, Range = 9 });
            if (!MagicExists(Spell.OneWithNature)) MagicInfoList.Add(new MagicInfo { Name = "OneWithNature", Spell = Spell.OneWithNature, Icon = 106, Level1 = 50, Level2 = 53, Level3 = 56, Need1 = 17000, Need2 = 19000, Need3 = 24000, BaseCost = 80, LevelCost = 15, MPowerBase = 75, MPowerBonus = 35, PowerBase = 30, PowerBonus = 20, Range = 9 });
            if (!MagicExists(Spell.MentalState)) MagicInfoList.Add(new MagicInfo { Name = "MentalState", Spell = Spell.MentalState, Icon = 81, Level1 = 11, Level2 = 15, Level3 = 22, Need1 = 500, Need2 = 900, Need3 = 1800, BaseCost = 1, LevelCost = 1, Range = 0 });

            //Custom
            if (!MagicExists(Spell.Portal)) MagicInfoList.Add(new MagicInfo { Name = "Portal", Spell = Spell.Portal, Icon = 1, Level1 = 7, Level2 = 11, Level3 = 14, Need1 = 150, Need2 = 350, Need3 = 700, BaseCost = 3, LevelCost = 2, Range = 9 });
            if (!MagicExists(Spell.BattleCry)) MagicInfoList.Add(new MagicInfo { Name = "BattleCry", Spell = Spell.BattleCry, Icon = 42, Level1 = 48, Level2 = 51, Level3 = 55, Need1 = 8000, Need2 = 11000, Need3 = 15000, BaseCost = 22, LevelCost = 10, Range = 0 });
        }

        private void UpdateMagicInfo()
        {
            for (var i = 0; i < MagicInfoList.Count; i++)
            {
                switch (MagicInfoList[i].Spell)
                {
                    //warrior
                    case Spell.Thrusting:
                        MagicInfoList[i].MultiplierBase = 0.25f;
                        MagicInfoList[i].MultiplierBonus = 0.25f;
                        break;
                    case Spell.HalfMoon:
                        MagicInfoList[i].MultiplierBase = 0.3f;
                        MagicInfoList[i].MultiplierBonus = 0.1f;
                        break;
                    case Spell.ShoulderDash:
                        MagicInfoList[i].MPowerBase = 4;
                        break;
                    case Spell.TwinDrakeBlade:
                        MagicInfoList[i].MultiplierBase = 0.8f;
                        MagicInfoList[i].MultiplierBonus = 0.1f;
                        break;
                    case Spell.FlamingSword:
                        MagicInfoList[i].MultiplierBase = 1.4f;
                        MagicInfoList[i].MultiplierBonus = 0.4f;
                        break;
                    case Spell.CrossHalfMoon:
                        MagicInfoList[i].MultiplierBase = 0.4f;
                        MagicInfoList[i].MultiplierBonus = 0.1f;
                        break;
                    case Spell.BladeAvalanche:
                        MagicInfoList[i].MultiplierBase = 1f;
                        MagicInfoList[i].MultiplierBonus = 0.4f;
                        break;
                    case Spell.SlashingBurst:
                        MagicInfoList[i].MultiplierBase = 3.25f;
                        MagicInfoList[i].MultiplierBonus = 0.25f;
                        break;
                    //wiz
                    case Spell.Repulsion:
                        MagicInfoList[i].MPowerBase = 4;
                        break;
                    //tao
                    case Spell.Poisoning:
                        MagicInfoList[i].MPowerBase = 0;
                        break;
                    case Spell.Curse:
                        MagicInfoList[i].MPowerBase = 20;
                        break;
                    case Spell.Plague:
                        MagicInfoList[i].MPowerBase = 0;
                        MagicInfoList[i].PowerBase = 0;
                        break;
                    //sin
                    case Spell.FatalSword:
                        MagicInfoList[i].MPowerBase = 20;
                        break;
                    case Spell.DoubleSlash:
                        MagicInfoList[i].MultiplierBase = 0.8f;
                        MagicInfoList[i].MultiplierBonus = 0.1f;
                        break;
                    case Spell.FireBurst:
                        MagicInfoList[i].MPowerBase = 4;
                        break;
                    case Spell.MoonLight:
                        MagicInfoList[i].MPowerBase = 20;
                        break;
                    case Spell.DarkBody:
                        MagicInfoList[i].MPowerBase = 20;
                        break;
                    case Spell.Hemorrhage:
                        MagicInfoList[i].MultiplierBase = 0.2f;
                        MagicInfoList[i].MultiplierBonus = 0.05f;
                        break;
                    case Spell.CrescentSlash:
                        MagicInfoList[i].MultiplierBase = 1f;
                        MagicInfoList[i].MultiplierBonus = 0.4f;
                        break;
                        //archer
                        //no changes :p
                }
            }
        }
        #endregion

        // TODO 从数据库中读取数据
        public void loadFromDatabase()
        {
        }

        // 将数据保存到数据库
        public void saveDataToDatabase()
        {
            // 保存基础信息
            var basicModel = new BasicModel()
            {
                Version = LoadVersion,
                CustomVersion = LoadCustomVersion,
                MapIndex = MapIndex,
                ItemIndex = ItemIndex,
                MonsterIndex = MonsterIndex,
                NPCIndex = NPCIndex,
                QuestIndex = QuestIndex,
                GameshopIndex = GameshopIndex,
                ConquestIndex = ConquestIndex,
                RespawnIndex = RespawnIndex
            };
            Manager.DB.Insertable(basicModel).ExecuteCommand();

            // 保存地图信息
            for (var i = 0; i < MapInfoList.Count; i++)
                MapInfoList[i].Save();

            // 保存物品信息
            for (var i = 0; i < ItemInfoList.Count; i++)
                ItemInfoList[i].Save();

            // 保存怪物信息
            for (var i = 0; i < MonsterInfoList.Count; i++)
                MonsterInfoList[i].Save();

            // 保存NPC信息
            for (var i = 0; i < NPCInfoList.Count; i++)
                NPCInfoList[i].Save();

            // 保存任务信息
            for (var i = 0; i < QuestInfoList.Count; i++)
                QuestInfoList[i].Save();

            // FIXME 保存龙信息 感觉是后来维护者乱 JB 加的, 懒得保存了
            // DragonInfo.Save();

            // 保存魔法信息
            for (var i = 0; i < MagicInfoList.Count; i++)
                MagicInfoList[i].Save();

            // 保存商店信息 ?不确定是不是商城
            for (var i = 0; i < GameShopItemList.Count; i++)
                GameShopItemList[i].Save();

            // FIXME 从文件读出来是空的, 目前看来没什么用
            // 保存土城攻城信息
            // for (var i = 0; i < ConquestInfos.Count; i++)
            //     ConquestInfos[i].Save();

            // FIXME 从文件读出来是空的, 目前看来没什么用
            // 保存刷怪时钟信息
            // RespawnTick.Save();
        }
    }
}