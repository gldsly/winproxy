package main

import (
	"fmt"
	"github.com/getlantern/systray"
	"github.com/gldsly/winproxy/icon"
	"github.com/skratchdot/open-golang/open"
)

var (
	quickButton *systray.MenuItem
	settingButton *systray.MenuItem
)

func initMenu() {
	settingButton = systray.AddMenuItem("控制台", "查看和修改转发参数")
	quickButton = systray.AddMenuItem("退出", "退出")
}

func onExit() {
	// clean up here
}

func onReady() {
	systray.SetIcon(icon.Data)
	// windows 无法显示标题
	//systray.SetTitle("端口转发")
	systray.SetTooltip("Windows Port Proxy")

	initMenu()

	for {
		select {
		case <-quickButton.ClickedCh:
			systray.Quit()
			fmt.Println("退出")
		case <-settingButton.ClickedCh:
			open.Run("http://127.0.0.1:57391")
		}
	}
}

func main() {
	go StartService()
	systray.Run(onReady, onExit)
}
