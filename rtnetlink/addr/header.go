package addr

/*
  Copyright (c) 2011, Abneptis LLC. All rights reserved.
  Original Author: James D. Nurmi <james@abneptis.com>

  See LICENSE for details
*/

import (
	"encoding/binary"
	"errors"
)
import "github.com/vishvananda/go-netlink/rtnetlink"
import "github.com/vishvananda/go-netlink"

const HEADER_LENGTH = 8

type Header [HEADER_LENGTH]byte

func NewHeader(afam rtnetlink.Family, pl uint8, fl Flags, scope rtnetlink.Scope, ifindex uint32) *Header {
	hdr := Header{byte(afam), pl, byte(fl), byte(scope)}
	binary.LittleEndian.PutUint32(hdr[4:8], ifindex)
	return &hdr
}

func (self Header) Len() int                        { return HEADER_LENGTH }
func (self Header) AddressFamily() rtnetlink.Family { return rtnetlink.Family(self[0]) }
func (self Header) PrefixLength() uint8             { return self[1] }
func (self Header) Flags() Flags                    { return Flags(self[2]) }
func (self Header) Scope() rtnetlink.Scope          { return rtnetlink.Scope(self[3]) }
func (self Header) InterfaceIndex() uint32          { return binary.LittleEndian.Uint32(self[4:8]) }

func (self *Header) UnmarshalNetlink(in []byte) (err error) {
	if len(in) != HEADER_LENGTH {
		err = errors.New("Wrong length for Header")
	} else {
		copy(self[0:HEADER_LENGTH], in[0:HEADER_LENGTH])
	}
	return
}

func (self Header) MarshalNetlink() (out []byte, err error) {
	out = netlink.Padded(self[0:HEADER_LENGTH])
	return
}
