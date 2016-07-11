package main

import (
	&#34;fmt&#34;
	&#34;log&#34;
	&#34;syscall&#34;
	&#34;unsafe&#34;
)

var (
	user32             = syscall.MustLoadDLL(&#34;user32.dll&#34;)
	procEnumWindows    = user32.MustFindProc(&#34;EnumWindows&#34;)
	procGetWindowTextW = user32.MustFindProc(&#34;GetWindowTextW&#34;)
)

func EnumWindows(enumFunc uintptr, lparam uintptr) (err error) {
	r1, _, e1 := syscall.Syscall(procEnumWindows.Addr(), 2, uintptr(enumFunc), uintptr(lparam), 0)
	if r1 == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func GetWindowText(hwnd syscall.Handle, str *uint16, maxCount int32) (len int32, err error) {
	r0, _, e1 := syscall.Syscall(procGetWindowTextW.Addr(), 3, uintptr(hwnd), uintptr(unsafe.Pointer(str)), uintptr(maxCount))
	len = int32(r0)
	if len == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func FindWindow(title string) (syscall.Handle, error) {
	var hwnd syscall.Handle
	cb := syscall.NewCallback(func(h syscall.Handle, p uintptr) uintptr {
		b := make([]uint16, 200)
		_, err := GetWindowText(h, &amp;b[0], int32(len(b)))
		if err != nil {
			// ignore the error
			return 1 // continue enumeration
		}
		if syscall.UTF16ToString(b) == title {
			// note the window
			hwnd = h
			return 0 // stop enumeration
		}
		return 1 // continue enumeration
	})
	EnumWindows(cb, 0)
	if hwnd == 0 {
		return 0, fmt.Errorf(&#34;No window with title &#39;%s&#39; found&#34;, title)
	}
	return hwnd, nil
}

func main() {
	const title = &#34;Game&#34;
	h, err := FindWindow(title)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf(&#34;Found &#39;%s&#39; window: handle=0x%x\n&#34;, title, h)
}
