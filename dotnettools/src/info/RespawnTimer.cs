using System;
using System.IO;
using System.Collections.Generic;

namespace dotnettools
{
    public class RespawnTimer
    {
        public byte BaseSpawnRate = 20;//amount of minutes between respawnticks (with no bonus)
        public ulong CurrentTickcounter = 0; //counter used to respawn everything
        public long LastTick = 0; //what 'time' was the last tick?
        public int LastUsercount = 0; //stops it from having to check delay each time
        public long CurrentDelay = 0;
        public List<RespawnTickOption> Respawn = new List<RespawnTickOption>();

        public RespawnTimer(BinaryReader reader, Manager manager)
        {
            Manager Envir = manager;
            BaseSpawnRate = reader.ReadByte();
            CurrentTickcounter = reader.ReadUInt64();
            LastTick = Envir.Time;
            Respawn.Clear();
            int Optioncount = reader.ReadInt32();
            for (int i = 0; i < Optioncount; i++)
            {
                RespawnTickOption Option = new RespawnTickOption(reader);
                Respawn.Add(Option);
            }
            CurrentDelay = (long)Math.Round((double)BaseSpawnRate * (double)60000);
        }
    }

    public class RespawnTickOption
    {
        public int UserCount = 1;
        public double DelayLoss = 1.0;

        public RespawnTickOption(BinaryReader reader)
        {
            UserCount = reader.ReadInt32();
            DelayLoss = reader.ReadDouble();
        }
    }
}