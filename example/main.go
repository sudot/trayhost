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
			"Ime",
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

	trayhost.Initialize("Neki", iconData, menuItems)

	go func() {
		for now := range time.Tick(1 * time.Second) {
			text := fmt.Sprintf("%v", now)
			trayhost.UpdateMenuItem(99, trayhost.MenuItem{
				text,
				true,
				func() {
					fmt.Println("new item", text)
				},
			})
		}
	}()

	// Enter the host system's event loop
	trayhost.EnterLoop()

	// This is only reached once the user chooses the Exit menu item
	fmt.Println("Exiting")
}
