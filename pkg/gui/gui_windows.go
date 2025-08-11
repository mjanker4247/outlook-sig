//go:build windows
// +build windows

package gui

import "fmt"

// ShowGUI displays a message that GUI is not available during cross-compilation
func ShowGUI() {
	fmt.Println("GUI mode is not available during cross-compilation.")
	fmt.Println("Please use CLI mode with appropriate flags:")
	fmt.Println("  SignatureInstaller.exe --name \"Your Name\" --email \"your.email@example.com\" --phone \"+49 123 456789\"")
	fmt.Println("")
	fmt.Println("Or build on Windows to enable GUI mode.")
}
