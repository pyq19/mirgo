using System;
using System.IO;

namespace dotnettools
{
    public class MagicInfo
    {
        public string Name;
        public Spell Spell;
        public byte BaseCost;
        public byte LevelCost;
        public byte Icon;
        public byte Level1;
        public byte Level2;
        public byte Level3;
        public ushort Need1;
        public ushort Need2;
        public ushort Need3;
        public uint DelayBase;
        public uint DelayReduction;
        public ushort PowerBase;
        public ushort PowerBonus;
        public ushort MPowerBase;
        public ushort MPowerBonus;
        public byte Range;
        public float MultiplierBase;
        public float MultiplierBonus;

        public MagicInfo(BinaryReader reader, int version = int.MaxValue, int Customversion = int.MaxValue)
        {
            Name = reader.ReadString();
            Spell = (Spell)reader.ReadByte();
            BaseCost = reader.ReadByte();
            LevelCost = reader.ReadByte();
            Icon = reader.ReadByte();
            Level1 = reader.ReadByte();
            Level2 = reader.ReadByte();
            Level3 = reader.ReadByte();
            Need1 = reader.ReadUInt16();
            Need2 = reader.ReadUInt16();
            Need3 = reader.ReadUInt16();
            DelayBase = reader.ReadUInt32();
            DelayReduction = reader.ReadUInt32();
            PowerBase = reader.ReadUInt16();
            PowerBonus = reader.ReadUInt16();
            MPowerBase = reader.ReadUInt16();
            MPowerBonus = reader.ReadUInt16();

            if (version > 66)
                Range = reader.ReadByte();
            if (version > 70)
            {
                MultiplierBase = reader.ReadSingle();
                MultiplierBonus = reader.ReadSingle();
            }
        }
    }

}