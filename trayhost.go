package trayhost

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
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
	WINDOWS      = iota
	OSX          = iota
	GNOME        = iota
	UNITY        = iota
	LinuxGeneric = iota
)

type MenuItem struct {
	Title    string
	Disabled bool
	Handler  func()
}

type MenuItems []MenuItem

var clickHandler func()
var Debug = false
var prefixPath string
var trayHostLog = log.New(os.Stdout, "TrayHost:", log.LstdFlags)
var menuItems MenuItems

func NewMenuItem(title string, handler func()) MenuItem {
	return MenuItem{Title: title, Handler: handler}
}

func NewMenuItemDisabled(title string) MenuItem {
	return MenuItem{Title: title, Disabled: true}
}

func NewMenuItemDivided() MenuItem {
	return MenuItem{}
}

// Run the host system's event loop
// 初始化系统托盘
// example/main.go 和 examplepath/main.go 都是很好的例子,你通过这写例子应该能了解到其所有使用方法.
//
// 如果设置 Debug = true，则会在操作过程中显示一些日志信息
//
// 你可以通过 title 设置鼠标停留在托盘图标上是显示的文字,如果你不想在鼠标停留时出现此文字,可以设置为空字符串 "".
//
// 托盘图标当然少不了要一个图片,你可以通过 SetIconPath 指定图片路径,也可以通过 SetIconData 指定图片的字节数组数据.
// 若你通过 SetIconPath 设置图片,在初始化的时候给 workDir 的值则表示图片所在的目录,当然你也可以让他为空字符串 "".
// 比如图片路径是 examplepath/icons/icon-1-256.ico,你可以通过设置 workDir 为 examplepath/icons,
// 然后 SetIconPath 就可以写成 icon-1-256.ico.
// 若你通过 SetIconData 设置图片,在初始化的时候给 workDir 设置任意值都不影响,因为此时不需要从任何地方读取图片文件.
//
// 一些程序的托盘图标在接收到鼠标的点击事件时,可以执行一个任意的操作,你也可以实现这一的功能,通过 handler 来实现即可.
// 当然 handler 只能接收到鼠标的左键点击响应,因为鼠标的右键点击已经被占用,用来显示 items 里的每一个操作了
//
// 也许你想要有一个强大的菜单操作,通过 SetMenu 可以帮助你实现这个想法.
// items 中的每一项菜单也可以响应鼠标的点击事件,但是只能响应鼠标左键的点击事件,因为一般情况下这样就够了.
//
// 在程序处理好所有初始化操作后,需要调用 EnterLoop 才能让系统托盘响应鼠标的各种点击事件,而且此调用不能在协程里执行,
// 如果你没有任何办法改变在主线程里执行,那么你可以通过 runtime.LockOSThread() 来操作
// 如果你需要在 EnterLoop 后执行类似 http.ListenAndServe(":1234", nil) 需要系统挂起的操作,直接写在 EnterLoop 后是不可行的
// 因为 EnterLoop 函数执行后,程序就会挂起在当前位置, EnterLoop 后的所有操作都无法执行,触发程序退出
// 所以你需要另起一个协程处理诸如 http.ListenAndServe(":1234", nil) 此类的操作,就像下面这样:
//
// go func() {
// 	http.ListenAndServe(":1234", nil)
// }()
//
// trayhost.EnterLoop()
//
// title   是鼠标在托盘图标上停留一段时间后显示的文字
// workDir 是托盘图标的路径或者说是通过图片数据写入图片的临时目录,可以是相对路径,也可以是绝对路径
// handler 是在托盘图标上点击鼠标左键时触发的回调
func Initialize(title string, workDir string, handler func()) {
	if !Debug {
		trayHostLog.SetOutput(ioutil.Discard)
	}
	prefixPath = workDir
	clickHandler = handler
	initialize(title)
}
func EnterLoop() {
	C.native_loop()
}

func Exit() {
	C.exit_loop()
}

// 更新系统托盘的图标
func SetIconPath(iconPth string) {
	iconPth = filepath.Join(prefixPath, iconPth)
	trayHostLog.Printf("Setting icon %s", iconPth)
	setIcon(iconPth)
}

func SetIconData(imageData []byte) {
	iconPth, err := createTempFile(imageData)
	if err != nil {
		return
	}
	defer os.Remove(iconPth)
	trayHostLog.Printf("Setting icon %s", iconPth)
	setIcon(iconPth)
}

// 此操作将会把托盘右键弹出的所有菜单项全部替换掉
func SetMenu(menu MenuItems) {
	menuItems = menu
	for index, item := range menuItems {
		_ = setMenuItem(index, item)
	}
}

// 更新指定的菜单
func UpdateMenuItem(index int, item MenuItem) {
	menuItems[index] = item
	_ = setMenuItem(index, item)
}

func createTempFile(iconData []byte) (filename string, err error) {
	file, err := ioutil.TempFile(os.TempDir(), "trayhosticon")
	if err != nil {
		return
	}
	defer file.Close()
	_, err = file.Write(iconData)
	filename = file.Name()
	return
}

func cBool(b bool) C.int {
	if b {
		return 1
	} else {
		return 0
	}
}
