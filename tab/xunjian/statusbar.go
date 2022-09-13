package xunjian

import (
	"fmt"
	"image"
	"log"
	"os"
	"time"

	ui "github.com/gizak/termui/v3"
)

// StatusBar status bar
type StatusBar struct {
	ui.Block
}

var (
	appName = ""
)

// NewStatusBar status bar
func NewStatusBar(name string) *StatusBar {
	appName = name
	self := &StatusBar{*ui.NewBlock()}
	self.Border = false
	return self
}

// Draw status bar
func (ctx *StatusBar) Draw(buf *ui.Buffer) {
	ctx.Block.Draw(buf)

	hostname, err := os.Hostname()
	if err != nil {
		log.Printf("could not get hostname: %v", err)
		return
	}
	buf.SetString(
		hostname,
		ui.NewStyle(ui.ColorWhite),
		image.Pt(ctx.Inner.Min.X, ctx.Inner.Min.Y+(ctx.Inner.Dy()/2)),
	)

	currentTime := time.Now()
	formattedTime := currentTime.Format("15:04:05")
	formattedTime = fmt.Sprintf("%v  %v", formattedTime, "Tab 切换标签/Esc 返回上一级")
	buf.SetString(
		formattedTime,
		ui.NewStyle(ui.ColorWhite),
		image.Pt(
			ctx.Inner.Min.X+(ctx.Inner.Dx()/2)-len(formattedTime)/2,
			ctx.Inner.Min.Y+(ctx.Inner.Dy()/2),
		),
	)

	buf.SetString(
		appName,
		ui.NewStyle(ui.ColorWhite),
		image.Pt(
			ctx.Inner.Max.X-6,
			ctx.Inner.Min.Y+(ctx.Inner.Dy()/2),
		),
	)
}
