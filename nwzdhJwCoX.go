func ptraceWait(pid int) (wpid int, status syscall.WaitStatus, err error) {
	wpid, err = syscall.Wait4(pid, &status, syscall.WALL, nil)
	if err != nil {
		return 0, 0, err
	}
	fmt.Printf("\t\twait: wpid=%5d, status=0x%06x\n", wpid, status)
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
		return fmt.Errorf("peek: got %d bytes, want %d", len(data))
	}
	return nil
}

func ptracePoke(pid int, addr uint64, data []byte) error {
	n, err := syscall.PtracePokeText(pid, uintptr(addr), data)
	if err != nil {
		return err
	}
	if n != len(data) {
		return fmt.Errorf("poke: got %d bytes, want %d", len(data))
	}
	return nil
}

func ptraceStep(c *cpu.CPU, m *memory.Memory, p *system.Process) {

	runtime.LockOSThread()

	args := []string{p.Executable}
	args = append(args, strings.Fields(p.CommandLine)...)
	proc, err := os.StartProcess(p.Executable, args, &os.ProcAttr{
		Files: []*os.File{
			os.Stdin,
			os.Stdout,
			os.Stderr,
		},
		Sys: &syscall.SysProcAttr{
			Ptrace:    true,
			Pdeathsig: syscall.SIGKILL,
		},
	})

	if err != nil {
		log.Fatalf("StartProcess: %v", err)
	}

	fmt.Printf("\tproc.Pid=%d\n", proc.Pid)

	_, status, err := ptraceWait(proc.Pid)
	if err != nil {
		log.Fatalf("wait0: %v", err)
	}
	if status != 0x00057f { // 0x05=SIGTRAP, 0x7f=stopped.
		log.Fatalf("status: got %#x, want %#x", status, 0x57f)
	}
	err = syscall.PtraceSetOptions(proc.Pid, syscall.PTRACE_O_TRACECLONE|syscall.PTRACE_O_TRACEEXIT)
	if err != nil {
		log.Fatalf("PtraceSetOptions: %v", err)
	}

	var buf [4]byte
	if err := ptracePeek(proc.Pid, p.Entry, buf[:]); err != nil {
		log.Fatalf("peek: %v", err)
	}
	breakpoints := map[uint64]breakpoint{
		p.Entry: {pc: p.Entry, origInstr: buf},
	}
	buf = [4]byte{0x0, 0x0, 0x20, 0xd4} // ARMv8 breakpoint op.
	if err := ptracePoke(proc.Pid, p.Entry, buf[:]); err != nil {
		log.Fatalf("poke: %v", err)
	}

	err = syscall.PtraceCont(proc.Pid, 0)
	if err != nil {
		log.Fatalf("PtraceCont: %v", err)
	}

	for {
		pid, status, err := ptraceWait(-1)
		if err != nil {
			log.Fatalf("wait1: %v", err)
		}

		switch status {
		case 0x00057f: // 0x05=SIGTRAP, 0x7f=stopped.
			regs := syscall.PtraceRegs{}

			/* *** THIS CALL RETURNS AN I/O ERROR *** */
			if err := syscall.PtraceGetRegs(pid, &regs); err != nil {
				log.Printf("PtraceGetRegs: %v", err)
			}

			log.Println(regs)

			// if a breakpoint, replay the instruction we replaced.
			if bp, ok := breakpoints[regs.Pc-uint64(4)]; ok {
				regs.Pc -= uint64(4)
				if err := syscall.PtraceSetRegs(pid, &regs); err != nil {
					log.Fatalf("PtraceSetRegs: %v", err)
				}
				if !ok {
					log.Fatalf("no breakpoint for address %#x\n", regs.Pc)
				}
				buf = bp.origInstr
				if err := ptracePoke(pid, p.Entry, buf[:]); err != nil {
					log.Fatalf("poke: %v", err)
				}
				fmt.Printf("\thit breakpoint at %#x, pid=%5d\n", regs.PC, pid)
				if err := syscall.PtraceSingleStep(pid); err != nil {
					log.Fatalf("PtraceSingleStep: %v", err)
				}
				_, status, err := ptraceWait(pid)
				if err != nil {
					log.Fatalf("wait2: %v", err)
				}
				if status != 0x00057f {
					log.Fatalf("PtraceSingleStep: unexpected status %#x\n", status)
				}
			}

			// single step the execution
			if err := syscall.PtraceSingleStep(pid); err != nil {
				log.Fatalf("PtraceSingleStep: %v", err)
			}
			_, status, err := ptraceWait(pid)
			if err != nil {
				log.Fatalf("wait3: %v", err)
			}
			if status != 0x00057f {
				log.Fatalf("PtraceSingleStep: unexpected status %#x\n", status)
			}

		default:
			log.Fatalf("Status unhandled %d\n", status)
			break
		}
	}
}