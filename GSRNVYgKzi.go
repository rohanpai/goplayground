//
// demonstration that the cols and rows reported on windows console app will overflow to negative if they're too big and don't cast them
//
package main

import (
	"fmt"
	"syscall"
	"unsafe"
)

type (
	SHORT int16
	WORD  uint16

	SMALL_RECT struct {
		Left   SHORT
		Top    SHORT
		Right  SHORT
		Bottom SHORT
	}

	COORD struct {
		X SHORT
		Y SHORT
	}

	CONSOLE_SCREEN_BUFFER_INFO struct {
		Size              COORD
		CursorPosition    COORD
		Attributes        WORD
		Window            SMALL_RECT
		MaximumWindowSize COORD
	}
)

var kernel32DLL = syscall.NewLazyDLL("kernel32.dll")
var getConsoleScreenBufferInfoProc = kernel32DLL.NewProc("GetConsoleScreenBufferInfo")

func main() {
	stdoutHandle := getStdHandle(syscall.STD_OUTPUT_HANDLE)

	var info, err = GetConsoleScreenBufferInfo(stdoutHandle)

	if err != nil {
		panic("could not get console screen buffer info")
	}

	fmt.Printf("max x: %i max y: %i\n", info.MaximumWindowSize.X, info.MaximumWindowSize.Y)
	fmt.Printf("max x * max y = %i\n", info.MaximumWindowSize.X*info.MaximumWindowSize.Y)
	fmt.Printf("max x * max y = %i\n", int(info.MaximumWindowSize.X)*int(info.MaximumWindowSize.Y))

}

func getError(r1, r2 uintptr, lastErr error) error {
	// If the function fails, the return value is zero.
	if r1 == 0 {
		if lastErr != nil {
			return lastErr
		}
		return syscall.EINVAL
	}
	return nil
}

func getStdHandle(stdhandle int) uintptr {
	handle, err := syscall.GetStdHandle(stdhandle)
	if err != nil {
		panic(fmt.Errorf("could not get standard io handle %d", stdhandle))
	}
	return uintptr(handle)
}

// GetConsoleScreenBufferInfo retrieves information about the specified console screen buffer.
// http://msdn.microsoft.com/en-us/library/windows/desktop/ms683171(v=vs.85).aspx
func GetConsoleScreenBufferInfo(handle uintptr) (*CONSOLE_SCREEN_BUFFER_INFO, error) {
	var info CONSOLE_SCREEN_BUFFER_INFO
	if err := getError(getConsoleScreenBufferInfoProc.Call(handle, uintptr(unsafe.Pointer(&info)), 0)); err != nil {
		return nil, err
	}
	return &info, nil
}
