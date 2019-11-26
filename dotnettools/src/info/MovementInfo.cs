using System;
using System.IO;
using System.Drawing;

namespace dotnettools
{
    public class MovementInfo
    {
        public int MapIndex;
        public Point Source;
        public Point Destination;
        public bool NeedHole;
        public bool NeedMove;
        public int ConquestIndex;

        public MovementInfo(BinaryReader reader, Manager manager)
        {
            Manager Envir = manager;
            MapIndex = reader.ReadInt32();
            Source = new Point(reader.ReadInt32(), reader.ReadInt32());
            Destination = new Point(reader.ReadInt32(), reader.ReadInt32());

            if (Envir.LoadVersion < 16) return;
            NeedHole = reader.ReadBoolean();

            if (Envir.LoadVersion < 48) return;
            NeedMove = reader.ReadBoolean();

            if (Envir.LoadVersion < 69) return;
            ConquestIndex = reader.ReadInt32();
        }

        // TODO
        public void Save(int mapIndex)
        {

        }
    }
}