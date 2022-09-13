package ui

import (
	"github.com/gizak/termui/v3"
	"image"
)

// GridWrap grid
type GridWrap struct {
	grid *termui.Grid
}

// GetRect grid
func (grid *GridWrap) GetRect() image.Rectangle {
	return grid.grid.GetRect()
}

// SetRect grid
func (grid *GridWrap) SetRect(i int, i2 int, i3 int, i4 int) {
	grid.grid.SetRect(i, i2, i3, i4)
}

// Lock grid
func (grid *GridWrap) Lock() {
	grid.grid.Lock()
}

// Unlock grid
func (grid *GridWrap) Unlock() {
	grid.grid.Unlock()
}

// Draw grid
func (grid *GridWrap) Draw(buf *termui.Buffer) {
	grid.grid.Draw(buf)
}

// Clear grid
func (grid *GridWrap) Clear() {
	grid.grid.Items = []*termui.GridItem{}
}

// NewGridWrap aaa
func NewGridWrap() *GridWrap {
	wrap := GridWrap{
		grid: termui.NewGrid(),
	}
	return &wrap
}

// Set aaa
func (grid *GridWrap) Set(entries ...interface{}) {
	grid.grid.Set(entries...)
}
