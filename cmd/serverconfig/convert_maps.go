package main

import (
	"bufio"
	"fmt"
	"path/filepath"
	"strings"
)

// MapInfo 表示地图定义。
// Index 为字符串：地图编号可以是数字（"0"）或字母（"D001"、"dygw"）。
type MapInfo struct {
	Index string   `json:"index"`
	Name  string   `json:"name"`
	Flags int      `json:"flags"`
	Props []string `json:"props,omitempty"`
}

// MapRoute 表示地图间的传送路线。
type MapRoute struct {
	SrcMap string `json:"srcMap"`
	SrcX   int    `json:"srcX"`
	SrcY   int    `json:"srcY"`
	DstMap string `json:"dstMap"`
	DstX   int    `json:"dstX"`
	DstY   int    `json:"dstY"`
}

// MiniMap 表示小地图映射。
type MiniMap struct {
	MapName   string `json:"mapName"`
	MiniMapID int    `json:"miniMapId"`
}

// StartPoint 表示安全区/重生点。
type StartPoint struct {
	MapName string `json:"mapName"`
	X       int    `json:"x"`
	Y       int    `json:"y"`
	Range   int    `json:"range"`
}

// ConvertMaps 转换地图配置文件。
func ConvertMaps(inputDir, outputDir string) error {
	envirDir := filepath.Join(inputDir, "Envir")

	// Convert mapinfo.txt
	if err := convertMapInfo(envirDir, outputDir); err != nil {
		return fmt.Errorf("converting mapinfo.txt: %w", err)
	}

	// Copy .map files
	if err := copyMapFiles(inputDir, outputDir); err != nil {
		return fmt.Errorf("copying map files: %w", err)
	}

	// 转换 MiniMap.txt
	if err := convertMiniMap(envirDir, outputDir); err != nil {
		return fmt.Errorf("converting MiniMap.txt: %w", err)
	}

	// 转换 StartPoint.txt
	if err := convertStartPoint(envirDir, outputDir); err != nil {
		return fmt.Errorf("converting StartPoint.txt: %w", err)
	}

	return nil
}

func convertMapInfo(envirDir, outputDir string) error {
	mapInfoFile := filepath.Join(envirDir, "mapinfo.txt")
	data, err := ReadGBKFile(mapInfoFile)
	if err != nil {
		return err
	}

	var maps []MapInfo
	var routes []MapRoute

	scanner := bufio.NewScanner(strings.NewReader(string(data)))
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || line[0] == ';' {
			continue
		}

		// 路线定义: srcMap x,y -> dstMap x,y
		if strings.Contains(line, "->") {
			parts := strings.SplitN(line, "->", 2)
			srcParts := strings.Fields(strings.TrimSpace(parts[0]))
			dstParts := strings.Fields(strings.TrimSpace(parts[1]))
			if len(srcParts) >= 2 && len(dstParts) >= 2 {
				var route MapRoute
				route.SrcMap = srcParts[0]
				fmt.Sscanf(srcParts[1], "%d,%d", &route.SrcX, &route.SrcY)
				route.DstMap = dstParts[0]
				fmt.Sscanf(dstParts[1], "%d,%d", &route.DstX, &route.DstY)
				routes = append(routes, route)
			}
			continue
		}

		// 地图定义: [index name flags] [PROPS...]
		// 属性后缀出现在右括号之后（如 [D001 兽人古墓一层 0] DARK NORECALL），
		// flags 可省略（如 [6 魔龙城]）。
		openIdx := strings.Index(line, "[")
		closeIdx := strings.Index(line, "]")
		if openIdx >= 0 && closeIdx > openIdx {
			content := line[openIdx+1 : closeIdx]
			propsStr := strings.TrimSpace(line[closeIdx+1:])
			parts := strings.Fields(content)
			if len(parts) >= 2 {
				mi := MapInfo{Index: parts[0], Name: parts[1]}
				if len(parts) >= 3 {
					fmt.Sscanf(parts[2], "%d", &mi.Flags)
				}
				if propsStr != "" {
					mi.Props = strings.Fields(propsStr)
				}
				maps = append(maps, mi)
			}
		}
	}

	// 地图定义 → maps/map_info.jsonc
	mapResult := map[string]interface{}{
		"maps": maps,
	}
	mapJSON, err := MarshalJSONC(mapResult, mapInfoComments, false)
	if err != nil {
		return err
	}
	mapComment := fmt.Sprintf("地图信息\n来源: asset/server/Envir/mapinfo.txt\n数量: %d 个地图\n字段说明:\n  index  地图编号（对应 Map 目录下 .map 文件名，可为数字或字母）\n  name   地图标题（中文名）\n  flags  原版数字标志位（Go 服务端暂未消费，属性以 props 为准）\n  props  地图属性 token（SAFE/DARK/MINE/NORECONNECT(0)/EXPRATE(3) 等，见服务端 parseMapFlag；可缺省）", len(maps))
	if err := WriteJSONC(filepath.Join(outputDir, "maps", "map_info.jsonc"), mapJSON, mapComment); err != nil {
		return err
	}

	// 传送路线 → maps/map_routes.jsonc（与 server MapManager.InitRoutes 约定一致）
	routeResult := map[string]interface{}{
		"routes": routes,
	}
	routeJSON, err := MarshalJSONC(routeResult, mapRouteComments, false)
	if err != nil {
		return err
	}
	routeComment := fmt.Sprintf("传送路线\n来源: asset/server/Envir/mapinfo.txt\n数量: %d 条路线\n字段说明:\n  srcMap/srcX/srcY  出发点地图编号与坐标（踩到该坐标触发传送）\n  dstMap/dstX/dstY  目标点地图编号与坐标", len(routes))
	return WriteJSONC(filepath.Join(outputDir, "maps", "map_routes.jsonc"), routeJSON, routeComment)
}

