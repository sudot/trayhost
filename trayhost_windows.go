// +build windows

package trayhost

import "C"

import (
	"syscall"
	"unsafe"
)

func setMenuItem(id int, item MenuItem) {
	// log.Printf("setMenuItem: %d %s", id, item.Title)
	// ignore errors
	titlePtr, _ := syscall.UTF16PtrFromString(item.Title)
	cSetMenuItem((C.int)(id), (*C.char)(unsafe.Pointer(titlePtr)), cbool(item.Disabled))
}
