package main

import (
	&#34;fmt&#34;
	&#34;syscall&#34;
	&#34;unsafe&#34;

	&#34;github.com/lxn/win&#34;
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
		fmt.Printf(&#34;%d:%s &#34;, win.pid, win.process)
		fmt.Printf(&#34;[%X]&#34;, win.hwnd)
		fmt.Printf(&#34; %s (%s) %d,%d %dx%d\n&#34;,
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
	user32        = syscall.NewLazyDLL(&#34;user32.dll&#34;)
	getWindowText = user32.NewProc(&#34;GetWindowTextW&#34;)
)

const bufSiz = 128 // Max length I want to see

func getName(hwnd win.HWND, get *syscall.LazyProc) string {
	var buf [bufSiz]uint16
	siz, _, _ := get.Call(uintptr(hwnd), uintptr(unsafe.Pointer(&amp;buf[0])), uintptr(len(buf)))
	if siz == 0 {
		return &#34;&#34;
	}
	name := syscall.UTF16ToString(buf[:siz])
	if siz == bufSiz-1 {
		name = name &#43; &#34;\u22EF&#34;
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
	win.GetWindowRect(hwnd, &amp;w.r)
	w.name = getName(hwnd, getWindowText)
	w.hasChild = win.GetWindow(hwnd, win.GW_CHILD) != 0
	d.list = append(d.list, w)
	return 1
}

func listWindows(hwnd win.HWND) []window {
	var d cbData
	d.list = make([]window, 0)
	d.pid = make(map[uint32]string)
	win.EnumChildWindows(hwnd, syscall.NewCallback(perWindow), uintptr(unsafe.Pointer(&amp;d)))
	return d.list
}
