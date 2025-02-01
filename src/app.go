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
	pageIns.build()
	pageMainMessage.build()
	pageConfirm.build()
	application.registerGlobalShortcuts()

	if err := app.SetRoot(application.pages, true).EnableMouse(true).EnablePaste(true).Run(); err != nil {
		panic(err)
	}
}

func (application *applicationType) registerGlobalShortcuts() {
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyCtrlC:
			application.ConfirmQuit()
		default:
			return event
		}
		return nil
	})
}

func (application *applicationType) ConfirmQuit() {
	pageConfirm.show("Are you sure you want to exit?", application.Quit)
}

func (application *applicationType) Quit() {
	if database.DB != nil {
		database.DB.Close()
	}
	app.Stop()
}
