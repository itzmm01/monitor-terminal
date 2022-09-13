package comm

import (
	"fmt"
	"io/ioutil"
	"monitor-ter/utils"
)

// SetQrcode set qrcode
func SetQrcode() {
	cmdStr := fmt.Sprintf(
		"python %s/setup-tools/tools/qrcode_tool.py -t '%s' -p > %s/tmp/report.txt",
		utils.TcsTools, utils.GetData(false), utils.Home,
	)
	res, err := utils.CmdRun(cmdStr)
	if err != nil {
		utils.Log.Printf("[ERROR]: " + cmdStr + res)
	}
}

// GetQrcode get arcode
func GetQrcode() string {
	SetQrcode()
	filePath := utils.Home + "/tmp/report.txt"
	File, err := ioutil.ReadFile(filePath)
	if err != nil {
		utils.Log.Printf("二维码文件失败 %v", err)
	}
	return string(File)
}
