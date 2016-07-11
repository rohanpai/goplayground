package main

import (
        "github.com/mikioh/tunnel"
        "log"
        "net"
        "os"
        "os/exec"
)

func main() {
        if os.Getuid() != 0 {
                log.Fatalf("need administrator privlege")
        }

        c, err := tunnel.New()
        if err != nil {
                log.Fatalf("tunnel.New failed: %v", err)
        }
        ifi, err := c.Interface()
        if err != nil {
                log.Fatalf("tunnel.Interface failed: %v", err)
        }
        if err := setup(ifi.Name); err != nil {
                log.Fatalf("platform dependent setup failed: %v", err)
        }
        defer func() {
                c.Close()
                teardown(ifi.Name)
        }()

        ifas, err := ifi.Addrs()
        if err != nil {
                log.Fatalf("Interface.Addrs failed: %v", err)
        }
        for _, ifa := range ifas {
                if ip := ifa.(*net.IPNet).IP.String(); ip == src {
                        log.Printf("fixed: %v, %v", ip, src)
                        return
                }
        }
        log.Println("still need investigation")
}

var (
        src = "169.254.0.1"
        dst = "169.254.0.254"
)

func setup(name string) error {
        cmd := exec.Command("ifconfig", name, "inet", src, "dstaddr", dst)
        if err := cmd.Run(); err != nil {
                return err
        }
        return nil
}

func teardown(name string) error {
        cmd := exec.Command("ifconfig", name, "delete")
        if err := cmd.Run(); err != nil {
                return err
        }
        return nil
}