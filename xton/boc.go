package xton

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/xssnick/tonutils-go/tvm/cell"
)

func DecodeMsgHashFromBoc(boc string) (string, string, error) {

	a, e := base64.StdEncoding.DecodeString(boc)
	fmt.Println(a, e)
	if e != nil {
		return "", "", e
	}

	c, e := cell.FromBOC(a)
	fmt.Println(c, e)

	if e != nil {
		return "", "", e
	}

	hashBytes := c.Hash()
	hashHex := hex.EncodeToString(hashBytes)
	fmt.Println(hashHex)

	return hashHex, base64.StdEncoding.EncodeToString(hashBytes), nil

	//curl -X 'GET' \
	//  'https://toncenter.com/api/v3/transactionsByMessage?direction=in&msg_hash=AmFKajlhOZtrbRZgXu_PXNCGUrppWvRLRV6Cp5cjnpk%3D&limit=128&offset=0' \
	//  -H 'accept: application/json'
}
