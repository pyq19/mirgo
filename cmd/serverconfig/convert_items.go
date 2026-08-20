package main

import (
	"bufio"
	"fmt"
	"path/filepath"
	"strings"
)

// FilterItem 表示过滤列表中的物品。
type FilterItem struct {
	Type  int    `json:"type"`
	Name  string `json:"name"`
	Props []int  `json:"props"`
}

// ItemRule 表示物品规则。
type ItemRule struct {
	Name  string `json:"name"`
	Rules string `json:"rules"`
}

// GroupItem 表示物品套装/分组。
type GroupItem struct {
	ID    int      `json:"id"`
	Count int      `json:"count"`
	Name  string   `json:"name"`
	Items []string `json:"items"`
	Bonus string   `json:"bonus"`
	Desc  string   `json:"desc"`
}

// MakeItem 表示合成配方。
// schema 与服务端 cmd/server/makedrug.go 的 DrugRecipe 对齐（product/materials）。
type MakeItem struct {
	Product   string         `json:"product"`
	Materials []MakeItemPart `json:"materials"`
}

// MakeItemPart 表示合成配方中的材料。
type MakeItemPart struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

// ConvertItems 转换物品配置文件。
func ConvertItems(inputDir, outputDir string) error {
	envirDir := filepath.Join(inputDir, "Envir")

	// 转换 FilterItemList.txt
	if err := convertFilterItems(envirDir, outputDir); err != nil {
		return fmt.Errorf("converting FilterItemList.txt: %w", err)
	}

	// 转换 ItemRuleList.txt
	if err := convertItemRules(envirDir, outputDir); err != nil {
		return fmt.Errorf("converting ItemRuleList.txt: %w", err)
	}

	// 转换 GroupItemList.txt
	if err := convertGroupItems(envirDir, outputDir); err != nil {
		return fmt.Errorf("converting GroupItemList.txt: %w", err)
	}

	// 转换 UnbindList.txt
	if err := convertUnbindList(envirDir, outputDir); err != nil {
		return fmt.Errorf("converting UnbindList.txt: %w", err)
	}

	// 转换 MakeItem.txt
	if err := convertMakeItems(envirDir, outputDir); err != nil {
		return fmt.Errorf("converting MakeItem.txt: %w", err)
	}

	// 转换 DisableTakeOffList.txt
	if err := convertDisableTakeOffList(envirDir, outputDir); err != nil {
		return fmt.Errorf("converting DisableTakeOffList.txt: %w", err)
	}

	return nil
}

// convertDisableTakeOffList 转换禁止脱下物品列表。
// Delphi 格式 "物品名 索引"（M2Share.pas:4578-4609），索引 = wIndex-1，
// 即 0 基 DB 索引；Go 的 WIndex 同为 0 基，直接保存。
func convertDisableTakeOffList(envirDir, outputDir string) error {
	listFile := filepath.Join(envirDir, "DisableTakeOffList.txt")
	data, err := ReadGBKFile(listFile)
	if err != nil {
		return err
	}

	type TakeOffItem struct {
		Name string `json:"name"`
		Idx  int    `json:"idx"`
	}

	var items []TakeOffItem
	scanner := bufio.NewScanner(strings.NewReader(string(data)))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || line[0] == ';' {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) >= 2 {
			var item TakeOffItem
			item.Name = parts[0]
			fmt.Sscanf(parts[1], "%d", &item.Idx)
			items = append(items, item)
		}
	}

	result := map[string]interface{}{
		"items": items,
	}

	jsonData, err := MarshalJSONC(result, takeoffComments, true)
	if err != nil {
		return err
	}

	outputFile := filepath.Join(outputDir, "items", "disable_takeoff_list.jsonc")
	comment := fmt.Sprintf("禁止脱下物品列表\n来源: asset/server/Envir/DisableTakeOffList.txt\n数量: %d 个物品\n字段说明:\n  name  物品名（仅备注用，匹配以 idx 为准）\n  idx   物品 0 基数据库索引（即 StdItems 行号；Delphi 原文件记为 wIndex-1）", len(items))

	return WriteJSONC(outputFile, jsonData, comment)
}

func convertFilterItems(envirDir, outputDir string) error {
	filterFile := filepath.Join(envirDir, "FilterItemList.txt")
	data, err := ReadGBKFile(filterFile)
	if err != nil {
		return err
	}

	var items []FilterItem
	scanner := bufio.NewScanner(strings.NewReader(string(data)))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || line[0] == ';' {
			continue
		}

		parts := strings.Fields(line)
		if len(parts) >= 2 {
			var item FilterItem
			fmt.Sscanf(parts[0], "%d", &item.Type)
			item.Name = parts[1]

			// Parse properties
			for i := 2; i < len(parts); i++ {
				var prop int
				fmt.Sscanf(parts[i], "%d", &prop)
				item.Props = append(item.Props, prop)
			}

			items = append(items, item)
		}
	}

	result := map[string]interface{}{
		"items": items,
	}

	jsonData, err := MarshalJSONC(result, filterItemComments, false)
	if err != nil {
		return err
	}

	outputFile := filepath.Join(outputDir, "items", "filter_items.jsonc")
	comment := fmt.Sprintf("物品过滤列表\n来源: asset/server/Envir/FilterItemList.txt\n数量: %d 个物品\n字段说明:\n  type   过滤分类编号（原文件第 1 列）\n  name   物品名\n  props  附加数值参数（原文件第 3 列起）", len(items))

	return WriteJSONC(outputFile, jsonData, comment)
}

