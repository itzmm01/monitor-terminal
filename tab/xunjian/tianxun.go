package xunjian

import (
	"fmt"
	"github.com/gizak/termui/v3/widgets"
	"io/ioutil"
	tx2 "monitor-ter/tab/xunjian/tx"
	"strings"

	"monitor-ter/utils"
	"strconv"
	"time"

	"github.com/gizak/termui/v3"
	"monitor-ter/ui"
)

type TXMethod string

var Log = utils.GetLogger("comm")
var TxDir = utils.TcsTools + "/csv"

// TXWidget tianxun
type TXWidget struct {
	*ui.Table
	updateInterval time.Duration
	Result         utils.TaskResult
	Focus          bool
	ColResizer     func(w, h int)
}

// Task tianxun
type Task struct {
	TaskName string
	Host     string
	Status   string
	Logs     string
	Desc     string
}

var (
	color         = ui.Default
	InspectResult *TXWidget
	History       *tx2.History
	XjItem        *tx2.History
	xjView        *tx2.ExecWidget
	rootView      *ui.GridWrap
	checkItem     *utils.CheckItem
	progressFrame *widgets.Paragraph
	successFrame  *widgets.Paragraph
	warnFrame     *widgets.Paragraph
	errorFrame    *widgets.Paragraph

	InspectResultIndex int
)

// InitTxWidgets tianxun
func InitTxWidgets() {
	rootView = ui.NewGridWrap()

	InspectResult = &TXWidget{
		Table:          ui.NewTable(),
		updateInterval: time.Second,
		Focus:          false,
	}
	InspectResult.Title = " 异常结果(↑/↓ 上/下 enter 选中 ) "
	InspectResult.ShowCursor = true
	InspectResult.ShowLocation = true
	InspectResult.ColGap = 1
	InspectResult.PadLeft = 0
	InspectResult.Header = []string{"No.", "检查项", "主机", "状态"}

	InspectResult.ColResizer = func(w, h int) {
		InspectResult.ColWidths = []int{
			(w / 30) * 1, w / 30 * 17, w / 30 * 8, w / 30 * 3,
		}
	}

	History = tx2.NewHistory()
	History.Header = []string{"No.", "名称"}
	History.Title = "排障历史(enter 选中)"
	History.ShowCursor = true
	History.ShowLocation = true
	History.ColGap = 1
	History.PadLeft = 0
	History.UniqueCol = 0
	History.ColResizer = func(w, h int) {
		History.ColWidths = []int{
			w / 6, w * 4 / 6,
		}
	}

	InspectResult.CursorColor = termui.Color(color.ProcCursor)
	path := History.Update(TxDir, "csv")
	//path := utils.GetLastFile(TxDir, "csv")
	InspectResult.update(path)

	xjView = tx2.NewExecWidget(OnRefresh)
	XjItem = tx2.NewHistory()
	XjItem.Header = []string{"巡检项", "标识"}
	XjItem.Title = "选择巡检项"
	XjItem.ShowCursor = true
	XjItem.ShowLocation = true
	XjItem.ColGap = 1
	XjItem.PadLeft = 0
	XjItem.UniqueCol = 0
	XjItem.ColResizer = func(w, h int) {
		XjItem.ColWidths = []int{
			w / 6, w * 4 / 6,
		}
	}
	res, err := utils.CmdRun(utils.TcsTools + "/tasks/command/list.sh")
	if err != nil {
		XjItem.Rows = [][]string{
			{"1", "all"},
			{"2", "iaas"},
			{"3", "kubelet"},
			{"4", "docker"},
		}
		Log.Println(res, err)
	} else {
		lines := strings.Split(res, "\n")
		XjItem.Rows = [][]string{}
		for index, line := range lines {
			if line != "" {
				var item = []string{strconv.Itoa(index), line}
				XjItem.Rows = append(XjItem.Rows, item)
			}
		}
	}
	xjView = tx2.NewExecWidget(OnRefresh)

}

