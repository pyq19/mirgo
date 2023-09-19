package game

var PacketNameMap = make(map[int]string)

func initClientPacketMap() {
	PacketNameMap[CM_CLIENT_VERSION] = "CLIENT_VERSION"
	PacketNameMap[CM_DISCONNECT] = "DISCONNECT"
	PacketNameMap[CM_KEEP_ALIVE] = "KEEP_ALIVE"
	PacketNameMap[CM_NEW_ACCOUNT] = "NEW_ACCOUNT"
	PacketNameMap[CM_CHANGE_PASSWORD] = "CHANGE_PASSWORD"
	PacketNameMap[CM_LOGIN] = "LOGIN"
	PacketNameMap[CM_NEW_CHARACTER] = "NEW_CHARACTER"
	PacketNameMap[CM_DELETE_CHARACTER] = "DELETE_CHARACTER"
	PacketNameMap[CM_START_GAME] = "START_GAME"
	PacketNameMap[CM_LOG_OUT] = "LOG_OUT"
	PacketNameMap[CM_TURN] = "TURN"
	PacketNameMap[CM_WALK] = "WALK"
	PacketNameMap[CM_RUN] = "RUN"
	PacketNameMap[CM_CHAT] = "CHAT"
	PacketNameMap[CM_MOVE_ITEM] = "MOVE_ITEM"
	PacketNameMap[CM_STORE_ITEM] = "STORE_ITEM"
	PacketNameMap[CM_TAKE_BACK_ITEM] = "TAKE_BACK_ITEM"
	PacketNameMap[CM_MERGE_ITEM] = "MERGE_ITEM"
	PacketNameMap[CM_EQUIP_ITEM] = "EQUIP_ITEM"
	PacketNameMap[CM_REMOVE_ITEM] = "REMOVE_ITEM"
	PacketNameMap[CM_REMOVE_SLOT_ITEM] = "REMOVE_SLOT_ITEM"
	PacketNameMap[CM_SPLIT_ITEM] = "SPLIT_ITEM"
	PacketNameMap[CM_USE_ITEM] = "USE_ITEM"
	PacketNameMap[CM_DROP_ITEM] = "DROP_ITEM"
	PacketNameMap[CM_DEPOSIT_REFINE_ITEM] = "DEPOSIT_REFINE_ITEM"
	PacketNameMap[CM_RETRIEVE_REFINE_ITEM] = "RETRIEVE_REFINE_ITEM"
	PacketNameMap[CM_REFINE_CANCEL] = "REFINE_CANCEL"
	PacketNameMap[CM_REFINE_ITEM] = "REFINE_ITEM"
	PacketNameMap[CM_CHECK_REFINE] = "CHECK_REFINE"
	PacketNameMap[CM_REPLACE_WED_RING] = "REPLACE_WED_RING"
	PacketNameMap[CM_DEPOSIT_TRADE_ITEM] = "DEPOSIT_TRADE_ITEM"
	PacketNameMap[CM_RETRIEVE_TRADE_ITEM] = "RETRIEVE_TRADE_ITEM"
	PacketNameMap[CM_DROP_GOLD] = "DROP_GOLD"
	PacketNameMap[CM_PICK_UP] = "PICK_UP"
	PacketNameMap[CM_INSPECT] = "INSPECT"
	PacketNameMap[CM_CHANGE_A_MODE] = "CHANGE_A_MODE"
	PacketNameMap[CM_CHANGE_P_MODE] = "CHANGE_P_MODE"
	PacketNameMap[CM_CHANGE_TRADE] = "CHANGE_TRADE"
	PacketNameMap[CM_ATTACK] = "ATTACK"
	PacketNameMap[CM_RANGE_ATTACK] = "RANGE_ATTACK"
	PacketNameMap[CM_HARVEST] = "HARVEST"
	PacketNameMap[CM_CALL_NPC] = "CALL_NPC"
	PacketNameMap[CM_TALK_MONSTER_NPC] = "TALK_MONSTER_NPC"
	PacketNameMap[CM_BUY_ITEM] = "BUY_ITEM"
	PacketNameMap[CM_SELL_ITEM] = "SELL_ITEM"
	PacketNameMap[CM_CRAFT_ITEM] = "CRAFT_ITEM"
	PacketNameMap[CM_REPAIR_ITEM] = "REPAIR_ITEM"
	PacketNameMap[CM_BUY_ITEM_BACK] = "BUY_ITEM_BACK"
	PacketNameMap[CM_S_REPAIR_ITEM] = "S_REPAIR_ITEM"
	PacketNameMap[CM_MAGIC_KEY] = "MAGIC_KEY"
	PacketNameMap[CM_MAGIC] = "MAGIC"
	PacketNameMap[CM_SWITCH_GROUP] = "SWITCH_GROUP"
	PacketNameMap[CM_ADD_MEMBER] = "ADD_MEMBER"
	PacketNameMap[CM_DEL_MEMBER] = "DEL_MEMBER"
	PacketNameMap[CM_GROUP_INVITE] = "GROUP_INVITE"
	PacketNameMap[CM_TOWN_REVIVE] = "TOWN_REVIVE"
	PacketNameMap[CM_SPELL_TOGGLE] = "SPELL_TOGGLE"
	PacketNameMap[CM_CONSIGN_ITEM] = "CONSIGN_ITEM"
	PacketNameMap[CM_MARKET_SEARCH] = "MARKET_SEARCH"
	PacketNameMap[CM_MARKET_REFRESH] = "MARKET_REFRESH"
	PacketNameMap[CM_MARKET_PAGE] = "MARKET_PAGE"
	PacketNameMap[CM_MARKET_BUY] = "MARKET_BUY"
	PacketNameMap[CM_MARKET_GET_BACK] = "MARKET_GET_BACK"
	PacketNameMap[CM_REQUEST_USER_NAME] = "REQUEST_USER_NAME"
	PacketNameMap[CM_REQUEST_CHAT_ITEM] = "REQUEST_CHAT_ITEM"
	PacketNameMap[CM_EDIT_GUILD_MEMBER] = "EDIT_GUILD_MEMBER"
	PacketNameMap[CM_EDIT_GUILD_NOTICE] = "EDIT_GUILD_NOTICE"
	PacketNameMap[CM_GUILD_INVITE] = "GUILD_INVITE"
	PacketNameMap[CM_GUILD_NAME_RETURN] = "GUILD_NAME_RETURN"
	PacketNameMap[CM_REQUEST_GUILD_INFO] = "REQUEST_GUILD_INFO"
	PacketNameMap[CM_GUILD_STORAGE_GOLD_CHANGE] = "GUILD_STORAGE_GOLD_CHANGE"
	PacketNameMap[CM_GUILD_STORAGE_ITEM_CHANGE] = "GUILD_STORAGE_ITEM_CHANGE"
	PacketNameMap[CM_GUILD_WAR_RETURN] = "GUILD_WAR_RETURN"
	PacketNameMap[CM_MARRIAGE_REQUEST] = "MARRIAGE_REQUEST"
	PacketNameMap[CM_MARRIAGE_REPLY] = "MARRIAGE_REPLY"
	PacketNameMap[CM_CHANGE_MARRIAGE] = "CHANGE_MARRIAGE"
	PacketNameMap[CM_DIVORCE_REQUEST] = "DIVORCE_REQUEST"
	PacketNameMap[CM_DIVORCE_REPLY] = "DIVORCE_REPLY"
	PacketNameMap[CM_ADD_MENTOR] = "ADD_MENTOR"
	PacketNameMap[CM_MENTOR_REPLY] = "MENTOR_REPLY"
	PacketNameMap[CM_ALLOW_MENTOR] = "ALLOW_MENTOR"
	PacketNameMap[CM_CANCEL_MENTOR] = "CANCEL_MENTOR"
	PacketNameMap[CM_TRADE_REQUEST] = "TRADE_REQUEST"
	PacketNameMap[CM_TRADE_REPLY] = "TRADE_REPLY"
	PacketNameMap[CM_TRADE_GOLD] = "TRADE_GOLD"
	PacketNameMap[CM_TRADE_CONFIRM] = "TRADE_CONFIRM"
	PacketNameMap[CM_TRADE_CANCEL] = "TRADE_CANCEL"
	PacketNameMap[CM_EQUIP_SLOT_ITEM] = "EQUIP_SLOT_ITEM"
	PacketNameMap[CM_FISHING_CAST] = "FISHING_CAST"
	PacketNameMap[CM_FISHING_CHANGE_AUTOCAST] = "FISHING_CHANGE_AUTOCAST"
	PacketNameMap[CM_ACCEPT_QUEST] = "ACCEPT_QUEST"
	PacketNameMap[CM_FINISH_QUEST] = "FINISH_QUEST"
	PacketNameMap[CM_ABANDON_QUEST] = "ABANDON_QUEST"
	PacketNameMap[CM_SHARE_QUEST] = "SHARE_QUEST"
	PacketNameMap[CM_ACCEPT_REINCARNATION] = "ACCEPT_REINCARNATION"
	PacketNameMap[CM_CANCEL_REINCARNATION] = "CANCEL_REINCARNATION"
	PacketNameMap[CM_COMBINE_ITEM] = "COMBINE_ITEM"
	PacketNameMap[CM_SET_CONCENTRATION] = "SET_CONCENTRATION"
	PacketNameMap[CM_AWAKENING_NEED_MATERIALS] = "AWAKENING_NEED_MATERIALS"
	PacketNameMap[CM_AWAKENING_LOCKED_ITEM] = "AWAKENING_LOCKED_ITEM"
	PacketNameMap[CM_AWAKENING] = "AWAKENING"
	PacketNameMap[CM_DISASSEMBLE_ITEM] = "DISASSEMBLE_ITEM"
	PacketNameMap[CM_DOWNGRADE_AWAKENING] = "DOWNGRADE_AWAKENING"
	PacketNameMap[CM_RESET_ADDED_ITEM] = "RESET_ADDED_ITEM"
	PacketNameMap[CM_SEND_MAIL] = "SEND_MAIL"
	PacketNameMap[CM_READ_MAIL] = "READ_MAIL"
	PacketNameMap[CM_COLLECT_PARCEL] = "COLLECT_PARCEL"
	PacketNameMap[CM_DELETE_MAIL] = "DELETE_MAIL"
	PacketNameMap[CM_LOCK_MAIL] = "LOCK_MAIL"
	PacketNameMap[CM_MAIL_LOCKED_ITEM] = "MAIL_LOCKED_ITEM"
	PacketNameMap[CM_MAIL_COST] = "MAIL_COST"
	PacketNameMap[CM_UPDATE_INTELLIGENT_CREATURE] = "UPDATE_INTELLIGENT_CREATURE"
	PacketNameMap[CM_INTELLIGENT_CREATURE_PICKUP] = "INTELLIGENT_CREATURE_PICKUP"
	PacketNameMap[CM_ADD_FRIEND] = "ADD_FRIEND"
	PacketNameMap[CM_REMOVE_FRIEND] = "REMOVE_FRIEND"
	PacketNameMap[CM_REFRESH_FRIENDS] = "REFRESH_FRIENDS"
	PacketNameMap[CM_ADD_MEMO] = "ADD_MEMO"
	PacketNameMap[CM_GUILD_BUFF_UPDATE] = "GUILD_BUFF_UPDATE"
	PacketNameMap[CM_NPC_CONFIRM_INPUT] = "NPC_CONFIRM_INPUT"
	PacketNameMap[CM_GAMESHOP_BUY] = "GAMESHOP_BUY"
	PacketNameMap[CM_REPORT_ISSUE] = "REPORT_ISSUE"
	PacketNameMap[CM_GET_RANKING] = "GET_RANKING"
	PacketNameMap[CM_OPENDOOR] = "OPENDOOR"
	PacketNameMap[CM_GET_RENTED_ITEMS] = "GET_RENTED_ITEMS"
	PacketNameMap[CM_ITEM_RENTAL_REQUEST] = "ITEM_RENTAL_REQUEST"
	PacketNameMap[CM_ITEM_RENTAL_FEE] = "ITEM_RENTAL_FEE"
	PacketNameMap[CM_ITEM_RENTAL_PERIOD] = "ITEM_RENTAL_PERIOD"
	PacketNameMap[CM_DEPOSIT_RENTAL_ITEM] = "DEPOSIT_RENTAL_ITEM"
	PacketNameMap[CM_RETRIEVE_RENTAL_ITEM] = "RETRIEVE_RENTAL_ITEM"
	PacketNameMap[CM_CANCEL_ITEM_RENTAL] = "CANCEL_ITEM_RENTAL"
	PacketNameMap[CM_ITEM_RENTAL_LOCK_FEE] = "ITEM_RENTAL_LOCK_FEE"
	PacketNameMap[CM_ITEM_RENTAL_LOCK_ITEM] = "ITEM_RENTAL_LOCK_ITEM"
	PacketNameMap[CM_CONFIRM_ITEM_RENTAL] = "CONFIRM_ITEM_RENTAL"
}