func convertItemRules(envirDir, outputDir string) error {
	rulesFile := filepath.Join(envirDir, "ItemRuleList.txt")
	data, err := ReadGBKFile(rulesFile)
	if err != nil {
		return err
	}

	var rules []ItemRule
	scanner := bufio.NewScanner(strings.NewReader(string(data)))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || line[0] == ';' {
			continue
		}

		parts := strings.SplitN(line, "\t", 2)
		if len(parts) >= 2 {
			rules = append(rules, ItemRule{
				Name:  strings.TrimSpace(parts[0]),
				Rules: strings.TrimSpace(parts[1]),
			})
		}
	}

	result := map[string]interface{}{
		"rules": rules,
	}

	jsonData, err := MarshalJSONC(result, itemRuleComments, false)
	if err != nil {
		return err
	}

	outputFile := filepath.Join(outputDir, "items", "item_rules.jsonc")
	comment := fmt.Sprintf("物品规则列表\n来源: asset/server/Envir/ItemRuleList.txt\n数量: %d 条规则\n字段说明:\n  name   物品名\n  rules  规则标志串（原文件制表符后的数值序列，逐位对应各项规则开关）", len(rules))

	return WriteJSONC(outputFile, jsonData, comment)
}

func convertGroupItems(envirDir, outputDir string) error {
	groupFile := filepath.Join(envirDir, "GroupItemList.txt")
	data, err := ReadGBKFile(groupFile)
	if err != nil {
		return err
	}

	var groups []GroupItem
	scanner := bufio.NewScanner(strings.NewReader(string(data)))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || line[0] == ';' {
			continue
		}

		parts := strings.Split(line, "\t")
		if len(parts) >= 5 {
			var group GroupItem
			fmt.Sscanf(parts[0], "%d", &group.ID)
			fmt.Sscanf(parts[1], "%d", &group.Count)
			group.Name = parts[2]
			group.Items = strings.Split(parts[3], "|")
			group.Bonus = parts[4]
			if len(parts) > 5 {
				group.Desc = parts[5]
			}
			groups = append(groups, group)
		}
	}

	result := map[string]interface{}{
		"groups": groups,
	}

	jsonData, err := MarshalJSONC(result, groupItemComments, true)
	if err != nil {
		return err
	}

	outputFile := filepath.Join(outputDir, "items", "group_items.jsonc")
	comment := fmt.Sprintf("套装定义\n来源: asset/server/Envir/GroupItemList.txt\n数量: %d 套\n字段说明:\n  id     套装编号\n  count  激活套装所需件数\n  name   套装名称\n  items  套装包含的物品名列表\n  bonus  套装加成属性值（竖线分隔数值）\n  desc   套装生效提示文本", len(groups))

	return WriteJSONC(outputFile, jsonData, comment)
}

func convertUnbindList(envirDir, outputDir string) error {
	unbindFile := filepath.Join(envirDir, "UnbindList.txt")
	data, err := ReadGBKFile(unbindFile)
	if err != nil {
		return err
	}

	type UnbindItem struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}

	var items []UnbindItem
	scanner := bufio.NewScanner(strings.NewReader(string(data)))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || line[0] == ';' {
			continue
		}

		parts := strings.Fields(line)
		if len(parts) >= 2 {
			var item UnbindItem
			fmt.Sscanf(parts[0], "%d", &item.ID)
			item.Name = parts[1]
			items = append(items, item)
		}
	}

	result := map[string]interface{}{
		"items": items,
	}

	jsonData, err := MarshalJSONC(result, unbindComments, true)
	if err != nil {
		return err
	}

	outputFile := filepath.Join(outputDir, "items", "unbind_list.jsonc")
	comment := fmt.Sprintf("解绑物品列表\n来源: asset/server/Envir/UnbindList.txt\n数量: %d 个物品\n字段说明:\n  id    Shape 编号（对应 StdMode 31 包裹类物品的 Shape 值）\n  name  解包后得到的物品名", len(items))

	return WriteJSONC(outputFile, jsonData, comment)
}

