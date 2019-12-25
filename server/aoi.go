package main

import (
	"sync"
)

// Zinx应用-MMO游戏案例-(2)AOI兴趣点算法 https://www.jianshu.com/p/e5b5db9fa6fe

// AOIManager 每一张地图对应一个 AOI 管理模块
type AOIManager struct {
	MinX  int       // 区域左边界坐标
	MaxX  int       // 区域右边界坐标
	CntsX int       // x方向格子的数量
	MinY  int       // 区域上边界坐标
	MaxY  int       // 区域下边界坐标
	CntsY int       // y方向的格子数量
	grids *sync.Map // map[int]*Grid 当前区域中都有哪些格子，key=格子ID， value=格子对象
}

// NewAOIManager 初始化一个AOI区域
func NewAOIManager(minX, maxX, cntsX, minY, maxY, cntsY int) *AOIManager {
	aoiMgr := &AOIManager{
		MinX:  minX,
		MaxX:  maxX,
		CntsX: cntsX,
		MinY:  minY,
		MaxY:  maxY,
		CntsY: cntsY,
		grids: new(sync.Map),
	}

	// 给AOI初始化区域中所有的格子
	for y := 0; y < cntsY; y++ {
		for x := 0; x < cntsX; x++ {
			// 计算格子ID
			// 格子编号：id = idy *nx + idx  (利用格子坐标得到格子编号)
			gridID := y*cntsX + x

			// 初始化一个格子放在AOI中的map里，key是当前格子的ID
			grid := NewGrid(gridID,
				aoiMgr.MinX+x*aoiMgr.gridWidth(),
				aoiMgr.MinX+(x+1)*aoiMgr.gridWidth(),
				aoiMgr.MinY+y*aoiMgr.gridLength(),
				aoiMgr.MinY+(y+1)*aoiMgr.gridLength())
			aoiMgr.grids.Store(gridID, grid)
		}
	}

	return aoiMgr
}

// gridWidth 得到每个格子在x轴方向的宽度
func (m *AOIManager) gridWidth() int {
	return (m.MaxX - m.MinX) / m.CntsX
}

// gridLength 得到每个格子在x轴方向的长度
func (m *AOIManager) gridLength() int {
	return (m.MaxY - m.MinY) / m.CntsY
}

/*
// String 打印信息方法
func (m *AOIManager) String() string {
	s := fmt.Sprintf("AOIManagr:\nminX:%d, maxX:%d, cntsX:%d, minY:%d, maxY:%d, cntsY:%d\n Grids in AOI Manager:\n",
		m.MinX, m.MaxX, m.CntsX, m.MinY, m.MaxY, m.CntsY)
	for _, grid := range m.grids {
		s += fmt.Sprintln(grid)
	}
	return s
}
*/
