package main

import(
	"fmt"
	"os"
	"time"
	"syscall"
)

type FileTime struct {
	MTime time.Time
	CTime time.Time
	ATime time.Time
}

func main(){

	file, err := FTime(os.TempDir())
	if err == nil {
		fmt.Println("Mo Time", file.MTime)
		fmt.Println("Ac Time", file.ATime)
		fmt.Println("Cr Time", file.CTime)
	}
}

// Gets the Modified, Create and Access time of a file
func FTime(file string) (t *FileTime, err error) {
	fileinfo, err := os.Stat(file)
	if err != nil {
	    return
	}
	t = new(FileTime)
	var stat = fileinfo.Sys().(*syscall.Stat_t)
	t.ATime = time.Unix(stat.Atim.Sec, stat.Atim.Nsec)
	t.CTime = time.Unix(stat.Ctim.Sec, stat.Ctim.Nsec)
	t.MTime = time.Unix(stat.Mtim.Sec, stat.Mtim.Nsec)
	return
}
