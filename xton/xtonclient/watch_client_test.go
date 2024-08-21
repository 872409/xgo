package tonclient

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/xssnick/tonutils-go/address"
	"github.com/xssnick/tonutils-go/tlb"
	"math/big"
	"testing"
	"time"
)

func TestNewAccountWatcher(t *testing.T) {
	got, err := NewAccountWatcher(context.Background(), 46330160000003, "UQAZ9CBAryyCtrwvm4kkNYx69hf8mL4M7FkBoIj0iui94S5s", 1*time.Second, order1)
	if err != nil {
		fmt.Println(err)
		return
	}
	got.Start()
	time.Sleep(3 * time.Second)
	got.Stop()
	time.Sleep(3 * time.Second)
}

func order1(tx *tlb.Transaction, fromWallet string, toWallet string, amount *big.Int, comment string) {
	fmt.Println("order1: ", fromWallet, toWallet, amount, toWallet, comment)
}

func TestName2(t *testing.T) {
	_, api, err := NewClient(context.Background())
	fmt.Println(err)
	add := address.MustParseAddr("UQAZ9CBAryyCtrwvm4kkNYx69hf8mL4M7FkBoIj0iui94S5s")

	minLT := uint64(46122579000001)
	aborted := 0
	total := 0
	txMap := make(map[uint64]int)
	maxTx, err := MapTransactions(api, add, 3*time.Second, func(tx *tlb.Transaction) bool {

		if _, found := txMap[tx.LT]; found {
			txMap[tx.LT] = txMap[tx.LT] + 1
		} else {
			txMap[tx.LT] = 1
			total++
			fmt.Println("total: ", total)
		}

		//fmt.Println("tx: ", tx.LT, base64.StdEncoding.EncodeToString(tx.Hash), hex.EncodeToString(tx.Hash))

		d := tx.Description.Description.(tlb.TransactionDescriptionOrdinary)
		computePh, ok := d.ComputePhase.Phase.(tlb.ComputePhaseVM)
		if ok {
			if computePh.Details.ExitCode < 0 {
				aborted++
				fmt.Println("d.Aborted", tx.LT, base64.StdEncoding.EncodeToString(tx.Hash), hex.EncodeToString(tx.Hash))
			}
		}

		//t := time.Unix(int64(tx.Now), 0)
		//fmt.Println("Now: ", tx.LT, t)

		//if tx.IO.In.MsgType == tlb.MsgTypeInternal {
		//
		//	//fmt.Println(tx.LT, tx.IO.In.AsInternal().Amount, tx.IO.In.AsInternal().Comment())
		//
		//}

		if tx.LT == minLT {
			//d := tx.Description.Description.(tlb.TransactionDescriptionOrdinary)
			//fmt.Println("minTxLT: ", tx, d)
			return true
		}
		return false
	})

	fmt.Println("total", total)
	fmt.Println("aborted", aborted)
	fmt.Println("maxTx", maxTx, err)
}
func TestName(t *testing.T) {

	_, api, err := NewClient(context.Background())
	fmt.Println(err)
	add := address.MustParseAddr("UQAZ9CBAryyCtrwvm4kkNYx69hf8mL4M7FkBoIj0iui94S5s")

	acct, err := GetAccount(context.Background(), api, add)

	getList := func(minTxLT uint64, limit uint32, txLT uint64, txHash []byte) (bool, uint64, []byte, error) {
		list, err := api.ListTransactions(context.Background(), add, limit, txLT, txHash)
		fmt.Println(err)
		if err != nil {
			return false, 0, nil, err
		}
		listLastIdx := len(list) - 1

		if listLastIdx < 0 {
			return false, 0, nil, fmt.Errorf("listLastIdx is negative")
		}

		for _, tx := range list {
			if tx.IO.In.MsgType == tlb.MsgTypeInternal {
				d := tx.Description.Description.(tlb.TransactionDescriptionOrdinary)
				fmt.Println("d.Aborted", tx.LT, d.Aborted)
				if d.Aborted {

				}
				//fmt.Println(tx.LT, tx.IO.In.AsInternal().Amount, tx.IO.In.AsInternal().Comment(), d)
			}

			if tx.LT == minTxLT {
				d := tx.Description.Description.(tlb.TransactionDescriptionOrdinary)
				fmt.Println("minTxLT: ", tx, d)
				return true, tx.LT, tx.Hash, nil
			}
		}
		return false, list[0].LT, list[0].Hash, nil
	}

	txLT := acct.LastTxLT
	txHash := acct.LastTxHash
	minLT := uint64(45758534000001)
	ok := false
	for !ok {
		ok, txLT, txHash, err = getList(minLT, 100, txLT, txHash)
		fmt.Println(ok, txLT, txHash, err)
		time.Sleep(3 * time.Second)
	}

}
