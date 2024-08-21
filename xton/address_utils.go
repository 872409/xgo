package xton

import (
	"encoding/hex"
	"fmt"
	"github.com/xssnick/tonutils-go/address"
	"regexp"
	"strconv"
	"strings"
)

func AddressToRaw(addr *address.Address) string {
	out := hex.EncodeToString(addr.Data())
	return out
	//return hex.EncodeToString([]byte(receiver.String()))
}

func AddressRawParse(addStr string) string {
	add := address.MustParseRawAddr(addStr)
	//add.SetBounce(false)
	return add.String()
}

func AddressToRawStr(addStr string) string {
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
