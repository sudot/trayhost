TrayHost
========

__TrayHost__ is a library for placing a __Go__ application in the task bar (system tray, notification area, or dock) in a consistent manner across multiple platforms. Currently, there is built-in support for __Windows__, __Mac OSX__, and __Linux__ systems that support GTK+ 3 status icons (including Gnome 2, KDE 4, Cinnamon, MATE and other desktop environments).

The intended usage is for applications that utilize web technology for the user interface, but require access to the client system beyond what is offered in a browser sandbox (for instance, an application that requires access to the user's file system).

The library places a tray icon on the host system's task bar that can be used to open a URL, giving users easy access to the web-based user interface. 

API docs can be found [here](http://godoc.org/github.com/cratonica/trayhost)

The Interesting Part
----------------------
```go
package main

import (
	"fmt"
	"github.com/sudot/trayhost"
	"os"
	"runtime"
	"time"
)

func main() {
	// EnterLoop must be called on the OS's main thread
	runtime.LockOSThread()

	// 构造菜单项,最终会按 key 的数值从小到大-从上到下排列
	menuItems := trayhost.MenuItems{
		0: trayhost.NewMenuItemDisabled("TrayHost"),
		// Title 为空则为分割线
		1: trayhost.NewMenuItemDivided(),
		2: trayhost.NewMenuItem("Item A", func() {
			fmt.Println("item A")
		}),
		3: trayhost.NewMenuItem("Item B", nil),
		4: trayhost.NewMenuItem(fmt.Sprintf("Time: %v", time.Now()), nil),
		// Title 为空则为分割线
		5: trayhost.NewMenuItemDivided(),
		6: trayhost.NewMenuItem("Exit", trayhost.Exit),
	}

	trayhost.Debug = true
	_ = trayhost.Initialize("TrayHost example", iconData, menuItems, os.TempDir())
	trayhost.SetClickHandler(onClick)
	_ = trayhost.SetIconImage(trayhost.IconAlternative, iconData2)
	_ = trayhost.SetIconImage(trayhost.IconAttention, iconData3)

	go func() {
		// 更改菜单内容
		for now := range time.Tick(1 * time.Second) {
			trayhost.UpdateCh <- trayhost.MenuItemUpdate{
				ItemId: 4,
				Item: trayhost.MenuItem{
					Title: fmt.Sprintf("Time: %v", now),
				},
			}
		}
	}()

	go func() {
		// 更换托盘图标
		for _ = range time.Tick(10 * time.Second) {
			_ = trayhost.SetIcon(trayhost.IconAlternative)
			time.Sleep(5 * time.Second)
			_ = trayhost.SetIcon(trayhost.IconAttention)
		}
	}()

	// Enter the host system's event loop
	trayhost.EnterLoop()

	// This is only reached once the user chooses the Exit menu item
	fmt.Println("Exiting")
}

func onClick() {
	fmt.Println("You clicked tray icon")
}

```

Build Environment
--------------------------
Before continuing, make sure that your GOPATH environment variable is set, and that you have Git and Mercurial installed and that __go__, __git__, and __hg__ are in your PATH.

Cross-compilation is not currently supported, so you will need access to a machine running the platform that you wish to target. 

Generally speaking, make sure that your system is capable of doing [cgo](http://golang.org/doc/articles/c_go_cgo.html) builds.

#### Linux
In addition to the essential GNU build tools, you will need to have the GTK+ 3.0 development headers installed.

#### Windows
To do cgo builds, you will need to install [MinGW](http://www.mingw.org/). In order to prevent the terminal window from appearing when your application runs, you'll need access to a copy of [editbin.exe](http://msdn.microsoft.com/en-us/library/xd3shwhf.aspx) which comes packaged with Microsoft's C/C++ build tools.

#### Mac OSX
__Note__: TrayHost requires __Go 1.1__ when targetting Mac OSX, or linking will fail due to issues with previous versions of Go and Mach-O binaries.

You'll need the "Command Line Tools for Xcode", which can be installed using Xcode. You should be able to run the __cc__ command from a terminal window.

Installing
-----------
Once your build environment is configured, go get the library:

```bash
go get github.com/sudot/trayhost
```

If all goes well, you shouldn't get any errors.

Using
-----
Use the included __example/main.go__ file as a template to get going.  OSX will throw a runtime error if __EnterLoop__ is called on a child thread, so the first thing you must do is lock the OS thread. Your application code will need to run on a child goroutine. __SetUrl__ can be called lazily if you need to take some time to determine what port you are running on. 

Before it will build, you will need to pick an icon for display in the system tray.

#### Generating the Tray Icon
Included in the project is a tool for generating the icon that gets displayed in the system tray. An icon sized 64x64 pixels should suffice, but there aren't any restrictions here as the system will take care of fitting it (just don't get carried away). 

Icons are embedded into the application by generating a Go array containing the byte data using the [2goarray](http://github.com/cratonica/2goarray) tool, which will automatically be installed if it is missing. The generated .go file will be compiled into the output program, so there is no need to distribute the icon with the program. If you want to embed more resources, check out the [embed](http://github.com/cratonica/embed) project.

#### Linux/OSX
From your project root, run __make_icon.sh__, followed by the path to a __PNG__ file to use. For example:

    make_icon.sh example/icons/icon-1-256.png

This will generate a file called __iconunix.go__ and set its build options so it won't be built in Windows.

#### Windows
From the project root, run __make_icon.bat__, followed by the path to a __Windows ICO__ file to use. If you need to create an ICO file, the online tool [ConvertICO](http://convertico.com/) can do this painlessly. 

Example:

    make_icon.bat example/icons/icon-1-256.ico

This will generate a file called __iconwin.go__ and set its build options so it will only be built in Windows.
    
#### Disabling the Command Prompt Window on Windows
The [editbin](http://msdn.microsoft.com/en-us/library/xd3shwhf.aspx) tool will allow you to change the subsystem of the output executable so that users won't see the command window while your application is running. The easiest way to do this is to open the Visual Studio Command Prompt from the start menu (or, alternatively, find __vcvarsall.bat__ in your Visual Studio installation directory and CALL it passing the __x86__ argument). Once you are in this environment, issue the command:

    editbin.exe /SUBSYSTEM:WINDOWS path\to\program.exe

Now when you run the program, you won't see a terminal window.
