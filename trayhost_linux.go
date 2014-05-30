// +build linux

package trayhost

import (
	"os"
)

func getDesktop() int {
	currentDesktop := os.Getenv("XDG_CURRENT_DESKTOP")

	switch currentDesktop {
	case "Unity":
		return UNITY
	case "GNOME":
		return GNOME
	default:
		return LINUX_GENERIC
	}
}
