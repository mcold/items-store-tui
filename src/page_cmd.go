package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"strconv"
)

type pageCmdType struct {
	cmds      *tview.List
	descrs    *tview.TextArea
	filterFrm *tview.Form
	saveFrm   *tview.Form
	mIdDescr  map[int]string
	mPosId    map[int]int
}

var pageCmd pageCmdType

func (pageCmd *pageCmdType) build() {
	pageCmd.cmds = tview.NewList()
	pageCmd.cmds.Box.SetBorder(true)
	pageCmd.cmds.Box.SetTitle("command (alt+q)")

	pageCmd.descrs = tview.NewTextArea()
	pageCmd.descrs.Box.SetBorder(true)
	pageCmd.descrs.Box.SetTitle("description (alt+w)")

	pageCmd.filterFrm = tview.NewForm().
		AddInputField("COMMAND", "", 20, nil, nil)

	pageCmd.filterFrm.Box.SetBorder(true)
	pageCmd.filterFrm.Box.SetTitle("find (alt+f)")

	flexCmd := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(pageCmd.filterFrm, 0, 1, true).
		AddItem(pageCmd.cmds, 0, 10, false)

	pageCmd.mIdDescr = make(map[int]string)
	pageCmd.mPosId = make(map[int]int)
	err := database.Connect()
	if err != nil {
		return
	} else {
		query := "SELECT id, command, descr FROM cmd order by command"
		cmds, err := database.Query(query)
		check(err)

		rowCount := 1
		for cmds.Next() {
			id := 0
			cmd := ""
			descr := ""
			cmds.Scan(&id, &cmd, &descr)

			pageCmd.cmds.AddItem(cmd, "", rune(0), func() {})

			pageCmd.mIdDescr[id] = descr
			pageCmd.mPosId[rowCount-1] = id
			rowCount++
		}
	}

	pageCmd.saveFrm = tview.NewForm().AddButton("Save", func() {
		saveDescr()
	})

	flexDescr := tview.NewFlex().
		AddItem(pageCmd.descrs, 0, 10, false).
		AddItem(pageCmd.saveFrm, 0, 2, false)

	flexDescr.SetDirection(tview.FlexRow)

	flexCmplx := tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(flexCmd, 0, 1, true).
		AddItem(flexDescr, 0, 4, false).
		SetFullScreen(true)

	pageCmd.filterFrm.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEnter {
			app.SetFocus(pageCmd.cmds)
			refreshCmdList()
			if pageCmd.cmds.GetItemCount() > 0 {
				pageCmd.cmds.SetCurrentItem(0)
				pageCmd.descrs.SetText(pageCmd.mIdDescr[pageCmd.mPosId[pageCmd.cmds.GetCurrentItem()]]+"\n", true)
			}
			return nil
		}

		return event
	})

	pageCmd.cmds.SetFocusFunc(func() {
		pageCmd.descrs.SetText(pageCmd.mIdDescr[pageCmd.mPosId[pageCmd.cmds.GetCurrentItem()]], false)
		//pageCmd.descrs.SetText(strconv.Itoa(pageCmd.cmds.GetCurrentItem()), true)
	})

	pageCmd.cmds.SetSelectedBackgroundColor(tcell.ColorGreen)

	flexCmplx.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 'q' && event.Modifiers() == tcell.ModAlt {
			app.SetFocus(pageCmd.cmds)
			return nil
		}
		if event.Rune() == 'w' && event.Modifiers() == tcell.ModAlt {
			app.SetFocus(pageCmd.descrs)

			return nil
		}
		if event.Rune() == 'f' && event.Modifiers() == tcell.ModAlt {
			app.SetFocus(pageCmd.filterFrm)
			return nil
		}
		if event.Key() == tcell.KeyInsert && event.Modifiers() == tcell.ModCtrl {
			application.pages.SwitchToPage("new command")
			return nil
		}
		if event.Rune() == 'h' && event.Modifiers() == tcell.ModAlt {
			application.pages.ShowPage("help")
			return nil
		}
		if event.Key() == tcell.KeyEsc {
			application.ConfirmQuit()
		}
		return event
	})

	pageCmd.descrs.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 's' && event.Modifiers() == tcell.ModAlt {
			saveDescr()
			return nil
		}
		return event
	})

	pageCmd.saveFrm.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyUp {
			app.SetFocus(pageCmd.descrs)
			return nil
		}
		return event
	})

	pageCmd.cmds.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {

		if event.Key() == tcell.KeyEnter {
			app.SetFocus(pageCmd.descrs)
			return nil
		}
		return event
	})

	flexCmd.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyTab {
			curId := pageCmd.cmds.GetCurrentItem() + 1
			cnt := pageCmd.cmds.GetItemCount()
			if curId < cnt {
				pageCmd.cmds.SetCurrentItem(curId)
				pageCmd.descrs.SetText(pageCmd.mIdDescr[pageCmd.mPosId[pageCmd.cmds.GetCurrentItem()]]+"\n", true)
			}

			return nil
		}

		if event.Key() == tcell.KeyDown {
			curId := pageCmd.cmds.GetCurrentItem() + 1
			cnt := pageCmd.cmds.GetItemCount()
			if curId < cnt {
				pageCmd.cmds.SetCurrentItem(curId)
				pageCmd.descrs.SetText(pageCmd.mIdDescr[pageCmd.mPosId[pageCmd.cmds.GetCurrentItem()]]+"\n", true)
			}

			return nil
		}
		if event.Key() == tcell.KeyBacktab {
			curId := pageCmd.cmds.GetCurrentItem() - 1
			if curId > -1 {
				pageCmd.cmds.SetCurrentItem(curId)
				pageCmd.descrs.SetText(pageCmd.mIdDescr[pageCmd.mPosId[pageCmd.cmds.GetCurrentItem()]]+"\n", true)
			}

			return nil
		}
		if event.Key() == tcell.KeyUp {
			curId := pageCmd.cmds.GetCurrentItem() - 1
			if curId > -1 {
				pageCmd.cmds.SetCurrentItem(curId)
				pageCmd.descrs.SetText(pageCmd.mIdDescr[pageCmd.mPosId[pageCmd.cmds.GetCurrentItem()]]+"\n", true)
			}

			return nil
		}
		if event.Key() == tcell.KeyRight {
			app.SetFocus(pageCmd.descrs)
		}
		if event.Key() == tcell.KeyLeft {
			app.SetFocus(pageCmd.filterFrm)
		}

		return event
	})

	pageCmd.filterFrm.SetFocus(1)

	application.pages.AddPage("commands", flexCmplx, true, true)
}

