package main

import (
  "github.com/AllenDang/w32"
  "syscall"
  "time"
  "unsafe"
)

var (
  moduser32               = syscall.NewLazyDLL("user32.dll")
  procSetForegroundWindow = moduser32.NewProc("SetForegroundWindow")
  procFindWindowW         = moduser32.NewProc("FindWindowW")
)

const (
  KEYEVENTF_KEYDOWN     = 0 //key UP
  KEYEVENTF_EXTENDEDKEY = 0x0001
  KEYEVENTF_KEYUP       = 0x0002 //key UP
  KEYEVENTF_UNICODE     = 0x0004
  KEYEVENTF_SCANCODE    = 0x0008 // scancode
)

type HWND uintptr

func auto() {
  // hwnd, err := FindWindow("Notepad", "Untitled - Notepad")
  // if err != nil {
    // log.Fatal(err)
  // }

  // SetForegroundWindow(hwnd)

  time.Sleep(time.Second * 3)

  sendkeys("Hello there!")
  sendkey(w32.VK_RETURN)

  time.Sleep(time.Second * 1)
}

func main() {
  auto()
}

func sendkey(vk uint16) {
  var inputs []w32.INPUT
  inputs = append(inputs, w32.INPUT{
    Type: w32.INPUT_KEYBOARD,
    Ki: w32.KEYBDINPUT{
      WVk:         vk,
      WScan:       0,
      DwFlags:     KEYEVENTF_KEYDOWN,
      Time:        0,
      DwExtraInfo: 0,
    },
  })

  inputs = append(inputs, w32.INPUT{
    Type: w32.INPUT_KEYBOARD,
    Ki: w32.KEYBDINPUT{
      WVk:         vk,
      WScan:       0,
      DwFlags:     KEYEVENTF_KEYUP,
      Time:        0,
      DwExtraInfo: 0,
    },
  })
  w32.SendInput(inputs)
}

func sendkeys(str string) {
  var inputs []w32.INPUT
  for _, s := range str {
    inputs = append(inputs, w32.INPUT{
      Type: w32.INPUT_KEYBOARD,
      Ki: w32.KEYBDINPUT{
        WVk:         0,
        WScan:       uint16(s),
        DwFlags:     KEYEVENTF_KEYDOWN | KEYEVENTF_UNICODE,
        Time:        0,
        DwExtraInfo: 0,
      },
    })

    inputs = append(inputs, w32.INPUT{
      Type: w32.INPUT_KEYBOARD,
      Ki: w32.KEYBDINPUT{
        WVk:         0,
        WScan:       uint16(s),
        DwFlags:     KEYEVENTF_KEYUP | KEYEVENTF_UNICODE,
        Time:        0,
        DwExtraInfo: 0,
      },
    })
  }

  w32.SendInput(inputs)
}

func SetForegroundWindow(hwnd HWND) bool {
  ret, _, _ := procSetForegroundWindow.Call(
    uintptr(hwnd))

  return ret != 0
}

func FindWindow(cls string, win string) (ret HWND, err error) {
  lpszClass := syscall.StringToUTF16Ptr(cls)
  lpszWindow := syscall.StringToUTF16Ptr(win)

  r0, _, e1 := syscall.Syscall(procFindWindowW.Addr(), 2, uintptr(unsafe.Pointer(lpszClass)), uintptr(unsafe.Pointer(lpszWindow)), 0)
  ret = HWND(r0)
  if ret == 0 {
    if e1 != 0 {
      err = error(e1)
    } else {
      err = syscall.EINVAL
    }
  }
  return
}
