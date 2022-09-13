package jobs

import (
	utils "monitor-ter/utils"
	"strconv"
	"sync"
	"time"
)

func workFunc(
	host map[string]string, taskStruct TaskStruct, job Job,
	response chan []string, limiter chan bool, wg *sync.WaitGroup,
) {
	defer wg.Done()
	res, err := taskStruct.RunTask(host, job)
	timeObj := time.Now()
	NowTime := timeObj.Format("2006-01-02 15:04:05")
	if taskStruct.IgnoreErrors {
		if err != nil {
			response <- []string{taskStruct.Name, host["ip"], "1", res, taskStruct.Threshold, NowTime, err.Error()}
		} else {
			response <- []string{taskStruct.Name, host["ip"], "0", res, taskStruct.Threshold, NowTime}
		}

	} else {
		Console.ErrorFile("Error that cannot be skipped: " + err.Error())
	}
	<-limiter
}

// Pool poll
func Pool(hostList []map[string]string, taskStruct TaskStruct, job Job) [][]string {

	wg := &sync.WaitGroup{}
	// 控制并发数为cpu核数
	limiter := make(chan bool, 2)
	defer close(limiter)

	// 函数内的局部变量channel, 专门用来接收函数内所有goroutine的结果
	responseChannel := make(chan []string, 6)
	// 为读取结果控制器创建新的WaitGroup, 需要保证控制器内的所有值都已经正确处理完毕, 才能结束
	wgResponse := &sync.WaitGroup{}

	var result [][]string
	go func() {
		// wgResponse计数器+1
		wgResponse.Add(1)

		for response := range responseChannel {
			result = append(result, response)
		}
		// 当 responseChannel被关闭时且channel中所有的值都已经被处理完毕后, 将执行到这一行
		wgResponse.Done()
	}()

	for _, host := range hostList {
		// 计数器+1
		wg.Add(1)
		limiter <- true
		go workFunc(host, taskStruct, job, responseChannel, limiter, wg)
	}

	// 等待所以协程执行完毕
	wg.Wait()
	// 关闭接收结果channel
	close(responseChannel)
	// 等待wgResponse的计数器归零
	wgResponse.Wait()
	// 返回聚合后结果
	return result
}

// CommPool comm pool
func CommPool(hostList []map[string]string, args string) [][]string {

	wg := &sync.WaitGroup{}
	// 控制并发数为cpu核数
	limiter := make(chan bool, 5)
	defer close(limiter)

	// 函数内的局部变量channel, 专门用来接收函数内所有goroutine的结果
	responseChannel := make(chan []string, 5)
	// 为读取结果控制器创建新的WaitGroup, 需要保证控制器内的所有值都已经正确处理完毕, 才能结束
	wgResponse := &sync.WaitGroup{}

	var result [][]string
	go func() {
		// wgResponse计数器+1
		wgResponse.Add(1)

		for response := range responseChannel {
			result = append(result, response)
		}
		// 当 responseChannel被关闭时且channel中所有的值都已经被处理完毕后, 将执行到这一行
		wgResponse.Done()
	}()

	for _, host := range hostList {
		// 计数器+1
		wg.Add(1)
		limiter <- true
		go commWorkFunc(host, responseChannel, limiter, wg, args)
	}

	// 等待所以协程执行完毕
	wg.Wait()
	// 关闭接收结果channel
	close(responseChannel)
	// 等待wgResponse的计数器归零
	wgResponse.Wait()
	// 返回聚合后结果
	return result
}
func commWorkFunc(host map[string]string, response chan []string, limiter chan bool, wg *sync.WaitGroup, args string) {
	defer wg.Done()
	timeObj := time.Now()
	NowTime := timeObj.Format("2006-01-02 15:04:05")
	if host["ip"] == "127.0.0.1" {
		response <- []string{"check_host", host["ip"], "0", "ok", NowTime}
	} else {
		sshPort, _ := strconv.ParseInt(host["port"], 10, 64)
		sshClient := utils.ClientConfig{
			Host:     host["ip"],
			Port:     sshPort,
			Username: host["user"],
			Password: host["password"],
		}
		connectErr := sshClient.CreateClient()
		if connectErr != nil {
			response <- []string{"check_host", host["ip"], "1", "connect failed", NowTime, connectErr.Error()}
		} else {
			res, err := sshClient.RunShell(args)
			if err != nil {
				response <- []string{"check_host", host["ip"], "1", res, NowTime, err.Error()}
			} else {
				response <- []string{"check_host", host["ip"], "0", res, NowTime}
			}
		}
	}

	<-limiter
}
