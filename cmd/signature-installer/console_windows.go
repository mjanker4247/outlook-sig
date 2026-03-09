//go:build windows

package main

import (
	"os"
	"syscall"
)

func init() {
	// Built with -H windowsgui so Windows never opens a console window on launch.
	// When CLI flags are present we re-attach to the parent process's console
	// (cmd.exe / PowerShell) so that output and stdin work normally from a terminal.
	// With no arguments the GUI starts silently without any console window.
	if len(os.Args) <= 1 {
		return
	}

	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	attachConsole := kernel32.NewProc("AttachConsole")

	const attachParentProcess = ^uintptr(0) // ATTACH_PARENT_PROCESS = -1
	ret, _, _ := attachConsole.Call(attachParentProcess)
	if ret == 0 {
		// No parent console (e.g. launched from Explorer with arguments).
		// Do not call AllocConsole — that would open a new window.
		return
	}

	// Reopen standard handles so fmt.Print / fmt.Scan reach the terminal.
	if f, err := os.OpenFile("CONIN$", os.O_RDWR, 0); err == nil {
		os.Stdin = f
	}
	if f, err := os.OpenFile("CONOUT$", os.O_RDWR, 0); err == nil {
		os.Stdout = f
		os.Stderr = f
	}
}
