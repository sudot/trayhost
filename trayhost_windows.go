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
	C.init((*C.char)(unsafe.Pointer(titlePtr)))
}

func setMenuItem(id int, item MenuItem) (err error) {
	titlePtr, err := syscall.UTF16PtrFromString(item.Title)
	if err != nil {
		return
	}
	C.set_menu_item((C.int)(id), (*C.char)(unsafe.Pointer(titlePtr)), cbool(item.Disabled))
	return
}

func setIcon(iconPth string) {
	cIconPth, _ := syscall.UTF16PtrFromString(iconPth)
	C.set_icon((*C.char)(unsafe.Pointer(cIconPth)))
}
