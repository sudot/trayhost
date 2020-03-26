package trayhost

import "C"

//export trayCallback
func trayCallback(itemId C.int) {
	id := int(itemId)
	if id == -1 {
		trayHostLog.Println("Tray click")
		if clickHandler != nil {
			clickHandler()
		}
		return
	}

	menuItem, hasItem := menuItems[id]
	if hasItem {
		if menuItem.Handler != nil {
			menuItem.Handler()
		} else {
			trayHostLog.Printf("Item %s has no handler", menuItem.Title)
		}
	}
}

//export goLog
func goLog(msg *C.char) {
	trayHostLog.Printf("cgo: %s", C.GoString(msg))
}
