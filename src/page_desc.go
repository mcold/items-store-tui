package main

import (
	"github.com/atotto/clipboard"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type pageDescType struct {
	desc     *tview.TextArea
	btnFrm   *tview.Form
	mIdDesc  map[int]string
	mIdTrans map[int]string
	*tview.Flex
}

var pageDesc pageDescType

func (pageDesc *pageDescType) build() {

	pageDesc.desc = tview.NewTextArea()
	pageDesc.desc.Box.SetBorder(true).
		SetTitle("F5").
		SetTitleAlign(tview.AlignLeft).
		SetBorderPadding(1, 0, 1, 0)

	pageDesc.mIdDesc = make(map[int]string)

	pageDesc.btnFrm = tview.NewForm().
		AddButton("Save", func() { save() }).
		AddButton("Copy", func() { copyDescr() })

	pageDesc.Flex = tview.NewFlex().
		AddItem(pageDesc.desc, 0, 15, false).
		AddItem(pageDesc.btnFrm, 0, 1, false)

	pageDesc.desc.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 's' && event.Modifiers() == tcell.ModAlt {
			save()
			return nil
		}
		if event.Rune() == 'v' && event.Modifiers() == tcell.ModAlt {

			clipBoardContent, err := clipboard.ReadAll()
			check(err)
			pageDesc.desc.SetText(pageDesc.desc.GetText()+"\n"+clipBoardContent, true)
			return nil
		}
		return event
	})

	pageDesc.btnFrm.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyUp {
			app.SetFocus(pageDesc.desc)
			return nil
		}
		return event
	})

	pageDesc.Flex.SetDirection(tview.FlexRow)

	pageInfo.pages.AddPage("desc", pageDesc.Flex, true, false)
}
