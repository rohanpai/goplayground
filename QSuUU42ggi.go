package main

// // Flags are based on python3.3-config --cflags and --ldflags
// #cgo darwin CFLAGS:-I/usr/local/Cellar/python3/3.3.0/Frameworks/Python.framework/Versions/3.3/include/python3.3m -I/usr/local/Cellar/python3/3.3.0/Frameworks/Python.framework/Versions/3.3/include/python3.3m -fno-common -dynamic -DNDEBUG -g -fwrapv -O3 -Wall -Wstrict-prototypes -isysroot /Applications/Xcode.app/Contents/Developer/Platforms/MacOSX.platform/Developer/SDKs/MacOSX10.8.sdk -arch x86_64 -isysroot /Applications/Xcode.app/Contents/Developer/Platforms/MacOSX.platform/Developer/SDKs/MacOSX10.8.sdk -I/usr/local/include -I/usr/local/opt/sqlite/include
// #cgo darwin LDFLAGS: -L. -L/usr/local/Cellar/python3/3.3.0/Frameworks/Python.framework/Versions/3.3/lib/python3.3/config-3.3m -ldl -framework CoreFoundation -lpython3.3m
// #include <Python.h>
import "C"

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
	"unsafe"
)

type syscall_Termios struct {
	Iflag     uint64
	Oflag     uint64
	Cflag     uint64
	Lflag     uint64
	Cc        [20]uint8
	Pad_cgo_0 [4]byte
	Ispeed    uint64
	Ospeed    uint64
}

const (
	syscall_IGNBRK = 0x1
	syscall_BRKINT = 0x2
	syscall_PARMRK = 0x8
	syscall_ISTRIP = 0x20
	syscall_INLCR  = 0x40
	syscall_IGNCR  = 0x80
	syscall_ICRNL  = 0x100
	syscall_IXON   = 0x200
	syscall_OPOST  = 0x1
	syscall_ECHO   = 0x8
	syscall_ECHONL = 0x10
	syscall_ICANON = 0x100
	syscall_ISIG   = 0x80
	syscall_IEXTEN = 0x400
	syscall_CSIZE  = 0x300
	syscall_PARENB = 0x1000
	syscall_CS8    = 0x300
	syscall_VMIN   = 0x10
	syscall_VTIME  = 0x11

	syscall_TCGETS = 0x40487413
	syscall_TCSETS = 0x80487414
)

func fcntl(fd int, cmd int, arg int) (val int, err error) {
	r, _, e := syscall.Syscall(syscall.SYS_FCNTL, uintptr(fd), uintptr(cmd), uintptr(arg))
	val = int(r)
	if e != 0 {
		panic(e)
	}
	return
}

func tcsetattr(fd int, termios *syscall_Termios) {
	r, _, e := syscall.Syscall(syscall.SYS_IOCTL, uintptr(fd), uintptr(syscall_TCSETS), uintptr(unsafe.Pointer(termios)))
	if r != 0 {
		panic(os.NewSyscallError("SYS_IOCTL", e))
	}
}

func tcgetattr(fd int, termios *syscall_Termios) {
	r, _, e := syscall.Syscall(syscall.SYS_IOCTL, uintptr(fd), uintptr(syscall_TCGETS), uintptr(unsafe.Pointer(termios)))
	if r != 0 {
		panic(os.NewSyscallError("SYS_IOCTL", e))
	}
}

func main() {
	C.Py_InitializeEx(0)
	defer C.Py_Finalize()
	var (
		in        int
		err       error
		sigio     = make(chan os.Signal)
		orig_tios syscall_Termios
	)
	in, err = syscall.Open("/dev/tty", syscall.O_RDONLY, 0)
	if err != nil {
		panic(err)
	}

	signal.Notify(sigio, syscall.SIGIO)

	fcntl(in, syscall.F_SETFL, syscall.O_ASYNC|syscall.O_NONBLOCK)
	tcgetattr(in, &orig_tios)
	tios := orig_tios
	tios.Iflag &^= syscall_BRKINT | syscall_IXON
	tios.Lflag &^= syscall_ECHO | syscall_ICANON |
		syscall_ISIG | syscall_IEXTEN
	tios.Cflag &^= syscall_CSIZE | syscall_PARENB
	tios.Cflag |= syscall_CS8
	tios.Cc[syscall_VMIN] = 1
	tios.Cc[syscall_VTIME] = 0

	tcsetattr(in, &tios)
	defer func() {
		tcsetattr(in, &orig_tios)
		syscall.Close(in)
	}()
	buf := make([]byte, 128)
	for {
		log.Println("Waiting on signal")
		select {
		case <-sigio:
			log.Println("Read..")
			n, err := syscall.Read(in, buf)
			log.Println("Received:", n, err)
			if err == syscall.EAGAIN || err == syscall.EWOULDBLOCK {
				break
			}
			if n > 0 {
				if buf[0] == 17 {
					return
				}
				log.Println("buf:", buf[:n])
			}
		case <-time.After(time.Second * 5):
			log.Println("Timed out")
			return
		}
	}
}
