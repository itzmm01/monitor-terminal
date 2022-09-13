package baseMonitor

import (
	"github.com/gizak/termui/v3"
	"monitor-ter/tab/comm"
	"monitor-ter/tab/xunjian"
	"monitor-ter/ui"
	"monitor-ter/utils"
)

var (
	baseDashboard []*xunjian.DashboardWidget
	color         = ui.Default
	Log           = utils.GetLogger("base monitor")
)

// InitWidgets init
func InitWidgets(init bool) {
	baseDashboard = comm.InitWidgets("baseItem", init)
}

// SetGrid init grid
func SetGrid(grid *termui.Grid) {
	allRows := comm.SetGrid(baseDashboard)
	grid.Set(allRows...)
}

// TxEventLoop 键盘事件监听
func TxEventLoop(id string) bool {
	if utils.Active.Tab != 0 {
		return false
	}
	if len(baseDashboard) == 0 {
		return false
	}
	Log.Printf("%v-%v", baseDashboard[utils.Active.Index].Title, baseDashboard[utils.Active.Index].GetCurIndex())
	switch id {
	case "<Up>":
		baseDashboard[utils.Active.Index].ScrollUp()
		return true
	case "<Down>":
		baseDashboard[utils.Active.Index].ScrollDown()
		return true
	case "<Right>":
		if utils.Active.Index < len(baseDashboard)-1 {
			utils.Active.Index++
		}
		return true

	case "<Left>":
		if utils.Active.Index > 0 {
			utils.Active.Index--
		}
		return true
	}
	return false
}
