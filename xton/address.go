package xton

import (
	"encoding/hex"
	"fmt"
	"github.com/xssnick/tonutils-go/address"
	"regexp"
	"strconv"
	"strings"
)

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

func IsRaw(source string) bool {
	parts := strings.Split(source, ":")
	if len(parts) != 2 {
		return false
	}

	wc, hash := parts[0], parts[1]

	// Check if wc can be parsed as an integer
	if _, err := strconv.Atoi(wc); err != nil {
		return false
	}

	// Check if hash is a valid hex string of length 64
	if len(hash) != 64 {
		return false
	}

	validHexRegex := regexp.MustCompile(`^[a-f0-9]+$`)
	if !validHexRegex.MatchString(strings.ToLower(hash)) {
		return false
	}

	return true
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

func AddressRawParse(addStr string) string {
	add := address.MustParseRawAddr(addStr)
	//add.SetBounce(false)
	return add.String()
}

func AddressToRaw(addStr string) string {
	add := address.MustParseAddr(addStr)
	//add.SetBounce(false)
	return add.String()
}

func AddressRawParseToBounceable(addStr string) string {
	add := address.MustParseRawAddr(addStr)

	add.SetBounce(true)
	//var dst []byte
	//var addr []byte
	//add.StringToBytes(dst, addr)
	//fmt.Println("AddressRawParseToBounceable", string(dst), string(addr))
	return hex.EncodeToString([]byte(add.String()))
}
func AddressParseToBounceable(addStr string) string {
	add := address.MustParseAddr(addStr)

	if add.IsBounceable() {
		fmt.Println(addStr, addStr)
		return addStr
	}

	add.SetBounce(true)
	fmt.Println(addStr, add.String())
	return add.String()
}
