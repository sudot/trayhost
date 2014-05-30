// +build linux darwin

package trayhost

import "C"

func setMenuItem(id int, item MenuItem) {
	cSetMenuItem((C.int)(id), C.CString(item.Title), cbool(item.Disabled))
}
