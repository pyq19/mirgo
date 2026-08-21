package main

import (
	"bufio"
	"fmt"
	"path/filepath"
	"strings"
)

// NpcDef 表示 NPC 定义。
type NpcDef struct {
	Name    string `json:"name"`
	Race    int    `json:"race"`
	MapName string `json:"mapName"`
	X       int    `json:"x"`
	Y       int    `json:"y"`
	Face    int    `json:"face"`
	Body    int    `json:"body"`
}

// MerchantDef 表示商人/NPC 位置。
type MerchantDef struct {
	ID      string `json:"id"`
	MapName string `json:"mapName"`
	X       int    `json:"x"`
	Y       int    `json:"y"`
	Name    string `json:"name"`
	Face    int    `json:"face"`
	Body    int    `json:"body"`
	Castle  int    `json:"castle"`
}

// ConvertNPCs 转换 NPC 配置文件。
func ConvertNPCs(inputDir, outputDir string) error {
	envirDir := filepath.Join(inputDir, "Envir")

	// 转换 Npcs.txt
	if err := convertNpcList(envirDir, outputDir); err != nil {
		return fmt.Errorf("converting Npcs.txt: %w", err)
	}

	// 转换 merchant.txt
	if err := convertMerchantList(envirDir, outputDir); err != nil {
		return fmt.Errorf("converting merchant.txt: %w", err)
	}

	// Copy market_def/*.txt
	if err := copyNpcScripts(envirDir, outputDir, "market_def", "merchant_scripts"); err != nil {
		return fmt.Errorf("copying market_def: %w", err)
	}

	// Copy Npc_def/*.txt
	if err := copyNpcScripts(envirDir, outputDir, "Npc_def", "npc_scripts"); err != nil {
		return fmt.Errorf("copying Npc_def: %w", err)
	}

	// Copy MapQuest_def/*.txt
	if err := copyNpcScripts(envirDir, outputDir, "MapQuest_def", "map_quest_scripts"); err != nil {
		return fmt.Errorf("copying MapQuest_def: %w", err)
	}

	return nil
}

func convertNpcList(envirDir, outputDir string) error {
	npcsFile := filepath.Join(envirDir, "Npcs.txt")
	data, err := ReadGBKFile(npcsFile)
	if err != nil {
		return err
	}

	var npcs []NpcDef
	scanner := bufio.NewScanner(strings.NewReader(string(data)))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || line[0] == ';' {
			continue
		}

		parts := strings.Fields(line)
		if len(parts) >= 6 {
			var npc NpcDef
			npc.Name = parts[0]
			fmt.Sscanf(parts[1], "%d", &npc.Race)
			npc.MapName = parts[2]
			fmt.Sscanf(parts[3], "%d", &npc.X)
			fmt.Sscanf(parts[4], "%d", &npc.Y)
			fmt.Sscanf(parts[5], "%d", &npc.Face)
			if len(parts) > 6 {
				fmt.Sscanf(parts[6], "%d", &npc.Body)
			}
			npcs = append(npcs, npc)
		}
	}

	result := map[string]interface{}{
		"npcs": npcs,
	}

	jsonData, err := MarshalJSONC(result, npcComments, false)
	if err != nil {
		return err
	}

	outputFile := filepath.Join(outputDir, "npcs", "npc_list.jsonc")
	comment := fmt.Sprintf("NPC定义\n来源: asset/server/Envir/Npcs.txt\n数量: %d 个NPC\n字段说明:\n  name     NPC 名字（也用于匹配脚本 npcs/npc_scripts/<name>-<map>.txt）\n  race     NPC 类型（原文件列：0=店主 1=国王 2=沙巴克官员；服务端仅对 15 用和平NPC图像）\n  mapName  所在地图编号\n  x/y      所在坐标\n  face     脸部图像号\n  body     身体图像号（外观 Appr）", len(npcs))

	return WriteJSONC(outputFile, jsonData, comment)
}

