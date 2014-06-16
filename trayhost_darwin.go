// +build darwin

package trayhost

func initialize(title string) {
	titlePtr := C.CString(title)
	defer C.free(unsafe.Pointer(titlePtr))
	C.init((*C.char)(unsafe.Pointer(titlePtr)))
}
