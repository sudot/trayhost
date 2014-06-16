// +build linux darwin

package trayhost

import (
	"unsafe"
)

/*
#include <stdlib.h>
#include "platform/trayhost.h"
*/
import "C"

func setMenuItem(id int, item MenuItem) {
	cTitle := C.CString(item.Title)
	defer C.free(unsafe.Pointer(cTitle))
	C.set_menu_item((C.int)(id), cTitle, cbool(item.Disabled))
}

func setIcon(iconPth string) {
	cIconPth := C.CString(iconPth)
	defer C.free(unsafe.Pointer(cIconPth))
	C.set_icon(cIconPth)
}
