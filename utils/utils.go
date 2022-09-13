package utils

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/wxnacy/wgo/arrays"
	"golang.org/x/text/encoding/simplifiedchinese"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"math"
	"net/http"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strings"
	"unsafe"
)

var (
	KB = uint64(math.Pow(2, 10))
	MB = uint64(math.Pow(2, 20))
	GB = uint64(math.Pow(2, 30))
	TB = uint64(math.Pow(2, 40))
)

var (
	Home       = GetCurDir()
	CpuNum     = runtime.NumCPU()
	UserName   = GetCurUser()
	ErrCode    = map[string]string{}
	TcsTools   = Home + "/tcs-tools"
	SetupTools = Home + "/tcs-tools/setup-tools"
	Temp       = "/tmp/" + UserName
	HostYml    = TcsTools + "/conf/host.yml"
	confData   = GetConf()
	ParamData  = map[string]string{
		"Home":    Home,
		"LogDir":  Home + "/logs",
		"LogFile": fmt.Sprintf("%v/logs/monitor.log", Home),
		"Conf":    Home + "/conf",
	}
	UpdateCode = "no"
	Tabs       = []string{
		"基础容量", "业务容量", "巡检排障", "扫描播报",
	}
	InitData = map[string]map[string][][]string{"获取监控数据": {
		"加载数据中": {
			{"加载数据中", "local", "0", "加载数据中", "0", "0"},
		}},
	}
	InitMonitorData = map[string]string{
		"baseItem":     "init",
		"businessItem": "init",
	}
)

const (
	WAIT    = "WAIT"
	RUNNING = "RUNNING"
	FAIL    = "FAIL"
	OK      = "OK"
	ERROR   = "ERROR"
	WARN    = "WARN"
)

type activeObj struct {
	Tab   int
	Form  string
	Index int
}

var (
	Active = activeObj{
		Tab:   0,
		Form:  "",
		Index: 0,
	}
)

// HostInfo aaa
type HostInfo struct {
	Ip       string `yaml:"ip"`
	User     string `yaml:"user"`
	Port     int    `yaml:"port"`
	Password string `yaml:"password"`
}

// HostMonitorInfo aaa
type HostMonitorInfo struct {
	Time string `json:"time"`
	Ip   string `json:"ip"`
	Cpu  string `json:"cpu"`
	Mem  string `json:"mem"`
	Swap string `json:"swap"`
	Disk []struct {
		Disk string `json:"disk"`
		Used string `json:"used"`
	} `json:"disk"`
}

// MonitorInfo aaa
type MonitorInfo struct {
	Host     string              `json:"ip"`
	Time     string              `json:"time"`
	ItemList []map[string]string `json:"itemList"`
}

// GetConf get config
func GetConf() map[string]string {
	File, err := ioutil.ReadFile(Home + "/conf/conf.yml")
	if err != nil {
		Log.Println("read host.yml fail: " + Home + "/conf/conf.yml")
		os.Exit(1)
	}
	var config map[string]string
	yaml.Unmarshal(File, &config)
	return config
}

// GetCurUser aaa
func GetCurUser() string {
	u, _ := user.Current()
	return u.Username
}

// BytesToKB aaa
func BytesToKB(b uint64) float64 {
	return float64(b) / float64(KB)
}

// BytesToMB aaa
func BytesToMB(b uint64) float64 {
	return float64(b) / float64(MB)
}

// BytesToGB aaa
func BytesToGB(b uint64) float64 {
	return float64(b) / float64(GB)
}

// BytesToTB aaa
func BytesToTB(b uint64) float64 {
	return float64(b) / float64(TB)
}

// ConvertBytes aaa
func ConvertBytes(b uint64) (float64, string) {
	switch {
	case b < KB:
		return float64(b), "B"
	case b < MB:
		return BytesToKB(b), "KB"
	case b < GB:
		return BytesToMB(b), "MB"
	case b < TB:
		return BytesToGB(b), "GB"
	default:
		return BytesToTB(b), "TB"
	}
}

// GetCurDir aaa
func GetCurDir() string {
	execFile, _ := filepath.Abs(os.Args[0])
	execFileTmp1 := strings.Split(execFile, `/`)
	basedir := strings.Join(execFileTmp1[0:len(execFileTmp1)-1], `/`)
	if runtime.GOOS == "windows" {
		return "./"
	}
	return basedir
}

