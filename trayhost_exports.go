package trayhost

import "C"
import (
	"log"
)

//export tray_callback
func tray_callback(itemId C.int) {

	id := int(itemId)
	menuItem, hasItem := menuItems[id]

	if id == -1 {
		log.Println("Tray click")
		if clickHandler != nil {
			clickHandler()
		}
	}

	if hasItem {
		if menuItem.Handler != nil {
			menuItem.Handler()
		} else {
			log.Printf("Item %s has no handler", menuItem.Title)
		}

	}
}

//export go_log
func go_log(msg *C.char) {
	log.Printf("cgo: %s", C.GoString(msg))
}
