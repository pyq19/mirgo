--
-- 由SQLiteStudio v3.2.1 产生的文件 周二 3月 3 16:06:43 2020
--
-- 文本编码：UTF-8
--
PRAGMA foreign_keys = off;
BEGIN TRANSACTION;

-- 表：basic
CREATE TABLE `basic` (
  `id` integer  NOT NULL PRIMARY KEY AUTOINCREMENT
,  `game_version` integer DEFAULT NULL
,  `custom_version` integer DEFAULT NULL
,  `map_index` integer DEFAULT NULL
,  `item_index` integer DEFAULT NULL
,  `monster_index` integer DEFAULT NULL
,  `npc_index` integer DEFAULT NULL
,  `quest_index` integer DEFAULT NULL
,  `gameshop_index` integer DEFAULT NULL
,  `conquest_index` integer DEFAULT NULL
,  `respawn_index` integer DEFAULT NULL
);


-- 表：game_shop_item
CREATE TABLE "game_shop_item"
(
	id integer not null
		primary key autoincrement,
	item_id integer default NULL,
	gold_price integer default NULL,
	credit_price integer default NULL,
	count integer default NULL,
	class varchar(200) default NULL,
	category varchar(200) default NULL,
	stock integer default NULL,
	i_stock integer default NULL,
	deal integer default NULL,
	top_item integer default NULL,
	create_date date default NULL
);

-- 表：item
CREATE TABLE "item"
(
	id integer not null
		primary key autoincrement,
	name varchar(200) default NULL,
	type integer default NULL,
	grade integer default NULL,
	required_type integer default NULL,
	required_class integer default NULL,
	required_gender integer default NULL,
	item_set integer default NULL,
	shape integer default NULL,
	weight integer default NULL,
	light integer default NULL,
	required_amount integer default NULL,
	image integer default NULL,
	durability integer default NULL,
	stack_size integer default NULL,
	price integer default NULL,
	min_ac integer default NULL,
	max_ac integer default NULL,
	min_mac integer default NULL,
	max_mac integer default NULL,
	min_dc integer default NULL,
	max_dc integer default NULL,
	min_mc integer default NULL,
	max_mc integer default NULL,
	min_sc integer default NULL,
	max_sc integer default NULL,
	hp integer default NULL,
	mp integer default NULL,
	accuracy integer default NULL,
	agility integer default NULL,
	luck integer default NULL,
	attack_speed integer default NULL,
	start_item integer default NULL,
	bag_weight integer default NULL,
	hand_weight integer default NULL,
	wear_weight integer default NULL,
	effect integer default NULL,
	strong integer default NULL,
	magic_resist integer default NULL,
	poison_resist integer default NULL,
	health_recovery integer default NULL,
	spell_recovery integer default NULL,
	poison_recovery integer default NULL,
	hp_rate integer default NULL,
	mp_rate integer default NULL,
	critical_rate integer default NULL,
	critical_damage integer default NULL,
	bools integer default NULL,
	max_ac_rate integer default NULL,
	max_mac_rate integer default NULL,
	holy integer default NULL,
	freezing integer default NULL,
	poison_attack integer default NULL,
	bind integer default NULL,
	reflect integer default NULL,
	hp_drain_rate integer default NULL,
	unique_item integer default NULL,
	random_stats_id integer default NULL,
	can_fast_run integer default NULL,
	can_awakening integer default NULL,
	tool_tip varchar(2000) default NULL
);

-- 表：magic
CREATE TABLE "magic" (
  `id` integer  NOT NULL PRIMARY KEY AUTOINCREMENT
,  `name` varchar(200) DEFAULT NULL
,  `spell` integer DEFAULT NULL
,  `base_cost` integer DEFAULT NULL
,  `level_cost` integer DEFAULT NULL
,  `icon` integer DEFAULT NULL
,  `level_1` integer DEFAULT NULL
,  `level_2` integer DEFAULT NULL
,  `level_3` integer DEFAULT NULL
,  `need_1` integer DEFAULT NULL
,  `need_2` integer DEFAULT NULL
,  `need_3` integer DEFAULT NULL
,  `delay_base` integer DEFAULT NULL
,  `delay_reduction` integer DEFAULT NULL
,  `power_base` integer DEFAULT NULL
,  `power_bonus` integer DEFAULT NULL
,  `m_power_base` integer DEFAULT NULL
,  `m_power_bonus` integer DEFAULT NULL
,  `magic_range` integer DEFAULT NULL
,  `multiplier_base` float(5,3) DEFAULT NULL
,  `multiplier_bonus` float(5,3) DEFAULT NULL
);

