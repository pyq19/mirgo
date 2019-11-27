using System;
using System.IO;
using System.Drawing;
using System.Collections.Generic;

namespace dotnettools
{

    public class DragonInfo
    {

        public bool Enabled;
        public string MapFileName, MonsterName, BodyName;
        public Point Location, DropAreaTop, DropAreaBottom;
        public List<DragonDropInfo>[] Drops = new List<DragonDropInfo>[13];
        public long[] Exps = new long[12];
        public byte Level;
        public long Experience;

        public DragonInfo() { }

        public DragonInfo(BinaryReader reader)
        {
            Enabled = reader.ReadBoolean();
            MapFileName = reader.ReadString();
            MonsterName = reader.ReadString();
            BodyName = reader.ReadString();
            Location = new Point(reader.ReadInt32(), reader.ReadInt32());
            DropAreaTop = new Point(reader.ReadInt32(), reader.ReadInt32());
            DropAreaBottom = new Point(reader.ReadInt32(), reader.ReadInt32());

            Level = 1;

            for (int i = 0; i < Exps.Length; i++)
            {
                Exps[i] = reader.ReadInt64();
            }
            for (int i = 0; i < Drops.Length; i++)
            {
                Drops[i] = new List<DragonDropInfo>();
            }
        }

        // FIXME 我觉得 DragonInfo.cs 这是 c# mir2 代码后来的维护者乱鸡儿加的, 懒得存了
        public void Save()
        {

        }

    }

    public class DragonDropInfo
    {
        public int Chance;
        public ItemInfo Item;
        public uint Gold;
        public byte level;
    }
}