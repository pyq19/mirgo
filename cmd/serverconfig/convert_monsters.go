package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// MonGen 表示怪物刷怪点。
type MonGen struct {
	MapName  string `json:"mapName"`
	X        int    `json:"x"`
	Y        int    `json:"y"`
	Name     string `json:"name"`
	Range    int    `json:"range"`
	Count    int    `json:"count"`
	Interval int    `json:"interval"`
}

// MonItem 表示怪物掉落物品。
type MonItem struct {
	Prob  string `json:"prob"`
	Name  string `json:"name"`
	Count int    `json:"count,omitempty"`
}

// ConvertMonsters 转换怪物配置文件。
func ConvertMonsters(inputDir, outputDir string) error {
	envirDir := filepath.Join(inputDir, "Envir")

	// Convert mongen.txt
	if err := convertMonGen(envirDir, outputDir); err != nil {
		return fmt.Errorf("converting mongen.txt: %w", err)
	}

	// 转换 MonItems/*.txt
	if err := convertMonItems(envirDir, outputDir); err != nil {
		return fmt.Errorf("converting MonItems: %w", err)
	}

	// Copy SmartMonster/*.ini
	if err := copySmartMonster(envirDir, outputDir); err != nil {
		return fmt.Errorf("copying SmartMonster: %w", err)
	}

	// Copy MonUseItems/*.txt（人形怪装备配置）
	if err := copyMonUseItems(envirDir, outputDir); err != nil {
		return fmt.Errorf("copying MonUseItems: %w", err)
	}

	return nil
}

func convertMonGen(envirDir, outputDir string) error {
	mongenFile := filepath.Join(envirDir, "mongen.txt")
	data, err := ReadGBKFile(mongenFile)
	if err != nil {
		return err
	}

	var gens []MonGen
	scanner := bufio.NewScanner(strings.NewReader(string(data)))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || line[0] == ';' {
			continue
		}

		parts := strings.Fields(line)
		if len(parts) >= 6 {
			var gen MonGen
			gen.MapName = parts[0]
			fmt.Sscanf(parts[1], "%d", &gen.X)
			fmt.Sscanf(parts[2], "%d", &gen.Y)
			gen.Name = parts[3]
			fmt.Sscanf(parts[4], "%d", &gen.Range)
			fmt.Sscanf(parts[5], "%d", &gen.Count)
			if len(parts) > 6 {
				fmt.Sscanf(parts[6], "%d", &gen.Interval)
			}
			gens = append(gens, gen)
		}
	}

	result := map[string]interface{}{
		"spawns": gens,
	}

	jsonData, err := MarshalJSONC(result, monGenComments, false)
	if err != nil {
		return fmt.Errorf("marshaling mon gen: %w", err)
	}

	outputFile := filepath.Join(outputDir, "monsters", "mon_gen.jsonc")
	comment := fmt.Sprintf("怪物刷新点（怪物刷新点定义）\n来源: asset/server/Envir/mongen.txt\n数量: %d 个刷新点\n", len(gens)) + monGenFieldDocs

	return WriteJSONC(outputFile, jsonData, comment)
}

// monGenFieldDocs 是 mon_gen.jsonc 头部注释追加的字段说明表，
// 按 MonGen 结构体字段顺序，每个字段一行。
const monGenFieldDocs = `字段说明（按结构体字段顺序）:
mapName  地图编号（刷新点所在地图，不在地图列表中的条目加载时跳过）
x        刷新点中心 X 坐标
y        刷新点中心 Y 坐标
name     怪物名称（引用 Monster 表名称）
range    刷新半径（格；在中心 ±range 内随机选可走点）
count    维持存活怪物数量（按 game.monGenRate 缩放）
interval 重生间隔（分钟；Delphi LocalDB.pas:1237 按分钟×60×1000 转 ms，空值按 10 分钟处理）`

// monGenComments 是 mon_gen.jsonc 的中文注释表。
// 键为点号路径；commentAll=false 时元素注释只渲染在首条记录。
// 文案来源：cmd/server/mongen.go loadMonGenFromFile/SpawnMonster 消费逻辑。
var monGenComments = map[string]string{
	"spawns":            "怪物刷新点（mongen.txt 每行一条）",
	"spawns[].mapName":  "地图编号（不在地图列表则跳过）",
	"spawns[].x":        "刷新点 X 坐标",
	"spawns[].y":        "刷新点 Y 坐标",
	"spawns[].name":     "怪物名称（引用 Monster 表）",
	"spawns[].range":    "刷新半径（格）",
	"spawns[].count":    "维持存活怪物数量",
	"spawns[].interval": "重生间隔（分钟，服务端 ×60×1000 转 ms）",
}

