package comm

import (
	"github.com/gizak/termui/v3"
	jobTools "monitor-ter/job-tools"
	"monitor-ter/tab/xunjian"
	"monitor-ter/ui"
	"monitor-ter/utils"
)

var (
	color = ui.Default
)

// SetGrid init gird comm
func SetGrid(form []*xunjian.DashboardWidget) []interface{} {
	var allCol1 []interface{}
	var allCol2 []interface{}
	var allCol3 []interface{}
	for index, v := range form {
		if index <= 3 {
			allCol1 = append(allCol1, termui.NewCol(1.0/4, v))
		} else if index >= 4 && index <= 7 {
			allCol2 = append(allCol2, termui.NewCol(1.0/4, v))
		} else if index >= 8 && index <= 11 {
			allCol3 = append(allCol3, termui.NewCol(1.0/4, v))
		}
		
	}
	var allRows []interface{}
	if len(allCol1) > 0 {
		allRows = append(allRows, termui.NewRow(1.0/3, allCol1...))
	}
	if len(allCol2) > 0 {
		allRows = append(allRows, termui.NewRow(1.0/3, allCol2...))
	}
	if len(allCol3) > 0 {
		allRows = append(allRows, termui.NewRow(1.0/3, allCol3...))
	}
	return allRows
}

// InitWidgets init widget comm
func InitWidgets(formName string, init bool) []*xunjian.DashboardWidget {
	var formList []*xunjian.DashboardWidget
	termWidth, termHeight := termui.TerminalDimensions()
	res := jobTools.GetMonitorData(init, formName)

	for _, v := range res {
		for name, value := range v {
			var form *xunjian.DashboardWidget
			params := map[string]string{}
			params["desc"] = name
			var dataList [][]string
			for _, hostInfo := range value {
				dataList = append(dataList, []string{hostInfo[5], hostInfo[3], hostInfo[4]})
			}
			utils.SortDataList(dataList)
			newDataList := [][]string{{"名称", "指标", "阈值"}}
			for index1 := range dataList {
				if index1 >= 4 {
					break
				}
				newDataList = append(newDataList, dataList[index1])
			}
			form = xunjian.NewDashboardWidget(params, newDataList, termWidth, termHeight)
			form.CursorColor = termui.Color(color.ProcCursor)
			form.SelectedRow = 0
			formList = append(formList, form)
		}
	}
	return formList
}
