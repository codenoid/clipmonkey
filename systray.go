package main

import (
	"clipmonkey/icon"
	"os"
	"strings"

	"github.com/getlantern/systray"
	"golang.design/x/clipboard"
)

func onReady() {
	systray.SetIcon(icon.GetIcon())
	systray.SetTitle("Clip Monkey")
	systray.SetTooltip("Clip Monkey")

	mQuit := systray.AddMenuItem("Quit", "Quit the whole app")
	go func() {
		for range mQuit.ClickedCh {
			systray.Quit()
			os.Exit(0)
		}
	}()

	for _, item := range clipboardHistory {
		showNewItem(item)
	}
}

func onExit() {
	// clean up here
}

func showNewItem(clip string) {
	label := strings.Join(strings.Fields(strings.TrimSpace(clip)), " ")

	if len(label) > 20 {
		label = label[:20] + "..."
	}
	btn := systray.AddMenuItem(label, "Click to copy")
	go func(btn *systray.MenuItem, label, clip string) {
		for range btn.ClickedCh {
			clipboard.Write(clipboard.FmtText, []byte(clip))
		}
	}(btn, label, clip)
}