func copyMapFiles(inputDir, outputDir string) error {
	srcDir := filepath.Join(inputDir, "Map")
	dstDir := filepath.Join(outputDir, "maps")

	count, err := CopyDir(srcDir, dstDir, "*.map")
	if err != nil {
		return err
	}

	fmt.Printf("  复制了 %d 个地图文件\n", count)
	return nil
}

func convertMiniMap(envirDir, outputDir string) error {
	miniMapFile := filepath.Join(envirDir, "MiniMap.txt")
	data, err := ReadGBKFile(miniMapFile)
	if err != nil {
		return err
	}

	var miniMaps []MiniMap
	scanner := bufio.NewScanner(strings.NewReader(string(data)))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || line[0] == ';' {
			continue
		}

		parts := strings.Fields(line)
		if len(parts) >= 2 {
			var mm MiniMap
			mm.MapName = parts[0]
			fmt.Sscanf(parts[1], "%d", &mm.MiniMapID)
			miniMaps = append(miniMaps, mm)
		}
	}

	result := map[string]interface{}{
		"miniMaps": miniMaps,
	}

	jsonData, err := MarshalJSONC(result, miniMapComments, false)
	if err != nil {
		return err
	}

	outputFile := filepath.Join(outputDir, "maps", "mini_map.jsonc")
	comment := fmt.Sprintf("小地图映射\n来源: asset/server/Envir/MiniMap.txt\n数量: %d 个映射\n字段说明:\n  mapName     地图编号（对应 map_info 的 index）\n  miniMapId   小地图图像编号（客户端小地图资源索引，乱改将看不到小地图）", len(miniMaps))

	return WriteJSONC(outputFile, jsonData, comment)
}

