package tui

import "github.com/rivo/tview"

func CreateLayout() *tview.Grid {
	newPrimitive := func(text string) tview.Primitive {
		return tview.NewTextView().
			SetTextAlign(tview.AlignCenter).
			SetText(text)
	}
	// menu := newPrimitive("Menu")
	main := newPrimitive("Main content")
	sideBar := newPrimitive("Side Bar")

	grid := tview.NewGrid().
		SetRows(3, 0, 3).
		SetColumns(0, 30).
		SetBorders(true).
		AddItem(
			newPrimitive("Header"), 
			0, 0, 1, 3, 0, 0, false).
		AddItem(
			newPrimitive("Terminal"), 
			2, 0, 1, 3, 0, 0, false).
		AddItem(
			main, 
			1, 0, 1, 2, 0, 0, false).
		AddItem(
			sideBar,
			1, 2, 1, 1, 0, 0, false)

	// // Layout for screens wider than 100 cells.
	// grid.AddItem(main, 1, 0, 1, 1, 0, 100, false).
	// 	AddItem(sideBar, 1, 2, 1, 1, 0, 100, false)

	grid.SetTitle("dag-doc")
	return grid
}
