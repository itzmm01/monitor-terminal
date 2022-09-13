package tx

import (
	"encoding/json"
	"github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"io/ioutil"
	"monitor-ter/ui"
	"monitor-ter/utils"
)

//播报
var (
	rootView *ui.TabPane
	//history   *tx.History
	qrcodeNew *widgets.List
	DirPath   = utils.TcsTools + "/csv"
)

// InitWidget aaa
func InitWidget() {
	termWidth, termHeight := termui.TerminalDimensions()
	rootView = ui.NewTabPane()
	rootView.SetRect(0, 1, termWidth, termHeight)

	qrcodeNew = widgets.NewList()
	qrcodeNew.Title = "关注微信\"交付中心小助手\",使用扫码排障功能"
	qrcodeNew.TextStyle = termui.NewStyle(termui.ColorYellow)
	qrcodeNew.WrapText = false
}

// SetupGrid aaa
func SetupGrid(grid *termui.Grid) {
	grid.Set(
		termui.NewRow(1.0,
			termui.NewCol(1.0, qrcodeNew),
		),
	)
	rootView.Set(grid)
}

// UpdateCode aaa
func UpdateCode(filename string) {
	file := utils.TcsTools + "/csv/csv-" + filename + ".csv"
	qrcodeNew.Rows = []string{
		getQrcode(file),
	}
}

// SetQrcode aaa
func SetQrcode(resultPath string) {
	data := utils.ParseSetupToolsCsv(resultPath)
	var errCode []string
	for _, line := range data.Items {
		if line.Name == "ok" {
			errCode = append(errCode, "ok")
			continue
		}
		if line.StatusText != "OK" {
			if _, ok := utils.ErrCode[line.Name]; ok {
				errCode = append(errCode, utils.ErrCode[line.Name])
			} else {
				errCode = append(errCode, "UNKNOWN")
			}
		}
	}
	errCode = utils.UniqList(errCode)
	jsonByte, _ := json.Marshal(errCode)
	Log.Printf(string(jsonByte))
	utils.CreateQrCode(string(jsonByte))
}
func getQrcode(file string) string {
	var resultPath string
	if file != "" {
		resultPath = file
	} else {
		resultPath = utils.GetLastFile(DirPath, "csv")
	}
	SetQrcode(resultPath)
	filePath := utils.Home + "/tmp/qrcode.txt"
	File, err := ioutil.ReadFile(filePath)
	if err != nil {
		Log.Printf("二维码文件失败 %v", err)
	}
	return string(File)
}
