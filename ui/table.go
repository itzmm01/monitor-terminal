package ui

import (
	"fmt"
	"image"
	"log"
	"strings"

	"github.com/gizak/termui/v3"
)

// Table a
type Table struct {
	*termui.Block

	Header []string
	Rows   [][]string

	ColWidths []int
	ColGap    int
	PadLeft   int

	ShowCursor  bool
	CursorColor termui.Color

	ShowLocation bool

	UniqueCol    int    // the column used to uniquely identify each table row
	SelectedItem string // used to keep the cursor on the correct item if the data changes
	SelectedRow  int
	TopRow       int // used to indicate where in the table we are scrolled at

	ColResizer func()
}

// NewTable returns a new Table instance
func NewTable() *Table {
	return &Table{
		Block:       termui.NewBlock(),
		SelectedRow: 0,
		TopRow:      0,
		UniqueCol:   0,
		ColResizer:  func() {},
	}
}

// Draw a
func (cts *Table) Draw(buf *termui.Buffer) {
	cts.Block.Draw(buf)

	if cts.ShowLocation {
		cts.drawLocation(buf)
	}

	cts.ColResizer()

	// finds exact column starting position
	var colXPos []int
	cur := 1 + cts.PadLeft
	for _, w := range cts.ColWidths {
		colXPos = append(colXPos, cur)
		cur += w
		cur += cts.ColGap
	}

	// prints header
	for i, h := range cts.Header {
		width := cts.ColWidths[i]
		if width == 0 {
			continue
		}
		// don't render column if it doesn't fit in widget
		if width > (cts.Inner.Dx()-colXPos[i])+1 {
			continue
		}
		buf.SetString(
			h,
			termui.NewStyle(termui.Theme.Default.Fg, termui.ColorClear, termui.ModifierBold),
			image.Pt(cts.Inner.Min.X+colXPos[i]-1, cts.Inner.Min.Y),
		)
	}

	if cts.TopRow < 0 {
		log.Printf("table widget TopRow value less than 0. TopRow: %v", cts.TopRow)
		return
	}

	// prints each row
	for rowNum := cts.TopRow; rowNum < cts.TopRow+cts.Inner.Dy()-1 && rowNum < len(cts.Rows); rowNum++ {
		row := cts.Rows[rowNum]
		y := (rowNum + 2) - cts.TopRow

		// prints cursor
		style := termui.NewStyle(termui.Theme.Default.Fg)
		if cts.ShowCursor {
			if (cts.SelectedItem == "" && rowNum == cts.SelectedRow) ||
				(cts.SelectedItem != "" && cts.SelectedItem == row[cts.UniqueCol]) {
				style.Fg = cts.CursorColor
				style.Modifier = termui.ModifierReverse
				for _, width := range cts.ColWidths {
					if width == 0 {
						continue
					}
					buf.SetString(
						strings.Repeat(" ", cts.Inner.Dx()),
						style,
						image.Pt(cts.Inner.Min.X, cts.Inner.Min.Y+y-1),
					)
				}
				cts.SelectedItem = row[cts.UniqueCol]
				cts.SelectedRow = rowNum
			}
		}

		// prints each col of the row
		for i, width := range cts.ColWidths {
			if width == 0 {
				continue
			}
			// don't render column if width is greater than distance to end of widget
			if width > (cts.Inner.Dx()-colXPos[i])+1 {
				continue
			}
			r := termui.TrimString(row[i], width)
			buf.SetString(
				r,
				style,
				image.Pt(cts.Inner.Min.X+colXPos[i]-1, cts.Inner.Min.Y+y-1),
			)
		}
	}
}

func (cts *Table) drawLocation(buf *termui.Buffer) {
	total := len(cts.Rows)
	topRow := cts.TopRow + 1
	bottomRow := cts.TopRow + cts.Inner.Dy() - 1
	if bottomRow > total {
		bottomRow = total
	}

	loc := fmt.Sprintf(" %d - %d of %d ", topRow, bottomRow, total)

	width := len(loc)
	buf.SetString(loc, cts.TitleStyle, image.Pt(cts.Max.X-width-2, cts.Min.Y))
}

// GetCurIndex a
func (cts *Table) GetCurIndex() int {
	return cts.SelectedRow
}

// GetAllData a
func (cts *Table) GetAllData() [][]string {

	return cts.Rows
}

// GetCurRow a
func (cts *Table) GetCurRow() []string {
	return cts.Rows[cts.GetCurIndex()]
}

// calcPos is used to calculate the cursor position and the current view into the table.
func (cts *Table) calcPos() {
	cts.SelectedItem = ""

	if cts.SelectedRow < 0 {
		cts.SelectedRow = 0
	}
	if cts.SelectedRow < cts.TopRow {
		cts.TopRow = cts.SelectedRow
	}

	if cts.SelectedRow > len(cts.Rows)-1 {
		cts.SelectedRow = len(cts.Rows) - 1
	}
	if cts.SelectedRow > cts.TopRow+(cts.Inner.Dy()-2) {
		cts.TopRow = cts.SelectedRow - (cts.Inner.Dy() - 2)
	}
}

// ScrollUp a
func (cts *Table) ScrollUp() {
	cts.SelectedRow--
	cts.calcPos()
}

// ScrollDown a
func (cts *Table) ScrollDown() {
	cts.SelectedRow++
	cts.calcPos()
}

// ScrollTop a
func (cts *Table) ScrollTop() {
	cts.SelectedRow = 0
	cts.calcPos()
}

// ScrollBottom a
func (cts *Table) ScrollBottom() {
	cts.SelectedRow = len(cts.Rows) - 1
	cts.calcPos()
}

// ScrollHalfPageUp a
func (cts *Table) ScrollHalfPageUp() {
	cts.SelectedRow = cts.SelectedRow - (cts.Inner.Dy()-2)/2
	cts.calcPos()
}

// ScrollHalfPageDown a
func (cts *Table) ScrollHalfPageDown() {
	cts.SelectedRow = cts.SelectedRow + (cts.Inner.Dy()-2)/2
	cts.calcPos()
}

// ScrollPageUp a
func (cts *Table) ScrollPageUp() {
	cts.SelectedRow -= cts.Inner.Dy() - 2
	cts.calcPos()
}

// ScrollPageDown a
func (cts *Table) ScrollPageDown() {
	cts.SelectedRow += cts.Inner.Dy() - 2
	cts.calcPos()
}

// HandleClick a
func (cts *Table) HandleClick(x, y int) {
	x = x - cts.Min.X
	y = y - cts.Min.Y
	if (x > 0 && x <= cts.Inner.Dx()) && (y > 0 && y <= cts.Inner.Dy()) {
		cts.SelectedRow = (cts.TopRow + y) - 2
		cts.calcPos()
	}
}
