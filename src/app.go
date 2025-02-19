package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"log"
	"os"
)

type applicationType struct {
	pages         *tview.Pages
	ListShortcuts []rune
}

var app *tview.Application

func (application *applicationType) init() {
	file, err := os.OpenFile("app.log", os.O_TRUNC|os.O_CREATE, 0666)
	check(err)
	log.SetOutput(file)

	app = tview.NewApplication()

	application.pages = tview.NewPages()
	pageInfo.build()
	pageItem.build()

	_, _, w, _ := pageItem.filterFrm.GetInnerRect()
	pageItem.filterFrm.GetFormItem(0).(*tview.InputField).SetFieldWidth(w * 2)

	pageIns.build()
	pageMainMessage.build()
	pageConfirm.build()
	application.registerGlobalShortcuts()
	app.SetFocus(pageItem.filterFrm)

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
