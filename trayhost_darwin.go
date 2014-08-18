// +build darwin

package trayhost

import (
	"unsafe"
)

/*
#include <stdlib.h>
#include "platform/trayhost_darwin.h"
*/
import "C"

func initialize(title string) {
	titlePtr := C.CString(title)
	defer C.free(unsafe.Pointer(titlePtr))
	C.init((*C.char)(unsafe.Pointer(titlePtr)))
}
