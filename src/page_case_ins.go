package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
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
	})

	pageCaseIns.Flex = tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(pageCaseIns.caseArea, 6, 1, true).
		AddItem(pageCaseIns.commentArea, 0, 5, false).
		AddItem(pageCaseIns.frmSave, 0, 1, false)

	pageInfo.pages.AddPage("caseIns", pageCaseIns.Flex, true, false)
}

func saveCase() {

	if len(strings.Trim(pageCaseIns.caseArea.GetText(), "")) > 0 {
		query := "INSERT INTO case(use_case, comment)" + "\n" +
			"VALUES( '" + pageCaseIns.caseArea.GetText() + "'," +
			"'" + pageCaseIns.commentArea.GetText() + "')" +
			"ON CONFLICT(use_case)" + "\n" +
			"DO UPDATE SET comment='" + pageCaseIns.commentArea.GetText() + "'"

		_, err := database.Exec(query)
		check(err)
	}
}
