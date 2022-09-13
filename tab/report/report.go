package report

import (
	"github.com/gizak/termui/v3"
	qrcode "monitor-ter/tab/report/qcode"
	"monitor-ter/tab/report/qcodeData"
	"monitor-ter/ui"
	"monitor-ter/utils"
)

type TXMethod string

var Log = utils.GetLogger("comm")

var (
	rootView *ui.GridWrap
)

// InitTxWidgets init
func InitTxWidgets() {
	rootView = ui.NewGridWrap()

	qrcodeUI := termui.NewGrid()
	qrcode.InitWidget()
	qrcode.SetupGrid(qrcodeUI)

	qrcodeDataUI := termui.NewGrid()
	qrcodeData.InitWidget()
	qrcodeData.SetupGrid(qrcodeDataUI)

	rootView.Set(
		termui.NewRow(1,
			termui.NewCol(1.0/10*3, qrcodeDataUI),
			termui.NewCol(1.0/10*7, qrcodeUI),
		),
	)
	termui.Render(rootView)
}

// SetGrid 渲染界面 ,会不停的调用
func SetGrid(grid *termui.Grid) {
	grid.Set(termui.NewRow(1, termui.NewCol(1, rootView)))
	RenderRoot(false)
}

// RenderRoot render
func RenderRoot(render bool) {
	rootView.Clear()
	updateUI()
	qrcodeUI := termui.NewGrid()
	qrcode.InitWidget()
	qrcode.SetupGrid(qrcodeUI)

	qrcodeDataUI := termui.NewGrid()
	qrcodeData.InitWidget()
	qrcodeData.SetupGrid(qrcodeDataUI)

	rootView.Set(
		termui.NewRow(1,
			termui.NewCol(1.0/10*2, qrcodeDataUI),
			termui.NewCol(1.0/10*8, qrcodeUI),
		),
	)

	if render {
		termui.Render(rootView)
	}
}

// TxEventLoop 键盘事件监听
func TxEventLoop(id string) bool {
	if utils.Active.Tab != 3 {
		return false
	}
	switch id {
	case "<C-r>":
		Log.Printf("ctrl -r")
		return true
	}
	return false
}

func updateUI() {
	termui.Render(rootView)
}
