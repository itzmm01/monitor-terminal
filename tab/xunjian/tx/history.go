package tx

import (
	"fmt"
	"github.com/gizak/termui/v3"
	"monitor-ter/ui"
	"monitor-ter/utils"
	"path/filepath"
	"strings"
)

// TaskHistory aaa
type TaskHistory struct {
	Name   string
	Status string
	Path   string
	Date   int64
}

// HistoryList aaa
type HistoryList []TaskHistory

// Len aaa
func (list HistoryList) Len() int {
	return len(list)
}

// Less aaa
func (list HistoryList) Less(i, j int) bool {
	return list[i].Date >= list[j].Date
}

// Swap aaa
func (list HistoryList) Swap(i, j int) {
	list[i], list[j] = list[j], list[i]
}

// History aaa
type History struct {
	*ui.Table
	List       []string
	Focus      bool
	ColResizer func(w, h int)
}

// NewHistory aaa
func NewHistory() *History {
	his := &History{
		Table: ui.NewTable(),
		Focus: true,
		ColResizer: func(w, h int) {

		},
	}
	his.ShowCursor = true
	his.ShowLocation = true
	his.ColGap = 1
	his.PadLeft = 0
	his.UniqueCol = 0
	his.CursorColor = termui.Color(ui.Default.ProcCursor)
	return his
}

// SetRect aaa
func (ctx *History) SetRect(x1, y1, x2, y2 int) {
	ctx.Table.SetRect(x1, y1, x2, y2)
	w := x2 - x1
	h := y2 - y1
	ctx.ColResizer(w, h)
}

// Update 加载本地的巡检历史
func (ctx *History) Update(dir string, filename string) string {

	var list []string
	//var listnew []string
	files, err := utils.GetAllFiles(dir)
	for _, file := range files {
		path, fileName := filepath.Split(file)
		resultPath := path + "/" + fileName
		match, _ := filepath.Match("*.csv", fileName)
		if !match {
			continue
		}
		list = append(list, resultPath)
	}
	if err != nil {
		Log.Printf("error %v", err)
		return ""
	}

	utils.SortFile(list)

	ctx.List = list
	if len(list) == 0 {
		list = append(list, "no history")
		ctx.List = list
	}
	rows := make([][]string, len(list))

	for i := range ctx.List {
		_, fileName := filepath.Split(ctx.List[i])
		rows[i] = make([]string, 2)
		rows[i][0] = fmt.Sprintf("[%v]", i+1)
		name := strings.Replace(fileName, "csv-", "", 1)
		rows[i][1] = strings.Replace(name, ".csv", "", 1)
	}
	ctx.Rows = rows
	var index = ctx.GetCurIndex()
	if index > -1 && index < len(ctx.List) {
		path := ctx.List[index]
		return path
	} else if len(ctx.List) > 0 {
		path := ctx.List[0]
		return path
	}
	return ""
}
