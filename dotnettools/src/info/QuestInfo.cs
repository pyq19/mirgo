using System;
using System.IO;
using System.Collections.Generic;

namespace dotnettools
{
    public class QuestInfo
    {

        public int QuestIndex;
        public uint NpcIndex;
        public NPCInfo NpcInfo;
        private uint _finishNpcIndex;
        public uint FinishNpcIndex
        {
            get { return _finishNpcIndex == 0 ? NpcIndex : _finishNpcIndex; }
            set { _finishNpcIndex = value; }
        }
        public string Name = string.Empty;
        public string Group = string.Empty;
        public string FileName = string.Empty;
        public string GotoMessage = string.Empty;
        public string KillMessage = string.Empty;
        public string ItemMessage = string.Empty;
        public string FlagMessage = string.Empty;

        public List<string> Description = new List<string>();
        public List<string> TaskDescription = new List<string>();
        public List<string> CompletionDescription = new List<string>();
        public int RequiredMinLevel, RequiredMaxLevel, RequiredQuest;
        public RequiredClass RequiredClass = RequiredClass.None;
        public QuestType Type;
        public List<QuestItemTask> CarryItems = new List<QuestItemTask>();
        public List<QuestKillTask> KillTasks = new List<QuestKillTask>();
        public List<QuestItemTask> ItemTasks = new List<QuestItemTask>();
        public List<QuestFlagTask> FlagTasks = new List<QuestFlagTask>();
        public uint GoldReward;
        public uint ExpReward;
        public uint CreditReward;
        public List<QuestItemReward> FixedRewards = new List<QuestItemReward>();
        public List<QuestItemReward> SelectRewards = new List<QuestItemReward>();

        public QuestInfo(BinaryReader reader, Manager manager)
        {
            Manager Envir = manager;

            QuestIndex = reader.ReadInt32();
            Name = reader.ReadString();
            Group = reader.ReadString();
            FileName = reader.ReadString();
            RequiredMinLevel = reader.ReadInt32();

            if (Envir.LoadVersion >= 38)
            {
                RequiredMaxLevel = reader.ReadInt32();
                if (RequiredMaxLevel == 0) RequiredMaxLevel = ushort.MaxValue;
            }

            RequiredQuest = reader.ReadInt32();
            RequiredClass = (RequiredClass)reader.ReadByte();
            Type = (QuestType)reader.ReadByte();
            GotoMessage = reader.ReadString();
            KillMessage = reader.ReadString();
            ItemMessage = reader.ReadString();
            if (Envir.LoadVersion >= 37) FlagMessage = reader.ReadString();

            // LoadInfo();
        }

        // TODO 从文件中读取任务
        public void LoadInfo(bool clear = false)
        {

        }

        public void Save()
        {
            var questInfoModel = new QuestInfoModel()
            {
                QuestIndex = QuestIndex,
                Name = Name,
                Group = Group,
                FileName = FileName,
                RequiredMinLevel = RequiredMinLevel,
                RequiredMaxLevel = RequiredMaxLevel,
                RequiredQuest = RequiredQuest,
                RequiredClass = (byte)RequiredClass,
                Type = (byte)Type,
                GotoMessage = GotoMessage,
                KillMessage = KillMessage,
                ItemMessage = ItemMessage,
                FlagMessage = FlagMessage,
            };
            Manager.DB.Insertable(questInfoModel).ExecuteCommand();
        }

    }

    public class QuestKillTask
    {
        public MonsterInfo Monster;
        public int Count;
        public string Message;
    }

    public class QuestItemTask
    {
        public ItemInfo Item;
        public uint Count;
        public string Message;
    }

    public class QuestFlagTask
    {
        public int Number;
        public string Message;
    }
}