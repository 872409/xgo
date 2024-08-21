package tonclient

import (
	"bytes"
	"context"
	"fmt"
	"github.com/xssnick/tonutils-go/address"
	"github.com/xssnick/tonutils-go/tlb"
	"github.com/xssnick/tonutils-go/ton"
	"sort"
	"time"
)

func getList(api *ton.APIClient, acctAddr *address.Address, limit uint32, txLT uint64, txHash []byte, txMap func(txItem *tlb.Transaction) bool) (bool, *tlb.Transaction, error) {
	list, err := api.ListTransactions(context.Background(), acctAddr, limit, txLT, txHash)
	if err != nil {
		return false, nil, err
	}

	listLastIdx := len(list) - 1

	if listLastIdx < 0 {
		return false, nil, fmt.Errorf("listLastIdx is negative")
	}

	if listLastIdx == 0 && bytes.Equal(list[listLastIdx].Hash, txHash) {
		return false, nil, err
	}

	sort.Slice(list, func(i, j int) bool {
		return list[i].LT > list[j].LT
	})

	for _, tx := range list {
		if bytes.Equal(tx.Hash, txHash) {
			continue
		}

		if txMap(tx) {
			return true, tx, nil
		}
	}
	return false, list[listLastIdx], nil
}

func MapTransactions(api *ton.APIClient, acctAddr *address.Address, interval time.Duration, txMapHandler func(tx *tlb.Transaction) bool) (maxTx *tlb.Transaction, err error) {

	acct, err := GetAccount(context.Background(), api, acctAddr)
	if err != nil {
		return nil, err
	}

	lastTx := &tlb.Transaction{
		LT:   acct.LastTxLT,
		Hash: acct.LastTxHash,
	}
	ok := false
	var _lastTx *tlb.Transaction
	for !ok {
		ok, _lastTx, err = getList(api, acctAddr, 1000, lastTx.LT, lastTx.Hash, func(tx *tlb.Transaction) bool {

			if maxTx == nil || tx.LT > maxTx.LT {
				maxTx = tx
			}

			return txMapHandler(tx)

		})

		if _lastTx != nil {
			lastTx = _lastTx
		}

		if err != nil {
			fmt.Println("getList:", err)
		}

		//fmt.Println(ok, lastTx, err)
		time.Sleep(interval)
		//time.Sleep(3 * time.Second)
	}
	fmt.Println("maxTx", maxTx)
	return maxTx, nil
}