// MaxInt aaa
func MaxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// ConvertLocalizedString aaa
func ConvertLocalizedString(s string) string {
	if strings.ContainsAny(s, ",") {
		return strings.Replace(s, ",", ".", 1)
	} else {
		return s
	}
}

// OpenFile 判断文件是否存在  存在则OpenFile 不存在则Create
func OpenFile(filename string) (*os.File, error) {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		os.Create(filename)
		return os.OpenFile(filename, os.O_WRONLY|os.O_APPEND, 0666)
	}
	return os.OpenFile(filename, os.O_WRONLY|os.O_APPEND, 0666) //打开文件
}

// 去除字符串两端无用字符集
func strip(s_ string, chars_ string) string {
	s, chars := []rune(s_), []rune(chars_)
	length := len(s)
	max := len(s) - 1
	l, r := true, true //标记当左端或者右端找到正常字符后就停止继续寻找
	start, end := 0, max
	tmpEnd := 0
	charset := make(map[rune]bool) //创建字符集，也就是唯一的字符，方便后面判断是否存在
	for i := 0; i < len(chars); i++ {
		charset[chars[i]] = true
	}
	for i := 0; i < length; i++ {
		if _, exist := charset[s[i]]; l && !exist {
			start = i
			l = false
		}
		tmpEnd = max - i
		if _, exist := charset[s[tmpEnd]]; r && !exist {
			end = tmpEnd
			r = false
		}
		if !l && !r {
			break
		}
	}
	if l && r { // 如果左端和右端都没找到正常字符，那么表示该字符串没有正常字符
		return ""
	}
	return string(s[start : end+1])
}

// CmdRun Run cmd
func CmdRun(command string) (resultStr string, status error) {
	var result []byte
	var err error

	sysType := runtime.GOOS
	if sysType == "windows" {
		result, err = exec.Command("cmd", "/c", command).CombinedOutput()
	} else if sysType == "linux" {
		result, err = exec.Command("/bin/sh", "-c", command).CombinedOutput()
	} else {
		Log.Printf(fmt.Sprintf("no support system: %v", sysType))
	}
	resultStr = ConvertByte2String(result, "GB18030")
	return strip(resultStr, "\n"), err
}

type Charset string

// ConvertByte2String 字符转换
func ConvertByte2String(byte []byte, charset Charset) string {
	UTF8 := Charset("UTF-8")
	GB18030 := Charset("GB18030")
	var str string
	switch charset {
	case GB18030:
		var decodeBytes, _ = simplifiedchinese.GB18030.NewDecoder().Bytes(byte)
		str = string(decodeBytes)
	case UTF8:
		fallthrough
	default:
		str = string(byte)
	}
	return str
}

// GetAllFiles 获取指定目录下的所有文件,包含子目录下的文件
func GetAllFiles(folder string) ([]string, error) {
	var result []string
	filepath.Walk(folder, func(path string, fi os.FileInfo, err error) error {
		if err != nil {
			Log.Printf(err.Error())
			return err
		}
		if !fi.IsDir() {
			//如果想忽略这个目录 return filepath.SkipDir
			result = append(result, path)
		}
		return nil
	})
	return result, nil
}

// GetLastFile aaa
func GetLastFile(dir, flag string) string {
	var resultPath string
	files, _ := GetAllFiles(dir)
	for _, file := range files {
		path, fileName := filepath.Split(file)
		match, _ := filepath.Match("*."+flag, fileName)
		if !match {
			continue
		}
		resultPath = path + "/" + fileName
	}
	return resultPath
}

// SortDataList aaa
func SortDataList(list [][]string) {
	sort.Slice(list, func(i, j int) bool {
		return list[i][1] > list[j][1] //按照每行的第一个元素排序
	})
}

// SortFile aaa
func SortFile(kArray []string) {
	sort.Slice(kArray, func(i, j int) bool {
		return kArray[i] > kArray[j]
	})
}

var (
	GroupHost = map[string][]map[string]string{}
)

