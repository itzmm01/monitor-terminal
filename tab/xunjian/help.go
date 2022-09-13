package xunjian

import (
	"image"
	"strings"

	ui "github.com/gizak/termui/v3"
)

// HelpMenu help
type HelpMenu struct {
	ui.Block
	content string
}

// NewHelpMenu help
func NewHelpMenu(content string) *HelpMenu {
	ctx := &HelpMenu{
		Block:   *ui.NewBlock(),
		content: " <ESC> 退出 \n" + content,
	}
	return ctx
}

// Resize help
func (cts *HelpMenu) Resize(termWidth, termHeight int) {
	var textWidth = 0
	for _, line := range strings.Split(cts.content, "\n") {
		textWidth = maxInt(len(line), textWidth)
	}
	textWidth += 20
	textHeight := 22
	x := (termWidth - textWidth) / 2
	y := (termHeight - textHeight) / 2

	cts.Block.SetRect(x, y, textWidth+x, textHeight+y)
}

// Draw help
func (cts *HelpMenu) Draw(buf *ui.Buffer) {
	cts.Block.Draw(buf)

	for y, line := range strings.Split(cts.content, "\n") {
		for x, rune := range line {
			buf.SetCell(
				ui.NewCell(rune, ui.NewStyle(7)),
				image.Pt(cts.Inner.Min.X+x, cts.Inner.Min.Y+y-1),
			)
		}
	}
}

func maxInt(a int, b int) int {
	if a > b {
		return a
	}
	return b
}
