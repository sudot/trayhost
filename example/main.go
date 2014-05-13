package main

import (
	"fmt"
	"github.com/overlordtm/trayhost"
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
			"Exit",
			false,
			trayhost.Exit,
		}}

	trayhost.Initialize("Some name", iconData, menuItems)

	go func() {
		for now := range time.Tick(3 * time.Second) {
			trayhost.UpdateCh <- trayhost.MenuItemUpdate{2, trayhost.MenuItem{
				fmt.Sprintf("zoki %v", now),
				false,
				func() {
					fmt.Println("zoki")
				},
			}}

			trayhost.UpdateCh <- trayhost.MenuItemUpdate{3, trayhost.MenuItem{
				fmt.Sprintf("boki %v", now),
				false,
				func() {
					fmt.Println("boki")
				},
			}}
		}
	}()

	// Enter the host system's event loop
	trayhost.EnterLoop()

	// This is only reached once the user chooses the Exit menu item
	fmt.Println("Exiting")
}