func mapInterfaceChange(src map[interface{}]interface{}) map[string]string {
	mapTmp := map[string]string{}
	for k, v := range src {
		mapTmp[k.(string)] = fmt.Sprintf("%v", v)
	}
	return mapTmp
}
func checkMapKey(map1 map[string]string, key string) bool {
	if _, ok := map1[key]; ok {
		return true
	} else {
		return false
	}
}

func getGroupParam(allParams map[string]map[string]string, host map[string]string, group, key, value string) {
	if !checkMapKey(host, key) {
		if _, ok := allParams[group]; ok {
			if checkMapKey(allParams[group], key) {
				host[key] = allParams[group][key]
				return
			}
		}
		if _, ok := allParams["ALL"]; ok {
			host[key] = allParams["ALL"][key]
		} else {
			host[key] = value
		}
	}
}
func checkHostExist(hostList []map[string]string, dstHost map[string]string) bool {
	for _, host := range hostList {
		if host["ip"] == dstHost["ip"] {
			return true
		}
	}
	return false
}

// ReadSrcHost aaa
func ReadSrcHost(filePath string) map[string][]map[string]string {
	File, err := ioutil.ReadFile(filePath)
	if err != nil {
		Log.Println("read host.yml fail: " + filePath)
		os.Exit(1)
	}
	var config interface{}
	yaml.Unmarshal(File, &config)
	allParams := map[string]map[string]string{}
	allHostsTmp := map[string][]map[string]string{}
	for key, value := range config.(map[interface{}]interface{}) {
		groupName := strings.Replace(key.(string), "[vars]", "", -1)
		groupName = strings.Replace(groupName, "[var]", "", -1)
		if strings.Contains(key.(string), "[vars]") || strings.Contains(key.(string), "[var]") {
			params := mapInterfaceChange(value.(map[interface{}]interface{}))
			allParams[groupName] = params
		} else {
			var tmpHosts []map[string]string
			for _, host := range value.([]interface{}) {
				hostInfo := mapInterfaceChange(host.(map[interface{}]interface{}))
				if !checkHostExist(tmpHosts, hostInfo) {
					tmpHosts = append(tmpHosts, hostInfo)
				}
			}
			allHostsTmp[groupName] = tmpHosts
		}
	}
	allHostsTmp["local"] = []map[string]string{{"ip": "127.0.0.1"}}
	for group, hostList := range allHostsTmp {
		for _, host := range hostList {
			if _, ok := host["instance_key"]; ok {
				password, _ := CmdRun(fmt.Sprintf(
					"python %v/scheduler/keygen.py -m decrypt -k '%v' -s '%v' ",
					SetupTools, host["instance_key"], host["password"],
				))
				host["password"] = password
			}
			getGroupParam(allParams, host, group, "user", "root")
			getGroupParam(allParams, host, group, "port", "22")
			getGroupParam(allParams, host, group, "password", "root")
			GroupHost[group] = append(GroupHost[group], host)
			GroupHost["ALL"] = append(GroupHost[group], host)
			if _, ok := host["name"]; ok {
				nameList := strings.Split(host["name"], `,`)
				if len(nameList) == 0 {
					continue
				}
				for _, name := range nameList {
					GroupHost[name] = append(GroupHost[name], host)
				}
			}

		}
	}
	return GroupHost
}

// CreateQrCode aaa
func CreateQrCode(data string) {
	var cmdStr string
	if data == "[\"ok\"]" {
		cmdStr = fmt.Sprintf(
			"echo > %s/tmp/qrcode.txt", Home,
		)
	} else {
		cmdStr = fmt.Sprintf(
			"python %s/tools/qrcode_tool.py -t '%s' -p > %s/tmp/qrcode.txt",
			SetupTools, data, Home,
		)
	}

	res, err := CmdRun(cmdStr)
	if err != nil {
		Log.Printf("[ERROR]: " + cmdStr)
		Log.Printf("[ERROR]: " + res)
	}

}

// CheckListExist aaa
func CheckListExist(list []string, value string) bool {
	index := arrays.ContainsString(list, value)
	if index == -1 {
		return false
	} else {
		return true
	}
}

