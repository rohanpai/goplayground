package main

import (
	&#34;fmt&#34;
	&#34;os&#34;
	&#34;runtime&#34;
	&#34;syscall&#34;

	&#34;github.com/hjr265/ptrace.go/ptrace&#34;
)

type RunningObject struct {
	Time        syscall.Timeval
	TimeLimit   int64
	MemoryLimit int64
	Memory      int64
}

func (r *RunningObject) Millisecond() int64 {
	return r.Time.Sec*1000 &#43; r.Time.Usec/1000
}

func Run(src string, args []string) *RunningObject {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	var rusage syscall.Rusage
	var incall = true
	var runningObject RunningObject
	proc, err := os.StartProcess(src, args, &amp;os.ProcAttr{Sys: &amp;syscall.SysProcAttr{
		Ptrace: true},
	})
	if err != nil {
		panic(err)
	}
	//set CPU time limit
	var rlimit syscall.Rlimit
	rlimit.Cur = 1
	rlimit.Max = 1 &#43; 1
	err = prLimit(proc.Pid, syscall.RLIMIT_CPU, &amp;rlimit)
	if err != nil {
		fmt.Println(err)
		return &amp;runningObject
	}
	rlimit.Cur = 1024
	rlimit.Max = 1024 &#43; 1024
	err = prLimit(proc.Pid, syscall.RLIMIT_DATA, &amp;rlimit)
	if err != nil {
		fmt.Println(err)
		return &amp;runningObject
	}
	err = prLimit(proc.Pid, syscall.RLIMIT_STACK, &amp;rlimit)
	if err != nil {
		fmt.Println(err)
		return &amp;runningObject
	}
	tracer, err := ptrace.Attach(proc)
	if err != nil {
		panic(err)
	}
	for {
		status := syscall.WaitStatus(0)
		_, err := syscall.Wait4(proc.Pid, &amp;status, syscall.WSTOPPED, &amp;rusage)
		if err != nil {
			panic(err)
		}
		if status.Exited() {
			fmt.Println(&#34;exit&#34;)
			fmt.Println(rusage.Stime)
			return &amp;runningObject
		}
		if status.CoreDump() {
			fmt.Println(&#34;CoreDump&#34;)
			return &amp;runningObject
		}
		if status.Continued() {
			fmt.Println(&#34;Continued&#34;)
			return &amp;runningObject
		}
		if status.Signaled() {
			return &amp;runningObject
		}
		if status.Stopped() &amp;&amp; status.StopSignal() != syscall.SIGTRAP {
			switch status.StopSignal() {
			case syscall.SIGXCPU:
				fmt.Println(&#34;SIGXCPU&#34;)
				runningObject.Time = rusage.Utime
				fmt.Println(runningObject.Millisecond())
				return &amp;runningObject
			case syscall.SIGSEGV:
				fmt.Println(&#34;SIGSEGV&#34;)
				runningObject.Memory = rusage.Minflt
				fmt.Println(runningObject.Memory)
				return &amp;runningObject
			default:
				fmt.Println(&#34;default&#34;)
			}
			return &amp;runningObject
		} else {
			regs, err := tracer.GetRegs()
			if err != nil {
				panic(err)
			}
			if regs.Orig_rax == syscall.SYS_WRITE {
				if incall {
					incall = false

					_, err = tracer.GetRegs()
					if err != nil {
						panic(err)
					}
					fmt.Printf(&#34;The child made a system call with, %d,%d,%d \n&#34;, regs.Rdi, regs.Rsi, regs.Rdx)
				} else {
					incall = true
					regs, err := tracer.GetRegs()
					if err != nil {
						panic(err)
					}
					fmt.Printf(&#34;write returned %v\n&#34;, regs.Rax)
					fmt.Printf(&#34;call %d\n&#34;, regs.Rdi)
				}
			}
		}
		//0表示不发出信号
		err = tracer.Syscall(syscall.Signal(0))
		if err != nil {
			fmt.Println(err)
		}
	}
}

func main() {
	Run(&#34;/bin/sleep&#34;, []string{&#34;sleep&#34;, 5})
}
