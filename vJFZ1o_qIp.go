package main

import(
	&#34;fmt&#34;
	&#34;os&#34;
	&#34;time&#34;
	&#34;syscall&#34;
)

type FileTime struct {
	MTime time.Time
	CTime time.Time
	ATime time.Time
}

func main(){

	file, err := FTime(os.TempDir())
	if err == nil {
		fmt.Println(&#34;Mo Time&#34;, file.MTime)
		fmt.Println(&#34;Ac Time&#34;, file.ATime)
		fmt.Println(&#34;Cr Time&#34;, file.CTime)
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