// ReadERRCode aaa
func ReadERRCode() {
	f, err := os.Open(TcsTools + "/conf/err_code.csv")
	if err != nil {
		return
	}
	reader := csv.NewReader(f)
	preData, err := reader.ReadAll() // preData 数据格式为 [][]string
	if err != nil {
		return
	}
	for _, line := range preData {
		ErrCode[line[2]] = line[0]
	}
}

// UniqList aaa
func UniqList(list []string) []string {
	// 创建一个临时map用来存储数组元素
	temp := make(map[string]bool)
	var newList []string
	for _, v := range list {
		// 遍历数组元素，判断此元素是否已经存在map中
		if _, ok := temp[v]; ok {
			continue
		} else {
			temp[v] = true
			newList = append(newList, v)
		}
	}
	return newList

}

// PathExists 判断文件是否存在
func PathExists(path string) bool {
	_, err := os.Stat(path)
	//当为空文件或文件夹存在
	if err == nil {
		return true
	}
	//os.IsNotExist(err)为true，文件或文件夹不存在
	if os.IsNotExist(err) {
		return false
	}
	//其它类型，不确定是否存在
	return false
}

// WriteFile aaa
func WriteFile(context, filePath string) {
	err := ioutil.WriteFile(filePath, []byte(context), 0755)
	if err != nil {
		Log.Println("文件写入失败", err)
		return
	}
}

// CodeInfo aaa
type CodeInfo struct {
	LtcId           string                `json:"ltc_id"`
	Jid             string                `json:"jid"`
	BaseMonitor     map[string][][]string `json:"base_monitor"`
	BusinessMonitor map[string][][]string `json:"business_monitor"`
	LastResult      []int                 `json:"last_result"`
}

// ReadJsonFileData aaa
func ReadJsonFileData(filepath string) map[string][][]string {
	jsonByte, err := ioutil.ReadFile(filepath)
	if err != nil {
		jsonByte = []byte("{}")
	}
	var jsonData map[string][][]string
	_ = json.Unmarshal(jsonByte, &jsonData)

	return jsonData
}
func monitorFormat(monitor map[string][][]string) string {
	var monitorStr = "[\n"
	for k, v := range monitor {
		var tmpList1 []string
		for _, v1 := range v {
			v1Str, _ := json.Marshal(v1)
			tmpList1 = append(tmpList1, string(v1Str))
		}
		tmpList1Str := strings.Join(tmpList1, ",\n  ")
		tmpStr1 := fmt.Sprintf(" {\"%v\":\n  [%v]\n }", k, tmpList1Str)
		monitorStr = monitorStr + tmpStr1 + ",\n "
	}
	re, _ := regexp.Compile("},\n $")
	monitorStrNew := re.ReplaceAllString(monitorStr, "}\n")

	return monitorStrNew + "]"
}

// JsonFormat aaa
func JsonFormat(data CodeInfo) string {
	result, _ := json.Marshal(data.LastResult)
	dataFormat := fmt.Sprintf(`{
"ltc_id": "%v",
"jid": "%v",
"base_monitor": %v,
"business_monitor": %v,
"inspect_result": %v
}`, data.LtcId, data.Jid, monitorFormat(data.BaseMonitor), monitorFormat(data.BusinessMonitor), string(result))
	return dataFormat
}

// GetData aaa
func GetData(format bool) string {
	resultPath := GetLastFile(TcsTools+"/csv", "csv")
	data := ParseSetupToolsCsv(resultPath)
	lastResult := []int{data.Success, data.Warn, data.Error}
	baseReport := ReadJsonFileData(Home + "/tmp/baseItem.json")
	businessReport := ReadJsonFileData(Home + "/tmp/businessItem.json")

	reportData := CodeInfo{
		LtcId:           confData["ltc_id"],
		Jid:             confData["jid"],
		BaseMonitor:     baseReport,
		BusinessMonitor: businessReport,
		LastResult:      lastResult,
	}

	reportStr, _ := json.Marshal(reportData)

	if format {
		return JsonFormat(reportData)
	} else {
		return string(reportStr)
	}

}

// CheckListExists aaa
func CheckListExists(value string, list []string) bool {
	for _, v := range list {
		if v == value {
			return true
		}
	}
	return false
}

