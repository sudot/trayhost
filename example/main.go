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
		trayhost.MenuItem{
			"Ime",
			true,
			nil,
		},
		trayhost.MenuItem{
			"",
			true,
			nil,
		},
		trayhost.MenuItem{
			"Item A",
			false,
			func() {
				fmt.Println("item A")
			},
		},
		trayhost.MenuItem{
			"Item B",
			false,
			nil,
		},
		trayhost.MenuItem{
			"Exit",
			false,
			trayhost.Exit,
		}}

	trayhost.Initialize("Neki", iconData, menuItems)

	time.AfterFunc(10*time.Second, func() {
		trayhost.SetMenu(trayhost.MenuItems{trayhost.MenuItem{
			"ne item",
			false,
			func() {
				fmt.Println("new item")
			},
		}})
	})

	// Enter the host system's event loop
	trayhost.EnterLoop()

	// This is only reached once the user chooses the Exit menu item
	fmt.Println("Exiting")
}
