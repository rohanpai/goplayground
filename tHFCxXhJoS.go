package main

import (
        "net"
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
        if remote == "" {
                remote = "127.0.0.1:2947"
        }
        conn, err := net.Dial("tcp", remote)
        gps := new(GPS)
        gps.conn = conn
        gps.Remote = remote
        return gps, err
}

func main () {
}
