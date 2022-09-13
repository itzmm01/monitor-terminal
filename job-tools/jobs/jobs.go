package jobs

import (
	utils2 "monitor-ter/utils"
	"strings"
	"time"
)

var (
	GlobalParams = map[string]string{}
)

// Job job
type Job struct {
	Name     string
	Host     string
	When     string `json:"condition"`
	Parallel bool   `json:"node_parallel"`
	PrintLog bool   `json:"print_log"`
	Tasks    []TaskStruct
	Params   map[string]string
}

func checkHostReady(hostList []map[string]string, init bool) []map[string]string {
	resList := CommPool(hostList, "date > /dev/null")
	var readyHost []map[string]string
	var existList []string
	for _, res := range resList {
		if res[2] == "0" {
			host := utils2.GetHostInfo(hostList, res[1])
			if host != nil && !utils2.CheckListExists(host["ip"], existList) {
				existList = append(existList, host["ip"])
				readyHost = append(readyHost, host)
			}
		} else {
			Console.ErrorFile(strings.Join(res, " "))
		}
	}
	return readyHost

}

// ReplaceParam replace param
func (job Job) ReplaceParam() {
	for _, params := range []map[string]string{job.Params, GlobalParams} {
		job.Name = utils2.ReplaceParam(job.Name, params)
		job.Host = utils2.ReplaceParam(job.Host, params)
		job.When = utils2.ReplaceParam(job.When, params)
	}

}

// RunJob run job
func (job Job) RunJob() map[string][][]string {
	jobRes := map[string][][]string{}
	job.ReplaceParam()
	Console.InfoFile("----------------Job execution starts: " + job.Name)
	if _, ok := GroupHost[job.Host]; !ok {
		Console.ErrorFile("Host mismatch: " + job.Host)
		return jobRes
	} else {
		if len(GroupHost[job.Host]) == 0 {
			Console.ErrorFile("Host is empty: " + job.Host)
			return jobRes
		}

		var readyHost []map[string]string
		if job.Host == "local" || job.Host == "127.0.0.1" {
			readyHost = []map[string]string{{"ip": "127.0.0.1"}}
		} else {
			readyHost = checkHostReady(GroupHost[job.Host], false)
		}
		if len(readyHost) == 0 {
			Console.ErrorFile("ready host is empty")

		} else {
			for _, task := range job.Tasks {
				Console.InfoFile("----Task: " + task.Name)
				if job.PrintLog {
					task.PrintLog = true
				}
				timeObj := time.Now()
				NowTime := timeObj.Format("2006-01-02 15:04:05")
				if job.Parallel || task.Parallel {
					taskResult := Pool(readyHost, task, job)
					jobRes[task.Name] = taskResult
				} else {
					for _, host := range readyHost {
						runTaskRes, err := task.RunTask(host, job)
						if err != nil {
							jobRes[task.Name] = append(jobRes[task.Name], []string{
								task.Name, host["ip"], "1", runTaskRes, task.Threshold, NowTime, err.Error(),
							})
						} else {
							jobRes[task.Name] = append(jobRes[task.Name], []string{
								task.Name, host["ip"], "0", runTaskRes, task.Threshold, NowTime,
							})
						}
					}
				}
			}
			//}
		}
		Console.InfoFile("----------------Job execution completed: " + job.Name)
	}
	return jobRes
}
