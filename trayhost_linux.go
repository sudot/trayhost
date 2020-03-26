// +build linux

package trayhost

import (
	"os"
)

import (
	"unsafe"
)

/*
#include <stdlib.h>
#include "platform/trayhost_linux.h"
*/
import "C"

func getDesktop() int {
	currentDesktop := os.Getenv("XDG_CURRENT_DESKTOP")

	switch currentDesktop {
	case "Unity":
		return UNITY
	case "GNOME":
		return GNOME
	default:
		return LinuxGeneric
	}
}

func initialize(title string) {
	titlePtr := C.CString(title)
	defer C.free(unsafe.Pointer(titlePtr))
	C.init((*C.char)(unsafe.Pointer(titlePtr)), (C.int)(getDesktop()))
}
