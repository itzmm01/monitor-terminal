package business

import (
	"github.com/gizak/termui/v3"
	"monitor-ter/tab/comm"
	"monitor-ter/tab/xunjian"
	"monitor-ter/utils"
)

var (
	businessForm []*xunjian.DashboardWidget
)

// InitWidgets init
func InitWidgets(init bool) {
	businessForm = comm.InitWidgets("businessItem", init)
}

// SetGrid init grid
func SetGrid(grid *termui.Grid) {
	allRows := comm.SetGrid(businessForm)
	grid.Set(allRows...)
}

// TxEventLoop 键盘事件监听
func TxEventLoop(id string) bool {
	if utils.Active.Tab != 1 {
		return false
	}
	switch id {
	case "<Up>":
		businessForm[utils.Active.Index].ScrollUp()
		return true
	case "<Down>":
		businessForm[utils.Active.Index].ScrollDown()
		return true
	case "<Right>":
		if utils.Active.Index < len(businessForm)-1 {
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
