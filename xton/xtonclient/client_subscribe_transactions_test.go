package tonclient

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/xssnick/tonutils-go/address"
	"github.com/xssnick/tonutils-go/tlb"
	"log"
	"testing"
)

func TestLoopTransactionsV2(t *testing.T) {
	ctx := context.Background()
	_, api, err := NewClient(ctx)
	if err != nil {
		log.Fatalln("get masterchain info err: ", err.Error())
		return
	}

	fmt.Println(err)
	watchAddress := address.MustParseAddr("UQAJWmeBW1--sX5BAoCrw-2QlVJh_IiDSCSdI5snJwDht6bP")
	//lastLt := 50844256000001
	//lastLt := 50689999000004
	lastLt := 0
	//lastLt := 46429334000003
	WatchTransactionsV2(ctx, api, watchAddress, uint64(lastLt), func(tx *tlb.Transaction) {
		fmt.Println("tx:", tx.LT, tx.IO.In.MsgType)
	}, func(eventType EventType, lastProcessedLT uint64, account *tlb.Account) {
		fmt.Println(eventType, lastProcessedLT)
	})
}

func TestNewJettonClient2(t *testing.T) {
	r, e := base64.StdEncoding.DecodeString("te6cckEBAQEACgAAEAAAAAB0ZXN0cbBecQ==")
	fmt.Println(string(r), e)
}

func TestNewJettonClient(t *testing.T) {
	ctx := context.Background()
	_, api, err := NewClient(ctx)
	if err != nil {
		log.Fatalln("get masterchain info err: ", err.Error())
		return
	}

	accountAddress := address.MustParseAddr("UQAJWmeBW1--sX5BAoCrw-2QlVJh_IiDSCSdI5snJwDht6bP")
	jettonMaster := address.MustParseAddr("EQCxE6mUtQJKFnGfaROTKOt1lZbDiiX1kCixRv7Nw2Id_sDs")
	jettonClient, err := NewJettonClient(ctx, api, accountAddress, jettonMaster)

	if err != nil {
		fmt.Println("new jetton client err: ", err.Error())
		return
	}

	lastLt := 0
	WatchTransactionsV2(ctx, api, accountAddress, uint64(lastLt), func(tx *tlb.Transaction) {
		//fmt.Println("tx:", tx.LT, tx.IO.In.MsgType)
		transfer, err := jettonClient.HandleInTx(tx)
		if err != nil {
			return
		}
		amt := tlb.MustFromNano(transfer.TransferNotification.Amount.Nano(), 6)
		fmt.Println("Jetton transfer", transfer.Comment, transfer.TransferNotification.Amount, amt, err)

	}, func(eventType EventType, lastProcessedLT uint64, account *tlb.Account) {
		fmt.Println(eventType, lastProcessedLT)
	})
}

func TestNewJettonClient_gala(t *testing.T) {
	ctx := context.Background()
	_, api, err := NewClient(ctx)
	if err != nil {
		log.Fatalln("get masterchain info err: ", err.Error())
		return
	}

	accountAddress := address.MustParseAddr("UQDUaIr_R6tcJEZSSOA2aBYw6ZpDP78KR22jb-jRndIWepWy")
	jettonMaster := address.MustParseAddr("EQBadmOayy7_bD18skopfOZw2kmTgDdBhXPVsuTQq1lalaBV")
	jettonClient, err := NewJettonClient(ctx, api, accountAddress, jettonMaster)

	if err != nil {
		fmt.Println("new jetton client err: ", err.Error())
		return
	}

	lastLt := 0
	WatchTransactionsV2(ctx, api, accountAddress, uint64(lastLt), func(tx *tlb.Transaction) {
		//fmt.Println("tx:", tx.LT, tx.IO.In.MsgType)
		transfer, err := jettonClient.HandleInTx(tx)
		if err != nil {
			fmt.Println("handle tx err: ", err.Error())
			return
		}

		amt := tlb.MustFromNano(transfer.TransferNotification.Amount.Nano(), 8)
		fmt.Println("Jetton transfer", transfer.Comment, transfer.TransferNotification.Amount, amt, err)

	}, func(eventType EventType, lastProcessedLT uint64, account *tlb.Account) {
		fmt.Println(eventType, lastProcessedLT)
	})
}
