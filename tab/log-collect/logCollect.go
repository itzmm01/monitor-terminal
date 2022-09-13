package logCollect

import (
	"github.com/gizak/termui/v3"
	ui "monitor-ter/ui"
	"monitor-ter/utils"
	"time"
)

var (
	color   = ui.Default
	LogItem *TXWidget
)

// TXWidget aaa
type TXWidget struct {
	*ui.Table
	updateInterval time.Duration
	Result         utils.TaskResult
	Focus          bool
	ColResizer     func(w, h int)
}

// InitTxWidgets aaa
func InitTxWidgets(grid *ui.GridWrap) {
	LogItem = &TXWidget{
		Table:          ui.NewTable(),
		updateInterval: time.Second,
		Focus:          false,
	}
	LogItem.Title = " ↑/↓ 上/下选择 "
	LogItem.ShowCursor = true
	LogItem.ShowLocation = true
	LogItem.ColGap = 1
	LogItem.PadLeft = 0
	LogItem.Header = []string{"No.", "日志模块名"}

	LogItem.ColResizer = func(w, h int) {
		LogItem.ColWidths = []int{
			(w / 30) * 1, w / 30 * 29,
		}
	}
	LogItem.CursorColor = termui.Color(color.ProcCursor)
	LogItem.Rows = [][]string{
		[]string{"1", "log1"},
		[]string{"2", "log2"},
		[]string{"3", "log3"},
	}
}

// TxEventLoop aaa
func TxEventLoop(id string) bool {
	if utils.Active.Tab != 3 {
		return false
	}
	switch id {
	case "Enter":

	}
	return true
}
