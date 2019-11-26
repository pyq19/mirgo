using System;
using System.IO;
using System.Drawing;
using System.Collections.Generic;

namespace dotnettools
{
    public class ConquestInfo
    {
        public int Index;
        public bool FullMap;
        public Point Location;
        public ushort Size;
        public string Name;
        public int MapIndex;
        public int PalaceIndex;
        public List<int> ExtraMaps = new List<int>();
        public List<ConquestArcherInfo> ConquestGuards = new List<ConquestArcherInfo>();
        public List<ConquestGateInfo> ConquestGates = new List<ConquestGateInfo>();
        public List<ConquestWallInfo> ConquestWalls = new List<ConquestWallInfo>();
        public List<ConquestSiegeInfo> ConquestSieges = new List<ConquestSiegeInfo>();
        public List<ConquestFlagInfo> ConquestFlags = new List<ConquestFlagInfo>();

        public int GuardIndex;
        public int GateIndex;
        public int WallIndex;
        public int SiegeIndex;
        public int FlagIndex;

        public byte StartHour = 0;
        public int WarLength = 60;

        private int counter;

        public ConquestType Type = ConquestType.Request;
        public ConquestGame Game = ConquestGame.CapturePalace;

        public bool Monday;
        public bool Tuesday;
        public bool Wednesday;
        public bool Thursday;
        public bool Friday;
        public bool Saturday;
        public bool Sunday;

        //King of the hill
        public Point KingLocation;
        public ushort KingSize;

        //Control points
        public List<ConquestFlagInfo> ControlPoints = new List<ConquestFlagInfo>();
        public int ControlPointIndex;

        public ConquestInfo(BinaryReader reader, Manager manager)
        {
            Manager Envir = manager;

            Index = reader.ReadInt32();

            if (Envir.LoadVersion > 73)
            {
                FullMap = reader.ReadBoolean();
            }

            Location = new Point(reader.ReadInt32(), reader.ReadInt32());
            Size = reader.ReadUInt16();
            Name = reader.ReadString();
            MapIndex = reader.ReadInt32();
            PalaceIndex = reader.ReadInt32();
            GuardIndex = reader.ReadInt32();
            GateIndex = reader.ReadInt32();
            WallIndex = reader.ReadInt32();
            SiegeIndex = reader.ReadInt32();

            if (Envir.LoadVersion > 72)
            {
                FlagIndex = reader.ReadInt32();
            }

            counter = reader.ReadInt32();
            for (int i = 0; i < counter; i++)
            {
                ConquestGuards.Add(new ConquestArcherInfo(reader));
            }
            counter = reader.ReadInt32();
            for (int i = 0; i < counter; i++)
            {
                ExtraMaps.Add(reader.ReadInt32());
            }
            counter = reader.ReadInt32();
            for (int i = 0; i < counter; i++)
            {
                ConquestGates.Add(new ConquestGateInfo(reader));
            }
            counter = reader.ReadInt32();
            for (int i = 0; i < counter; i++)
            {
                ConquestWalls.Add(new ConquestWallInfo(reader));
            }
            counter = reader.ReadInt32();
            for (int i = 0; i < counter; i++)
            {
                ConquestSieges.Add(new ConquestSiegeInfo(reader));
            }

            if (Envir.LoadVersion > 72)
            {
                counter = reader.ReadInt32();
                for (int i = 0; i < counter; i++)
                {
                    ConquestFlags.Add(new ConquestFlagInfo(reader));
                }
            }

            StartHour = reader.ReadByte();
            WarLength = reader.ReadInt32();
            Type = (ConquestType)reader.ReadByte();
            Game = (ConquestGame)reader.ReadByte();

            Monday = reader.ReadBoolean();
            Tuesday = reader.ReadBoolean();
            Wednesday = reader.ReadBoolean();
            Thursday = reader.ReadBoolean();
            Friday = reader.ReadBoolean();
            Saturday = reader.ReadBoolean();
            Sunday = reader.ReadBoolean();

            KingLocation = new Point(reader.ReadInt32(), reader.ReadInt32());
            KingSize = reader.ReadUInt16();

            if (Envir.LoadVersion > 74)
            {
                ControlPointIndex = reader.ReadInt32();
                counter = reader.ReadInt32();
                for (int i = 0; i < counter; i++)
                {
                    ControlPoints.Add(new ConquestFlagInfo(reader));
                }
            }
        }

    }

    public class ConquestSiegeInfo
    {
        public int Index;
        public Point Location;
        public int MobIndex;
        public string Name;
        public uint RepairCost;
        public ConquestSiegeInfo(BinaryReader reader)
        {
            Index = reader.ReadInt32();
            Location = new Point(reader.ReadInt32(), reader.ReadInt32());
            MobIndex = reader.ReadInt32();
            Name = reader.ReadString();
            RepairCost = reader.ReadUInt32();
        }
    }

    public class ConquestWallInfo
    {
        public int Index;
        public Point Location;
        public int MobIndex;
        public string Name;
        public uint RepairCost;
        public ConquestWallInfo(BinaryReader reader)
        {
            Index = reader.ReadInt32();
            Location = new Point(reader.ReadInt32(), reader.ReadInt32());
            MobIndex = reader.ReadInt32();
            Name = reader.ReadString();
            RepairCost = reader.ReadUInt32();
        }
    }

    public class ConquestGateInfo
    {
        public int Index;
        public Point Location;
        public int MobIndex;
        public string Name;
        public uint RepairCost;

        public ConquestGateInfo(BinaryReader reader)
        {
            Index = reader.ReadInt32();
            Location = new Point(reader.ReadInt32(), reader.ReadInt32());
            MobIndex = reader.ReadInt32();
            Name = reader.ReadString();
            RepairCost = reader.ReadUInt32();
        }
    }

    public class ConquestArcherInfo
    {
        public int Index;
        public Point Location;
        public int MobIndex;
        public string Name;
        public uint RepairCost;

        public ConquestArcherInfo(BinaryReader reader)
        {
            Index = reader.ReadInt32();
            Location = new Point(reader.ReadInt32(), reader.ReadInt32());
            MobIndex = reader.ReadInt32();
            Name = reader.ReadString();
            RepairCost = reader.ReadUInt32();
        }

    }

    public class ConquestFlagInfo
    {
        public int Index;
        public Point Location;
        public string Name;
        public string FileName = string.Empty;

        public ConquestFlagInfo(BinaryReader reader)
        {
            Index = reader.ReadInt32();
            Location = new Point(reader.ReadInt32(), reader.ReadInt32());
            Name = reader.ReadString();
            FileName = reader.ReadString();
        }
    }

}