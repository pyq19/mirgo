package main

import (
	"fmt"
	"path/filepath"
	"strconv"
)

// ExpTable 表示经验表配置。
type ExpTable struct {
	HighLevel          int              `json:"highLevel"`
	KillMonExpMultiple int              `json:"killMonExpMultiple"`
	BaseExp            int64            `json:"baseExp"`
	Levels             map[string]int64 `json:"levels"`
}

// StringsConfig 表示服务端字符串配置。
type StringsConfig struct {
	Messages map[string]string `json:"messages"`
}

// GlobalVars 表示全局变量配置。
type GlobalVars struct {
	Variables map[string]int `json:"variables"`
}

// ConvertExpTable 转换经验表配置。
func ConvertExpTable(inputDir, outputDir string) error {
	expFile := filepath.Join(inputDir, "Exps.ini")
	expINI, err := ParseINI(expFile)
	if err != nil {
		return fmt.Errorf("parsing Exps.ini: %w", err)
	}

	expTable := ExpTable{
		HighLevel:          getINIInt(expINI, "Exp", "HighLevel", 1000),
		KillMonExpMultiple: getINIInt(expINI, "Exp", "KillMonExpMultiple", 1),
		Levels:             make(map[string]int64),
	}

	// Parse base exp
	if baseExp, ok := expINI["Exp"]["BaseExp"]; ok {
		if n, err := strconv.ParseInt(baseExp, 10, 64); err == nil {
			expTable.BaseExp = n
		}
	}

	// Parse level exp values
	if expSection, ok := expINI["Exp"]; ok {
		for k, v := range expSection {
			if len(k) > 5 && k[:5] == "Level" {
				if n, err := strconv.ParseInt(v, 10, 64); err == nil {
					expTable.Levels[k] = n
				}
			}
		}
	}

	// 序列化为 JSONC（带字段级注释）
	data, err := MarshalJSONC(expTable, expComments, false)
	if err != nil {
		return fmt.Errorf("marshaling exp table: %w", err)
	}

	// Write output
	outputFile := filepath.Join(outputDir, "exp_table.jsonc")
	comment := "经验值表\n来源: asset/server/Exps.ini\n说明: 各等级所需经验值"

	return WriteJSONC(outputFile, data, comment)
}

// ConvertStrings 转换服务端字符串配置。
func ConvertStrings(inputDir, outputDir string) error {
	strFile := filepath.Join(inputDir, "String.ini")
	strINI, err := ParseINI(strFile)
	if err != nil {
		return fmt.Errorf("parsing String.ini: %w", err)
	}

	config := StringsConfig{
		Messages: make(map[string]string),
	}

	// Get all messages from [String] section
	if strSection, ok := strINI["String"]; ok {
		for k, v := range strSection {
			config.Messages[k] = v
		}
	}

	// 序列化为 JSONC（带字段级注释）
	data, err := MarshalJSONC(config, stringsComments, false)
	if err != nil {
		return fmt.Errorf("marshaling strings: %w", err)
	}

	// Write output
	outputFile := filepath.Join(outputDir, "strings.jsonc")
	comment := "服务端提示文本\n来源: asset/server/String.ini\n说明: 游戏内提示消息（约210条）"

	return WriteJSONC(outputFile, data, comment)
}

// ConvertGlobalVars 转换全局变量配置。
func ConvertGlobalVars(inputDir, outputDir string) error {
	gvFile := filepath.Join(inputDir, "GlobalVal.ini")
	gvINI, err := ParseINI(gvFile)
	if err != nil {
		return fmt.Errorf("parsing GlobalVal.ini: %w", err)
	}

	config := GlobalVars{
		Variables: make(map[string]int),
	}

	// Get all variables from any section
	for _, section := range gvINI {
		for k, v := range section {
			if n, err := strconv.Atoi(v); err == nil {
				config.Variables[k] = n
			}
		}
	}

	// 序列化为 JSONC（带字段级注释）
	data, err := MarshalJSONC(config, globalVarsComments, false)
	if err != nil {
		return fmt.Errorf("marshaling global vars: %w", err)
	}

	// Write output
	outputFile := filepath.Join(outputDir, "global_vars.jsonc")
	comment := "全局变量\n来源: asset/server/GlobalVal.ini\n说明: 2000个全局变量（GlobalVal0~GlobalVal1999）"

	return WriteJSONC(outputFile, data, comment)
}

// expComments 是 exp_table.jsonc 的中文注释表。
// 文案来源：Delphi M2Share.pas/ObjBase.pas 的经验逻辑。
var expComments = map[string]string{
	"highLevel":          "最高等级上限（Delphi HighLevel；超过后按封顶经验处理）",
	"killMonExpMultiple": "杀怪经验倍数，以怪物数据库经验值为基准（Delphi dwKillMonExpMultiple，默认1）",
	"baseExp":            "基础经验值（Exps.ini BaseExp）",
	// levels 为动态键 map，整体说明放键上方
	"levels": "各等级所需累计经验值，键为 LevelN（N=等级），值为升到该等级所需总经验",
}

// stringsComments 是 strings.jsonc 的中文注释表。
var stringsComments = map[string]string{
	// messages 为动态键 map，整体说明放键上方
	"messages": "游戏内提示消息映射，键为消息标识（如 ClientSoftVersionError、AttackModeOfAll），值为提示文本",
}

// globalVarsComments 是 global_vars.jsonc 的中文注释表。
var globalVarsComments = map[string]string{
	// variables 为动态键 map，整体说明放键上方
	"variables": "全局变量表，键为 GlobalValN（N=0~1999），值为整型变量，供 NPC 脚本读写",
}
