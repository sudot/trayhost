package main

import (
	"fmt"
	"github.com/overlordtm/trayhost"
	"os"
	"runtime"
	"time"
)

func main() {
	// EnterLoop must be called on the OS's main thread
	runtime.LockOSThread()

	menuItems := trayhost.MenuItems{
		0: trayhost.MenuItem{
			"Trayhost",
			true,
			nil,
		},
		1: trayhost.MenuItem{
			"",
			true,
			nil,
		},
		2: trayhost.MenuItem{
			"Item A",
			false,
			func() {
				fmt.Println("item A")
			},
		},
		3: trayhost.MenuItem{
			"Item B",
			false,
			nil,
		},
		4: trayhost.MenuItem{
			fmt.Sprintf("Time: %v", time.Now()),
			false,
			nil,
		},
		5: trayhost.MenuItem{
			"Exit",
			false,
			trayhost.Exit,
		}}

	trayhost.Initialize("Trayhost example", iconData, menuItems, os.TempDir())
	trayhost.SetClickHandler(onClick)
	trayhost.SetIconImage(trayhost.ICON_ALTERNATIVE, iconData2)
	trayhost.SetIconImage(trayhost.ICON_ATTENTION, iconData3)

	go func() {
		for now := range time.Tick(1 * time.Second) {
			trayhost.UpdateCh <- trayhost.MenuItemUpdate{4, trayhost.MenuItem{
				fmt.Sprintf("Time: %v", now),
				false,
				nil,
			},
			}
		}
	}()

	go func() {
		for _ = range time.Tick(10 * time.Second) {
			trayhost.SetIcon(trayhost.ICON_ALTERNATIVE)
			time.Sleep(5 * time.Second)
			trayhost.SetIcon(trayhost.ICON_ATTENTION)
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
