using System.IO;
using System.Drawing;
using System.Collections.Generic;

namespace dotnettools
{
    public class MapInfo
    {
        public Manager manager;
        public int MapIndex;
        public string FileName;
        public string Title;
        public ushort MiniMap;
        public ushort BigMap;
        public ushort Music;
        public LightSetting Light;
        public byte MapDarkLight;
        public byte MineIndex;
        public bool NoTeleport;
        public bool NoReconnect;
        public bool NoRandom;
        public bool NoEscape;
        public bool NoRecall;
        public bool NoDrug;
        public bool NoPosition;
        public bool NoFight;
        public bool NoThrowItem;
        public bool NoDropPlayer;
        public bool NoDropMonster;
        public bool NoNames;
        public bool NoMount;
        public bool NeedBridle;
        public bool Fight;
        public bool NeedHole;
        public bool Fire;
        public bool Lightning;
        public bool NoTownTeleport;
        public bool NoReincarnation;
        public string NoReconnectMap;
        public int FireDamage;
        public int LightningDamage;
        public List<SafeZoneInfo> SafeZones = new List<SafeZoneInfo>();
        public List<MovementInfo> Movements = new List<MovementInfo>();
        public List<RespawnInfo> Respawns = new List<RespawnInfo>();
        public List<NPCInfo> NPCs = new List<NPCInfo>();
        public List<MineZone> MineZones = new List<MineZone>();
        public List<Point> ActiveCoords = new List<Point>();
        public InstanceInfo Instance;

        public MapInfo(BinaryReader reader, Manager manager)
        {
            Manager Envir = manager;
            this.manager = manager;
            MapIndex = reader.ReadInt32();
            FileName = reader.ReadString();
            Title = reader.ReadString();
            MiniMap = reader.ReadUInt16();
            Light = (LightSetting)reader.ReadByte();

            if (Envir.LoadVersion >= 3) BigMap = reader.ReadUInt16();

            int count = reader.ReadInt32();
            for (int i = 0; i < count; i++)
                SafeZones.Add(new SafeZoneInfo(reader) { Info = this });

            count = reader.ReadInt32();
            for (int i = 0; i < count; i++)
                Respawns.Add(new RespawnInfo(reader, Envir.LoadVersion, Envir.LoadCustomVersion, Envir));

            if (Envir.LoadVersion <= 33)
            {
                count = reader.ReadInt32();
                for (int i = 0; i < count; i++)
                    NPCs.Add(new NPCInfo(reader, Envir));
            }

            count = reader.ReadInt32();
            for (int i = 0; i < count; i++)
                Movements.Add(new MovementInfo(reader, Envir, MapIndex));

            if (Envir.LoadVersion < 14) return;

            NoTeleport = reader.ReadBoolean();
            NoReconnect = reader.ReadBoolean();
            NoReconnectMap = reader.ReadString();

            NoRandom = reader.ReadBoolean();
            NoEscape = reader.ReadBoolean();
            NoRecall = reader.ReadBoolean();
            NoDrug = reader.ReadBoolean();
            NoPosition = reader.ReadBoolean();
            NoThrowItem = reader.ReadBoolean();
            NoDropPlayer = reader.ReadBoolean();
            NoDropMonster = reader.ReadBoolean();
            NoNames = reader.ReadBoolean();
            Fight = reader.ReadBoolean();
            if (Envir.LoadVersion == 14) NeedHole = reader.ReadBoolean();
            Fire = reader.ReadBoolean();
            FireDamage = reader.ReadInt32();
            Lightning = reader.ReadBoolean();
            LightningDamage = reader.ReadInt32();
            if (Envir.LoadVersion < 23) return;
            MapDarkLight = reader.ReadByte();
            if (Envir.LoadVersion < 26) return;
            count = reader.ReadInt32();
            for (int i = 0; i < count; i++)
                MineZones.Add(new MineZone(reader));
            if (Envir.LoadVersion < 27) return;
            MineIndex = reader.ReadByte();

            if (Envir.LoadVersion < 33) return;
            NoMount = reader.ReadBoolean();
            NeedBridle = reader.ReadBoolean();

            if (Envir.LoadVersion < 42) return;
            NoFight = reader.ReadBoolean();

            if (Envir.LoadVersion < 53) return;
            Music = reader.ReadUInt16();
            if (Envir.LoadVersion < 78) return;
            NoTownTeleport = reader.ReadBoolean();
            if (Envir.LoadVersion < 79) return;
            NoReincarnation = reader.ReadBoolean();
        }

        public void Save()
        {
            var mapInfoModel = new MapInfoModel()
            {
                Id = MapIndex,
                FileName = FileName,
                Title = Title,
                MiniMap = MiniMap,
                Light = (byte)Light,
                BigMap = BigMap,
                NoTeleport = NoTeleport,
                NoReconnect = NoReconnect,
                NoReconnectMap = NoReconnectMap,
                NoRandom = NoRandom,
                NoEscape = NoEscape,
                NoRecall = NoRecall,
                NoDrug = NoDrug,
                NoPosition = NoPosition,
                NoThrowItem = NoThrowItem,
                NoDropPlayer = NoDropPlayer,
                NoDropMonster = NoDropMonster,
                NoNames = NoNames,
                Fight = Fight,
                Fire = Fire,
                FireDamage = FireDamage,
                Lightning = Lightning,
                LightningDamage = LightningDamage,
                MapDarkLight = MapDarkLight,
                MineIndex = MineIndex,
                NoMount = NoMount,
                NeedBridle = NeedBridle,
                NoFight = NoFight,
                Music = Music,
                NoTownTeleport = NoTownTeleport,
                NoReincarnation = NoReincarnation,
            };
            Manager.DB.Insertable(mapInfoModel).ExecuteCommand();

            for (int i = 0; i < SafeZones.Count; i++)
                SafeZones[i].Save(MapIndex);

            for (int i = 0; i < Respawns.Count; i++)
                Respawns[i].Save(MapIndex);

            for (int i = 0; i < Movements.Count; i++)
                Movements[i].Save(MapIndex);

            for (int i = 0; i < MineZones.Count; i++)
                MineZones[i].Save(MapIndex);
        }
    }

    // TODO
    public class InstanceInfo
    {
    }
}