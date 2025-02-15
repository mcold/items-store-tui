package main

import (
	"log"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type pageInsType struct {
	item  *tview.TextArea
	trans *tview.TextArea
	descr *tview.TextArea
}

var pageIns pageInsType

func (pageIns *pageInsType) build() {
	pageIns.item = tview.NewTextArea()
	pageIns.item.SetBorder(true)
	pageIns.item.SetTitle("ITEM")
	pageIns.item.SetBorderPadding(1, 1, 1, 0)

	pageIns.trans = tview.NewTextArea()
	pageIns.trans.SetBorder(true)
	pageIns.trans.SetTitle("TRANSCRIPTION")
	pageIns.trans.SetBorderPadding(1, 1, 1, 0)

	pageIns.descr = tview.NewTextArea()
	pageIns.descr.SetBorder(true)
	pageIns.descr.SetTitle("DESCRIPTION")
	pageIns.descr.SetBorderPadding(1, 0, 1, 0)

	frmSave := tview.NewForm().AddButton("Save", func() {
		saveItem()
	})

	pageIns.item.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyTab {
			app.SetFocus(pageIns.trans)
			return nil
		}
		return event
	})

	pageIns.trans.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
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
			app.SetFocus(pageIns.item)
			return nil
		}
		return event
	})

	frmSave.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyTab {
			app.SetFocus(pageIns.item)
			return nil
		}
		if event.Key() == tcell.KeyBacktab {
			app.SetFocus(pageIns.descr)
			return nil
		}
		if event.Key() == tcell.KeyEnter {
			saveItem()
		}
		return event
	})

	flexIns := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(pageIns.item, 0, 1, true).
		AddItem(pageIns.trans, 0, 1, true).
		AddItem(pageIns.descr, 0, 6, true).
		AddItem(frmSave, 0, 1, false)

	flexIns.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 's' && event.Modifiers() == tcell.ModAlt {
			saveItem()
			return nil
		}
		if event.Key() == tcell.KeyEsc {
			application.pages.SwitchToPage("items")
			return nil
		}

		return event
	})

	application.pages.AddPage("new item", flexIns, true, false)
}

func saveItem() {
	saveItemDB()
	pageIns.item.SetText("", true)
	pageIns.descr.SetText("", true)
	refreshItemList()
	application.pages.SwitchToPage("items")
	app.SetFocus(pageItem.items)
}

func saveItemDB() {
	log.Println("-------------------------------")
	log.Println("saveItemDB")
	log.Println("--------------------")

	log.Println("itemName: " + pageIns.item.GetText())
	log.Println("trans: " + pageIns.trans.GetText())
	log.Println("itemDesc: " + pageIns.item.GetText())

	log.Println(len(strings.Trim(pageIns.item.GetText(), "")))

	if len(strings.Trim(pageIns.item.GetText(), "")) > 0 {
		query := "INSERT INTO item(name, trans, descr)" + "\n" +
			"VALUES( '" + pageIns.item.GetText() + "'," +
			"'" + pageIns.trans.GetText() + "'" + "," +
			"'" + pageIns.descr.GetText() + "')" + "\n" +
			"ON CONFLICT(name)" + "\n" +
			"DO UPDATE SET descr='" + pageIns.descr.GetText() + "'"

		log.Println(query)

		_, err := database.Exec(query)
		check(err)
	}

	clearInsFields()

	log.Println("-------------------------------")
}

func clearInsFields() {
	pageIns.item.SetText("", true)
	pageIns.trans.SetText("", true)
	pageIns.descr.SetText("", true)
}
