package utils

import (
	"encoding/csv"
	"os"
)

const (
	CheckFail = 1
	CheckOk   = 2
	CheckWarn = 3
)

// CheckItem aaa
type CheckItem struct {
	Name       string
	Datetime   string
	Status     int
	StatusText string
	Msg        string
	Doc        string
	Type       string
	Host       string
	Top        string
}

// TaskResult aaa
type TaskResult struct {
	Items   []CheckItem
	Success int
	Error   int
	Warn    int
	Msg     string
}

func emptyData(msg string) CheckItem {
	return CheckItem{
		Name:       msg,
		Datetime:   "",
		Status:     0,
		StatusText: "",
		Msg:        "",
		Doc:        "",
		Type:       "",
		Host:       "",
		Top:        "",
	}
}

// ParseSetupToolsCsv TASK_NAME,HOST,STATUS,LOGS,DESC
func ParseSetupToolsCsv(path string) TaskResult {
	var result = TaskResult{
		Items:   []CheckItem{},
		Success: 0,
		Error:   0,
		Warn:    0,
	}
	if len(path) == 0 || !PathExists(path) || path == "no history" {
		result.Msg = "No Data"
		result.Items = []CheckItem{
			emptyData("No Data"),
		}
		return result
	}
	file, err := os.Open(path)
	if err != nil {
		result.Msg = "文件打开异常"
		result.Items = []CheckItem{
			emptyData("文件打开异常"),
		}
		Log.Printf("open file(%v) Error : %v", path, err)
		return result
	}
	defer file.Close()
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		result.Msg = "文件内容读取异常，请确保是CSV文件"
		result.Items = []CheckItem{
			emptyData("文件内容读取异常，请确保是CSV文件"),
		}
		Log.Printf("open file(%v) Error : %v", path, err)
		return result
	}

	for _, record := range records {
		if len(record) < 5 {
			continue
		}
		name := record[0]
		if name == "TASK_NAME" {
			continue
		}
		status := record[2]
		item := CheckItem{
			Name:       record[0],
			Status:     0,
			StatusText: status,
			Msg:        record[3],
			Doc:        record[4],
			Host:       record[1],
		}
		if status == OK {
			result.Success++
		} else if status == FAIL || status == ERROR {
			result.Items = append(result.Items, item)
			result.Error++
		} else if status == WARN {
			result.Warn++
		} else {
			Log.Printf("record : %v", record)
		}
	}
	if len(result.Items) == 0 {
		result.Items = append(result.Items, CheckItem{
			Name:       "ok",
			Status:     0,
			StatusText: "ok",
			Msg:        "ok",
			Doc:        "ok",
			Host:       "ok",
		})
	}
	return result
}