func convertMakeItems(envirDir, outputDir string) error {
	makeFile := filepath.Join(envirDir, "MakeItem.txt")
	data, err := ReadGBKFile(makeFile)
	if err != nil {
		return err
	}

	var recipes []MakeItem
	var currentRecipe *MakeItem

	scanner := bufio.NewScanner(strings.NewReader(string(data)))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || line[0] == ';' {
			continue
		}

		// Recipe header: [item name]
		if line[0] == '[' && line[len(line)-1] == ']' {
			if currentRecipe != nil {
				recipes = append(recipes, *currentRecipe)
			}
			name := line[1 : len(line)-1]
			currentRecipe = &MakeItem{Product: name}
			continue
		}

		// Ingredient: name count
		if currentRecipe != nil {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				var part MakeItemPart
				part.Name = parts[0]
				fmt.Sscanf(parts[1], "%d", &part.Count)
				currentRecipe.Materials = append(currentRecipe.Materials, part)
			}
		}
	}

	if currentRecipe != nil {
		recipes = append(recipes, *currentRecipe)
	}

	result := map[string]interface{}{
		"recipes": recipes,
	}

	jsonData, err := MarshalJSONC(result, makeItemComments, false)
	if err != nil {
		return err
	}

	outputFile := filepath.Join(outputDir, "items", "make_items.jsonc")
	comment := fmt.Sprintf("合成配方\n来源: asset/server/Envir/MakeItem.txt\n数量: %d 个配方\n字段说明:\n  product    产物物品名\n  materials  所需材料列表（name 材料名，count 数量）\n说明: 服务端 makedrug.go 消费；未列出 goldCost，缺省按金币模式处理", len(recipes))

	return WriteJSONC(outputFile, jsonData, comment)
}

// filterItemComments 是 filter_items.jsonc 的中文注释表（条目多，仅首条记录渲染，commentAll=false）。
var filterItemComments = map[string]string{
	// items 为非空数组，注释放键上方
	"items":         "物品过滤列表（GEEM2 FilterItemList，Go 服务端暂未消费，保留）",
	"items[].type":  "过滤分类编号（原文件第 1 列）",
	"items[].name":  "物品名",
	"items[].props": "附加数值参数（原文件第 3 列起）",
}

// itemRuleComments 是 item_rules.jsonc 的中文注释表（条目多，仅首条记录渲染，commentAll=false）。
var itemRuleComments = map[string]string{
	// rules 为非空数组，注释放键上方
	"rules":         "物品规则列表（逐物品行为开关，Go 服务端暂未消费，保留）",
	"rules[].name":  "物品名",
	"rules[].rules": "规则标志串（原文件制表符后的数值序列，逐位对应各项规则开关）",
}

// groupItemComments 是 group_items.jsonc 的中文注释表（条目少，逐元素渲染，commentAll=true）。
var groupItemComments = map[string]string{
	// groups 为非空数组，注释放键上方
	"groups":         "套装列表（Go 服务端暂未消费，保留）",
	"groups[].id":    "套装编号",
	"groups[].count": "激活套装所需件数",
	"groups[].name":  "套装名称",
	"groups[].items": "套装包含的物品名列表",
	"groups[].bonus": "套装加成属性值（竖线分隔数值）",
	"groups[].desc":  "套装生效提示文本",
}

// unbindComments 是 unbind_list.jsonc 的中文注释表（条目少，逐元素渲染，commentAll=true）。
// 文案来源：cmd/server/itemdb.go LoadUnbindList（StdMode 31 解包）。
var unbindComments = map[string]string{
	// items 为非空数组，注释放键上方
	"items":        "解绑（解包）映射列表（服务端 ItemDB.LoadUnbindList 消费）",
	"items[].id":   "Shape 编号（对应 StdMode 31 包裹类物品的 Shape 值）",
	"items[].name": "解包后得到的物品名",
}

// takeoffComments 是 disable_takeoff_list.jsonc 的中文注释表（条目少，逐元素渲染，commentAll=true）。
// 文案来源：cmd/server/itemdb.go LoadDisableTakeOffList。
var takeoffComments = map[string]string{
	// items 为非空数组，注释放键上方
	"items":        "禁止脱下物品列表（服务端 ItemDB.LoadDisableTakeOffList 消费，禁止脱下/死亡掉落）",
	"items[].name": "物品名（仅备注用，匹配以 idx 为准）",
	"items[].idx":  "物品 0 基数据库索引（即 StdItems 行号；Delphi 原文件记为 wIndex-1）",
}

// makeItemComments 是 make_items.jsonc 的中文注释表（59 条 > 50，仅首条记录渲染，commentAll=false）。
// 文案来源：cmd/server/makedrug.go LoadDrugRecipes/DrugRecipe。
var makeItemComments = map[string]string{
	// recipes 为非空数组，注释放键上方
	"recipes":                     "制药/合成配方列表（服务端 LoadDrugRecipes 消费）",
	"recipes[].product":           "产物物品名",
	"recipes[].materials":         "所需材料列表",
	"recipes[].materials[].name":  "材料物品名",
	"recipes[].materials[].count": "所需材料数量",
}
