package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type pageMainMessageType struct {
	helpText        string
	textViewMessage *tview.TextView
}

var pageMainMessage pageMainMessageType

func (pageMainMessage *pageMainMessageType) build() {
	pageMainMessage.textViewMessage = tview.NewTextView()
	pageMainMessage.textViewMessage.SetBorder(true)
	pageMainMessage.textViewMessage.SetBorderPadding(1, 1, 1, 1)
	pageMainMessage.textViewMessage.SetTitleAlign(tview.AlignCenter)
	pageMainMessage.textViewMessage.SetRegions(true)
	pageMainMessage.textViewMessage.SetWordWrap(true)
	pageMainMessage.textViewMessage.SetWrap(true)
	pageMainMessage.textViewMessage.SetScrollable(true)

	pageMainMessage.helpText = `KEYBOARD SHORTCUTS

[green]Global:[-:-:-:-]
[yellow]esc[-:-:-:-]: Exit page/Quit application.
[yellow]ctrl + shift + v[-:-:-:-]: Paste text. // TODO
[yellow]ctrl + z[-:-:-:-]: Undo text. // TODO

[green]Main UI:[-:-:-:-]
[yellow]alt + f[-:-:-:-]: Focus on the filter item.
[yellow]alt + q[-:-:-:-]: Focus on the item list.
[yellow]alt + w[-:-:-:-]: Focus on the description area.
[yellow]ctrl + s[-:-:-:-]: Save.

[yellow]alt + v[-:-:-:-]: Paste clipboard buffer content (only description field).

ABOUT

ITEM-STORE is a free and open-source project maintained by contributors. Feel free report issues or submit feature requests on GitHub
`
	pageMainMessage.textViewMessage.SetText(pageMainMessage.helpText)

	pageMainMessage.textViewMessage.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEsc {
			application.pages.SwitchToPage("items")
			return nil
		}

		return event
	})

	application.pages.AddPage("help", pageMainMessage.textViewMessage, true, false)

}

func (pageMainMessage *pageMainMessageType) show(textAlign int, title, message string) {
	pageMainMessage.textViewMessage.Clear()

	switch message {
	case "helpText":
		pageMainMessage.textViewMessage.SetTextAlign(tview.AlignLeft)
		pageMainMessage.textViewMessage.SetTitle("Help (alt+h)")
		pageMainMessage.textViewMessage.SetText(pageMainMessage.helpText)
	default:
		pageMainMessage.textViewMessage.SetTextAlign(textAlign)
		pageMainMessage.textViewMessage.SetTitle(title)
		pageMainMessage.textViewMessage.SetText(message)
	}

	application.pages.SwitchToPage("message")
}

func (pageMainMessage *pageMainMessageType) focus() {
	app.SetFocus(pageMainMessage.textViewMessage)
}
