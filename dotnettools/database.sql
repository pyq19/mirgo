CREATE TABLE basic (
    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    game_version INT,
    custom_version INT,
    map_index INT,
    item_index INT,
    monster_index INT,
    npc_index INT,
    quest_index INT,
    gameshop_index INT,
    conquest_index INT,
    respawn_index INT
);

-- MapInfo
CREATE TABLE map_info (
    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    map_index INT,
    file_name VARCHAR(100),
    title VARCHAR(100),
    mini_map INT,
    big_map INT,
    music INT,
    light INT,
    map_dark_light INT,
    mine_index INT,
    no_teleport INT,
    no_reconnect INT,
    no_random INT,
    no_escape INT,
    no_recall INT,
    no_drug INT,
    no_position INT,
    no_fight INT,
    no_throw_item INT,
    no_drop_player INT,
    no_drop_monster INT,
    no_names INT,
    no_mount INT,
    need_bridle INT,
    fight INT,
    fire INT,
    lightning INT,
    no_town_teleport INT,
    no_reincarnation INT,
    no_reconnect_map VARCHAR(100),
    fire_damage INT,
    lightning_damage INT
);

CREATE TABLE safe_zone_info (
    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    map_index INT,
    location_x INT,
    location_y INT,
    size INT,
    start_point INT
);

CREATE TABLE movement_info (
    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    map_index INT,
    source_x INT,
    source_y INT,
    destination_x INT,
    destination_y INT,
    need_hole INT,
    need_move INT,
    conquest_index INT
);

CREATE TABLE respawn_info(
    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    map_index INT,
    monster_index INT,
    location_x INT,
    location_y INT,
    count INT,
    spread INT,
    delay INT,
    random_delay INT,
    direction INT,
    route_path VARCHAR(1000),
    respawn_index INT,
    save_respawn_time INT,
    respawn_ticks INT
);

CREATE TABLE npc_info();

CREATE TABLE mine_zone(
    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    map_index INT,
    mine INT,
    location_x INT,
    location_y INT,
    size INT
);

CREATE TABLE active_coords();

-- ItemInfo
-- MonsterInfo
-- NPCInfo
-- QuestInfo
-- MagicInfo
-- GameShop
-- Conquest
-- RespawnTick