using System.Drawing;
using System.IO;

namespace dotnettools
{
    public class RespawnInfo
    {
        public Manager Envir;

        public int MonsterIndex;
        public Point Location;
        public ushort Count;
        public ushort Spread;
        public ushort Delay;
        public ushort RandomDelay;
        public byte Direction;
        public string RoutePath;
        public int RespawnIndex;
        public bool SaveRespawnTime;
        public ushort RespawnTicks;

        public RespawnInfo(BinaryReader reader, int Version, int Customversion, Manager manager)
        {
            Envir = manager;

            MonsterIndex = reader.ReadInt32();
            Location = new Point(reader.ReadInt32(), reader.ReadInt32());

            Count = reader.ReadUInt16();
            Spread = reader.ReadUInt16();

            Delay = reader.ReadUInt16();
            Direction = reader.ReadByte();

            if (Envir.LoadVersion >= 36)
            {
                RoutePath = reader.ReadString();
            }

            if (Version > 67)
            {
                RandomDelay = reader.ReadUInt16();
                RespawnIndex = reader.ReadInt32();
                SaveRespawnTime = reader.ReadBoolean();
                RespawnTicks = reader.ReadUInt16();
            }
            else
            {
                RespawnIndex = ++Envir.RespawnIndex;
            }
        }
    }
}
