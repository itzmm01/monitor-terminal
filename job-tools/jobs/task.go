package jobs

import (
	"fmt"
	utils2 "monitor-ter/utils"
	"os"
	"strconv"
	"strings"
)

// TaskStruct task
type TaskStruct struct {
	Name         string
	Type         string
	Cmd          string
	When         string `json:"item_condition"`
	Expr         string `json:"expr"`
	Loop         string
	Threshold    string
	Register     string
	PrintLog     bool `json:"print_log"`
	Parallel     bool
	IgnoreErrors bool `json:"allow_failed"`
}

// RunTask run task
func (task TaskStruct) RunTask(hostInfo map[string]string, job Job) (string, error) {
	taskInfo := ReplaceTaskParam(task, job, hostInfo["ip"])

	if taskInfo.When != "" && !utils2.WhenExpr(taskInfo.When) {
		Console.SkipFile("Fail to match: " + taskInfo.When)
		Console.SkipFile("    " + hostInfo["ip"])
		return "", nil
	}
	var res string
	var err error

	if hostInfo["ip"] == "127.0.0.1" || hostInfo["ip"] == "local" {
		res, err = utils2.CmdRun(taskInfo.Cmd)
	} else {
		sshPort, _ := strconv.ParseInt(hostInfo["port"], 10, 64)
		sshClient := utils2.ClientConfig{
			Host:     hostInfo["ip"],
			Port:     sshPort,
			Username: hostInfo["user"],
			Password: hostInfo["password"],
		}
		if err1 := sshClient.CreateClient(); err1 != nil {
			return "", utils2.Error(fmt.Sprintf("connect failed: %v-%v", sshClient.Host, err1.Error()))
		}
		res, err = sshClient.RunShell(taskInfo.Cmd)
	}

	Console.InfoFile(hostInfo["ip"] + "Execute command: " + taskInfo.Cmd)
	taskInfo.Expr = strings.Replace(taskInfo.Expr, "result", res, -1)
	if taskInfo.Expr != "" && !utils2.WhenExpr(taskInfo.Expr) {
		err1 := utils2.Error("Fail to match: " + taskInfo.Expr)
		Console.WarnFile(err1.Error())
		Console.ErrorFile("cmd results: " + res)
		if !taskInfo.IgnoreErrors {
			os.Exit(1)
		} else {
			err = err1
		}
	} else {
		if err != nil {
			Console.ErrorFile("cmd results: " + res)
		} else {
			Console.SuccessFile("cmd results: " + res)
		}
	}
	if taskInfo.Register != "" {
		registerName := fmt.Sprintf("%v@%v", taskInfo.Register, hostInfo["ip"])
		GlobalParams[registerName] = res
		GlobalParams[taskInfo.Register] = res
	}

	return res, err
}
