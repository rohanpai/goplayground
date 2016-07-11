// Copyright 2012 The golibpcap Authors. All rights reserved.                      
// Use of this source code is governed by a BSD-style                              
// license that can be found in the LICENSE file.

package pkt

/*
#include <net/ethernet.h>
#include <netinet/ether.h>
#include <netinet/in.h>

// http://standards.ieee.org/getieee802/download/802.1Q-2011.pdf
#define _ETHERTYPE_VLAN_CTAG ETHERTYPE_VLAN
#define _ETHERTYPE_VLAN_STAG 0x88a8
#define _ETHERTYPE_VLAN_ITAG 0x88e7

struct __attribute__((packed)) vlan_ctag_stub {
	uint16_t tci;
	uint16_t ether_type;
};

struct __attribute__((packed)) vlan_stag_stub {
	uint16_t service_tci;
	uint16_t customer_tpid;
	uint16_t customer_tci;
	uint16_t ether_type;
};

struct __attribute__((packed)) vlan_itag_stub {
	uint8_t tci[16];
	uint16_t ether_type;
};
*/
import "C"

import (
	"fmt"
	"net"
	"unsafe"
)

// The EthHdr struct is a wrapper for the ether_header struct in <net/ethernet.h>.
type EthHdr struct {
	cptr      *C.struct_ether_header // C pointer to ether_header
	SrcAddr   net.HardwareAddr       // the sender's MAC address
	DstAddr   net.HardwareAddr       // the receiver's MAC address
	EtherType uint16                 // packet type ID field
	Tpid      uint16                 // IEEE 802.1Q Tag Protocol Identifier. Either 0 or one of _ETHERTYPE_VLAN_{C,S,I}TAG
	Tag       interface{}            // IEEE 802.1Q Tag if present, on of *VlanCTag, *VlanSTag, *VlanITag
	payload   unsafe.Pointer
}

// C-TAG format
type VlanCTag struct {
	Tci uint16
}

// S-TAG format
type VlanSTag struct {
	ServiceTci   uint16
	CustomerTpid uint16
	CustomerTci  uint16
}

// I-TAG format
// TODO: add single fields
type VlanITag struct {
	tci []byte
}

// With an unsafe.Pointer to the block of C memory NewEthHdr returns a filled in EthHdr struct.
func NewEthHdr(p unsafe.Pointer) (*EthHdr, unsafe.Pointer) {
	ethHdr := &EthHdr{
		cptr:    (*C.struct_ether_header)(p),
		payload: unsafe.Pointer(uintptr(p) + uintptr(C.ETHER_HDR_LEN)),
	}
	ethHdr.SrcAddr = net.HardwareAddr(C.GoBytes(unsafe.Pointer(&ethHdr.cptr.ether_shost), C.ETH_ALEN))
	ethHdr.DstAddr = net.HardwareAddr(C.GoBytes(unsafe.Pointer(&ethHdr.cptr.ether_dhost), C.ETH_ALEN))
	ethHdr.EtherType = uint16(C.ntohs(C.uint16_t(ethHdr.cptr.ether_type)))

	switch ethHdr.EtherType {
	case C._ETHERTYPE_VLAN_CTAG:
		ptr := (*C.struct_vlan_ctag_stub)(ethHdr.payload)

		ethHdr.EtherType = uint16(C.ntohs(ptr.ether_type))
		ethHdr.Tpid = C._ETHERTYPE_VLAN_CTAG
		ethHdr.Tag = &VlanCTag{
			Tci: uint16(C.ntohs(ptr.tci)),
		}
		ethHdr.payload = unsafe.Pointer(uintptr(ethHdr.payload) + C.sizeof_struct_vlan_ctag_stub)
	case C._ETHERTYPE_VLAN_STAG:
		ptr := (*C.struct_vlan_stag_stub)(ethHdr.payload)

		ethHdr.EtherType = uint16(C.ntohs(ptr.ether_type))
		ethHdr.Tpid = C._ETHERTYPE_VLAN_STAG
		ethHdr.Tag = VlanSTag{
			ServiceTci:   uint16(C.ntohs(ptr.service_tci)),
			CustomerTpid: uint16(C.ntohs(ptr.customer_tpid)),
			CustomerTci:  uint16(C.ntohs(ptr.customer_tci)),
		}
		ethHdr.payload = unsafe.Pointer(uintptr(ethHdr.payload) + C.sizeof_struct_vlan_stag_stub)
	case C._ETHERTYPE_VLAN_ITAG:
		ptr := (*C.struct_vlan_itag_stub)(ethHdr.payload)

		ethHdr.EtherType = uint16(C.ntohs(ptr.ether_type))
		ethHdr.Tpid = C._ETHERTYPE_VLAN_ITAG
		ethHdr.Tag = &VlanITag{
			tci: C.GoBytes(unsafe.Pointer(&ptr.tci), 16), // 16 length on octets of the TCI field in the I-Tag format
		}
		ethHdr.payload = unsafe.Pointer(uintptr(ethHdr.payload) + C.sizeof_struct_vlan_itag_stub)
	}
	return ethHdr, ethHdr.payload
}

// JsonElement returns a JSON encoding of the EthHdr struct.
func (h *EthHdr) JsonElement() string {
	return fmt.Sprintf("\"ether_header\":{\"ether_shost\":\"%s\",\"ether_dhost\":\"%s\",\"ether_type\":%d}",
		h.SrcAddr.String(),
		h.DstAddr.String(),
		h.EtherType)
}

// CsvElement returns a CSV encoding of the EthHdr struct.
// The string "ETH" signifies the beginning of the EthHdr.
func (h *EthHdr) CsvElement() string {
	return fmt.Sprintf("\"ETH\",\"%s\",\"%s\",%d",
		h.SrcAddr.String(),
		h.DstAddr.String(),
		h.EtherType)
}

// String returns a minimal encoding of the EthHdr struct.
func (h *EthHdr) String() string {
	return fmt.Sprintf("%s->%s %#x",
		h.SrcAddr.String(),
		h.DstAddr.String(),
		h.EtherType)
}
