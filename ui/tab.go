package ui

import (
	"github.com/gizak/termui/v3"
	"image"
)

// TabPane a
type TabPane struct {
	views     []termui.Drawable
	rectangle image.Rectangle
	lock      bool
}

// NewTabPane a
func NewTabPane() *TabPane {
	wrap := TabPane{
		views:     []termui.Drawable{},
		rectangle: image.Rect(0, 0, 0, 0),
	}
	return &wrap
}

// GetRect a
func (tab *TabPane) GetRect() image.Rectangle {
	return tab.rectangle
}

// SetRect a
func (tab *TabPane) SetRect(i int, i2 int, i3 int, i4 int) {
	tab.rectangle = image.Rect(i, i2, i3, i4)
}

// Lock a
func (tab *TabPane) Lock() {
	tab.lock = true
}

// Unlock a
func (tab *TabPane) Unlock() {
	tab.lock = false
}

// Draw a
func (tab *TabPane) Draw(buf *termui.Buffer) {
	rect := tab.GetRect()
	for _, view := range tab.views {
		view.SetRect(rect.Min.X, rect.Min.Y, rect.Max.X, rect.Max.Y)
		view.Lock()
		view.Draw(buf)
		view.Unlock()
	}
}

// Clear a
func (tab *TabPane) Clear() {
	tab.views = []termui.Drawable{}
}

// Set a
func (tab *TabPane) Set(entries ...termui.Drawable) {
	for _, item := range entries {
		tab.views = append(tab.views, item)
	}
}
