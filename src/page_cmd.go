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
	btnSave   *tview.Button
	mIdDescr  map[int]string
	mPosId    map[int]int
	//pages     *tview.Pages
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

	pageCmd.filterFrm.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEnter {
			refreshCmdList()
			//filterIndex, _ := pageCmd.filterFrm.GetFocusedItemIndex()
			//query := "SELECT id, command, descr FROM cmd WHERE lower(command) like lower('%" + pageCmd.filterFrm.GetFormItem(filterIndex).(*tview.InputField).GetText() + "%') order by command"
			//
			//cmdFind, err := database.Query(query)
			//check(err)
			//
			//pageCmd.mIdDescr = make(map[int]string)
			//pageCmd.mPosId = make(map[int]int)
			//pageCmd.cmds.Clear()
			//rowCount := 1
			//for cmdFind.Next() {
			//	id := 0
			//	cmd := ""
			//	descr := ""
			//	cmdFind.Scan(&id, &cmd, &descr)
			//
			//	pageCmd.cmds.AddItem(cmd, "", rune(0), func() {})
			//
			//	pageCmd.mIdDescr[id] = descr
			//	pageCmd.mPosId[rowCount-1] = id
			//	rowCount++
			//}
			return nil
		}

		return event
	})

	flexCmd := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(pageCmd.filterFrm, 0, 2, true).
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

	frmSave := tview.NewForm().AddButton("Save", func() {
		query := "UPDATE cmd" + "\n" +
			"SET descr = '" + pageCmd.descrs.GetText() + "'\n" +
			"WHERE id = " + strconv.Itoa(pageCmd.mPosId[pageCmd.cmds.GetCurrentItem()])

		pageCmd.mIdDescr[pageCmd.mPosId[pageCmd.cmds.GetCurrentItem()]] = pageCmd.descrs.GetText()

		_, err := database.Exec(query)
		check(err)
	})

	//flexDescrKeys := tview.NewFlex().
	//	AddItem(pageCmd.descrs, 0, 10, false).
	//	AddItem(pageCmd.cmds, 0, 10, false)

	flexDescr := tview.NewFlex().
		AddItem(pageCmd.descrs, 0, 10, false).
		AddItem(frmSave, 0, 2, false)

	flexDescr.SetDirection(tview.FlexRow)

	flexCmplx := tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(flexCmd, 0, 1, false).
		AddItem(flexDescr, 0, 4, false).
		SetFullScreen(true)

	flexCmplx.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 'q' && event.Modifiers() == tcell.ModAlt {
			app.SetFocus(pageCmd.cmds)
			//pageCmd.cmds.SetSelectable(true, false)
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
		//if event.Rune() == 's' && event.Modifiers() == tcell.ModCtrl {
		//	frmSave.GetFormItemByLabel("Save").PasteHandler()
		//	app.SetFocus(pageCmd.cmds)
		//	return nil
		//}
		if event.Rune() == 'h' && event.Modifiers() == tcell.ModAlt {
			application.pages.ShowPage("help")
			return nil
		}
		if event.Key() == tcell.KeyEsc {
			application.ConfirmQuit()
		}
		return event
	})

	// TODO: error!!!
	pageCmd.descrs.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 's' && event.Modifiers() == tcell.ModAlt {
			frmSave.GetFormItemByLabel("Save").PasteHandler()
			app.SetFocus(pageCmd.cmds)
			return nil
		}
		return event

	})

	pageCmd.cmds.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {

		if event.Key() == tcell.KeyEnter {
			pageCmd.descrs.SetText(pageCmd.mIdDescr[pageCmd.mPosId[pageCmd.cmds.GetCurrentItem()]]+"\n", true)
		}

		return event
	})

	pageCmd.filterFrm.SetFocus(1)

	//pageCmd.pages = tview.NewPages()
	//pageMainMessage.build()
	//pageMainMessage.show(tview.AlignCenter, "", "helpText")

	//pageCmd.pages.ShowPage("helpText")

	application.pages.AddPage("commands", flexCmplx, true, true)
}

func refreshCmdList() {
	query := "SELECT id, command, descr FROM cmd WHERE lower(command) like lower('%" + pageCmd.filterFrm.GetFormItem(0).(*tview.InputField).GetText() + "%') order by command"
	cmdFind, err := database.Query(query)
	check(err)

	pageCmd.mIdDescr = make(map[int]string)
	pageCmd.mPosId = make(map[int]int)
	pageCmd.cmds.Clear()
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