// OnRefresh tianxun
func OnRefresh() {
	path := History.Update(TxDir, "csv")
	//path := utils.GetLastFile(TxDir, "csv")
	InspectResult.update(path)
	checkItem = nil
	result := InspectResult.Result
	var progress string
	File, err := ioutil.ReadFile(utils.Temp + "/progress.txt")
	if err != nil {
		progress = "100%"
	} else {
		progress = string(File) + "%"
	}
	progressFrame = InspectCount("progress", progress)
	successFrame = InspectCount("success", strconv.Itoa(result.Success))
	warnFrame = InspectCount("warn", strconv.Itoa(result.Warn))
	errorFrame = InspectCount("error", strconv.Itoa(result.Error))

}
func updateXJData() {
	OnRefresh()
	rootView.Set(
		termui.NewRow(2.0/20,
			termui.NewCol(1.0/8*2, xjView),
			termui.NewCol(1.0/8, progressFrame),
			termui.NewCol(1.0/8, successFrame),
			termui.NewCol(1.0/8, warnFrame),
			termui.NewCol(1.0/8, errorFrame),
		),
	)
	termui.Render(rootView)
}

// SetGrid 渲染界面 ,会不停的调用
func SetGrid(grid *termui.Grid) {
	grid.Set(termui.NewRow(1, termui.NewCol(1, rootView)))
	RenderRoot(false, "")
}

// RenderRoot tianxun
func RenderRoot(render bool, flag string) {
	rootView.Clear()
	if flag == "XjItem" {
		rootView.Set(
			termui.NewRow(1,
				termui.NewCol(1.0, XjItem),
			),
		)
	} else if flag == "LogItem" {

	} else {
		if checkItem != nil {
			tx2.ShowDetailWidget(rootView, *checkItem)
		} else {
			OnRefresh()
			grid1 := termui.NewGrid()
			tx2.InitWidget()
			tx2.SetupGrid(grid1)
			fileList := History.Rows[History.GetCurIndex()]
			tx2.UpdateCode(fileList[1])

			InspectResult.SelectedRow = InspectResultIndex
			rootView.Set(
				termui.NewRow(2.0/20,
					termui.NewCol(1.0/8*2, xjView),
					termui.NewCol(1.0/8, progressFrame),
					termui.NewCol(1.0/8, successFrame),
					termui.NewCol(1.0/8, warnFrame),
					termui.NewCol(1.0/8, errorFrame),
				),
				termui.NewRow(18.0/20,
					termui.NewCol(1.0/20*4, History),
					termui.NewCol(1.0/20*6, InspectResult),
					termui.NewCol(1.0/20*10, grid1),
				),
			)
		}
	}
	if render {
		termui.Render(rootView)
	}
}

// TestUpdate tianxun
func TestUpdate() {
	if activeForm == "History" || activeForm == "InspectResult" {
		if tx2.GetStatus() == "RUNNING" {
			updateXJData()
		} else {
			if utils.UpdateCode != "no" {
				RenderRoot(true, "")
				utils.UpdateCode = "no"
			}
		}
	}
}

