package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type pageInsType struct {
	cmd   *tview.TextArea
	descr *tview.TextArea
}

var pageIns pageInsType

func (pageIns *pageInsType) build() {
	pageIns.cmd = tview.NewTextArea()
	pageIns.cmd.SetBorder(true)
	pageIns.cmd.SetTitle("COMMAND")

	pageIns.descr = tview.NewTextArea()
	pageIns.descr.SetBorder(true)
	pageIns.descr.SetTitle("DESCR")

	frmSave := tview.NewForm().AddButton("Save", func() {
		saveCmd()
		application.pages.ShowPage("commands")
	})

	flexIns := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(pageIns.cmd, 0, 1, true).
		AddItem(pageIns.descr, 0, 10, true).
		AddItem(frmSave, 0, 2, false)

	flexIns.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 's' && event.Modifiers() == tcell.ModAlt {
			saveCmd()
			application.pages.SwitchToPage("commands")
			// TODO: renew command list
			return nil
		}
		if event.Key() == tcell.KeyEsc {
			pageIns.cmd.SetText("", true)
			pageIns.descr.SetText("", true)
			application.pages.SwitchToPage("commands")
			// TODO: renew command list
			return nil
		}

		return event
	})

	application.pages.AddPage("new command", flexIns, true, false)
}

func saveCmd() {
	// TODO: upsert
	query := "INSERT INTO cmd(command, descr)" + "\n" +
		"VALUES( '" + pageIns.cmd.GetText() + "'," +
		"'" + pageIns.descr.GetText() + "')"

	_, err := database.Exec(query)
	check(err)
}
