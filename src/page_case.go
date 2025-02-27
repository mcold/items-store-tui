package main

import (
	"database/sql"
	"github.com/atotto/clipboard"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"log"
	"strconv"
)

type pageCaseType struct {
	caseList     *tview.List
	descCaseArea *tview.TextArea
	mPosCaseID   map[int]int
	*tview.Flex
}

var pageCase pageCaseType

func (pageCase *pageCaseType) build() {

	pageCase.caseList = tview.NewList()
	pageCase.mPosCaseID = make(map[int]int)

	pageCase.caseList.Box.SetBorder(true).
		SetTitle("F6").
		SetTitleAlign(tview.AlignLeft).
		SetBorderPadding(1, 0, 1, 0)

	pageCase.descCaseArea = tview.NewTextArea()
	pageCase.descCaseArea.SetBorder(true).
		SetTitleAlign(tview.AlignLeft).
		SetBackgroundColor(tcell.ColorDarkBlue).
		SetBorderPadding(1, 1, 1, 1).
		SetTitle("F7").
		SetTitleAlign(tview.AlignLeft)

	pageCase.caseList.SetSelectedFunc(func(pos int, _ string, _ string, _ rune) {
		setCaseComment()
	})

	pageCase.Flex = tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(pageCase.caseList, 0, 1, true).
		AddItem(pageCase.descCaseArea, 0, 1, false)

	pageCase.Flex.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {

		if event.Key() == tcell.KeyInsert || event.Key() == tcell.KeyF8 {

			clearCaseInsFields()
			pageInfo.pages.SwitchToPage("caseIns")
			app.SetFocus(pageCaseIns.caseArea)
			return nil
		}
		if event.Key() == tcell.KeyDelete {
			deleteCase()
			setCases()
			return nil
		}
		if event.Key() == tcell.KeyF7 {
			app.SetFocus(pageCase.descCaseArea)
		}
		return event
	})

	pageCase.caseList.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {

		if event.Key() == tcell.KeyInsert || event.Key() == tcell.KeyF8 {

			clearCaseInsFields()
			pageInfo.pages.SwitchToPage("caseIns")
			app.SetFocus(pageCaseIns.caseArea)
			return nil
		}
		if event.Key() == tcell.KeyDelete {
			deleteCase()
			setCases()
			return nil
		}
		if event.Rune() == 'c' && event.Modifiers() == tcell.ModAlt {
			caseText, _ := pageCase.caseList.GetItemText(pageCase.caseList.GetCurrentItem())
			err := clipboard.WriteAll(caseText)
			check(err)
		}
		return event
	})

	pageCase.descCaseArea.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {

		if event.Rune() == 'v' && event.Modifiers() == tcell.ModAlt {

			clipBoardContent, err := clipboard.ReadAll()
			check(err)
			pageCase.descCaseArea.SetText(pageCase.descCaseArea.GetText()+"\n"+clipBoardContent, true)
			return nil
		}
		return event
	})

	pageInfo.pages.AddPage("case", pageCase.Flex, true, false)
}

func setCases() {
	log.Println("-------------------------------")
	log.Println("setCases")
	log.Println("--------------------")
	curItemID := pageItem.mPosId[pageItem.items.GetCurrentItem()]

	if curItemID > 0 {
		query := `SELECT id, use_case, comment
					FROM cases 
				   WHERE id_item = ` + strconv.Itoa(curItemID) +
			` ORDER BY id ASC`

		pageCase.caseList.Clear()

		caseFind, err := database.Query(query)
		check(err)

		rowCount := 0
		for caseFind.Next() {

			var id sql.NullInt64
			var useCase, comment sql.NullString

			err := caseFind.Scan(&id, &useCase, &comment)
			check(err)

			pageCase.caseList.AddItem(useCase.String, comment.String, rune(0), func() {})
			pageCase.descCaseArea.SetText(comment.String, true)
			pageCase.mPosCaseID[rowCount] = int(id.Int64)
			rowCount++
		}

		pageCase.caseList.SetCurrentItem(0)

	}

	log.Println("-------------------------------")
}

func setCaseComment() {

	log.Println("-------------------------------")
	log.Println("setCaseComment")
	log.Println("--------------------")

	log.Println("cur case id: " + strconv.Itoa(pageCase.caseList.GetCurrentItem()))
	curCaseID := pageCase.mPosCaseID[pageCase.caseList.GetCurrentItem()]
	query := `SELECT comment
					FROM cases 
				   WHERE id = ` + strconv.Itoa(curCaseID)

	log.Println(query)

	caseFind, err := database.Query(query)
	check(err)

	for caseFind.Next() {

		var comment sql.NullString

		err := caseFind.Scan(&comment)
		check(err)

		pageCase.descCaseArea.SetText(comment.String, true)
		log.Println(comment.String)

		return
	}

	log.Println("-------------------------------")

}

func deleteCase() {
	query := "DELETE FROM cases WHERE id = " + strconv.Itoa(pageCase.mPosCaseID[pageCase.caseList.GetCurrentItem()])
	log.Println(query)

	_, err := database.Exec(query)
	check(err)
}

func clearCaseInsFields() {
	pageCaseIns.caseArea.SetText("", true)
	pageCaseIns.commentArea.SetText("", true)
}
