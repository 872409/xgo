package tonclient

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/xssnick/tonutils-go/address"
	"github.com/xssnick/tonutils-go/tlb"
	"github.com/xssnick/tonutils-go/ton/jetton"
	"log"
	"math/big"
	"testing"
	"time"
)

type name struct {
	a string
}

func Test1111(t *testing.T) {
	var x = &name{
		//a: "a",
	}
	var c = &name{
		//a: "a",
	}
	fmt.Printf("%p,%p", x, c)
}
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

func TestJetton(t *testing.T) {
	ctx := context.Background()
	_, api, err := NewClient(ctx)
	master, err := api.CurrentMasterchainInfo(context.Background()) // we fetch block just to trigger chain proof check
	if err != nil {
		log.Fatalln("get masterchain info err: ", err.Error())
		return
	}

	fmt.Println(err)
	watchAddress := address.MustParseAddr("UQAJWmeBW1--sX5BAoCrw-2QlVJh_IiDSCSdI5snJwDht6bP")
	//watchAddress := address.MustParseAddr("EQB88BBo29q9CeUR4I8JQuOId3zqMJiGAHyrhTnmaGdLOX3d")
	//watchAddress := address.MustParseAddr("UQAZ9CBAryyCtrwvm4kkNYx69hf8mL4M7FkBoIj0iui94S5s")
	// address on which we are accepting payments
	treasuryAddress := watchAddress

	usdt := jetton.NewJettonMasterClient(api, address.MustParseAddr("EQCxE6mUtQJKFnGfaROTKOt1lZbDiiX1kCixRv7Nw2Id_sDs"))
	// get our jetton wallet address
	treasuryJettonWallet, err := usdt.GetJettonWalletAtBlock(context.Background(), treasuryAddress, master)

	//txMap := make(map[uint64]int)
	WatchTransactions(ctx, api, watchAddress, 46429334000003, func(tx *tlb.Transaction) {
		//fmt.Println("tx:", tx.LT, tx.IO.In.MsgType)
		//d := tx.Description.Description.(tlb.TransactionDescriptionOrdinary)
		//computePh, ok := d.ComputePhase.Phase.(tlb.ComputePhaseVM)
		//if !ok || computePh.Details.ExitCode < 0 {
		//	fmt.Println("d.Aborted", tx.LT, base64.StdEncoding.EncodeToString(tx.Hash), hex.EncodeToString(tx.Hash))
		//	return
		//}

		//if (tx.IO.In.MsgType !== 'internal' || tx.description.type !== 'generic' || tx.description.computePhase?.type !== 'vm') {
		//  return;
		//}
		//if tx.IO.In.MsgType == tlb.MsgTypeExternalIn {
		//	fmt.Println("tx:", tx.LT, hex.EncodeToString(tx.IO.In.AsExternalIn().Payload().Hash()))
		//}
		//if tx.IO.In.MsgType == tlb.MsgTypeInternal {
		//	inMsg := tx.IO.In.AsInternal()
		//	comment := tx.IO.In.AsInternal().Comment()
		//	fmt.Println("tx:", inMsg.SrcAddr, inMsg.DstAddr, comment)
		//}

		if tx.IO.In != nil && tx.IO.In.MsgType == tlb.MsgTypeInternal {
			ti := tx.IO.In.AsInternal()
			src := ti.SrcAddr

			fmt.Println("src", src)

			// verify that sender is our jetton wallet
			if ti.SrcAddr.Equals(treasuryJettonWallet.Address()) {
				var transfer jetton.TransferNotification
				if err = tlb.LoadFromCell(&transfer, ti.Body.BeginParse()); err == nil {

					payloadSlice := transfer.ForwardPayload.BeginParse()
					payloadOpcode, er := payloadSlice.LoadUInt(32)
					fmt.Println("payloadOpcode", payloadOpcode, er)
					if payloadOpcode != 0 {
						return
					}
					payloadComment := payloadSlice.MustLoadStringSnake()
					fmt.Println("payloadComment", payloadComment)
					// convert decimals to 6 for USDT (it can be fetched from jetton details too), default is 9
					amt := tlb.MustFromNano(transfer.Amount.Nano(), 6)

					// reassign sender to real jetton sender instead of its jetton wallet contract
					src = transfer.Sender
					log.Println("received", amt.String(), "USDT from", src.String())
				}
			}

			// show received ton amount
			//log.Println("received", ti.Amount.String(), "TON from", src.String())
		}
	})

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
