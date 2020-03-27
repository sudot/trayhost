package main

import (
	"fmt"
	"github.com/sudot/trayhost"
	"net/http"
	"os/exec"
	"runtime"
	"time"
)

func main() {
	// EnterLoop must be called on the OS's main thread
	runtime.LockOSThread()

	// Debug默认是false,在你的实际使用中不需要这一行代码
	trayhost.Debug = true
	trayhost.Initialize("TrayHost", func() {
		fmt.Println("You clicked tray icon")
		openUrl()
	})
	// 通过图片字节数组设置托盘图标
	trayhost.SetIconData(iconData)
	trayhost.SetMenu(trayhost.MenuItems{
		trayhost.NewMenuItemDisabled("TrayHost"),
		trayhost.NewMenuItemDivided(),
		trayhost.NewMenuItem("在浏览器打开", openUrl),
		trayhost.NewMenuItem("Item B", nil),
		trayhost.NewMenuItem(fmt.Sprintf("Time: %v", time.Now()), nil),
		trayhost.NewMenuItemDivided(),
		trayhost.NewMenuItem("Exit", trayhost.Exit),
	})

	go func() {
		// 启动一个 http 服务器
		_ = http.ListenAndServe(":1234", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte("TrayHost"))
		}))
	}()

	go func() {
		// 更新菜单项
		for now := range time.Tick(1 * time.Second) {
			trayhost.UpdateMenuItem(4, trayhost.NewMenuItem(fmt.Sprintf("Time: %v", now), nil))
		}
	}()
	go func() {
		// 每10秒更换两次托盘图标
		for _ = range time.Tick(10 * time.Second) {
			// 通过图片路径设置托盘图标
			trayhost.SetIconPath("example/icons/go.ico")
			time.Sleep(5 * time.Second)
			// 通过图片字节数组设置托盘图标
			trayhost.SetIconData(iconData)
		}
	}()

	// Enter the host system's event loop
	trayhost.EnterLoop()

	fmt.Println("Exiting")
}

func openUrl() {
	var commands = map[string]string{
		"windows": "explorer",
		"darwin":  "open",
		"linux":   "xdg-open",
	}
	run, _ := commands[runtime.GOOS]
	_ = exec.Command(run, "http://127.0.0.1:1234").Start()
}
