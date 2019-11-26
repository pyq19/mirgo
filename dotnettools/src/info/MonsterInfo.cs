using System;
using System.IO;
using System.Collections.Generic;

namespace dotnettools
{
    public class MonsterInfo
    {

        public int Index;
        public string Name = string.Empty;

        public Monster Image;
        public byte AI, Effect, ViewRange = 7, CoolEye;
        public ushort Level;

        public uint HP;
        public byte Accuracy, Agility, Light;
        public ushort MinAC, MaxAC, MinMAC, MaxMAC, MinDC, MaxDC, MinMC, MaxMC, MinSC, MaxSC;

        public ushort AttackSpeed = 2500, MoveSpeed = 1800;
        public uint Experience;

        public List<MonsterDropInfo> Drops = new List<MonsterDropInfo>();

        public bool CanTame = true, CanPush = true, AutoRev = true, Undead = false;

        public bool HasSpawnScript;
        public bool HasDieScript;

        public MonsterInfo(BinaryReader reader, Manager manager)
        {
            Manager Envir = manager;

            Index = reader.ReadInt32();
            Name = reader.ReadString();

            Image = (Monster)reader.ReadUInt16();
            AI = reader.ReadByte();
            Effect = reader.ReadByte();
            if (Envir.LoadVersion < 62)
            {
                Level = (ushort)reader.ReadByte();
            }
            else
            {
                Level = reader.ReadUInt16();
            }

            ViewRange = reader.ReadByte();
            if (Envir.LoadVersion >= 3) CoolEye = reader.ReadByte();

            HP = reader.ReadUInt32();

            if (Envir.LoadVersion < 62)
            {
                MinAC = (ushort)reader.ReadByte();
                MaxAC = (ushort)reader.ReadByte();
                MinMAC = (ushort)reader.ReadByte();
                MaxMAC = (ushort)reader.ReadByte();
                MinDC = (ushort)reader.ReadByte();
                MaxDC = (ushort)reader.ReadByte();
                MinMC = (ushort)reader.ReadByte();
                MaxMC = (ushort)reader.ReadByte();
                MinSC = (ushort)reader.ReadByte();
                MaxSC = (ushort)reader.ReadByte();
            }
            else
            {
                MinAC = reader.ReadUInt16();
                MaxAC = reader.ReadUInt16();
                MinMAC = reader.ReadUInt16();
                MaxMAC = reader.ReadUInt16();
                MinDC = reader.ReadUInt16();
                MaxDC = reader.ReadUInt16();
                MinMC = reader.ReadUInt16();
                MaxMC = reader.ReadUInt16();
                MinSC = reader.ReadUInt16();
                MaxSC = reader.ReadUInt16();
            }

            Accuracy = reader.ReadByte();
            Agility = reader.ReadByte();
            Light = reader.ReadByte();

            AttackSpeed = reader.ReadUInt16();
            MoveSpeed = reader.ReadUInt16();
            Experience = reader.ReadUInt32();

            if (Envir.LoadVersion < 6)
            {
                reader.BaseStream.Seek(8, SeekOrigin.Current);

                int count = reader.ReadInt32();
                reader.BaseStream.Seek(count * 12, SeekOrigin.Current);
            }

            CanPush = reader.ReadBoolean();
            CanTame = reader.ReadBoolean();

            if (Envir.LoadVersion < 18) return;
            AutoRev = reader.ReadBoolean();
            Undead = reader.ReadBoolean();
        }

        // TODO
        public void Save()
        {

        }
    }

    public class MonsterDropInfo
    {
        public int Chance;
        public ItemInfo Item;
        public uint Gold;

        public byte Type;
        public bool QuestRequired;
    }

}