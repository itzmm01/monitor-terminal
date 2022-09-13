package jobTools

import (
	"encoding/json"
	"flag"
	"fmt"
	jobs2 "monitor-ter/job-tools/jobs"
	utils2 "monitor-ter/utils"
)

type arrayFlags []string

// Value ...
func (i *arrayFlags) String() string {
	return fmt.Sprint(*i)
}

// Set 方法是flag.Value接口, 设置flag Value的方法.通过多个flag指定的值， 所以我们追加到最终的数组上.
func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

var paramsList arrayFlags

func args() map[string]*string {
	argsMap := map[string]*string{}
	argsMap["hostYml"] = flag.String("i", "./host.yml", "host.yml")
	argsMap["jobYml"] = flag.String("c", "./job.yml", "job.yml")
	argsMap["paramYml"] = flag.String("p", "", "params.yml")
	flag.Var(&paramsList, "g", "params")
	flag.Parse()
	return argsMap
}

// GoJob run job
func GoJob(host, job string) map[string]map[string][][]string {
	//argsMap := args()
	AllJobRes := map[string]map[string][][]string{}
	utils2.ParseParams(paramsList, jobs2.GlobalParams)
	//jobs2.ReadSrcHost(*argsMap["hostYml"])
	//jobs2.ReadJobYml(*argsMap["jobYml"], map[string]string{})
	utils2.ReadSrcHost(host)
	jobs2.ReadJobYml(job, map[string]string{}, AllJobRes)
	return AllJobRes
}

// GetMonitorData get monitor data
func GetMonitorData(init bool, monitorName string) map[string]map[string][][]string {
	var baseRes map[string]map[string][][]string
	if init {
		baseRes = utils2.InitData
		utils2.InitMonitorData[monitorName] = "ok"
	} else {
		baseRes = GoJob(utils2.HostYml, fmt.Sprintf("%v/%v.yml", utils2.ParamData["Conf"], monitorName))
	}
	result2 := map[string][][]string{}
	for _, v := range baseRes {
		for name, value := range v {
			params := map[string]string{}
			params["desc"] = name
			var dataList [][]string
			for _, hostInfo := range value {
				dataList = append(dataList, []string{hostInfo[1], hostInfo[3], hostInfo[4]})
			}
			utils2.SortDataList(dataList)
			newDataList := [][]string{{"名称", "指标", "阈值"}}
			for index1 := range dataList {
				result2[name] = append(result2[name], dataList[index1])
				if index1 >= 4 {
					break
				}
				newDataList = append(newDataList, dataList[index1])
			}
		}
	}
	jsonData, _ := json.Marshal(result2)
	utils2.WriteFile(string(jsonData), fmt.Sprintf("%v/tmp/%v.json", utils2.Home, monitorName))

	return baseRes
}
