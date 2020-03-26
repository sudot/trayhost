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

	if id < 0 || id >= len(menuItems) {
		return
	}
	menuItem := menuItems[id]
	if menuItem.Handler != nil {
		menuItem.Handler()
	} else {
		trayHostLog.Printf("Item %s has no handler", menuItem.Title)
	}
}

//export goLog
func goLog(msg *C.char) {
	trayHostLog.Printf("cgo: %s", C.GoString(msg))
}
