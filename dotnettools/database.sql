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

CREATE TABLE mine_zone(
    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    map_index INT,
    mine INT,
    location_x INT,
    location_y INT,
    size INT
);

-- ItemInfo
CREATE TABLE item_info (
    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    item_index INT,
    name VARCHAR(200),
    type INT,
    grade INT,
    required_type INT,
    required_class INT,
    required_gender INT,
    item_set INT,
    shape INT,
    weight INT,
    light INT,
    required_amount INT,
    image INT,
    durability INT,
    stack_size INT,
    price INT,
    min_ac INT,
    max_ac INT,
    min_mac INT,
    max_mac INT,
    min_dc INT,
    max_dc INT,
    min_mc INT,
    max_mc INT,
    min_sc INT,
    max_sc INT,
    hp INT,
    mp INT,
    accuracy INT,
    agility INT,
    luck INT,
    attack_speed INT,
    start_item INT,
    bag_weight INT,
    hand_weight INT,
    wear_weight INT,
    effect INT,
    strong INT,
    magic_resist INT,
    poison_resist INT,
    health_recovery INT,
    spell_recovery INT,
    poison_recovery INT,
    hp_rate INT,
    mp_rate INT,
    critical_rate INT,
    critical_damage INT,
    bools INT,
    max_ac_rate INT,
    max_mac_rate INT,
    holy INT,
    freezing INT,
    poison_attack INT,
    bind INT,
    reflect INT,
    hp_drain_rate INT,
    unique_item INT,
    random_stats_id INT,
    can_fast_run INT,
    can_awakening INT,
    tool_tip VARCHAR(2000)
);

-- MonsterInfo
CREATE TABLE monster_info (
    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    monster_index INT,
    name VARCHAR(200),
    image INT,
    ai INT,
    effect INT,
    level INT,
    view_range INT,
    cool_eye INT,
    hp INT,
    min_ac INT,
    max_ac INT,
    min_mac INT,
    max_mac INT,
    min_dc INT,
    max_dc INT,
    min_mc INT,
    max_mc INT,
    min_sc INT,
    max_sc INT,
    accuracy INT,
    agility INT,
    light INT,
    attack_speed INT,
    move_speed INT,
    experience INT,
    can_push INT,
    can_tame INT,
    auto_rev INT,
    undead INT
);

-- NPCInfo
CREATE TABLE npc_info (
    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    map_index INT,
    npc_index INT,
    file_name VARCHAR(200),
    name VARCHAR(200),
    location_x INT,
    location_y INT,
    rate INT,
    image INT,
    time_visible INT,
    hour_start INT,
    minute_start INT,
    hour_end INT,
    minute_end INT,
    min_lev INT,
    max_lev INT,
    day_of_week VARCHAR(200),
    class_required VARCHAR(200),
    flag_needed INT,
    conquest INT
);

-- QuestInfo
CREATE TABLE quest_info (
    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    quest_index INT,
    name VARCHAR(200),
    quest_group VARCHAR(200),
    file_name VARCHAR(200),
    required_min_level INT,
    required_max_level INT,
    required_quest INT,
    required_class INT,
    quest_type INT,
    goto_message VARCHAR(2000),
    kill_message VARCHAR(2000),
    item_message VARCHAR(2000),
    flag_message VARCHAR(2000)
);

-- MagicInfo
CREATE TABLE magic_info (
    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(200),
    spell INT,
    base_cost INT,
    level_cost INT,
    icon INT,
    level_1 INT,
    level_2 INT,
    level_3 INT,
    need_1 INT,
    need_2 INT,
    need_3 INT,
    delay_base INT,
    delay_reduction INT,
    power_base INT,
    power_bonus INT,
    m_power_base INT,
    m_power_bonus INT,
    magic_range INT,
    multiplier_base FLOAT(5, 3),
    multiplier_bonus FLOAT(5, 3)
);

-- GameShop
CREATE TABLE game_shop_item (
    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    game_shop_item_index INT,
    item_index INT,
    gold_price INT,
    credit_price INT,
    count INT,
    class VARCHAR(200),
    category VARCHAR(200),
    stock INT,
    i_stock INT,
    deal INT,
    top_item INT,
    create_date DATE
);

-- Conquest
-- RespawnTick
-- CREATE TABLE active_coords();