func convertMerchantList(envirDir, outputDir string) error {
	merchantFile := filepath.Join(envirDir, "merchant.txt")
	data, err := ReadGBKFile(merchantFile)
	if err != nil {
		return err
	}

	var merchants []MerchantDef
	scanner := bufio.NewScanner(strings.NewReader(string(data)))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || line[0] == ';' {
			continue
		}

		parts := strings.Fields(line)
		if len(parts) >= 7 {
			var m MerchantDef
			m.ID = parts[0]
			m.MapName = parts[1]
			fmt.Sscanf(parts[2], "%d", &m.X)
			fmt.Sscanf(parts[3], "%d", &m.Y)
			m.Name = parts[4]
			fmt.Sscanf(parts[5], "%d", &m.Face)
			fmt.Sscanf(parts[6], "%d", &m.Body)
			if len(parts) > 7 {
				fmt.Sscanf(parts[7], "%d", &m.Castle)
			}
			merchants = append(merchants, m)
		}
	}

	result := map[string]interface{}{
		"merchants": merchants,
	}

	jsonData, err := MarshalJSONC(result, merchantComments, false)
	if err != nil {
		return err
	}

	outputFile := filepath.Join(outputDir, "npcs", "merchant_list.jsonc")
	comment := fmt.Sprintf("商人/NPC位置\n来源: asset/server/Envir/merchant.txt\n数量: %d 个商人\n字段说明:\n  id       商人脚本 ID（匹配 npcs/merchant_scripts/<id>-<map>.txt）\n  mapName  所在地图编号\n  x/y      所在坐标\n  name     商店显示名称\n  face     脸部图像号\n  body     身体图像号（外观 Appr）\n  castle   是否属于沙巴克城堡（非 0 = 城堡商人：城战受限、城主行会折扣、收城堡税）", len(merchants))

	return WriteJSONC(outputFile, jsonData, comment)
}

// npcComments 是 npc_list.jsonc 的中文注释表（条目多，仅首条记录渲染，commentAll=false）。
// 文案来源：cmd/server/npclist.go、mongen.go LoadNpcs、npcobject.go Feature。
var npcComments = map[string]string{
	// npcs 为非空数组，注释放键上方
	"npcs":           "NPC 定义列表（服务端 LoadNpcList → LoadNpcs 消费）",
	"npcs[].name":    "NPC 名字（也用于匹配脚本 npcs/npc_scripts/<name>-<map>.txt）",
	"npcs[].race":    "NPC 类型（原文件列：0=店主 1=国王 2=沙巴克官员；服务端仅对 15 用和平NPC图像）",
	"npcs[].mapName": "所在地图编号",
	"npcs[].x":       "所在 X 坐标",
	"npcs[].y":       "所在 Y 坐标",
	"npcs[].face":    "脸部图像号",
	"npcs[].body":    "身体图像号（外观 Appr）",
}

// merchantComments 是 merchant_list.jsonc 的中文注释表（条目多，仅首条记录渲染，commentAll=false）。
// 文案来源：cmd/server/npclist.go、mongen.go LoadNpcs、merchantsystem.go（城堡折扣/税）。
var merchantComments = map[string]string{
	// merchants 为非空数组，注释放键上方
	"merchants":           "商人定义列表（服务端 LoadMerchantList → LoadNpcs 消费，IsMerchant=true）",
	"merchants[].id":      "商人脚本 ID（匹配 npcs/merchant_scripts/<id>-<map>.txt）",
	"merchants[].mapName": "所在地图编号",
	"merchants[].x":       "所在 X 坐标",
	"merchants[].y":       "所在 Y 坐标",
	"merchants[].name":    "商店显示名称",
	"merchants[].face":    "脸部图像号",
	"merchants[].body":    "身体图像号（外观 Appr）",
	"merchants[].castle":  "是否属于沙巴克城堡（非 0 = 城堡商人：城战受限、城主行会折扣、收城堡税）",
}

func copyNpcScripts(envirDir, outputDir, srcSubdir, dstSubdir string) error {
	srcDir := filepath.Join(envirDir, srcSubdir)
	dstDir := filepath.Join(outputDir, "npcs", dstSubdir)

	if !DirExists(srcDir) {
		fmt.Printf("  跳过 %s (目录不存在)\n", srcSubdir)
		return nil
	}

	// 递归复制以保留子目录（如 market_def/魔龙城、market_def/backup）。
	count, err := CopyDirRecursive(srcDir, dstDir, "*.txt")
	if err != nil {
		return err
	}

	fmt.Printf("  复制了 %d 个 %s 脚本\n", count, srcSubdir)
	return nil
}
