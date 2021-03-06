package paint

import (
	"github.com/nsf/termbox-go"

	"github.com/sebashwa/vixl44/modes"
	"github.com/sebashwa/vixl44/state"
)

func AdjustColor(diff int) {
	newIndex := int(state.SelectedColor) + diff

	if newIndex < 1 {
		state.SelectedColor = 256
	} else if newIndex > 256 {
		state.SelectedColor = 1
	} else {
		state.SelectedColor = termbox.Attribute(newIndex)
	}
}

func SelectColor() {
	position := state.Cursor.Position

	if state.CurrentMode == modes.PaletteMode {
		state.SelectedColor = state.Palette[position.X][position.Y]
	} else {
		state.SelectedColor = state.Canvas.Values[position.X][position.Y]
	}
}

func fillPixel(color termbox.Attribute) {
	position := state.Cursor.Position

	state.Canvas.Values[position.X][position.Y] = color
	state.Canvas.Values[position.X+1][position.Y] = color

	state.History.AddCanvasState(state.Canvas.GetValuesCopy())
}

func FillPixel() {
	fillPixel(state.SelectedColor)
}

func KillPixel() {
	fillPixel(termbox.ColorDefault)
}

func fillArea(color termbox.Attribute) {
	xMin, xMax, yMin, yMax := state.Cursor.GetVisualModeArea()

	for x := xMin; x <= xMax+1; x++ {
		for y := yMin; y <= yMax; y++ {
			state.Canvas.Values[x][y] = color
		}
	}

	state.History.AddCanvasState(state.Canvas.GetValuesCopy())
}

func FillArea() {
	fillArea(state.SelectedColor)
}

func KillArea() {
	fillArea(termbox.ColorDefault)
}

func floodFill(x, y int, targetColor termbox.Attribute, replacementColor termbox.Attribute) {
	if targetColor == replacementColor {
		return
	}
	if state.Canvas.Values[x][y] != targetColor {
		return
	}

	state.Canvas.Values[x][y] = replacementColor

	if x > 0 {
		floodFill(x-1, y, targetColor, replacementColor)
	}
	if x < state.Canvas.Columns-1 {
		floodFill(x+1, y, targetColor, replacementColor)
	}
	if y > 0 {
		floodFill(x, y-1, targetColor, replacementColor)
	}
	if y < state.Canvas.Rows-1 {
		floodFill(x, y+1, targetColor, replacementColor)
	}
}

func FloodFill() {
	x := state.Cursor.Position.X
	y := state.Cursor.Position.Y

	floodFill(x, y, state.Canvas.Values[x][y], state.SelectedColor)
	state.History.AddCanvasState(state.Canvas.GetValuesCopy())
}
