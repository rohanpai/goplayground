// Copyright 2015 go-smpp authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package smpptest

import (
	"net"
	"strconv"
	"sync/atomic"
	"testing"

	"github.com/fiorix/go-smpp/smpp/pdu"
	"github.com/fiorix/go-smpp/smpp/pdu/pdufield"
	"github.com/fiorix/go-smpp/smpp/pdu/pdutext"
)

var msgIDcounter int64 = 0

func submitSM(m pdu.Body) (string, error) {
	// submit message using fields from m, return message ID...
	return strconv.FormatInt(atomic.AddInt64(&msgIDcounter, 1), 10), nil
}

// SubmitSMHandler handles SubmitSM and returns SubmitSMResp.
func SubmitSMHandler(cli Conn, m pdu.Body) {
	msgID, _ := submitSM(m)
	resp := pdu.NewSubmitSMResp()
	resp.Header().Seq = m.Header().Seq
	resp.Fields().Set(pdufield.MessageID, msgID)
	cli.Write(resp)
}

func TestFakeServer(t *testing.T) {
	s := NewServer()
	s.Handler = SubmitSMHandler
	defer s.Close()
	c, err := net.Dial("tcp", s.Addr())
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()
	rw := newConn(c)
	// bind
	p := pdu.NewBindTransmitter()
	f := p.Fields()
	f.Set(pdufield.SystemID, "client")
	f.Set(pdufield.Password, "secret")
	f.Set(pdufield.InterfaceVersion, 0x34)
	if err = rw.Write(p); err != nil {
		t.Fatal(err)
	}
	// bind resp
	resp, err := rw.Read()
	if err != nil {
		t.Fatal(err)
	}
	id, ok := resp.Fields()[pdufield.SystemID]
	if !ok {
		t.Fatalf("missing system_id field: %#v", resp)
	}
	if id.String() != "smpptest" {
		t.Fatalf("unexpected system_id: want smpptest, have %q", id)
	}
	// submit_sm
	p = pdu.NewSubmitSM()
	f = p.Fields()
	f.Set(pdufield.SourceAddr, "foobar")
	f.Set(pdufield.DestinationAddr, "bozo")
	f.Set(pdufield.ShortMessage, pdutext.Latin1("Lorem ipsum"))
	if err = rw.Write(p); err != nil {
		t.Fatal(err)
	}
	// same submit_sm
	resp, err = rw.Read()
	if err != nil {
		t.Fatal(err)
	}
	f = resp.Fields()
	if id := f[pdufield.MessageID]; id.String() != "1" {
		t.Fatalf("unexpected messageID field data: want 1, have %q", id)
	}
}