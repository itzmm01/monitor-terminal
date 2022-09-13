package qrcodeData

import (
	"github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"monitor-ter/ui"
	"monitor-ter/utils"
	"strings"
)

//播报
var (
	rootView   *ui.TabPane
	qrcodeData *widgets.Paragraph
)

// InitWidget init
func InitWidget() {
	//
}

// SetupGrid init grid
func SetupGrid(grid *termui.Grid) {
	title := widgets.NewParagraph()
	title.Title = "上报源数据"
	text := []string{
		utils.GetData(true),
	}
	title.Text = strings.Join(text, "")
	grid.Set(
		termui.NewRow(8.0/8,
			termui.NewCol(4.0/4, title),
		),
	)
}
