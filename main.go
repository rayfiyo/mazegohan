package main

import (
	"fmt"
	"github.com/rivo/tview"
)

func main() {
	fmt.Println("Let's play")
	
	app := tview.NewApplication()

	a := tview.NewTextView()
	a.SetText("textarea(a)")

	flex := tview.NewFlex().
		AddItem(a, 0, 1, false)

	if err := app.SetRoot(flex, true).Run(); err != nil {
		panic(err)
	}
}