func refreshCmdList() {
	query := "SELECT id, command, descr FROM cmd WHERE lower(command) like lower('%" + pageCmd.filterFrm.GetFormItem(0).(*tview.InputField).GetText() + "%') order by command"
	cmdFind, err := database.Query(query)
	check(err)

	pageCmd.mIdDescr = make(map[int]string)
	pageCmd.mPosId = make(map[int]int)
	pageCmd.cmds.Clear()
	pageCmd.descrs.SetText("", false)
	rowCount := 1
	for cmdFind.Next() {
		id := 0
		cmd := ""
		descr := ""
		cmdFind.Scan(&id, &cmd, &descr)

		pageCmd.cmds.AddItem(cmd, "", rune(0), func() {})

		pageCmd.mIdDescr[id] = descr
		pageCmd.mPosId[rowCount-1] = id
		rowCount++
	}
}

func saveDescr() {
	query := "UPDATE cmd" + "\n" +
		"SET descr = '" + pageCmd.descrs.GetText() + "'\n" +
		"WHERE id = " + strconv.Itoa(pageCmd.mPosId[pageCmd.cmds.GetCurrentItem()])

	pageCmd.mIdDescr[pageCmd.mPosId[pageCmd.cmds.GetCurrentItem()]] = pageCmd.descrs.GetText()

	_, err := database.Exec(query)
	check(err)
}
