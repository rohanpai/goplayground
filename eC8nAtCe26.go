package main

import (
        &#34;github.com/mikioh/tunnel&#34;
        &#34;log&#34;
        &#34;net&#34;
        &#34;os&#34;
        &#34;os/exec&#34;
)

func main() {
        if os.Getuid() != 0 {
                log.Fatalf(&#34;need administrator privlege&#34;)
        }

        c, err := tunnel.New()
        if err != nil {
                log.Fatalf(&#34;tunnel.New failed: %v&#34;, err)
        }
        ifi, err := c.Interface()
        if err != nil {
                log.Fatalf(&#34;tunnel.Interface failed: %v&#34;, err)
        }
        if err := setup(ifi.Name); err != nil {
                log.Fatalf(&#34;platform dependent setup failed: %v&#34;, err)
        }
        defer func() {
                c.Close()
                teardown(ifi.Name)
        }()

        ifas, err := ifi.Addrs()
        if err != nil {
                log.Fatalf(&#34;Interface.Addrs failed: %v&#34;, err)
        }
        for _, ifa := range ifas {
                if ip := ifa.(*net.IPNet).IP.String(); ip == src {
                        log.Printf(&#34;fixed: %v, %v&#34;, ip, src)
                        return
                }
        }
        log.Println(&#34;still need investigation&#34;)
}

var (
        src = &#34;169.254.0.1&#34;
        dst = &#34;169.254.0.254&#34;
)

func setup(name string) error {
        cmd := exec.Command(&#34;ifconfig&#34;, name, &#34;inet&#34;, src, &#34;dstaddr&#34;, dst)
        if err := cmd.Run(); err != nil {
                return err
        }
        return nil
}

func teardown(name string) error {
        cmd := exec.Command(&#34;ifconfig&#34;, name, &#34;delete&#34;)
        if err := cmd.Run(); err != nil {
                return err
        }
        return nil
}