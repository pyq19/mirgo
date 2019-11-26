using System;
using System.Collections.Generic;
using System.Drawing;
using System.IO;
using System.Linq;
using System.Text;

namespace dotnettools
{
    public class SafeZoneInfo
    {
        public Point Location;
        public ushort Size;
        public bool StartPoint;
        public MapInfo Info;

        public SafeZoneInfo(BinaryReader reader)
        {
            Location = new Point(reader.ReadInt32(), reader.ReadInt32());
            Size = reader.ReadUInt16();
            StartPoint = reader.ReadBoolean();
        }

        public void Save(int mapIndex)
        {
            var safeZoneInfoModel = new SafeZoneInfoModel()
            {
                MapIndex = mapIndex,
                LocationX = Location.X,
                LocationY = Location.Y,
                Size = Size,
                StartPoint = this.StartPoint ? 1 : 0
            };
            Manager.DB.Insertable(safeZoneInfoModel).ExecuteCommand();
        }
    }
}