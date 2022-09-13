package tx

import (
	"fmt"
	"github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"monitor-ter/utils"
	"time"
)

// ExecWidget form obj
type ExecWidget struct {
	*widgets.Paragraph
	updateInterval time.Duration
	Status         string
}

// JobTask job
type JobTask struct {
	name        string `yaml:"name"`
	taskType    string `yaml:"type"`
	localScript string `yaml:"local_script"`
	cmd         string `yaml:"cmd"`
	allowFailed bool   `yaml:"allow_failed"`
	desc        string `yaml:"desc"`
}

// JobData job data
type JobData struct {
	name        string     `yaml:"name"`
	host        string     `yaml:"host"`
	tasks       *[]JobTask `yaml:"tasks"`
	allowFailed bool       `yaml:"allow_failed"`
}

// JobItem job item
type JobItem struct {
	job *JobData `yaml:"job"`
}

const (
	WAIT    = "WAIT"
	RUNNING = "RUNNING"
	FAIL    = "FAIL"
	OK      = "OK"
	ERROR   = "ERROR"
	WARN    = "WARN"
)

var CheckLog = utils.GetLogger("comm")

var (
	view      *ExecWidget
	colorDict = map[string]termui.Color{
		"WAIT":    termui.ColorWhite,
		"RUNNING": termui.ColorYellow,
		"FAIL":    termui.ColorRed,
		"OK":      termui.ColorGreen,
	}
)

// NewExecWidget init
func NewExecWidget(callback func()) *ExecWidget {
	view = &ExecWidget{
		Paragraph: widgets.NewParagraph(),
		Status:    WAIT,
	}
	view.Title = "Ctrl+R开始排障"
	view.TitleStyle.Modifier = termui.ModifierBold
	view.TextStyle.Modifier = termui.ModifierBold
	//textBox.TextStyle.Modifier = ui.ModifierUnderline
	view.BorderRight = true
	view.Border = true
	Update(true)
	return view
}

// Start starting
func Start(item string) {
	if view.Status == RUNNING {
		Log.Printf("tx is running, skip")
	} else {
		view.Status = RUNNING
		Update(true)
		go execTx(item)
	}
}

// GetStatus get status
func GetStatus() string {
	return view.Status
}

// Stop stopping
func Stop() {
	view.Status = WAIT
	Update(true)
}

// Update update
func Update(render bool) {
	var status = view.Status
	switch status {
	case WAIT:
		view.Text = "Ctrl+R开始排障"
	case RUNNING:
		view.Text = "排障中...(Ctrl+Q退出)"
	case FAIL:
		view.Text = "排障失败"
		utils.UpdateCode = "yes"
	case OK:
		view.Text = "排障成功"
		utils.UpdateCode = "yes"
	}
	if _, ok := colorDict[view.Status]; ok {
		view.TextStyle.Fg = colorDict[view.Status]
		view.TitleStyle.Fg = colorDict[view.Status]
	}
	termui.Render(view)
}

// execTx 执行巡检
func execTx(item string) {
	scanTime := time.Now().Format("2006-01-02-150405")
	CheckLog.Printf("execTx: %v", scanTime)
	RunCmd := fmt.Sprintf("cd %v && sh tcs-tools.sh check -mod %v", utils.TcsTools, item)
	Log.Printf("cmd : %v", RunCmd)
	_, err := utils.CmdRun(RunCmd)
	if err != nil {
		Log.Printf("execTx fail: %v", err)
		view.Status = FAIL
	} else {
		view.Status = OK
	}
	Update(true)
}
