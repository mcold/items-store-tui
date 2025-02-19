package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"log"
	"strconv"
	"strings"
)

type pageCaseInsType struct {
	caseArea    *tview.TextArea
	commentArea *tview.TextArea
	frmSave     *tview.Form
	*tview.Flex
}

var pageCaseIns pageCaseInsType

func (pageCaseIns *pageCaseInsType) build() {

	pageCaseIns.caseArea = tview.NewTextArea()
	pageCaseIns.caseArea.SetBorder(true).
		SetTitleAlign(tview.AlignLeft).
		SetBackgroundColor(tcell.ColorDarkBlue).
		SetBorderPadding(1, 1, 1, 1)

	pageCaseIns.commentArea = tview.NewTextArea()

	pageCaseIns.commentArea.SetBorder(true).
		SetTitleAlign(tview.AlignLeft).
		SetBackgroundColor(tcell.ColorDarkBlue).
		SetBorderPadding(1, 1, 1, 1)

	pageCaseIns.frmSave = tview.NewForm().AddButton("Save", func() {
		saveCase()
		pageInfo.pages.SwitchToPage("case")
		setCases()
		app.SetFocus(pageCase.caseList)
		pageCase.caseList.SetCurrentItem(0)
	})

	pageCaseIns.caseArea.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyTab {
			app.SetFocus(pageCaseIns.commentArea)
			return nil
		}
		return event
	})

	pageCaseIns.commentArea.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyTab {
			app.SetFocus(pageCaseIns.frmSave)
			return nil
		}
		return event
	})

	pageCaseIns.Flex = tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(pageCaseIns.caseArea, 6, 1, true).
		AddItem(pageCaseIns.commentArea, 0, 5, false).
		AddItem(pageCaseIns.frmSave, 0, 1, false)

	pageInfo.pages.AddPage("caseIns", pageCaseIns.Flex, true, false)
}

func saveCase() {
	log.Println("-------------------------------")
	log.Println("saveCase")
	log.Println("--------------------")

	if len(strings.Trim(pageCaseIns.caseArea.GetText(), "")) > 0 {
		query := "INSERT INTO cases(id_item, use_case, comment)" + "\n" +
			"VALUES( " + strings.ReplaceAll(strconv.Itoa(pageItem.mPosId[pageItem.items.GetCurrentItem()]), "'", "''") + ", '" + strings.ReplaceAll(pageCaseIns.caseArea.GetText(), "'", "''") + "'," +
			"'" + strings.ReplaceAll(pageCaseIns.commentArea.GetText(), "'", "''") + "')" +
			" ON CONFLICT(use_case)" + "\n" +
			" DO UPDATE SET comment='" + strings.ReplaceAll(pageCaseIns.commentArea.GetText(), "'", "''") + "'"

		log.Println(query)
		_, err := database.Exec(query)
		check(err)
	}

	log.Println("-------------------------------")
}
