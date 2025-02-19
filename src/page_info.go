package main

import (
	"github.com/rivo/tview"
)

type pageInfoType struct {
	*tview.Flex
	pages *tview.Pages
}

var pageInfo pageInfoType

func (pageInfo *pageInfoType) build() {

	pageInfo.pages = tview.NewPages()
	pageDesc.build()
	pageCase.build()
	pageCaseIns.build()

	pageInfo.Flex = tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(pageInfo.pages, 5, 0, true)

	pageInfo.pages.ShowPage("desc")
}
