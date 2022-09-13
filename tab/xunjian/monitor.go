package xunjian

import (
	"fmt"
	"github.com/gizak/termui/v3"
	"monitor-ter/ui"
	"time"
)

// DashboardWidget monitor
type DashboardWidget struct {
	*ui.Table
	updateInterval time.Duration
	title          string
}

// NewDashboardWidget monitor
func NewDashboardWidget(item map[string]string, itemData [][]string, termWidth, termHeight int) *DashboardWidget {
	ctx := &DashboardWidget{
		Table:          ui.NewTable(),
		updateInterval: time.Second,
	}
	ctx.title = item["desc"]
	ctx.Title = fmt.Sprintf(" %v ", item["desc"])
	ctx.ShowCursor = true
	ctx.ShowLocation = true
	ctx.ColGap = 3
	ctx.PadLeft = 1
	ctx.ColResizer = func() {
		ctx.ColWidths = []int{
			termWidth / 10, termWidth / 20, termWidth / 20,
		}

	}
	ctx.UniqueCol = 0
	ctx.Rows = itemData
	//ctx.update()

	//go func() {
	//	for range time.NewTicker(1 * time.Second).C {
	//		ctx.Lock()
	//		//ctx.update()
	//		ctx.resize()
	//		ctx.Unlock()
	//	}
	//}()

	return ctx
}

func (ctx *DashboardWidget) resize() {
	termWidth, _ := termui.TerminalDimensions()
	ctx.ColResizer = func() {
		ctx.ColWidths = []int{
			termWidth / 20, termWidth / 10, termWidth / 20,
		}
	}
}
func (ctx *DashboardWidget) update() {
}
