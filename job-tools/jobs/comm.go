package jobs

import (
	utils "monitor-ter/job-tools/util"
	utils2 "monitor-ter/utils"
	"strings"
)

var (
	Console   = utils.Console
	GroupHost = utils2.ReadSrcHost(utils2.HostYml)
)

// ReplaceTaskParam replace param
func ReplaceTaskParam(task TaskStruct, job Job, host string) TaskStruct {
	task.Name = strings.Replace(task.Name, `${IP}`, host, -1)
	task.Cmd = strings.Replace(task.Cmd, `${IP}`, host, -1)
	task.When = strings.Replace(task.When, `${IP}`, host, -1)
	task.Expr = strings.Replace(task.Expr, `${IP}`, host, -1)
	for _, params := range []map[string]string{job.Params, GlobalParams} {
		task.Name = utils2.ReplaceParam(task.Name, params)
		task.Cmd = utils2.ReplaceParam(task.Cmd, params)
		task.When = utils2.ReplaceParam(task.When, params)
		task.Expr = utils2.ReplaceParam(task.Expr, params)
	}
	return task
}
