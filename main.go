package main

import (
	"encoding/json"
	"flag"
	"fmt"
	_ "image/png"
	jobTools "monitor-ter/job-tools"
	"monitor-ter/tab"
	"monitor-ter/utils"
	"os"
	"runtime"
	"strings"
	"time"
)

import (
	_ "image/gif"
	_ "image/jpeg"
)

func argsParse() map[string]*string {
	argsMap := map[string]*string{}
	argsMap["mode"] = flag.String("m", "default", "default/get_data/add_cron/del_cron")
	argsMap["cron"] = flag.String("c",
		"* */1 * * * sh tcs-tools.sh check -mod all", "cron expr",
	)
	flag.Parse()
	return argsMap
}

//写入文件
func writeBytesToFile(filepath string, content []byte) {
	//打开文件，没有此文件则创建文件，将写入的内容append进去
	w1, err := os.OpenFile(filepath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return
	}
	w1.Write(content)
	w1.Close()

}
func getRes(name string) {
	timeObj := time.Now()
	NowTime := timeObj.Format("2006-01-02")
	jobRes := jobTools.GoJob(utils.HostYml, utils.ParamData["Conf"]+"/"+name+".yml")
	resBytes, _ := json.Marshal(jobRes)
	resBytes = append(resBytes, []byte("\n")...)
	savePath := fmt.Sprintf("%v/data/%v-%v.json", utils.Home, name, NowTime)
	writeBytesToFile(savePath, resBytes)
}
func getMonitorData() {
	getRes("businessItem")
	getRes("monitorItem")
}

func main() {
	args := argsParse()
	Log := utils.Log
	if *args["mode"] == "default" {
		Log.Printf("Starting %v", tab.AppName)
		runtime.GOMAXPROCS(runtime.NumCPU())
		os.Setenv("LANG", "en_US.UTF-8")
		os.Mkdir(utils.Home+"/tmp", 0755)
		os.Mkdir(utils.TcsTools+"/csv", 0755)
		//utils.CmdRun("sh" + utils.TcsTools + "/tcs-tools.sh init")
		utils.ReadERRCode()
		tab.Main()
	} else if *args["mode"] == "get_data" {
		getMonitorData()
	} else if *args["mode"] == "add_cron" {
		var cronStr string
		cronStr = strings.Replace(*args["cron"], "sh tcs-tools.sh",
			fmt.Sprintf("sh %v/tcs-tools.sh", utils.TcsTools), -1,
		)
		cronStr = cronStr + fmt.Sprintf("; %v/tianxun-lite -m report", utils.Home)
		utils.CrontabAdmin("add", cronStr)

	} else if *args["mode"] == "del_cron" {
		utils.CrontabAdmin("del", "")
	} else if *args["mode"] == "report" {
		jobTools.GetMonitorData(false, "baseItem")
		jobTools.GetMonitorData(false, "businessItem")
		reportData := utils.GetData(false)
		res, err := utils.PostData(reportData, "https://assist.tonyandmoney.cn/report/auto?user=test&api_secret=Ks6KqY61SE")
		if err != nil {
			fmt.Println(reportData)
			fmt.Println(res, err.Error())
			Log.Println(reportData)
			Log.Println(res, err.Error())
		}
		fmt.Println(res)
		Log.Println(res)
	} else {
		fmt.Println("no support: " + *args["mode"])
	}

}