func convertMonItems(envirDir, outputDir string) error {
	monItemsDir := filepath.Join(envirDir, "MonItems")
	dstDir := filepath.Join(outputDir, "monsters", "mon_items")

	// 读取目录并按扩展名不区分大小写过滤（.txt/.TXT 均命中）。
	entries, err := os.ReadDir(monItemsDir)
	if err != nil {
		return err
	}

	count := 0
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		base := entry.Name()
		ext := filepath.Ext(base)
		if !strings.EqualFold(ext, ".txt") {
			continue
		}
		srcFile := filepath.Join(monItemsDir, base)

		data, err := ReadGBKFile(srcFile)
		if err != nil {
			fmt.Printf("  警告: 读取 %s 失败: %v\n", srcFile, err)
			continue
		}

		monsterName := strings.TrimSuffix(base, ext)
		var items []MonItem

		scanner := bufio.NewScanner(strings.NewReader(string(data)))
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line == "" || line[0] == ';' {
				continue
			}

			// Format: 1/N itemname [count]
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				var item MonItem
				item.Prob = parts[0]
				item.Name = parts[1]
				if len(parts) > 2 {
					fmt.Sscanf(parts[2], "%d", &item.Count)
				}
				items = append(items, item)
			}
		}

		result := map[string]interface{}{
			"monster": monsterName,
			"items":   items,
		}

		jsonData, err := MarshalJSONC(result, monItemComments, true)
		if err != nil {
			continue
		}

		outputFile := filepath.Join(dstDir, monsterName+".jsonc")
		comment := fmt.Sprintf("怪物掉落表（%s 的掉落表）\n来源: asset/server/Envir/MonItems/%s\n怪物: %s\n", monsterName, base, monsterName) + monItemFieldDocs

		if err := WriteJSONC(outputFile, jsonData, comment); err != nil {
			fmt.Printf("  警告: 写入 %s 失败: %v\n", outputFile, err)
			continue
		}
		count++
	}

	fmt.Printf("  转换了 %d 个怪物掉落表\n", count)
	return nil
}

func copySmartMonster(envirDir, outputDir string) error {
	srcDir := filepath.Join(envirDir, "SmartMonster")
	dstDir := filepath.Join(outputDir, "monsters", "smart_monster")

	count, err := CopyDir(srcDir, dstDir, "*.ini")
	if err != nil {
		return err
	}

	fmt.Printf("  复制了 %d 个智能怪物配置\n", count)
	return nil
}

func copyMonUseItems(envirDir, outputDir string) error {
	srcDir := filepath.Join(envirDir, "MonUseItems")
	dstDir := filepath.Join(outputDir, "monsters", "mon_use_items")

	if !DirExists(srcDir) {
		fmt.Println("  跳过 MonUseItems (目录不存在)")
		return nil
	}

	count, err := CopyDir(srcDir, dstDir, "*.txt")
	if err != nil {
		return err
	}

	fmt.Printf("  复制了 %d 个怪物装备配置\n", count)
	return nil
}

// monItemFieldDocs 是 mon_items/*.jsonc 头部注释追加的字段说明表，
// 按 MonItem 结构体字段顺序，每个字段一行。
const monItemFieldDocs = `字段说明（按结构体字段顺序）:
prob  掉落概率 "1/N"（服务端取分母 N 作为判定基数，1/N 概率掉落）
name  物品名称（"金币" 表示金币掉落）
count 掉落数量（省略时按 1；金币条目为金币数量）`

// monItemComments 是 mon_items/*.jsonc 的中文注释表。
// 键为点号路径；小文件用 commentAll=true 逐元素渲染。
// 文案来源：cmd/server/droptable.go loadFile/parseProb 消费逻辑。
var monItemComments = map[string]string{
	"monster":       "怪物名（掉落表索引键，对应 Monster 表名称）",
	"items":         "掉落行（原文格式: prob name [count]）",
	"items[].prob":  "掉落概率 1/N",
	"items[].name":  "物品名（\"金币\" 为金币）",
	"items[].count": "掉落数量（省略按 1）",
}
