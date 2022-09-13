package tab

import (
	"monitor-ter/tab/base-monitor"
	"monitor-ter/tab/business"
	"monitor-ter/tab/report"
	XJ "monitor-ter/tab/xunjian"
	"os"
	"os/signal"
	"syscall"
	"time"

	"monitor-ter/ui"
	"monitor-ter/utils"

	"github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

const (
	AppName = "monitor-terminal"
)

var (
	color = ui.Default
	rate  = time.Second
	tab   *widgets.TabPane
	help  *XJ.HelpMenu
	grid  *termui.Grid
	bar   *XJ.StatusBar

	Log = utils.Log
)

func setupGrid(test string) {
	grid = termui.NewGrid()
	if test == "xj" {
		XJ.SetGrid(grid)
	} else if test == "business" {
		business.SetGrid(grid)
	} else if test == "report" {
		report.SetGrid(grid)
	} else {
		baseMonitor.SetGrid(grid)
	}
}

// SetDefaultColors tab
func SetDefaultColors() {
	termui.Theme.Default = termui.NewStyle(termui.Color(color.Fg), termui.Color(color.Bg))
	termui.Theme.Block.Title = termui.NewStyle(termui.Color(color.BorderLabel), termui.Color(color.Bg))
	termui.Theme.Block.Border = termui.NewStyle(termui.Color(color.BorderLine), termui.Color(color.Bg))
}

// InitWidgets tab
func InitWidgets() {
	baseMonitor.InitWidgets(true)
	business.InitWidgets(true)
	XJ.InitTxWidgets()
	report.InitTxWidgets()
	help = XJ.NewHelpMenu("")
	bar = XJ.NewStatusBar(AppName)

}

// RenderTab tab
func RenderTab() {
	termWidth, termHeight := termui.TerminalDimensions()
	setupGrid("")
	bar.SetRect(0, termHeight-1, termWidth, termHeight)
	grid.SetRect(1, 1, termWidth, termHeight-1)
	switch tab.ActiveTabIndex {
	case 0:
		termui.Render(tab, grid, bar)
	case 1:
		setupGrid("business")
		grid.SetRect(1, 1, termWidth, termHeight-1)
		termui.Render(tab, grid, bar)
	case 2:
		setupGrid("xj")
		grid.SetRect(1, 1, termWidth, termHeight-1)
		termui.Render(tab, grid, bar)
	case 3:
		setupGrid("report")
		grid.SetRect(1, 1, termWidth, termHeight-1)
		termui.Render(tab, grid, bar)
	}
}
func rendBar() {
	termui.Render(bar)
}

func timeNext() {
	rendBar()
	if tab.ActiveTabIndex == 2 {
		XJ.TestUpdate()
	} else if tab.ActiveTabIndex == 0 {
		initName := "baseItem"
		if utils.InitMonitorData[initName] != "init" && utils.InitMonitorData[initName] != "running" {
			utils.InitMonitorData[initName] = "running"
			go func() {
				baseMonitor.InitWidgets(false)
				RenderTab()
			}()
		}
	} else if tab.ActiveTabIndex == 1 {
		initName := "businessItem"
		if utils.InitMonitorData[initName] != "init" && utils.InitMonitorData[initName] != "running" {
			utils.InitMonitorData[initName] = "running"
			go func() {
				business.InitWidgets(false)
				RenderTab()
			}()
		}
	}
}

// EventLoop tab
func EventLoop() {
	setupGrid("")
	termWidth, termHeight := termui.TerminalDimensions()

	tab = widgets.NewTabPane(utils.Tabs...)
	tab.SetRect(0, 0, termWidth, 1)
	tab.Border = true

	//gridInspect.SetRect(1, 1, termWidth, termHeight-1)
	help.Resize(termWidth, termHeight)

	drawTicker := time.NewTicker(rate).C

	// handles kill signal
	sigTerm := make(chan os.Signal, 2)
	signal.Notify(sigTerm, os.Interrupt, syscall.SIGTERM)

	uiEvents := termui.PollEvents()

	previousKey := ""
	RenderTab()
	for {
		select {
		case <-sigTerm:
			return
		case <-drawTicker:
			timeNext()
		case e := <-uiEvents:
			utils.Active.Tab = tab.ActiveTabIndex
			var skip = false
			if tab.ActiveTabIndex == 0 {
				skip = baseMonitor.TxEventLoop(e.ID)
			} else if tab.ActiveTabIndex == 1 {
				skip = business.TxEventLoop(e.ID)
			} else if tab.ActiveTabIndex == 2 {
				skip = XJ.TxEventLoop(e.ID)
			} else if tab.ActiveTabIndex == 3 {
				skip = report.TxEventLoop(e.ID)
			}
			if skip {
				continue
			}
			switch e.ID {

			case "q", "<C-c>":
				return
			case "<Resize>":
				payload := e.Payload.(termui.Resize)
				help.Resize(payload.Width, payload.Height)
				grid.SetRect(1, 1, payload.Width, payload.Height)
				termui.Clear()
				RenderTab()
			case "?":
				RenderTab()
			case "<Tab>":
				if tab.ActiveTabIndex == 3 {
					tab.ActiveTabIndex = 0
				} else {
					tab.FocusRight()
				}
				termui.Clear()
				RenderTab()
			}
			if previousKey == e.ID {
				previousKey = ""
			} else {
				previousKey = e.ID
			}

		}
	}
}

// Main tab
func Main() {
	if err := termui.Init(); err != nil {
		Log.Fatalf("failed to initialize ui: %v", err)
	}
	defer termui.Close()
	SetDefaultColors()
	InitWidgets()
	EventLoop()
}
