func ptraceWait(pid int) (wpid int, status syscall.WaitStatus, err error) {
	wpid, err = syscall.Wait4(pid, &amp;status, syscall.WALL, nil)
	if err != nil {
		return 0, 0, err
	}
	fmt.Printf(&#34;\t\twait: wpid=%5d, status=0x%06x\n&#34;, wpid, status)
	return wpid, status, nil
}

type breakpoint struct {
	pc        uint64
	origInstr [4]byte
}

func ptracePeek(pid int, addr uint64, data []byte) error {
	n, err := syscall.PtracePeekText(pid, uintptr(addr), data)
	if err != nil {
		return err
	}
	if n != len(data) {
		return fmt.Errorf(&#34;peek: got %d bytes, want %d&#34;, len(data))
	}
	return nil
}

func ptracePoke(pid int, addr uint64, data []byte) error {
	n, err := syscall.PtracePokeText(pid, uintptr(addr), data)
	if err != nil {
		return err
	}
	if n != len(data) {
		return fmt.Errorf(&#34;poke: got %d bytes, want %d&#34;, len(data))
	}
	return nil
}

func ptraceStep(c *cpu.CPU, m *memory.Memory, p *system.Process) {

	runtime.LockOSThread()

	args := []string{p.Executable}
	args = append(args, strings.Fields(p.CommandLine)...)
	proc, err := os.StartProcess(p.Executable, args, &amp;os.ProcAttr{
		Files: []*os.File{
			os.Stdin,
			os.Stdout,
			os.Stderr,
		},
		Sys: &amp;syscall.SysProcAttr{
			Ptrace:    true,
			Pdeathsig: syscall.SIGKILL,
		},
	})

	if err != nil {
		log.Fatalf(&#34;StartProcess: %v&#34;, err)
	}

	fmt.Printf(&#34;\tproc.Pid=%d\n&#34;, proc.Pid)

	_, status, err := ptraceWait(proc.Pid)
	if err != nil {
		log.Fatalf(&#34;wait0: %v&#34;, err)
	}
	if status != 0x00057f { // 0x05=SIGTRAP, 0x7f=stopped.
		log.Fatalf(&#34;status: got %#x, want %#x&#34;, status, 0x57f)
	}
	err = syscall.PtraceSetOptions(proc.Pid, syscall.PTRACE_O_TRACECLONE|syscall.PTRACE_O_TRACEEXIT)
	if err != nil {
		log.Fatalf(&#34;PtraceSetOptions: %v&#34;, err)
	}

	var buf [4]byte
	if err := ptracePeek(proc.Pid, p.Entry, buf[:]); err != nil {
		log.Fatalf(&#34;peek: %v&#34;, err)
	}
	breakpoints := map[uint64]breakpoint{
		p.Entry: {pc: p.Entry, origInstr: buf},
	}
	buf = [4]byte{0x0, 0x0, 0x20, 0xd4} // ARMv8 breakpoint op.
	if err := ptracePoke(proc.Pid, p.Entry, buf[:]); err != nil {
		log.Fatalf(&#34;poke: %v&#34;, err)
	}

	err = syscall.PtraceCont(proc.Pid, 0)
	if err != nil {
		log.Fatalf(&#34;PtraceCont: %v&#34;, err)
	}

	for {
		pid, status, err := ptraceWait(-1)
		if err != nil {
			log.Fatalf(&#34;wait1: %v&#34;, err)
		}

		switch status {
		case 0x00057f: // 0x05=SIGTRAP, 0x7f=stopped.
			regs := syscall.PtraceRegs{}

			/* *** THIS CALL RETURNS AN I/O ERROR *** */
			if err := syscall.PtraceGetRegs(pid, &amp;regs); err != nil {
				log.Printf(&#34;PtraceGetRegs: %v&#34;, err)
			}

			log.Println(regs)

			// if a breakpoint, replay the instruction we replaced.
			if bp, ok := breakpoints[regs.Pc-uint64(4)]; ok {
				regs.Pc -= uint64(4)
				if err := syscall.PtraceSetRegs(pid, &amp;regs); err != nil {
					log.Fatalf(&#34;PtraceSetRegs: %v&#34;, err)
				}
				if !ok {
					log.Fatalf(&#34;no breakpoint for address %#x\n&#34;, regs.Pc)
				}
				buf = bp.origInstr
				if err := ptracePoke(pid, p.Entry, buf[:]); err != nil {
					log.Fatalf(&#34;poke: %v&#34;, err)
				}
				fmt.Printf(&#34;\thit breakpoint at %#x, pid=%5d\n&#34;, regs.PC, pid)
				if err := syscall.PtraceSingleStep(pid); err != nil {
					log.Fatalf(&#34;PtraceSingleStep: %v&#34;, err)
				}
				_, status, err := ptraceWait(pid)
				if err != nil {
					log.Fatalf(&#34;wait2: %v&#34;, err)
				}
				if status != 0x00057f {
					log.Fatalf(&#34;PtraceSingleStep: unexpected status %#x\n&#34;, status)
				}
			}

			// single step the execution
			if err := syscall.PtraceSingleStep(pid); err != nil {
				log.Fatalf(&#34;PtraceSingleStep: %v&#34;, err)
			}
			_, status, err := ptraceWait(pid)
			if err != nil {
				log.Fatalf(&#34;wait3: %v&#34;, err)
			}
			if status != 0x00057f {
				log.Fatalf(&#34;PtraceSingleStep: unexpected status %#x\n&#34;, status)
			}

		default:
			log.Fatalf(&#34;Status unhandled %d\n&#34;, status)
			break
		}
	}
}