// +build windows

package trayhost

import (
	"syscall"
	"unsafe"
)

/*
#include <stdlib.h>
#include "platform/trayhost.h"
#include "platform/trayhost_win.h"
*/
import "C"

func initialize(title string) {
	titlePtr, _ := syscall.UTF16PtrFromString(title)
	defer C.free(unsafe.Pointer(titlePtr))
	C.init((*C.char)(unsafe.Pointer(titlePtr)))
}

func setMenuItem(id int, item MenuItem) {
	titlePtr, _ := syscall.UTF16PtrFromString(item.Title)
	defer C.free(unsafe.Pointer(titlePtr))
	C.set_menu_item((C.int)(id), (*C.char)(unsafe.Pointer(titlePtr)), cbool(item.Disabled))
}

func setIcon(iconPth string) {
	cIconPth, _ := syscall.UTF16PtrFromString(iconPth)
	defer C.free(unsafe.Pointer(cIconPth))
	C.set_icon((*C.char)(unsafe.Pointer(cIconPth)))
}
