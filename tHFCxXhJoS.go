package main

import (
        &#34;net&#34;
)

type GPSFix struct {
        latitude        float64
        longitude       float64
        speed           float64
        climb           float64 
        time            float64
}

type GPS struct {
        conn            net.Conn
        Fix             GPSFix
        Remote          string
}

func NewGPS (remote string) (GPS, err error) {
        if remote == &#34;&#34; {
                remote = &#34;127.0.0.1:2947&#34;
        }
        conn, err := net.Dial(&#34;tcp&#34;, remote)
        gps := new(GPS)
        gps.conn = conn
        gps.Remote = remote
        return gps, err
}

func main () {
}
