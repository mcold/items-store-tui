package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type applicationType struct {
	pages         *tview.Pages
	ListShortcuts []rune
}

var app *tview.Application

func (application *applicationType) init() {
	app = tview.NewApplication()

	application.pages = tview.NewPages()
	pageCmd.build()
	application.registerGlobalShortcuts()

	if err := app.SetRoot(application.pages, true).EnableMouse(true).EnablePaste(true).Run(); err != nil {
		panic(err)
	}
}

func (application *applicationType) registerGlobalShortcuts() {
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyCtrlC:
			return nil
		//case tcell.KeyEsc:
		//	application.ConfirmQuit()
		default:
			return event
		}
		return nil
	})
}
