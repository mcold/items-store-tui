package main

import (
	"database/sql"
	"log"
	"strconv"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type pageItemType struct {
	//descrs   *tview.TextArea
	//btnFrm   *tview.Form
	//mIdDescr map[int]string
	mIdTrans   map[int]string
	items      *tview.List
	filterFrm  *tview.Form
	itemArea   *tview.TextArea
	mPosId     map[int]int
	curItemPos int
}

var pageItem pageItemType

func (pageItem *pageItemType) build() {

	pageItem.filterFrm = tview.NewForm().
		AddInputField("", "", 60, nil, nil)

	pageItem.filterFrm.Box.SetBorder(true)
	pageItem.filterFrm.Box.SetTitle("F2").
		SetTitleAlign(tview.AlignLeft)

	pageItem.items = tview.NewList()
	pageItem.items.Box.SetBorder(true).
		SetTitle("F3").
		SetTitleAlign(tview.AlignLeft).
		SetBorderPadding(1, 0, 1, 0)

	pageItem.itemArea = tview.NewTextArea()
	pageItem.itemArea.SetBorder(true).
		SetTitle("F4").
		SetTitleAlign(tview.AlignLeft).
		SetBackgroundColor(tcell.ColorDarkBlue).
		SetBorderPadding(1, 1, 1, 1)

	pageItem.filterFrm.GetFormItem(0).(*tview.InputField).SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyTab {
			app.SetFocus(pageItem.filterFrm.GetFormItem(1).(*tview.InputField))
			return nil
		}
		return event
	})

	pageItem.filterFrm.GetFormItem(0).(*tview.InputField).SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyTab {
			app.SetFocus(pageItem.filterFrm.GetFormItem(0).(*tview.InputField))
			return nil
		}
		return event
	})

	flexItem := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(pageItem.filterFrm, 5, 0, true).
		AddItem(pageItem.items, 0, 10, false).
		AddItem(pageItem.itemArea, 6, 0, false)

	pageItem.mIdTrans = make(map[int]string)
	pageItem.mPosId = make(map[int]int)
	err := database.Connect()
	if err != nil {
		return
	} else {
		query := "SELECT id, name, trans, descr FROM item order by name"
		items, err := database.Query(query)
		check(err)

		rowCount := 1
		for items.Next() {
			var id sql.NullInt64
			var item, trans, desc sql.NullString
			err := items.Scan(&id, &item, &trans, &desc)
			check(err)

			pageItem.items.AddItem(item.String, trans.String, rune(0), func() {})

			pageDesc.mIdDesc[int(id.Int64)] = desc.String
			pageItem.mPosId[rowCount-1] = int(id.Int64)
			rowCount++
		}
	}

	flexCmplx := tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(flexItem, 0, 3, true).
		AddItem(pageInfo.pages, 0, 8, false).
		SetFullScreen(true)

	pageItem.filterFrm.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEnter {
			app.SetFocus(pageItem.items)
			refreshItemList()
			if pageItem.items.GetItemCount() > 0 {
				pageItem.items.SetCurrentItem(0)
				pageDesc.desc.SetText(pageDesc.mIdDesc[pageItem.mPosId[pageItem.items.GetCurrentItem()]]+"\n", true)
			}
			return nil
		}

		return event
	})

	pageItem.items.SetFocusFunc(func() {
		itemText, _ := pageItem.items.GetItemText(pageItem.items.GetCurrentItem())
		pageItem.itemArea.SetText(itemText, true)
		pageDesc.desc.SetText(pageDesc.mIdDesc[pageItem.mPosId[pageItem.items.GetCurrentItem()]], false)
		pageInfo.pages.SwitchToPage("desc")
		pageCase.caseList.Clear()
		pageCase.descCaseArea.SetText("", true)
	})

	pageItem.items.SetSelectedBackgroundColor(tcell.ColorGreen)

	flexCmplx.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyF1 {
			application.pages.ShowPage("help")
			return nil
		}
		if event.Key() == tcell.KeyF2 {
			app.SetFocus(pageItem.filterFrm)
			return nil
		}
		if event.Key() == tcell.KeyF3 {
			app.SetFocus(pageItem.items)
			pageInfo.pages.SwitchToPage("desc")
			return nil
		}
		if event.Key() == tcell.KeyF4 {
			app.SetFocus(pageItem.itemArea)
			return nil
		}
		if event.Key() == tcell.KeyF5 {
			pageInfo.pages.SwitchToPage("desc")
			app.SetFocus(pageDesc.desc)
			return nil
		}
		if event.Key() == tcell.KeyF6 {
			setCases()
			pageCase.caseList.SetCurrentItem(0)
			pageInfo.pages.SwitchToPage("case")
			app.SetFocus(pageCase.caseList)

			return nil
		}
		if event.Key() == tcell.KeyInsert && event.Modifiers() == tcell.ModCtrl {
			application.pages.SwitchToPage("new item")
			return nil
		}
		if event.Key() == tcell.KeyEsc {
			application.ConfirmQuit()
		}
		return event
	})

	pageItem.items.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {

		if event.Key() == tcell.KeyEnter {
			app.SetFocus(pageDesc.desc)
			return nil
		}
		if event.Key() == tcell.KeyDelete {
			deleteItem()
			refreshItemList()

			return nil
		}
		return event
	})

	flexItem.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyTab {
			curId := pageItem.items.GetCurrentItem() + 1
			cnt := pageItem.items.GetItemCount()
			if curId < cnt {
				pageItem.items.SetCurrentItem(curId)
				pageDesc.desc.SetText(pageDesc.mIdDesc[pageItem.mPosId[pageItem.items.GetCurrentItem()]]+"\n", true)
			}

			return nil
		}

		if event.Key() == tcell.KeyDown {
			curId := pageItem.items.GetCurrentItem() + 1
			cnt := pageItem.items.GetItemCount()
			if curId < cnt {
				pageItem.items.SetCurrentItem(curId)
				itemText, _ := pageItem.items.GetItemText(pageItem.items.GetCurrentItem())
				pageItem.itemArea.SetText(itemText, true)
				pageDesc.desc.SetText(pageDesc.mIdDesc[pageItem.mPosId[pageItem.items.GetCurrentItem()]]+"\n", true)
			}

			return nil
		}
		if event.Key() == tcell.KeyBacktab {
			curId := pageItem.items.GetCurrentItem() - 1
			if curId > -1 {
				pageItem.items.SetCurrentItem(curId)
				pageDesc.desc.SetText(pageDesc.mIdDesc[pageItem.mPosId[pageItem.items.GetCurrentItem()]]+"\n", true)
			}

			return nil
		}
		if event.Key() == tcell.KeyUp {
			curId := pageItem.items.GetCurrentItem() - 1
			if curId > -1 {
				pageItem.items.SetCurrentItem(curId)
				itemText, _ := pageItem.items.GetItemText(pageItem.items.GetCurrentItem())
				pageItem.itemArea.SetText(itemText, true)
				pageDesc.desc.SetText(pageDesc.mIdDesc[pageItem.mPosId[pageItem.items.GetCurrentItem()]]+"\n", true)
			}

			return nil
		}
		if event.Key() == tcell.KeyRight {
			app.SetFocus(pageDesc.desc)
		}
		if event.Key() == tcell.KeyLeft {
			app.SetFocus(pageItem.filterFrm)
		}

		return event
	})

	pageItem.filterFrm.SetFocus(1)

	application.pages.AddPage("items", flexCmplx, true, true)
}