func keyUP() {
	if activeForm == "History" || activeForm == "" {
		History.ScrollUp()
		termui.Render(History)
	}
	if activeForm == "XjItem" {
		XjItem.ScrollUp()
		termui.Render(XjItem)
	}
	if activeForm == "InspectResult" {
		InspectResult.ScrollUp()
		termui.Render(InspectResult)
	}
}
func keyDown() {
	if activeForm == "XjItem" {
		XjItem.ScrollDown()
		termui.Render(XjItem)
	}
	if activeForm == "History" || activeForm == "" {
		History.ScrollDown()
		termui.Render(History)
	}
	if activeForm == "InspectResult" {
		InspectResult.ScrollDown()
		termui.Render(InspectResult)
	}
}
func keyEnter() {
	if activeForm == "History" {
		RenderRoot(true, "")
		if InspectResult.Rows[0][1] == "nodata" || InspectResult.Rows[0][1] == "ok" {
			return
		}
		activeForm = "InspectResult"
	} else if activeForm == "InspectResult" {
		InspectResultIndex = InspectResult.GetCurIndex()
		items := InspectResult.Result.Items
		if InspectResultIndex > -1 && InspectResultIndex < len(items) {
			checkItem = &items[InspectResultIndex]
		}
		RenderRoot(true, "")
		activeForm = "Result"
	} else if activeForm == "XjItem" {
		index := XjItem.GetCurIndex()
		items := XjItem.Rows
		var item []string
		if index > -1 && index < len(items) {
			item = items[index]
		}
		RenderRoot(true, "")
		tx2.Start(item[1])

		activeForm = "InspectResult"
	}
}
func keyEsc() {
	if activeForm == "XjItem" {
		activeForm = "InspectResult"
		RenderRoot(true, "")
	} else if activeForm == "Result" {
		activeForm = "InspectResult"
		checkItem = nil
		RenderRoot(true, "")
	} else {
		activeForm = "History"
		checkItem = nil
		RenderRoot(true, "")
	}
}

var activeForm = "History"

// TxEventLoop 键盘事件监听
func TxEventLoop(id string) bool {
	if utils.Active.Tab != 2 {
		return false
	}

	switch id {
	case "<C-r>":
		activeForm = "XjItem"
		RenderRoot(true, "XjItem")
	case "<C-e>":
		//activeForm = "LogItem"
		//RenderRoot(true, "LogItem")
	case "<C-q>":
		tx2.Stop()
	case "<Up>":
		keyUP()
		return true
	case "<Down>":
		keyDown()
		return true
	case "<Enter>":
		keyEnter()
		return true
	case "<Escape>":
		keyEsc()
		return true
	}
	return false
}

// InspectCount tianxun
func InspectCount(status, count string) *widgets.Paragraph {
	colorDict := map[string]termui.Color{
		"xj":       termui.ColorGreen,
		"progress": termui.ColorRed,
		"success":  termui.ColorGreen,
		"error":    termui.ColorRed,
		"warn":     termui.ColorYellow,
	}
	zhCNDict := map[string]string{
		"success":  "成功",
		"progress": "进度",
		"error":    "失败",
		"warn":     "警告",
		"xj":       "按Ctrl+r开始排障",
	}
	textBox := widgets.NewParagraph()
	textBox.Title = zhCNDict[status]
	textBox.TitleStyle.Modifier = termui.ModifierBold

	textBox.Text = count
	textBox.TextStyle.Modifier = termui.ModifierBold
	//textBox.TextStyle.Modifier = ui.ModifierUnderline
	textBox.BorderRight = true

	if _, ok := colorDict[status]; ok {
		textBox.TextStyle.Fg = colorDict[status]
		textBox.TitleStyle.Fg = colorDict[status]
	}
	textBox.Border = true
	return textBox
}

func (ctx *TXWidget) update(path string) {
	res := utils.ParseSetupToolsCsv(path)
	ctx.Result = res
	items := res.Items
	stringsList := make([][]string, len(items))
	for i, item := range items {
		if item.StatusText == "OK" {
			continue
		}
		stringsList[i] = make([]string, 4)
		stringsList[i][0] = fmt.Sprintf("%v", i+1)
		stringsList[i][1] = item.Name
		stringsList[i][2] = item.Host
		stringsList[i][3] = item.StatusText
	}
	ctx.Rows = stringsList
	ctx.ScrollTop()
}

// SetRect tianxun
func (ctx *TXWidget) SetRect(x1, y1, x2, y2 int) {
	ctx.Table.SetRect(x1, y1, x2, y2)
	w := x2 - x1
	h := y2 - y1
	ctx.ColResizer(w, h)
}
