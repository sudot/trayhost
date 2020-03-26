package main

import (
	"fmt"
	"github.com/sudot/trayhost"
	"os"
	"runtime"
	"time"
)

func main() {
	// EnterLoop must be called on the OS's main thread
	runtime.LockOSThread()

	// 构造菜单项,最终会按 key 的数值从小到大-从上到下排列
	menuItems := trayhost.MenuItems{
		0: trayhost.NewMenuItemDisabled("TrayHost"),
		// Title 为空则为分割线
		1: trayhost.NewMenuItemDivided(),
		2: trayhost.NewMenuItem("Item A", func() {
			fmt.Println("item A")
		}),
		3: trayhost.NewMenuItem("Item B", nil),
		4: trayhost.NewMenuItem(fmt.Sprintf("Time: %v", time.Now()), nil),
		// Title 为空则为分割线
		5: trayhost.NewMenuItemDivided(),
		6: trayhost.NewMenuItem("Exit", trayhost.Exit),
	}

	trayhost.Debug = true
	_ = trayhost.Initialize("TrayHost example", iconData, menuItems, os.TempDir())
	trayhost.SetClickHandler(onClick)
	_ = trayhost.SetIconImage(trayhost.IconAlternative, iconData2)
	_ = trayhost.SetIconImage(trayhost.IconAttention, iconData3)

	go func() {
		// 更改菜单内容
		for now := range time.Tick(1 * time.Second) {
			trayhost.UpdateCh <- trayhost.MenuItemUpdate{
				ItemId: 4,
				Item: trayhost.MenuItem{
					Title: fmt.Sprintf("Time: %v", now),
				},
			}
		}
	}()

	go func() {
		// 更换托盘图标
		for _ = range time.Tick(10 * time.Second) {
			_ = trayhost.SetIcon(trayhost.IconAlternative)
			time.Sleep(5 * time.Second)
			_ = trayhost.SetIcon(trayhost.IconAttention)
		}
	}()

	// Enter the host system's event loop
	trayhost.EnterLoop()

	// This is only reached once the user chooses the Exit menu item
	fmt.Println("Exiting")
}

func onClick() {
	fmt.Println("You clicked tray icon")
}
