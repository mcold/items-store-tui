package main

import (
	//"database/sql"

	"github.com/rivo/tview"
)

type pageCmdType struct {
	cmds      *tview.List
	descrs    *tview.List
	filterFrm *tview.Form
}

var pageCmd pageCmdType

func (pageCmd *pageCmdType) build() {
	pageCmd.cmds = tview.NewList()
	pageCmd.cmds.Box.SetBorder(true)
	pageCmd.cmds.Box.SetTitle("command (alt+q)")

	pageCmd.descrs = tview.NewList()
	pageCmd.descrs.Box.SetBorder(true)
	pageCmd.descrs.Box.SetTitle("description (alt+w)")

	//filter := tview.NewForm().AddInputField("cmd", "", 5, nil, nil).
	//	SetBorder(true)
	//filter := tview.NewForm().AddInputField("command", "", 10, nil, nil).
	//	SetBorder(true)
	//frm := tview.NewFlex().SetDirection(tview.FlexRow).
	//	AddItem(filter, 0, 10, false)
	//
	//fl := tview.NewFlex().SetDirection(tview.FlexRow).
	//	AddItem(frm, 0, 2, true).
	//	AddItem(pageCmd.cmds, 0, 10, false)

	//pageCmd.flex = tview.NewFlex().SetDirection(tview.FlexColumn).AddItem(fl, 0, 1, false).
	//	AddItem(pageCmd.descrs, 0, 1, false).
	//	AddItem(tview.NewInputField(), 0, 1, false).
	//	AddItem(filter, 0, 1, false)

	//pageCmd.flex = tview.NewFlex().SetDirection(tview.FlexColumn).
	//	AddItem(filter, 0, 1, false).
	//	AddItem(tview.NewInputField(), 0, 1, false)

	pageCmd.filterFrm = tview.NewForm().
		AddInputField("COMMAND", "", 30, nil, nil)

	flexCmd := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(pageCmd.filterFrm, 0, 1, true).
		AddItem(pageCmd.cmds, 0, 9, false)

	flexDescr := tview.NewFlex().
		AddItem(pageCmd.descrs, 0, 1, false)

	flexCmplx := tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(flexCmd, 0, 1, false).
		AddItem(flexDescr, 0, 4, false).
		SetFullScreen(true)

	// TODO: set focus on filter input field
	pageCmd.filterFrm.SetFocus(1)

	application.pages.AddPage("commands", flexCmplx, true, true)
}
