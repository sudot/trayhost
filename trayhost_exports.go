package trayhost

import (
	"fmt"
)

import "C"

//export tray_callback
func tray_callback(itemId C.int) {

	id := int(itemId)

	menuItem, has := menuItems[id]

	if id == -1 {
		fmt.Println("Tray click")
	}

	if has && menuItem.Handler != nil {
		menuItem.Handler()
	} else {
		fmt.Println("No handler")
	}
}
