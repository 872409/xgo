package xton

import (
	"encoding/hex"
	"github.com/xssnick/tonutils-go/address"
)

func WithAddress(addr *address.Address) *TonAddress {
	return &TonAddress{
		addr: addr,
	}
}

func MustParse(addr string) *TonAddress {
	add, _ := Parse(addr)
	return add
}

func Parse(addr string) (*TonAddress, error) {
	var add *address.Address
	var err error
	if IsRaw(addr) {
		add, err = address.ParseRawAddr(addr)
	} else {
		add, err = address.ParseAddr(addr)
	}
	if err != nil {
		return nil, err
	}

	return &TonAddress{
		addr: add,
	}, nil
}

type TonAddress struct {
	addr *address.Address
	//urlSafe bool
}

func (receiver *TonAddress) Bounce(bouncable bool) *TonAddress {
	receiver.addr.SetBounce(bouncable)
	return receiver
}

func (receiver *TonAddress) IsBounce() bool {
	return receiver.addr.IsBounceable()
}

func (receiver *TonAddress) TestnetOnly(testnetOnly bool) *TonAddress {
	receiver.addr.SetTestnetOnly(testnetOnly)
	return receiver
}

//	func (receiver *TonAddress) ToFriendly() string {
//		return receiver.String()
//		//return hex.EncodeToString([]byte(receiver.String()))
//	}
func (receiver *TonAddress) ToRaw() string {
	out := hex.EncodeToString(receiver.addr.Data())
	return out
	//return hex.EncodeToString([]byte(receiver.String()))
}

//	func (receiver *TonAddress) UrlSafe(urlSafe bool) *TonAddress {
//		receiver.urlSafe = urlSafe
//		return receiver
//		//return hex.EncodeToString([]byte(receiver.String()))
//	}
func (receiver *TonAddress) String() string {
	out := receiver.addr.String()

	//.replace(/\+/g, '-').replace(/\//g, '_');
	//if receiver.urlSafe {
	//	out = strings.ReplaceAll(out, "+", "-")
	//	out = strings.ReplaceAll(out, "/", "_")
	//}
	return out
	//return hex.EncodeToString([]byte(receiver.String()))
}
