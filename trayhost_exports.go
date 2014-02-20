package trayhost

import (
	"fmt"
)

import "C"

//export tray_callback
func tray_callback(itemId C.int) {
	item := menuItems[itemId]

	if item.Handler != nil {
		item.Handler()
	} else {
		fmt.Println("no handler")
	}
}
