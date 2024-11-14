package tonclient

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/xssnick/tonutils-go/address"
	"github.com/xssnick/tonutils-go/tlb"
	"github.com/xssnick/tonutils-go/ton"
	"log"
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

func WatchTransactions(ctx context.Context, api *ton.APIClient, acctAddr *address.Address, lastLT uint64, txMapHandler func(tx *tlb.Transaction)) {

	c := make(chan *tlb.Transaction)

	go SubscribeOnTransactions(ctx, api, acctAddr, lastLT, c)
	//go api.SubscribeOnTransactions(ctx, acctAddr, lastLT, c)
	//go watch(api, ctx, acctAddr, lastLT, c)

	for {
		select {
		case <-ctx.Done():
			return
		case tx := <-c:
			//fmt.Println("tx:", tx)
			txMapHandler(tx)
		}
	}

}

func watch(api ton.APIClientWrapped, ctx context.Context, addr *address.Address, redisLt uint64, channel chan<- *tlb.Transaction) {
	defer func() {
		close(channel)
	}()

	//logicCtx.Redis().Del(context.Background(), redisKey)
	// 定义睡眠时间
	for {
		// 每次循环，获取一次账号数据
		b, err := api.CurrentMasterchainInfo(ctx)
		if err != nil {
			log.Fatalln("get block err:", err.Error())
			return
		}

		res, err := api.GetAccount(ctx, b, addr)
		if err != nil {
			log.Fatalln("get account err:", err.Error())
			return
		}

		// 获取当前获取到的
		lastLt := res.LastTxLT
		lastHash := res.LastTxHash

		//redisLt, isGo := x.CheckTL(logicCtx, lastLt)
		//if !isGo {
		//	fmt.Println("你还没我大，不配走")
		//	time.Sleep(time.Duration(sleepSecond) * time.Second)
		//	continue
		//}

		//var redisLt uint64
		//redisLt, _ = logicCtx.Redis().Get(context.Background(), redisKey).Uint64()
		//fmt.Println("redis 获取的值：", redisLt)
		//if redisLt != 0 && lastLt <= redisLt {
		//	//fmt.Println("你还没我大，不配走")
		//	time.Sleep(time.Duration(sleepSecond) * time.Second)
		//	continue
		//}

		for {

			// load transactions in batches with size 15
			list, err := api.ListTransactions(ctx, addr, 10, lastLt, lastHash)
			if err != nil {
				log.Printf("send err: %s", err.Error())
				break
			}

			// 将数组倒叙
			sort.Slice(list, func(i, j int) bool {
				return list[i].LT > list[j].LT
			})

			lastLt = list[len(list)-1].PrevTxLT
			lastHash = list[len(list)-1].PrevTxHash

			for k, t := range list {

				if t.LT <= redisLt {
					fmt.Println("到这就结束了")
					goto endFor
				}

				fmt.Println("list[", k, "].LT", t.LT, t.String())
				channel <- t
				// 判断是否支付成功
				if t.OrigStatus != tlb.AccountStatusActive || t.EndStatus != tlb.AccountStatusActive {
					continue
				}

				//fmt.Println("t.EndStatus", t.EndStatus)
				//if t.EndStatus != tlb.AccountStatusActive {
				//	continue
				//}
				//
				//// 先获取报文
				//comment := t.IO.In.AsInternal().Comment()

			}
		}

	endFor:
		fmt.Println("记录了最大值")
		// 记录到数据库

	}
}

var ErrorNoTx = errors.New("no transactions were found")

func SubscribeOnTransactions(workerCtx context.Context, c *ton.APIClient, addr *address.Address, lastProcessedLT uint64, channel chan<- *tlb.Transaction) {
	defer func() {
		close(channel)
	}()

	wait := 0 * time.Second
	for {
		select {
		case <-workerCtx.Done():
			return
		case <-time.After(wait):
		}
		wait = 3 * time.Second

		ctx, cancel := context.WithTimeout(workerCtx, 10*time.Second)
		master, err := c.CurrentMasterchainInfo(ctx)
		cancel()
		if err != nil {
			fmt.Println("GetMasterchainInfo:", err)
			continue
		}

		ctx, cancel = context.WithTimeout(workerCtx, 10*time.Second)
		acc, err := c.GetAccount(ctx, master, addr)
		cancel()
		if err != nil {
			fmt.Println("GetAccount:", err)
			continue
		}
		if !acc.IsActive || acc.LastTxLT == 0 {
			// no transactions
			fmt.Println("no transactions:", acc.LastTxLT, acc.IsActive)
			continue
		}

		if lastProcessedLT == acc.LastTxLT {
			// already processed all
			//fmt.Println("already processed all", lastProcessedLT)
			continue
		}

		var transactions []*tlb.Transaction
		lastHash, lastLT := acc.LastTxHash, acc.LastTxLT

		waitList := 0 * time.Second
	list:
		for {
			select {
			case <-workerCtx.Done():
				return
			case <-time.After(waitList):
			}

			ctx, cancel = context.WithTimeout(workerCtx, 10*time.Second)
			//fmt.Println("lastLT, lastHash", lastLT, hex.EncodeToString(lastHash))
			res, err := c.ListTransactions(ctx, addr, 10, lastLT, lastHash)
			//fmt.Println("err,res", err, len(res), errors.Is(err, ton.ErrNoTransactionsWereFound))
			cancel()
			if err != nil {
				if errors.Is(err, ton.ErrNoTransactionsWereFound) {
					if len(transactions) > 0 {
						break
					}
				}
				//fmt.Println("ListTransactions:", err)
				if lsErr, ok := err.(ton.LSError); ok {
					fmt.Println("ListTransactions:", lsErr.Code, lsErr.Error())
					if lsErr.Code == -400 {
						break
					}
					// lt not in db error

					//lastLT = acc.LastTxLT
					//lastHash = acc.LastTxHash
					//lastProcessedLT = acc.LastTxLT
					break
				}

				waitList = 3 * time.Second
				continue
			}

			if len(res) == 0 {
				break
			}

			sort.Slice(res, func(i, j int) bool {
				return res[i].LT > res[j].LT
			})

			// reverse slice
			//for i, j := 0, len(res)-1; i < j; i, j = i+1, j-1 {
			//	res[i], res[j] = res[j], res[i]
			//}

			for i, tx := range res {
				//fmt.Println("range res:", i, tx.LT, tx.LT <= lastProcessedLT)
				if tx.LT <= lastProcessedLT {
					transactions = append(transactions, res[:i]...)
					break list
				}
			}

			lastLT, lastHash = res[len(res)-1].PrevTxLT, res[len(res)-1].PrevTxHash
			//fmt.Println("next,lastLT, lastHash", lastLT, hex.EncodeToString(lastHash))
			transactions = append(transactions, res...)
			waitList = 0 * time.Second
		}

		//fmt.Println("transactions:", len(transactions))
		if len(transactions) > 0 {
			lastProcessedLT = transactions[0].LT // mark last transaction as known to not trigger twice

			// reverse slice to send in correct time order (from old to new)
			//for i, j := 0, len(transactions)-1; i < j; i, j = i+1, j-1 {
			//	transactions[i], transactions[j] = transactions[j], transactions[i]
			//}
			//fmt.Println("range transactions:", len(transactions))
			for _, tx := range transactions {
				channel <- tx
			}

			wait = 0 * time.Second
		}
	}
}
