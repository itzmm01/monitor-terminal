package qrcode

import (
	"github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"monitor-ter/tab/report/comm"
	"monitor-ter/ui"
)

var (
	rootView  *ui.TabPane
	qrcodeNew *widgets.List
)

// InitWidget init
func InitWidget() {
	termWidth, termHeight := termui.TerminalDimensions()
	rootView = ui.NewTabPane()
	rootView.SetRect(0, 1, termWidth, termHeight)

	qrcodeNew = widgets.NewList()
	qrcodeNew.Title = "关注微信\"交付中心小助手\",使用扫码播报功能"
	qrcodeNew.Rows = []string{
		comm.GetQrcode(),
	}
	qrcodeNew.TextStyle = termui.NewStyle(termui.ColorYellow)
	qrcodeNew.WrapText = false

	grid := termui.NewGrid()
	grid.Set(
		termui.NewRow(1.0,
			termui.NewCol(1.0, qrcodeNew),
		),
	)
	rootView.Set(grid)
}

// SetupGrid set grid
func SetupGrid(grid *termui.Grid) {
	qrcodeNew.Rows = []string{
		comm.GetQrcode(),
	}
	grid.Set(
		termui.NewRow(1.0,
			termui.NewCol(1.0, qrcodeNew),
		),
	)
	rootView.Set(grid)
}