// PostData aaa
func PostData(data, url string) (res string, err error) {
	reader := bytes.NewReader([]byte(data))
	request, err := http.NewRequest("POST", url, reader)
	defer request.Body.Close() //程序在使用完回复后必须关闭回复的主体
	if err != nil {
		return "", err
	}
	request.Header.Set("Content-Type", "application/json;charset=UTF-8")
	request.Header.Set("'User-Agent", "Apipost client Runtime/+https://www.apipost.cn/")
	//必须设定该参数,POST参数才能正常提交，意思是以json串提交数据
	client := http.Client{}
	resp, err := client.Do(request) //Do 方法发送请求，返回 HTTP 回复
	if err != nil {
		return "", err
	}

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	//byte数组直接转成string，优化内存
	str := (*string)(unsafe.Pointer(&respBytes))
	return *str, nil
}

// CrontabAdmin aaa
func CrontabAdmin(op, cron string) {
	if op == "del" {
		cmd := "sed -i '/tianxun-lite -m report/d' /var/spool/cron/$USER"
		res, err := CmdRun(cmd)
		if err != nil {
			Log.Println("del crontab fail: " + res)
			fmt.Println("del crontab fail: " + res)
			return
		}
		Log.Println("del crontab ok")
		fmt.Println("del crontab ok")
		return
	}
	check, code := CmdRun("sed -n '/tianxun-lite -m report/p' /var/spool/cron/$USER")
	if code != nil {
		Log.Println("get crontab fail: " + check)
		fmt.Println("get crontab fail: " + check)
		return
	}
	var addCmd string
	if check == "" {
		addCmd = fmt.Sprintf("echo '%v' >> /var/spool/cron/$USER", cron)

	} else {
		addCmd = fmt.Sprintf("sed -i 's#.*tianxun-lite -m report#%v#g' /var/spool/cron/$USER", cron)
	}
	addRes, addCode := CmdRun(addCmd)
	Log.Println("cmd: " + addCmd)
	fmt.Println("cmd: " + addCmd)
	if addCode != nil {
		Log.Println("add crontab fail: " + addRes)
		fmt.Println("add crontab fail: " + addRes)
		return
	} else {
		Log.Println("add crontab ok")
		fmt.Println("add crontab ok")
	}
}

func runCmd(cmdStr string) string {
	list := strings.Split(cmdStr, ",")
	cmd := exec.Command(list[0], list[1:]...)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return ""
	}
	resultStr := ConvertByte2String([]byte(out.String()), "GB18030")
	return resultStr
}

// WhenExpr aaa
func WhenExpr(exprStr string) bool {
	cmd := fmt.Sprintf("python,pyeval.py,%v", exprStr)
	res := runCmd(cmd)
	if strings.Contains(res, "True") {
		return true
	} else {
		return false
	}
}

// ParseParams aaa
func ParseParams(paramsStr []string, globalParams map[string]string) {
	for _, params := range paramsStr {
		// 替换第一个等号 用作分隔符防止有多个等号时冲突
		newFlag := `==flag==`
		paramsNew := strings.Replace(params, `=`, newFlag, 1)
		param := strings.Split(paramsNew, newFlag)
		if len(param) != 2 {
			Log.Println("param format error: " + params)
			continue
		} else {
			globalParams[param[0]] = param[1]
		}
	}
}

// ReplaceParam aaa
func ReplaceParam(src string, param map[string]string) string {
	paramFlag := `\$\{\S+\}`
	re := regexp.MustCompile(paramFlag)
	res := re.FindAllStringSubmatch(src, -1)
	for _, variable := range res {
		variableName := strings.Replace(variable[0], `${`, ``, 1)
		variableName = strings.Replace(variableName, `}`, ``, 1)
		if _, ok := param[variableName]; ok {
			src = strings.Replace(src, `${`+variableName+`}`, param[variableName], -1)
		}
	}
	return src
}

// GetHostInfo aaa
func GetHostInfo(hostList []map[string]string, ip string) map[string]string {
	for _, hostInfo := range hostList {
		if hostInfo["ip"] == ip {
			return hostInfo
		}
	}
	return nil
}

// Error aaa
func Error(text string) error {
	return &errorString{text}
}

type errorString struct {
	s string
}

func (e *errorString) Error() string {
	return e.s
}
