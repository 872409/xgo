package tonclient

import (
	"context"
	"errors"
	"fmt"
	"github.com/xssnick/tonutils-go/address"
	"github.com/xssnick/tonutils-go/tlb"
	"github.com/xssnick/tonutils-go/ton"
	"sort"
	"time"
)

type EventType = string

const EventType_FetchEnd EventType = "end"
const EventType_ProcessStop EventType = "stop"
const EventType_Watch EventType = "watch"

type EventHandler = func(eventType EventType, lastProcessedLT uint64, account *tlb.Account)

func WatchTransactionsV2(ctx context.Context, api *ton.APIClient, acctAddr *address.Address, lastProcessedLT uint64, txMapHandler func(tx *tlb.Transaction), eventHandler EventHandler) {

	c := make(chan *tlb.Transaction)

	//go SubscribeOnTransactionsV2(ctx, api, acctAddr, lastLT, c)
	go SubscribeOnTransactionsV2(ctx, api, acctAddr, lastProcessedLT, c, eventHandler)

	for {
		select {
		case <-ctx.Done():
			return
		case tx := <-c:
			txMapHandler(tx)
		}
	}

}

func getListV2(workerCtx context.Context, api *ton.APIClient, acctAddr *address.Address, limit uint32, txLT uint64, txHash []byte) (result []*tlb.Transaction, err error) {

	ctx, cancel := context.WithTimeout(workerCtx, 10*time.Second)
	list, err := api.ListTransactions(ctx, acctAddr, limit, txLT, txHash)
	cancel()

	if errors.Is(err, ton.ErrNoTransactionsWereFound) {
		return []*tlb.Transaction{}, nil
	}

	if err != nil {
		var lsErr ton.LSError
		if errors.As(err, &lsErr) {
			if lsErr.Code == -400 {
				return []*tlb.Transaction{}, nil
			}
		}
		return nil, err
	}

	//1->0
	sort.Slice(list, func(i, j int) bool {
		return list[i].LT > list[j].LT
	})

	return list, nil
}

func getAccountV2(workerCtx context.Context, api ton.APIClientWrapped, addr *address.Address) (*tlb.Account, error) {
	ctx, cancel := context.WithTimeout(workerCtx, 10*time.Second)

	defer cancel()

	// 每次循环，获取一次账号数据
	b, err := api.CurrentMasterchainInfo(ctx)
	if err != nil {
		//fmt.Println("get block err:", err.Error())
		return nil, err
	}
	account, err := api.GetAccount(ctx, b, addr)
	if err != nil {
		//fmt.Println("get account err:", err.Error())
		return nil, err
	}
	return account, err
}

func SubscribeOnTransactionsV2(workerCtx context.Context, api *ton.APIClient, acctAddr *address.Address, lastProcessedLT uint64, channel chan<- *tlb.Transaction, eventHandler EventHandler) {

	lastTx := &tlb.Transaction{}
	maxTx := &tlb.Transaction{}
	var watchNew = false
	wait := time.Duration(0)
	var acct *tlb.Account
	for {
		select {
		case <-workerCtx.Done():
			return
		case <-time.After(wait):
		}

		if lastTx.LT == 0 || watchNew {
			_acct, err := getAccountV2(workerCtx, api, acctAddr)

			if err != nil {
				fmt.Println("getAccountV2 err", err)
				continue
			}

			acct = _acct

			lastTx = &tlb.Transaction{
				LT:   acct.LastTxLT,
				Hash: acct.LastTxHash,
			}

			//fmt.Println("GetAccount", lastTx.LT, lastProcessedLT)

			if lastTx.LT == lastProcessedLT {
				//fmt.Println("Wait Watch New")
				wait = time.Second * 3
				watchNew = true
				eventHandler(EventType_Watch, lastProcessedLT, acct)
				continue
			}
		}

		list, err := getListV2(workerCtx, api, acctAddr, 10, lastTx.LT, lastTx.Hash)
		if err != nil {
			//fmt.Println("getList err try:", err)
			wait = time.Second * 1
			continue
		}

		hasMore := list != nil && len(list) > 0

		if hasMore {

			fmt.Println("list", len(list))
			for _, tx := range list {
				//fmt.Println("PrevTxLT", tx.LT, tx.PrevTxLT)
				if tx.LT >= maxTx.LT {
					maxTx = tx
				}

				if tx.LT > lastProcessedLT {
					channel <- tx
				} else {
					watchNew = true
					lastProcessedLT = maxTx.LT
					fmt.Println("Stop! Watch New", lastTx.LT, tx.LT, lastProcessedLT)
					eventHandler(EventType_ProcessStop, lastProcessedLT, acct)
					break
				}
			}

			if watchNew {
				continue
			}

			lastTx = &tlb.Transaction{
				LT:   list[len(list)-1].PrevTxLT,
				Hash: list[len(list)-1].PrevTxHash,
			}
			if lastTx.LT == 0 {
				watchNew = true
				lastProcessedLT = maxTx.LT
				//fmt.Println("Stop! Watch New", lastTx.LT, lastProcessedLT)
				eventHandler(EventType_ProcessStop, lastProcessedLT, acct)
				continue
			}
			//fmt.Println("Transaction", lastTx.LT, lastProcessedLT)
			watchNew = false
		}

		if !hasMore {
			watchNew = true
			lastProcessedLT = maxTx.LT
			eventHandler(EventType_FetchEnd, lastProcessedLT, acct)
			//fmt.Println("End! Watch New", lastProcessedLT)
		}

	}
}
