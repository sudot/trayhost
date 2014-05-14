/*
Package trayhost is a library for placing a Go
application in the task bar (system tray,
notification area, or dock) in a consistent
manner across multiple platforms. Currently,
there is built-in support for Windows, Mac OSX,
and Linux systems that support GTK+ 3 status
icons (including Gnome 2, KDE 4, Cinnamon,
MATE and other desktop environments).

The indended usage is for applications that
utilize web technology for the user interface, but
require access to the client system beyond what
is offered in a browser sandbox (for instance,
an application that requires access to the user's
file system).

The library places a tray icon on the host system's
task bar that can be used to open a URL, giving users
easy access to the web-based user interface.

Further information can be found at the project's
home at http://github.com/cratonica/trayhost

Clint Caywood

http://github.com/cratonica/trayhost
*/
package trayhost

import (
	"reflect"
	"runtime"
	"sort"
	"unsafe"
)

/*
#cgo linux pkg-config: gtk+-2.0
#cgo linux CFLAGS: -DLINUX -I/usr/include/libappindicator-0.1/
#cgo linux LDFLAGS: -ldl
#cgo windows CFLAGS: -DWIN32
#cgo darwin CFLAGS: -DDARWIN -x objective-c
#cgo darwin LDFLAGS: -framework Cocoa
#include <stdlib.h>
#include "platform/platform.h"
*/
import "C"

var isExiting bool
var menuItems MenuItems
var UpdateCh = make(chan MenuItemUpdate, 99)

type MenuItem struct {
	Title    string
	Disabled bool
	Handler  func()
}

type MenuItemUpdate struct {
	ItemId int
	Item   MenuItem
}

type MenuItems map[int]MenuItem

// Run the host system's event loop
func Initialize(title string, imageData []byte, items MenuItems) {

	cTitle := C.CString(title)
	defer C.free(unsafe.Pointer(cTitle))

	// Copy the image data into unmanaged memory
	cImageData := C.malloc(C.size_t(len(imageData)))
	defer C.free(cImageData)
	var cImageDataSlice []C.uchar
	sliceHeader := (*reflect.SliceHeader)(unsafe.Pointer(&cImageDataSlice))
	sliceHeader.Cap = len(imageData)
	sliceHeader.Len = len(imageData)
	sliceHeader.Data = uintptr(cImageData)

	for i, v := range imageData {
		cImageDataSlice[i] = C.uchar(v)
	}

	// Initialize menu
	C.init(cTitle, &cImageDataSlice[0], C.uint(len(imageData)))

	setMenu(items)
}

func EnterLoop() {
	go updater()
	C.native_loop()
	// If reached, user clicked Exit
	isExiting = true
}

func Exit() {
	C.exit_loop()
}

func updater() {
	runtime.UnlockOSThread()
	for update := range UpdateCh {
		updateMenuItem(update.ItemId, update.Item)
	}
}

func cAddMenuItem(id C.int, title *C.char, disabled C.int) {
	defer C.free(unsafe.Pointer(title))
	C.add_menu_item(id, title, disabled)
}

func setMenu(menu MenuItems) {
	menuItems = menu

	menuItemOrder := make([]int, 0, len(menuItems))

	for key, _ := range menuItems {
		menuItemOrder = append(menuItemOrder, key)
	}

	sort.Ints(menuItemOrder)

	for id := range menuItemOrder {
		item := menuItems[id]
		addMenuItem(id, item)
	}
}

func updateMenuItem(id int, item MenuItem) {
	menuItems[id] = item
	addMenuItem(id, item)
}

func cbool(b bool) C.int {
	if b {
		return 1
	} else {
		return 0
	}
}