func refreshItemList() {
	log.Println("-------------------------------")
	log.Println("refreshItemList")
	log.Println("--------------------")

	pageItem.curItemPos = pageItem.items.GetCurrentItem()

	itemToken := strings.TrimSpace(pageItem.filterFrm.GetFormItem(0).(*tview.InputField).GetText())
	var query string

	log.Println("len itemToken: ", len(itemToken))

	if len(itemToken) > 0 {
		query = "SELECT id, name, trans, descr" +
			" FROM item" +
			" WHERE lower(name) like lower('%" + itemToken + "%')" +
			" or lower(descr) like lower('%" + itemToken + "%')" +
			" order by name"
	} else {
		query = "SELECT id, name, trans, descr" +
			" FROM item" +
			" order by name"
	}

	log.Println(query)

	itemFind, err := database.Query(query)
	check(err)

	pageDesc.mIdDesc = make(map[int]string)
	pageItem.mPosId = make(map[int]int)
	pageItem.items.Clear()
	rowCount := 1
	for itemFind.Next() {

		var id sql.NullInt64
		var item, trans, descr sql.NullString
		err := itemFind.Scan(&id, &item, &trans, &descr)
		check(err)

		pageItem.items.AddItem(item.String, trans.String, rune(0), func() {})

		pageItem.mIdTrans[int(id.Int64)] = trans.String
		pageDesc.mIdDesc[int(id.Int64)] = descr.String
		pageItem.mPosId[rowCount-1] = int(id.Int64)
		rowCount++
	}

	pageItem.items.SetCurrentItem(pageItem.curItemPos)

	log.Println("-------------------------------")
}

func save() {
	log.Println("-------------------------------")
	log.Println("save")
	log.Println("--------------------")

	query := "UPDATE item" + "\n" +
		"SET name = '" + pageItem.itemArea.GetText() + "',\n" +
		"descr = '" + pageDesc.desc.GetText() + "'\n" +
		"WHERE id = " + strconv.Itoa(pageItem.mPosId[pageItem.items.GetCurrentItem()])

	log.Println(query)

	pageDesc.mIdDesc[pageItem.mPosId[pageItem.items.GetCurrentItem()]] = pageDesc.desc.GetText()

	_, err := database.Exec(query)
	check(err)

	refreshItemList()
	log.Println("-------------------------------")
}

func copyDescr() {
	err := clipboard.WriteAll(pageDesc.desc.GetText())
	check(err)
}

func deleteItem() {

	queryCases := "DELETE FROM cases WHERE id_item = " + strconv.Itoa(pageItem.mPosId[pageItem.items.GetCurrentItem()])
	_, err := database.Exec(queryCases)
	check(err)

	queryItem := "DELETE FROM item WHERE id = " + strconv.Itoa(pageItem.mPosId[pageItem.items.GetCurrentItem()])
	_, err = database.Exec(queryItem)
	check(err)
}
