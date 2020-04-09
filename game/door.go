package game

import (
	"time"

	"github.com/yenkeia/mirgo/game/cm"
	"github.com/yenkeia/mirgo/game/proto/server"
)

type Door struct {
	Map      *Map
	Index    byte
	State    int //0: closed, 1: opening, 2: open, 3: closing
	LastTick time.Time
	Location cm.Point
}

func (d *Door) IsOpen() bool {
	return d.State == 2
}

func (d *Door) SetOpen(open bool) {
	if open {
		d.State = 2
	} else {
		d.State = 0
	}
	d.LastTick = time.Now()
	d.Map.BroadcastP(d.Location, &server.Opendoor{DoorIndex: d.Index, Close: !open}, nil)
}

func (d *Door) Tick(now time.Time) {
	if d.State == 2 {
		if now.Sub(d.LastTick) > 5*time.Second {
			d.SetOpen(false)
		}
	}
}

type Grid struct {
	W, H uint32
	Grid map[uint32]map[uint32]*Door // 二维map，节省内存。。
}

func NewGrid(w, h uint32) *Grid {
	return &Grid{
		W: w, H: h,
		Grid: map[uint32]map[uint32]*Door{},
	}
}

func (g *Grid) In(loc cm.Point) bool {
	return loc.X < g.W && loc.Y < g.H
}

func (g *Grid) Set(loc cm.Point, d *Door) {
	if g.In(loc) {
		if _, ok := g.Grid[loc.X]; !ok {
			g.Grid[loc.X] = map[uint32]*Door{}
		}

		g.Grid[loc.X][loc.Y] = d
	}
}

func (g *Grid) Get(loc cm.Point) *Door {
	if _, ok := g.Grid[loc.X]; !ok {
		return nil
	}
	if v, ok := g.Grid[loc.X][loc.Y]; !ok {
		return nil
	} else {
		return v
	}
}
