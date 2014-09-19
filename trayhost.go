/*
 */
package trayhost

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"
)

/*
#cgo linux pkg-config: gtk+-2.0
#cgo linux CFLAGS: -DLINUX -I/usr/include/libappindicator-0.1/
#cgo linux LDFLAGS: -ldl
#cgo windows CFLAGS: -DWIN32 -DUNICODE
#cgo darwin CFLAGS: -DDARWIN -x objective-c
#cgo darwin LDFLAGS: -framework Cocoa
#include <stdlib.h>
#include "platform/platform.h"
*/
import "C"

const (
	ICON_PRIMARY     = iota
	ICON_ALTERNATIVE = iota
	ICON_ATTENTION   = iota
)

const (
	WINDOWS       = iota
	OSX           = iota
	GNOME         = iota
	UNITY         = iota
	LINUX_GENERIC = iota
)

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

var isExiting bool = false
var menuItems MenuItems
var UpdateCh = make(chan MenuItemUpdate, 99)
var icons = map[int]string{}
var tmpFiles []string = make([]string, 0, 3)
var clickHandler func()
var Debug bool = false
var curIconId int = -1
var tmpDir string

// Run the host system's event loop
func Initialize(title string, imageData []byte, items MenuItems, tmpDirectory string) (err error) {
	if !Debug {
		log.SetOutput(ioutil.Discard)
	}
	tmpDir = tmpDirectory
	err = SetIconImage(ICON_PRIMARY, imageData)
	if err != nil {
		return
	}
	initialize(title)
	err = SetIcon(ICON_PRIMARY)
	if err != nil {
		return
	}
	setMenu(items)
	return
}

func SetIconImage(iconId int, imageData []byte) (err error) {
	iconPth, err := createTempFile(imageData)
	if err != nil {
		return
	}
	icons[iconId] = iconPth
	return
}

func EnterLoop() {
	go menuUpdateLoop()
	C.native_loop()
	isExiting = true
}

func Exit() {
	C.exit_loop()
	cleanup()
}

func SetIcon(iconId int) (err error) {

	if iconId != curIconId {
		iconPth, ok := icons[iconId]
		if !ok {
			err = fmt.Errorf("No icon with icon id %d", iconId)
			return
		}

		log.Printf("Setting icon %s (id: %d)", iconPth, iconId)
		setIcon(iconPth)
		curIconId = iconId
	}
	return
}

func SetClickHandler(handler func()) {
	clickHandler = handler
}

func menuUpdateLoop() {
	for update := range UpdateCh {
		updateMenuItem(update.ItemId, update.Item)
	}
}

func updateMenuItem(id int, item MenuItem) {
	menuItems[id] = item
	setMenuItem(id, item)
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
		setMenuItem(id, item)
	}
}

func createTempFile(iconData []byte) (filename string, err error) {
	file, err := ioutil.TempFile(tmpDir, "trayhosticon")
	if err != nil {
		return
	}
	defer file.Close()
	_, err = file.Write(iconData)
	filename = file.Name()
	tmpFiles = append(tmpFiles, filename)
	return
}

func cleanup() {
	for _, file := range tmpFiles {
		err := os.Remove(file)
		if err != nil {
			log.Printf("Failed to remove tmp file %s: %v\n", file, err)
		} else {
			log.Printf("Tmp file %s removed\n", file)
		}
	}
}

func cbool(b bool) C.int {
	if b {
		return 1
	} else {
		return 0
	}
}