func initServerPacketMap() {
	PacketNameMap[SM_CONNECTED] = "CONNECTED"
	PacketNameMap[SM_CLIENT_VERSION] = "CLIENT_VERSION"
	PacketNameMap[SM_DISCONNECT] = "DISCONNECT"
	PacketNameMap[SM_KEEP_ALIVE] = "KEEP_ALIVE"
	PacketNameMap[SM_NEW_ACCOUNT] = "NEW_ACCOUNT"
	PacketNameMap[SM_CHANGE_PASSWORD] = "CHANGE_PASSWORD"
	PacketNameMap[SM_CHANGE_PASSWORD_BANNED] = "CHANGE_PASSWORD_BANNED"
	PacketNameMap[SM_LOGIN] = "LOGIN"
	PacketNameMap[SM_LOGIN_BANNED] = "LOGIN_BANNED"
	PacketNameMap[SM_LOGIN_SUCCESS] = "LOGIN_SUCCESS"
	PacketNameMap[SM_NEW_CHARACTER] = "NEW_CHARACTER"
	PacketNameMap[SM_NEW_CHARACTER_SUCCESS] = "NEW_CHARACTER_SUCCESS"
	PacketNameMap[SM_DELETE_CHARACTER] = "DELETE_CHARACTER"
	PacketNameMap[SM_DELETE_CHARACTER_SUCCESS] = "DELETE_CHARACTER_SUCCESS"
	PacketNameMap[SM_START_GAME] = "START_GAME"
	PacketNameMap[SM_START_GAME_BANNED] = "START_GAME_BANNED"
	PacketNameMap[SM_START_GAME_DELAY] = "START_GAME_DELAY"
	PacketNameMap[SM_MAP_INFORMATION] = "MAP_INFORMATION"
	PacketNameMap[SM_USER_INFORMATION] = "USER_INFORMATION"
	PacketNameMap[SM_USER_LOCATION] = "USER_LOCATION"
	PacketNameMap[SM_OBJECT_PLAYER] = "OBJECT_PLAYER"
	PacketNameMap[SM_OBJECT_REMOVE] = "OBJECT_REMOVE"
	PacketNameMap[SM_OBJECT_TURN] = "OBJECT_TURN"
	PacketNameMap[SM_OBJECT_WALK] = "OBJECT_WALK"
	PacketNameMap[SM_OBJECT_RUN] = "OBJECT_RUN"
	PacketNameMap[SM_CHAT] = "CHAT"
	PacketNameMap[SM_OBJECT_CHAT] = "OBJECT_CHAT"
	PacketNameMap[SM_NEW_ITEM_INFO] = "NEW_ITEM_INFO"
	PacketNameMap[SM_MOVE_ITEM] = "MOVE_ITEM"
	PacketNameMap[SM_EQUIP_ITEM] = "EQUIP_ITEM"
	PacketNameMap[SM_MERGE_ITEM] = "MERGE_ITEM"
	PacketNameMap[SM_REMOVE_ITEM] = "REMOVE_ITEM"
	PacketNameMap[SM_REMOVE_SLOT_ITEM] = "REMOVE_SLOT_ITEM"
	PacketNameMap[SM_TAKE_BACK_ITEM] = "TAKE_BACK_ITEM"
	PacketNameMap[SM_STORE_ITEM] = "STORE_ITEM"
	PacketNameMap[SM_SPLIT_ITEM] = "SPLIT_ITEM"
	PacketNameMap[SM_SPLIT_ITEM1] = "SPLIT_ITEM1"
	PacketNameMap[SM_DEPOSIT_REFINE_ITEM] = "DEPOSIT_REFINE_ITEM"
	PacketNameMap[SM_RETRIEVE_REFINE_ITEM] = "RETRIEVE_REFINE_ITEM"
	PacketNameMap[SM_REFINE_CANCEL] = "REFINE_CANCEL"
	PacketNameMap[SM_REFINE_ITEM] = "REFINE_ITEM"
	PacketNameMap[SM_DEPOSIT_TRADE_ITEM] = "DEPOSIT_TRADE_ITEM"
	PacketNameMap[SM_RETRIEVE_TRADE_ITEM] = "RETRIEVE_TRADE_ITEM"
	PacketNameMap[SM_USE_ITEM] = "USE_ITEM"
	PacketNameMap[SM_DROP_ITEM] = "DROP_ITEM"
	PacketNameMap[SM_PLAYER_UPDATE] = "PLAYER_UPDATE"
	PacketNameMap[SM_PLAYER_INSPECT] = "PLAYER_INSPECT"
	PacketNameMap[SM_LOG_OUT_SUCCESS] = "LOG_OUT_SUCCESS"
	PacketNameMap[SM_LOG_OUT_FAILED] = "LOG_OUT_FAILED"
	PacketNameMap[SM_TIME_OF_DAY] = "TIME_OF_DAY"
	PacketNameMap[SM_CHANGE_A_MODE] = "CHANGE_A_MODE"
	PacketNameMap[SM_CHANGE_P_MODE] = "CHANGE_P_MODE"
	PacketNameMap[SM_OBJECT_ITEM] = "OBJECT_ITEM"
	PacketNameMap[SM_OBJECT_GOLD] = "OBJECT_GOLD"
	PacketNameMap[SM_GAINED_ITEM] = "GAINED_ITEM"
	PacketNameMap[SM_GAINED_GOLD] = "GAINED_GOLD"
	PacketNameMap[SM_LOSE_GOLD] = "LOSE_GOLD"
	PacketNameMap[SM_GAINED_CREDIT] = "GAINED_CREDIT"
	PacketNameMap[SM_LOSE_CREDIT] = "LOSE_CREDIT"
	PacketNameMap[SM_OBJECT_MONSTER] = "OBJECT_MONSTER"
	PacketNameMap[SM_OBJECT_ATTACK] = "OBJECT_ATTACK"
	PacketNameMap[SM_STRUCK] = "STRUCK"
	PacketNameMap[SM_OBJECT_STRUCK] = "OBJECT_STRUCK"
	PacketNameMap[SM_DAMAGE_INDICATOR] = "DAMAGE_INDICATOR"
	PacketNameMap[SM_DURA_CHANGED] = "DURA_CHANGED"
	PacketNameMap[SM_HEALTH_CHANGED] = "HEALTH_CHANGED"
	PacketNameMap[SM_DELETE_ITEM] = "DELETE_ITEM"
	PacketNameMap[SM_DEATH] = "DEATH"
	PacketNameMap[SM_OBJECT_DIED] = "OBJECT_DIED"
	PacketNameMap[SM_COLOUR_CHANGED] = "COLOUR_CHANGED"
	PacketNameMap[SM_OBJECT_COLOUR_CHANGED] = "OBJECT_COLOUR_CHANGED"
	PacketNameMap[SM_OBJECT_GUILD_NAME_CHANGED] = "OBJECT_GUILD_NAME_CHANGED"
	PacketNameMap[SM_GAIN_EXPERIENCE] = "GAIN_EXPERIENCE"
	PacketNameMap[SM_LEVEL_CHANGED] = "LEVEL_CHANGED"
	PacketNameMap[SM_OBJECT_LEVELED] = "OBJECT_LEVELED"
	PacketNameMap[SM_OBJECT_HARVEST] = "OBJECT_HARVEST"
	PacketNameMap[SM_OBJECT_HARVESTED] = "OBJECT_HARVESTED"
	PacketNameMap[SM_OBJECT_NPC] = "OBJECT_NPC"
	PacketNameMap[SM_NPC_RESPONSE] = "NPC_RESPONSE"
	PacketNameMap[SM_OBJECT_HIDE] = "OBJECT_HIDE"
	PacketNameMap[SM_OBJECT_SHOW] = "OBJECT_SHOW"
	PacketNameMap[SM_POISONED] = "POISONED"
	PacketNameMap[SM_OBJECT_POISONED] = "OBJECT_POISONED"
	PacketNameMap[SM_MAP_CHANGED] = "MAP_CHANGED"
	PacketNameMap[SM_OBJECT_TELEPORT_OUT] = "OBJECT_TELEPORT_OUT"
	PacketNameMap[SM_OBJECT_TELEPORT_IN] = "OBJECT_TELEPORT_IN"
	PacketNameMap[SM_TELEPORT_IN] = "TELEPORT_IN"
	PacketNameMap[SM_NPC_GOODS] = "NPC_GOODS"
	PacketNameMap[SM_NPC_SELL] = "NPC_SELL"
	PacketNameMap[SM_NPC_REPAIR] = "NPC_REPAIR"
	PacketNameMap[SM_NPC_S_REPAIR] = "NPC_S_REPAIR"
	PacketNameMap[SM_NPC_REFINE] = "NPC_REFINE"
	PacketNameMap[SM_NPC_CHECK_REFINE] = "NPC_CHECK_REFINE"
	PacketNameMap[SM_NPC_COLLECT_REFINE] = "NPC_COLLECT_REFINE"
	PacketNameMap[SM_NPC_REPLACE_WED_RING] = "NPC_REPLACE_WED_RING"
	PacketNameMap[SM_NPC_STORAGE] = "NPC_STORAGE"
	PacketNameMap[SM_SELL_ITEM] = "SELL_ITEM"
	PacketNameMap[SM_CRAFT_ITEM] = "CRAFT_ITEM"
	PacketNameMap[SM_REPAIR_ITEM] = "REPAIR_ITEM"
	PacketNameMap[SM_ITEM_REPAIRED] = "ITEM_REPAIRED"
	PacketNameMap[SM_NEW_MAGIC] = "NEW_MAGIC"
	PacketNameMap[SM_REMOVE_MAGIC] = "REMOVE_MAGIC"
	PacketNameMap[SM_MAGIC_LEVELED] = "MAGIC_LEVELED"
	PacketNameMap[SM_MAGIC] = "MAGIC"
	PacketNameMap[SM_MAGIC_DELAY] = "MAGIC_DELAY"
	PacketNameMap[SM_MAGIC_CAST] = "MAGIC_CAST"
	PacketNameMap[SM_OBJECT_MAGIC] = "OBJECT_MAGIC"
	PacketNameMap[SM_OBJECT_EFFECT] = "OBJECT_EFFECT"
	PacketNameMap[SM_RANGE_ATTACK] = "RANGE_ATTACK"
	PacketNameMap[SM_PUSHED] = "PUSHED"
	PacketNameMap[SM_OBJECT_PUSHED] = "OBJECT_PUSHED"
	PacketNameMap[SM_OBJECT_NAME] = "OBJECT_NAME"
	PacketNameMap[SM_USER_STORAGE] = "USER_STORAGE"
	PacketNameMap[SM_SWITCH_GROUP] = "SWITCH_GROUP"
	PacketNameMap[SM_DELETE_GROUP] = "DELETE_GROUP"
	PacketNameMap[SM_DELETE_MEMBER] = "DELETE_MEMBER"
	PacketNameMap[SM_GROUP_INVITE] = "GROUP_INVITE"
	PacketNameMap[SM_ADD_MEMBER] = "ADD_MEMBER"
	PacketNameMap[SM_REVIVED] = "REVIVED"
	PacketNameMap[SM_OBJECT_REVIVED] = "OBJECT_REVIVED"
	PacketNameMap[SM_SPELL_TOGGLE] = "SPELL_TOGGLE"
	PacketNameMap[SM_OBJECT_HEALTH] = "OBJECT_HEALTH"
	PacketNameMap[SM_MAP_EFFECT] = "MAP_EFFECT"
	PacketNameMap[SM_OBJECT_RANGE_ATTACK] = "OBJECT_RANGE_ATTACK"
	PacketNameMap[SM_ADD_BUFF] = "ADD_BUFF"
	PacketNameMap[SM_REMOVE_BUFF] = "REMOVE_BUFF"
	PacketNameMap[SM_OBJECT_HIDDEN] = "OBJECT_HIDDEN"
	PacketNameMap[SM_REFRESH_ITEM] = "REFRESH_ITEM"
	PacketNameMap[SM_OBJECT_SPELL] = "OBJECT_SPELL"
	PacketNameMap[SM_USER_DASH] = "USER_DASH"
	PacketNameMap[SM_OBJECT_DASH] = "OBJECT_DASH"
	PacketNameMap[SM_USER_DASH_FAIL] = "USER_DASH_FAIL"
	PacketNameMap[SM_OBJECT_DASH_FAIL] = "OBJECT_DASH_FAIL"
	PacketNameMap[SM_NPC_CONSIGN] = "NPC_CONSIGN"
	PacketNameMap[SM_NPC_MARKET] = "NPC_MARKET"
	PacketNameMap[SM_NPC_MARKET_PAGE] = "NPC_MARKET_PAGE"
	PacketNameMap[SM_CONSIGN_ITEM] = "CONSIGN_ITEM"
	PacketNameMap[SM_MARKET_FAIL] = "MARKET_FAIL"
	PacketNameMap[SM_MARKET_SUCCESS] = "MARKET_SUCCESS"
	PacketNameMap[SM_OBJECT_SIT_DOWN] = "OBJECT_SIT_DOWN"
	PacketNameMap[SM_IN_TRAP_ROCK] = "IN_TRAP_ROCK"
	PacketNameMap[SM_BASE_STATS_INFO] = "BASE_STATS_INFO"
	PacketNameMap[SM_USER_NAME] = "USER_NAME"
	PacketNameMap[SM_CHAT_ITEM_STATS] = "CHAT_ITEM_STATS"
	PacketNameMap[SM_GUILD_NOTICE_CHANGE] = "GUILD_NOTICE_CHANGE"
	PacketNameMap[SM_GUILD_MEMBER_CHANGE] = "GUILD_MEMBER_CHANGE"
	PacketNameMap[SM_GUILD_STATUS] = "GUILD_STATUS"
	PacketNameMap[SM_GUILD_INVITE] = "GUILD_INVITE"
	PacketNameMap[SM_GUILD_EXP_GAIN] = "GUILD_EXP_GAIN"
	PacketNameMap[SM_GUILD_NAME_REQUEST] = "GUILD_NAME_REQUEST"
	PacketNameMap[SM_GUILD_STORAGE_GOLD_CHANGE] = "GUILD_STORAGE_GOLD_CHANGE"
	PacketNameMap[SM_GUILD_STORAGE_ITEM_CHANGE] = "GUILD_STORAGE_ITEM_CHANGE"
	PacketNameMap[SM_GUILD_STORAGE_LIST] = "GUILD_STORAGE_LIST"
	PacketNameMap[SM_GUILD_REQUEST_WAR] = "GUILD_REQUEST_WAR"
	PacketNameMap[SM_DEFAULT_NPC] = "DEFAULT_NPC"
	PacketNameMap[SM_NPC_UPDATE] = "NPC_UPDATE"
	PacketNameMap[SM_NPC_IMAGE_UPDATE] = "NPC_IMAGE_UPDATE"
	PacketNameMap[SM_MARRIAGE_REQUEST] = "MARRIAGE_REQUEST"
	PacketNameMap[SM_DIVORCE_REQUEST] = "DIVORCE_REQUEST"
	PacketNameMap[SM_MENTOR_REQUEST] = "MENTOR_REQUEST"
	PacketNameMap[SM_TRADE_REQUEST] = "TRADE_REQUEST"
	PacketNameMap[SM_TRADE_ACCEPT] = "TRADE_ACCEPT"
	PacketNameMap[SM_TRADE_GOLD] = "TRADE_GOLD"
	PacketNameMap[SM_TRADE_ITEM] = "TRADE_ITEM"
	PacketNameMap[SM_TRADE_CONFIRM] = "TRADE_CONFIRM"
	PacketNameMap[SM_TRADE_CANCEL] = "TRADE_CANCEL"
	PacketNameMap[SM_MOUNT_UPDATE] = "MOUNT_UPDATE"
	PacketNameMap[SM_EQUIP_SLOT_ITEM] = "EQUIP_SLOT_ITEM"
	PacketNameMap[SM_FISHING_UPDATE] = "FISHING_UPDATE"
	PacketNameMap[SM_CHANGE_QUEST] = "CHANGE_QUEST"
	PacketNameMap[SM_COMPLETE_QUEST] = "COMPLETE_QUEST"
	PacketNameMap[SM_SHARE_QUEST] = "SHARE_QUEST"
	PacketNameMap[SM_NEW_QUEST_INFO] = "NEW_QUEST_INFO"
	PacketNameMap[SM_GAINED_QUEST_ITEM] = "GAINED_QUEST_ITEM"
	PacketNameMap[SM_DELETE_QUEST_ITEM] = "DELETE_QUEST_ITEM"
	PacketNameMap[SM_CANCEL_REINCARNATION] = "CANCEL_REINCARNATION"
	PacketNameMap[SM_REQUEST_REINCARNATION] = "REQUEST_REINCARNATION"
	PacketNameMap[SM_USER_BACK_STEP] = "USER_BACK_STEP"
	PacketNameMap[SM_OBJECT_BACK_STEP] = "OBJECT_BACK_STEP"
	PacketNameMap[SM_USER_DASH_ATTACK] = "USER_DASH_ATTACK"
	PacketNameMap[SM_OBJECT_DASH_ATTACK] = "OBJECT_DASH_ATTACK"
	PacketNameMap[SM_USER_ATTACK_MOVE] = "USER_ATTACK_MOVE"
	PacketNameMap[SM_COMBINE_ITEM] = "COMBINE_ITEM"
	PacketNameMap[SM_ITEM_UPGRADED] = "ITEM_UPGRADED"
	PacketNameMap[SM_SET_CONCENTRATION] = "SET_CONCENTRATION"
	PacketNameMap[SM_SET_OBJECT_CONCENTRATION] = "SET_OBJECT_CONCENTRATION"
	PacketNameMap[SM_SET_ELEMENTAL] = "SET_ELEMENTAL"
	PacketNameMap[SM_SET_OBJECT_ELEMENTAL] = "SET_OBJECT_ELEMENTAL"
	PacketNameMap[SM_REMOVE_DELAYED_EXPLOSION] = "REMOVE_DELAYED_EXPLOSION"
	PacketNameMap[SM_OBJECT_DECO] = "OBJECT_DECO"
	PacketNameMap[SM_OBJECT_SNEAKING] = "OBJECT_SNEAKING"
	PacketNameMap[SM_OBJECT_LEVEL_EFFECTS] = "OBJECT_LEVEL_EFFECTS"
	PacketNameMap[SM_SET_BINDING_SHOT] = "SET_BINDING_SHOT"
	PacketNameMap[SM_SEND_OUTPUT_MESSAGE] = "SEND_OUTPUT_MESSAGE"
	PacketNameMap[SM_NPC_AWAKENING] = "NPC_AWAKENING"
	PacketNameMap[SM_NPC_DISASSEMBLE] = "NPC_DISASSEMBLE"
	PacketNameMap[SM_NPC_DOWNGRADE] = "NPC_DOWNGRADE"
	PacketNameMap[SM_NPC_RESET] = "NPC_RESET"
	PacketNameMap[SM_AWAKENING_NEED_MATERIALS] = "AWAKENING_NEED_MATERIALS"
	PacketNameMap[SM_AWAKENING_LOCKED_ITEM] = "AWAKENING_LOCKED_ITEM"
	PacketNameMap[SM_AWAKENING] = "AWAKENING"
	PacketNameMap[SM_RECEIVE_MAIL] = "RECEIVE_MAIL"
	PacketNameMap[SM_MAIL_LOCKED_ITEM] = "MAIL_LOCKED_ITEM"
	PacketNameMap[SM_MAIL_SEND_REQUEST] = "MAIL_SEND_REQUEST"
	PacketNameMap[SM_MAIL_SENT] = "MAIL_SENT"
	PacketNameMap[SM_PARCEL_COLLECTED] = "PARCEL_COLLECTED"
	PacketNameMap[SM_MAIL_COST] = "MAIL_COST"
	PacketNameMap[SM_RESIZE_INVENTORY] = "RESIZE_INVENTORY"
	PacketNameMap[SM_RESIZE_STORAGE] = "RESIZE_STORAGE"
	PacketNameMap[SM_NEW_INTELLIGENT_CREATURE] = "NEW_INTELLIGENT_CREATURE"
	PacketNameMap[SM_UPDATE_INTELLIGENT_CREATURElIST] = "UPDATE_INTELLIGENT_CREATURElIST"
	PacketNameMap[SM_INTELLIGENT_CREATURE_ENABLE_RENAME] = "INTELLIGENT_CREATURE_ENABLE_RENAME"
	PacketNameMap[SM_INTELLIGENT_CREATURE_PICKUP] = "INTELLIGENT_CREATURE_PICKUP"
	PacketNameMap[SM_NPC_PEARL_GOODS] = "NPC_PEARL_GOODS"
	PacketNameMap[SM_TRANSFORM_UPDATE] = "TRANSFORM_UPDATE"
	PacketNameMap[SM_FRIEND_UPDATE] = "FRIEND_UPDATE"
	PacketNameMap[SM_LOVER_UPDATE] = "LOVER_UPDATE"
	PacketNameMap[SM_MENTOR_UPDATE] = "MENTOR_UPDATE"
	PacketNameMap[SM_GUILD_BUFF_LIST] = "GUILD_BUFF_LIST"
	PacketNameMap[SM_NPC_REQUEST_INPUT] = "NPC_REQUEST_INPUT"
	PacketNameMap[SM_GAME_SHOP_INFO] = "GAME_SHOP_INFO"
	PacketNameMap[SM_GAME_SHOP_STOCK] = "GAME_SHOP_STOCK"
	PacketNameMap[SM_RANKINGS] = "RANKINGS"
	PacketNameMap[SM_OPENDOOR] = "OPENDOOR"
	PacketNameMap[SM_GET_RENTED_ITEMS] = "GET_RENTED_ITEMS"
	PacketNameMap[SM_ITEM_RENTAL_REQUEST] = "ITEM_RENTAL_REQUEST"
	PacketNameMap[SM_ITEM_RENTAL_FEE] = "ITEM_RENTAL_FEE"
	PacketNameMap[SM_ITEM_RENTAL_PERIOD] = "ITEM_RENTAL_PERIOD"
	PacketNameMap[SM_DEPOSIT_RENTAL_ITEM] = "DEPOSIT_RENTAL_ITEM"
	PacketNameMap[SM_RETRIEVE_RENTAL_ITEM] = "RETRIEVE_RENTAL_ITEM"
	PacketNameMap[SM_UPDATE_RENTAL_ITEM] = "UPDATE_RENTAL_ITEM"
	PacketNameMap[SM_CANCEL_ITEM_RENTAL] = "CANCEL_ITEM_RENTAL"
	PacketNameMap[SM_ITEM_RENTAL_LOCK] = "ITEM_RENTAL_LOCK"
	PacketNameMap[SM_ITEM_RENTAL_PARTNER_LOCK] = "ITEM_RENTAL_PARTNER_LOCK"
	PacketNameMap[SM_CAN_CONFIRM_ITEM_RENTAL] = "CAN_CONFIRM_ITEM_RENTAL"
	PacketNameMap[SM_CONFIRM_ITEM_RENTAL] = "CONFIRM_ITEM_RENTAL"
	PacketNameMap[SM_NEW_RECIPE_INFO] = "NEW_RECIPE_INFO"
	PacketNameMap[SM_OPEN_BROWSER] = "OPEN_BROWSER"
}

func init() {
	initClientPacketMap()
	initServerPacketMap()
}