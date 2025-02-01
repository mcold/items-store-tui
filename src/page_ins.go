package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"strings"
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
	pageIns.descr.SetTitle("DESCRIPTION")

	frmSave := tview.NewForm().AddButton("Save", func() {
		saveCmd()
	})

	pageIns.cmd.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyTab {
			app.SetFocus(pageIns.descr)
			return nil
		}
		return event
	})

	pageIns.descr.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyTab {
			app.SetFocus(frmSave)
			return nil
		}
		if event.Key() == tcell.KeyBacktab {
			app.SetFocus(pageIns.cmd)
			return nil
		}
		return event
	})

	frmSave.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyTab {
			app.SetFocus(pageIns.cmd)
			return nil
		}
		if event.Key() == tcell.KeyBacktab {
			app.SetFocus(pageIns.descr)
			return nil
		}
		if event.Key() == tcell.KeyEnter {
			saveCmd()
		}
		return event
	})

	flexIns := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(pageIns.cmd, 0, 1, true).
		AddItem(pageIns.descr, 0, 12, true).
		AddItem(frmSave, 0, 1, false)

	flexIns.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 's' && event.Modifiers() == tcell.ModAlt {
			saveCmd()
			return nil
		}
		if event.Key() == tcell.KeyEsc {
			application.pages.SwitchToPage("commands")
			return nil
		}

		return event
	})

	application.pages.AddPage("new command", flexIns, true, false)
}

func saveCmd() {
	saveCmdDB()
	pageIns.cmd.SetText("", true)
	pageIns.descr.SetText("", true)
	refreshCmdList()
	application.pages.SwitchToPage("commands")
	app.SetFocus(pageCmd.cmds)
}

func saveCmdDB() {
	if len(strings.Trim(pageIns.cmd.GetText(), "")) > 0 {
		query := "INSERT INTO cmd(command, descr)" + "\n" +
			"VALUES( '" + pageIns.cmd.GetText() + "'," +
			"'" + pageIns.descr.GetText() + "')" + "\n" +
			"ON CONFLICT(command)" + "\n" +
			"DO UPDATE SET descr='" + pageIns.descr.GetText() + "'"

		_, err := database.Exec(query)
		check(err)
	}
}
