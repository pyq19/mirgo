package main

import (
	"github.com/yenkeia/mirgo/common"
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
				aoiMgr.MinY+y*aoiMgr.gridHeight(),
				aoiMgr.MinY+(y+1)*aoiMgr.gridHeight())
			aoiMgr.grids.Store(gridID, grid)
		}
	}

	return aoiMgr
}

// gridWidth 得到每个格子在x轴方向的宽度
func (m *AOIManager) gridWidth() int {
	return (m.MaxX - m.MinX) / m.CntsX
}

// gridLength 得到每个格子在y轴方向的高度
func (m *AOIManager) gridHeight() int {
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

// GetSurroundGridsByGridID 根据格子的gID得到当前周边的九宫格信息
func (m *AOIManager) GetSurroundGridsByGridID(gID int) (grids []*Grid) {
	// 判断gID是否存在
	v, ok := m.grids.Load(gID)
	if !ok {
		return
	}
	g := v.(*Grid)
	// 将当前gid添加到九宫格中
	grids = append(grids, g)
	// 根据gid得到当前格子所在的X轴编号
	idx := gID % m.CntsX
	// 判断当前idx左边是否还有格子
	if idx > 0 {
		v, _ := m.grids.Load(gID - 1)
		grids = append(grids, v.(*Grid))
	}
	// 判断当前的idx右边是否还有格子
	if idx < m.CntsX-1 {
		v, _ := m.grids.Load(gID + 1)
		grids = append(grids, v.(*Grid))
	}
	// 将x轴当前的格子都取出，进行遍历，再分别得到每个格子的上下是否有格子
	// 得到当前x轴的格子id集合
	gidsX := make([]int, 0, len(grids))
	for _, v := range grids {
		gidsX = append(gidsX, v.GID)
	}
	// 遍历x轴格子
	for _, v := range gidsX {
		// 计算该格子处于第几列
		idy := v / m.CntsX
		// 判断当前的idy上边是否还有格子
		if idy > 0 {
			v, _ := m.grids.Load(v - m.CntsX)
			grids = append(grids, v.(*Grid))
		}
		// 判断当前的idy下边是否还有格子
		if idy < m.CntsY-1 {
			v, _ := m.grids.Load(v + m.CntsX)
			grids = append(grids, v.(*Grid))
		}
	}
	return
}

// GetGridByCoordinate
func (m *AOIManager) GetGridByCoordinate(coordinate string) (grid *Grid) {
	p := common.NewPointByCoordinate(coordinate)
	x := int(p.X)
	y := int(p.Y)
	w := m.gridWidth()
	h := m.gridHeight()
	gridID := (y/h)*m.CntsX + (x / w)
	v, ok := m.grids.Load(gridID)
	if !ok {
		return nil
	}
	grid = v.(*Grid)
	return grid
}

// GetSurroundGridsByCoordinate 根据坐标求出周边九宫格
func (m *AOIManager) GetSurroundGridsByCoordinate(coordinate string) (grids []*Grid) {
	grid := m.GetGridByCoordinate(coordinate)
	return m.GetSurroundGridsByGridID(grid.GID)
}
