package tx

import (
	"github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"monitor-ter/ui"
	"monitor-ter/utils"
	"strings"
)

// ShowDetailWidget aaa
func ShowDetailWidget(grid *ui.GridWrap, data utils.CheckItem) {

	title := widgets.NewParagraph()
	title.Title = "异常详情"
	text := []string{
		"[1] 异常项 : " + data.Name,
		"[2] 异常主机: " + data.Host,
		"[3] 状态: " + data.StatusText,
		"[4] 详细信息: \n" + data.Msg,
		"[5] 说   明: \n" + data.Doc,
	}
	title.Text = strings.Join(text, "\n\n")
	grid.Clear()
	grid.Set(
		termui.NewRow(8.0/8,
			termui.NewCol(4.0/4, title),
		),
	)

}
