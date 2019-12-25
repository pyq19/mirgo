package main

import (
	"github.com/yenkeia/mirgo/common"
	"sync"
)

// Zinx应用-MMO游戏案例-(2)AOI兴趣点算法 https://www.jianshu.com/p/e5b5db9fa6fe

// AOIManager 每一张地图对应一个 AOI 管理模块
type AOIManager struct {
	Map   *Map      // 这个 AOI 属于哪张地图
	MinX  int       // 区域左边界坐标
	MaxX  int       // 区域右边界坐标
	CntsX int       // x方向区域的数量
	MinY  int       // 区域上边界坐标
	MaxY  int       // 区域下边界坐标
	CntsY int       // y方向的区域数量
	grids *sync.Map // map[int]*Grid 当前区域中都有哪些区域，key=区域ID， value=区域对象
}

// NewAOIManager 初始化一个AOI区域
func NewAOIManager(m *Map, minX, maxX, cntsX, minY, maxY, cntsY int) *AOIManager {
	aoi := &AOIManager{
		Map:   m,
		MinX:  minX,
		MaxX:  maxX,
		CntsX: cntsX,
		MinY:  minY,
		MaxY:  maxY,
		CntsY: cntsY,
		grids: new(sync.Map),
	}

	// 给AOI初始化区域中所有的区域
	for y := 0; y < cntsY; y++ {
		for x := 0; x < cntsX; x++ {
			// 计算区域ID
			// 区域编号：id = idy *nx + idx  (利用区域坐标得到区域编号)
			gridID := y*cntsX + x

			// 初始化一个区域放在AOI中的map里，key是当前区域的ID
			grid := NewGrid(aoi, gridID,
				aoi.MinX+x*aoi.gridWidth(),
				aoi.MinX+(x+1)*aoi.gridWidth(),
				aoi.MinY+y*aoi.gridHeight(),
				aoi.MinY+(y+1)*aoi.gridHeight())
			aoi.grids.Store(gridID, grid)
		}
	}

	return aoi
}

// gridWidth 得到每个区域在x轴方向的宽度
func (m *AOIManager) gridWidth() int {
	return (m.MaxX - m.MinX) / m.CntsX
}

// gridLength 得到每个区域在y轴方向的高度
func (m *AOIManager) gridHeight() int {
	return (m.MaxY - m.MinY) / m.CntsY
}

// GetSurroundGridsByGridID 根据区域的gID得到当前周边的九宫格信息
func (m *AOIManager) GetSurroundGridsByGridID(gID int) (grids []*Grid) {
	// 判断gID是否存在
	v, ok := m.grids.Load(gID)
	if !ok {
		return
	}
	g := v.(*Grid)
	// 将当前gid添加到九宫格中
	grids = append(grids, g)
	// 根据gid得到当前区域所在的X轴编号
	idx := gID % m.CntsX
	// 判断当前idx左边是否还有区域
	if idx > 0 {
		v, _ := m.grids.Load(gID - 1)
		grids = append(grids, v.(*Grid))
	}
	// 判断当前的idx右边是否还有区域
	if idx < m.CntsX-1 {
		v, _ := m.grids.Load(gID + 1)
		grids = append(grids, v.(*Grid))
	}
	// 将x轴当前的区域都取出，进行遍历，再分别得到每个区域的上下是否有区域
	// 得到当前x轴的区域id集合
	gidsX := make([]int, 0, len(grids))
	for _, v := range grids {
		gidsX = append(gidsX, v.GID)
	}
	// 遍历x轴区域
	for _, v := range gidsX {
		// 计算该区域处于第几列
		idy := v / m.CntsX
		// 判断当前的idy上边是否还有区域
		if idy > 0 {
			v, _ := m.grids.Load(v - m.CntsX)
			grids = append(grids, v.(*Grid))
		}
		// 判断当前的idy下边是否还有区域
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
