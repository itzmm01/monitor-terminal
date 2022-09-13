package ui

const (
	Bold int = 1 << (iota + 9)
	Underline
	Reverse
)

// Colorscheme coolor
type Colorscheme struct {
	Name   string
	Author string

	Fg int
	Bg int

	BorderLabel int
	BorderLine  int

	// should add at least 8 here
	CPULines []int

	BattLines []int

	MainMem int
	SwapMem int

	ProcCursor int

	Sparkline int

	DiskBar int

	// colors the temperature number a different color if it's over a certain threshold
	TempLow  int
	TempHigh int
}

var Default1 = Colorscheme{
	Fg: 249,
	Bg: -1,

	BorderLabel: 249,
	BorderLine:  239,

	CPULines: []int{81, 70, 208, 197, 249, 141, 221, 186},

	BattLines: []int{81, 70, 208, 197, 249, 141, 221, 186},

	MainMem: 208,
	SwapMem: 186,

	ProcCursor: 197,

	Sparkline: 81,

	DiskBar: 102,

	TempLow:  70,
	TempHigh: 208,
}

var Default = Colorscheme{
	Fg: 250,
	Bg: -1,

	BorderLabel: 250,
	BorderLine:  37,

	CPULines: []int{61, 33, 37, 64, 125, 160, 166, 136},

	BattLines: []int{61, 33, 37, 64, 125, 160, 166, 136},

	MainMem: 125,
	SwapMem: 166,

	ProcCursor: 136,

	Sparkline: 33,

	DiskBar: 245,

	TempLow:  64,
	TempHigh: 160,
}
