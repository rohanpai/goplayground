package main

import (
	"fmt"
	"syscall"
	"unsafe"

	"github.com/lxn/win"
)

type window struct {
	hwnd        win.HWND
	pid         uint32
	name, class string
	process     string
	r           win.RECT
	visible     bool
	hasChild    bool
}

func dump() {
	l := listWindows(win.HWND(0))
	for _, win := range l {
		fmt.Printf("%d:%s ", win.pid, win.process)
		fmt.Printf("[%X]", win.hwnd)
		fmt.Printf(" %s (%s) %d,%d %dx%d\n",
			win.name, win.class,
			win.r.Left, win.r.Top,
			win.r.Right-win.r.Left, win.r.Bottom-win.r.Top)
	}
}

func main() {
	dump()
	fmt.Println()
	dump()
}

var (
	user32        = syscall.NewLazyDLL("user32.dll")
	getWindowText = user32.NewProc("GetWindowTextW")
)

const bufSiz = 128 // Max length I want to see

func getName(hwnd win.HWND, get *syscall.LazyProc) string {
	var buf [bufSiz]uint16
	siz, _, _ := get.Call(uintptr(hwnd), uintptr(unsafe.Pointer(&buf[0])), uintptr(len(buf)))
	if siz == 0 {
		return ""
	}
	name := syscall.UTF16ToString(buf[:siz])
	if siz == bufSiz-1 {
		name = name + "\u22EF"
	}
	return name
}

type cbData struct {
	list []window
	pid  map[uint32]string
}

func perWindow(hwnd win.HWND, param uintptr) uintptr {
	d := (*cbData)(unsafe.Pointer(param))
	w := window{hwnd: hwnd}
	w.visible = win.IsWindowVisible(hwnd)
	win.GetWindowRect(hwnd, &w.r)
	w.name = getName(hwnd, getWindowText)
	w.hasChild = win.GetWindow(hwnd, win.GW_CHILD) != 0
	d.list = append(d.list, w)
	return 1
}

func listWindows(hwnd win.HWND) []window {
	var d cbData
	d.list = make([]window, 0)
	d.pid = make(map[uint32]string)
	win.EnumChildWindows(hwnd, syscall.NewCallback(perWindow), uintptr(unsafe.Pointer(&d)))
	return d.list
}