-- 表：map
CREATE TABLE "map"
(
	id integer not null
		primary key autoincrement,
	file_name varchar(100) default NULL,
	title varchar(100) default NULL,
	mini_map integer default NULL,
	big_map integer default NULL,
	music integer default NULL,
	light integer default NULL,
	map_dark_light integer default NULL,
	mine_index integer default NULL,
	no_teleport integer default NULL,
	no_reconnect integer default NULL,
	no_random integer default NULL,
	no_escape integer default NULL,
	no_recall integer default NULL,
	no_drug integer default NULL,
	no_position integer default NULL,
	no_fight integer default NULL,
	no_throw_item integer default NULL,
	no_drop_player integer default NULL,
	no_drop_monster integer default NULL,
	no_names integer default NULL,
	no_mount integer default NULL,
	need_bridle integer default NULL,
	fight integer default NULL,
	fire integer default NULL,
	lightning integer default NULL,
	no_town_teleport integer default NULL,
	no_reincarnation integer default NULL,
	no_reconnect_map varchar(100) default NULL,
	fire_damage integer default NULL,
	lightning_damage integer default NULL
);

-- 表：mine_zone
CREATE TABLE "mine_zone"
(
	id integer not null
		primary key autoincrement,
	map_id integer default NULL,
	mine integer default NULL,
	location_x integer default NULL,
	location_y integer default NULL,
	size integer default NULL
);

-- 表：monster
CREATE TABLE "monster"
(
	id integer not null
		primary key autoincrement,
	name varchar(200) default NULL,
	image integer default NULL,
	ai integer default NULL,
	effect integer default NULL,
	level integer default NULL,
	view_range integer default NULL,
	cool_eye integer default NULL,
	hp integer default NULL,
	min_ac integer default NULL,
	max_ac integer default NULL,
	min_mac integer default NULL,
	max_mac integer default NULL,
	min_dc integer default NULL,
	max_dc integer default NULL,
	min_mc integer default NULL,
	max_mc integer default NULL,
	min_sc integer default NULL,
	max_sc integer default NULL,
	accuracy integer default NULL,
	agility integer default NULL,
	light integer default NULL,
	attack_speed integer default NULL,
	move_speed integer default NULL,
	experience integer default NULL,
	can_push integer default NULL,
	can_tame integer default NULL,
	auto_rev integer default NULL,
	undead integer default NULL
);

-- 表：movement
CREATE TABLE "movement"
(
	id integer not null
		primary key autoincrement,
	source_map integer default NULL,
	source_x integer default NULL,
	source_y integer default NULL,
	destination_map integer default NULL,
	destination_x integer default NULL,
	destination_y integer default NULL,
	need_hole integer default NULL,
	need_move integer default NULL,
	conquest_index integer default NULL
);

-- 表：npc
CREATE TABLE "npc"
(
	id integer not null
		primary key autoincrement,
	map_id integer default NULL,
	file_name varchar(200) default NULL,
	name varchar(200) default NULL,
	location_x integer default NULL,
	location_y integer default NULL,
	rate integer default NULL,
	image integer default NULL,
	time_visible integer default NULL,
	hour_start integer default NULL,
	minute_start integer default NULL,
	hour_end integer default NULL,
	minute_end integer default NULL,
	min_lev integer default NULL,
	max_lev integer default NULL,
	day_of_week varchar(200) default NULL,
	class_required varchar(200) default NULL,
	flag_needed integer default NULL,
	conquest integer default NULL
);

-- 表：quest
CREATE TABLE "quest"
(
	id integer not null
		primary key autoincrement,
	name varchar(200) default NULL,
	quest_group varchar(200) default NULL,
	file_name varchar(200) default NULL,
	required_min_level integer default NULL,
	required_max_level integer default NULL,
	required_quest integer default NULL,
	required_class integer default NULL,
	quest_type integer default NULL,
	goto_message varchar(2000) default NULL,
	kill_message varchar(2000) default NULL,
	item_message varchar(2000) default NULL,
	flag_message varchar(2000) default NULL
);

-- 表：respawn
CREATE TABLE "respawn"
(
	id integer not null
		primary key autoincrement,
	map_id integer default NULL,
	monster_id integer default NULL,
	location_x integer default NULL,
	location_y integer default NULL,
	count integer default NULL,
	spread integer default NULL,
	delay integer default NULL,
	random_delay integer default NULL,
	direction integer default NULL,
	route_path varchar(1000) default NULL,
	respawn_index integer default NULL,
	save_respawn_time integer default NULL,
	respawn_ticks integer default NULL
);

-- 表：safe_zone
CREATE TABLE "safe_zone"
(
	id integer not null
		primary key autoincrement,
	map_id integer default NULL,
	location_x integer default NULL,
	location_y integer default NULL,
	size integer default NULL,
	start_point integer default NULL
);

COMMIT TRANSACTION;
PRAGMA foreign_keys = on;