func convertStartPoint(envirDir, outputDir string) error {
	startPointFile := filepath.Join(envirDir, "StartPoint.txt")
	data, err := ReadGBKFile(startPointFile)
	if err != nil {
		return err
	}

	var startPoints []StartPoint
	scanner := bufio.NewScanner(strings.NewReader(string(data)))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || line[0] == ';' {
			continue
		}

		parts := strings.Fields(line)
		if len(parts) >= 4 {
			var sp StartPoint
			sp.MapName = parts[0]
			fmt.Sscanf(parts[1], "%d", &sp.X)
			fmt.Sscanf(parts[2], "%d", &sp.Y)
			fmt.Sscanf(parts[3], "%d", &sp.Range)
			startPoints = append(startPoints, sp)
		}
	}

	result := map[string]interface{}{
		"startPoints": startPoints,
	}

	jsonData, err := MarshalJSONC(result, startPointComments, false)
	if err != nil {
		return err
	}

	outputFile := filepath.Join(outputDir, "maps", "start_points.jsonc")
	comment := fmt.Sprintf("安全区/复活点\n来源: asset/server/Envir/StartPoint.txt\n数量: %d 个点\n字段说明:\n  mapName  地图编号\n  x/y      安全区中心坐标\n  range    安全区半径（格）；<=0 时服务端按默认 5 处理", len(startPoints))

	return WriteJSONC(outputFile, jsonData, comment)
}

// mapInfoComments 是 map_info.jsonc 的中文注释表（条目多，仅首条记录渲染，commentAll=false）。
// 文案来源：cmd/server/mapmanager.go InitMapFlags/parseMapFlag。
var mapInfoComments = map[string]string{
	// maps 为非空数组，注释放键上方
	"maps":         "地图列表（服务端 MapManager.InitMapFlags 消费，按 index 匹配已加载的 .map 文件）",
	"maps[].index": "地图编号（对应 Map 目录下同名 .map 文件，可为数字或字母）",
	"maps[].name":  "地图标题（中文名）",
	"maps[].flags": "原版数字标志位（Delphi nServerIndex；Go 服务端暂未消费）",
	"maps[].props": "地图属性 token（SAFE/DARK/DAY/MINE/NORECONNECT(0)/EXPRATE(3)/DECHP(5/1000) 等，见服务端 parseMapFlag）",
}

// mapRouteComments 是 map_routes.jsonc 的中文注释表（条目多，仅首条记录渲染，commentAll=false）。
// 文案来源：cmd/server/mapmanager.go InitRoutes/FindRoute 与 playobject.go CheckMapRoute。
var mapRouteComments = map[string]string{
	// routes 为非空数组，注释放键上方
	"routes":          "传送路线列表（玩家踩到 srcX/srcY 即传送至 dstMap 的 dstX/dstY）",
	"routes[].srcMap": "出发点地图编号",
	"routes[].srcX":   "出发点 X 坐标",
	"routes[].srcY":   "出发点 Y 坐标",
	"routes[].dstMap": "目标地图编号",
	"routes[].dstX":   "目标点 X 坐标",
	"routes[].dstY":   "目标点 Y 坐标",
}

// miniMapComments 是 mini_map.jsonc 的中文注释表（条目多，仅首条记录渲染，commentAll=false）。
// 文案来源：cmd/server/mapmanager.go InitMiniMaps。
var miniMapComments = map[string]string{
	// miniMaps 为非空数组，注释放键上方
	"miniMaps":             "小地图映射列表（服务端 MapManager.InitMiniMaps 消费，写入 Environment.MinMap）",
	"miniMaps[].mapName":   "地图编号（对应 map_info 的 index）",
	"miniMaps[].miniMapId": "小地图图像编号（客户端小地图资源索引，乱改将看不到小地图）",
}

// startPointComments 是 start_points.jsonc 的中文注释表（条目多，仅首条记录渲染，commentAll=false）。
// 文案来源：cmd/server/safezone.go LoadSafeZones。
var startPointComments = map[string]string{
	// startPoints 为非空数组，注释放键上方
	"startPoints":           "安全区/复活点列表（服务端 LoadSafeZones 消费）",
	"startPoints[].mapName": "地图编号",
	"startPoints[].x":       "安全区中心 X 坐标",
	"startPoints[].y":       "安全区中心 Y 坐标",
	"startPoints[].range":   "安全区半径（格）；<=0 时服务端按默认 5 处理",
}
