package tonclient

import (
	"context"
	"fmt"
	"github.com/xssnick/tonutils-go/address"
	"github.com/xssnick/tonutils-go/tlb"
	"github.com/xssnick/tonutils-go/ton"
	"log"
	"main/comm/xutils"
	"math/big"
	"time"
)

//	type WatchClient struct {
//		client             ton.APIClientWrapped
//		pool               *liteclient.ConnectionPool
//		ctx                context.Context
//		interval           time.Duration
//		wallets            []string
//		accountWatcherList []*accountWatcher
//	}
//
//	func NewWatchClient(ctx context.Context, interval time.Duration, wallet ...string) *WatchClient {
//		pool := liteclient.NewConnectionPool()
//		api := ton.NewAPIClient(pool, ton.ProofCheckPolicyFast)
//
//		return &WatchClient{
//			wallets:  wallet,
//			pool:     pool,
//			interval: interval,
//			ctx:      pool.StickyContext(ctx),
//			client:   api.WithRetry(),
//		}
//	}
type TransactionHandler func(tx *tlb.Transaction, fromWallet string, toWallet string, amount *big.Int, comment string)

//
//func (r *WatchClient) Stop() {
//
//}
//
//func (r *WatchClient) Start(h TransactionHandler) {
//
//}
//
//func (r *WatchClient) transactionHandler(h TransactionHandler) {
//
//}

type accountWatcher struct {
	tickerChan  chan bool
	client      ton.APIClientWrapped
	ctx         context.Context
	interval    time.Duration
	addr        *address.Address
	addrAccount *tlb.Account
	handler     TransactionHandler
	lastLt      uint64
	lastHash    []byte
}

func NewAccountWatcher(ctx context.Context, lastLt uint64, addr string, interval time.Duration, handler TransactionHandler) (*accountWatcher, error) {
	pool, api, err := NewClient(ctx)
	if err != nil {
		log.Fatalln("connection err: ", err.Error())
		return nil, nil
	}

	return &accountWatcher{
		client:   api.WithRetry(),
		lastLt:   lastLt,
		ctx:      pool.StickyContext(ctx),
		addr:     address.MustParseAddr(addr),
		interval: interval,
		handler:  handler,
	}, nil
}

func (r *accountWatcher) initAccount() error {
	// 每次循环，获取一次账号数据
	b, err := r.client.CurrentMasterchainInfo(r.ctx)
	if err != nil {
		log.Fatalln("get block err:", err.Error())
		return err
	}
	account, err := r.client.GetAccount(r.ctx, b, r.addr)
	if err != nil {
		log.Fatalln("get account err:", err.Error())
		return err
	}
	r.addrAccount = account
	//
	//// 获取当前获取到的
	//lastLt := r.addrAccount.LastTxLT
	//lastHash := r.addrAccount.LastTxHash
	//
	////var redisLt uint64
	////redisLt, _ = logicCtx.Redis().Get(logicCtx.Ctx(), "TON_LT").Uint64()
	//
	//if lastLt < r.lastLt {
	//	fmt.Println("你还没我大，不配走")
	//	time.Sleep(r.interval * time.Second)
	//	r.lastLt = lastLt
	//	r.lastHash = lastHash
	//	return nil
	//}
	return nil
}
func (r *accountWatcher) getAccountLastLT() (uint64, []byte) {

	if r.lastHash != nil {
		return r.lastLt, r.lastHash
	}

	lastLt := r.addrAccount.LastTxLT
	lastHash := r.addrAccount.LastTxHash

	//var redisLt uint64
	//redisLt, _ = logicCtx.Redis().Get(logicCtx.Ctx(), "TON_LT").Uint64()

	if lastLt < r.lastLt {
		fmt.Println("你还没我大，不配走")
		time.Sleep(r.interval * time.Second)
		r.lastLt = lastLt
		r.lastHash = lastHash
		return 0, nil
	}

	return lastLt, lastHash
}

func (r *accountWatcher) Stop() {
	if r.tickerChan != nil {
		r.tickerChan <- true
	}
}

func (r *accountWatcher) Start() {
	r.initAccount()
	r.tickerChan = xutils.NewTicker(r.interval, func() {
		fmt.Println("tick")
		lastLt, lastHash := r.getAccountLastLT()
		if lastLt > 0 {
			r.goDo(lastLt, lastHash)
		}
	})

	//
	//go func(ticker *time.Ticker) {
	//	//defer r.ticker.Stop()
	//
	//	for range ticker.C {
	//		fmt.Println("tick")
	//		//lastLt, lastHash := r.getAccountLastLT()
	//		//if lastLt > 0 {
	//		//	//r.goDo(lastLt, lastHash)
	//		//}
	//	}
	//	fmt.Println("ticker out")
	//
	//	//for {
	//	//	select {
	//	//	case <-r.ticker.C:
	//	//		fmt.Println("ticker")
	//	//		lastLt, lastHash := r.getAccountLastLT()
	//	//		if lastLt > 0 {
	//	//			r.goDo(lastLt, lastHash)
	//	//		}
	//	//	}
	//	//}
	//
	//}(r.ticker)

}

//	func (r *accountWatcher) getTxList(minTxLT uint64, limit uint32, txLT uint64, txHash []byte) (bool, uint64, []byte, error) {
//		acct, err := GetAccount(context.Background(), api, add)
//
//		getList := func(minTxLT uint64, limit uint32, txLT uint64, txHash []byte) (bool, uint64, []byte, error) {
//			list, err := api.ListTransactions(context.Background(), add, limit, txLT, txHash)
//			fmt.Println(err)
//			if err != nil {
//				return false, 0, nil, err
//			}
//			listLastIdx := len(list) - 1
//
//			if listLastIdx < 0 {
//				return false, 0, nil, fmt.Errorf("listLastIdx is negative")
//			}
//
//			for _, tx := range list {
//				if tx.IO.In.MsgType == tlb.MsgTypeInternal {
//					d := tx.Description.Description.(tlb.TransactionDescriptionOrdinary)
//					fmt.Println("d.Aborted", tx.LT, d.Aborted)
//					if d.Aborted {
//
//					}
//					//fmt.Println(tx.LT, tx.IO.In.AsInternal().Amount, tx.IO.In.AsInternal().Comment(), d)
//				}
//
//				if tx.LT == minTxLT {
//					d := tx.Description.Description.(tlb.TransactionDescriptionOrdinary)
//					fmt.Println("minTxLT: ", tx, d)
//					return true, tx.LT, tx.Hash, nil
//				}
//			}
//			return false, list[0].LT, list[0].Hash, nil
//		}
//
//		txLT := acct.LastTxLT
//		txHash := acct.LastTxHash
//		minLT := uint64(45758534000001)
//		ok := false
//		for !ok {
//			ok, txLT, txHash, err = getList(minLT, 100, txLT, txHash)
//			fmt.Println(ok, txLT, txHash, err)
//			time.Sleep(3 * time.Second)
//		}
//	}
func (r *accountWatcher) goDo(lastLt uint64, lastHash []byte) {
	list, err := r.client.ListTransactions(r.ctx, r.addr, 10, 0, nil)
	if err != nil {
		log.Printf("send err: %s", err.Error())
		return
	}

	var maxLt uint64
	var maxHash []byte
	for _, tx := range list {
		if maxLt == 0 {
			maxLt = tx.LT
			maxHash = tx.Hash
			fmt.Println(maxLt, maxHash)
		}

		fmt.Println(tx.LT, tx)

		//in := tx.IO.In
		//msg := in.Msg
		//
		//fromWallet := msg.SenderAddr().String()
		//toWallet := msg.DestAddr().String()
		//amount := in.AsInternal().Amount.Nano()
		//snake := in.AsInternal().Comment()
		////fmt.Println("snake", snake)
		//r.handler(tx, fromWallet, toWallet, amount, snake)
	}

	//for i := len(list) - 1; i >= 0; i-- {
	//	//fmt.Println(i, list[i])
	//
	//	tx := list[i]
	//	in := tx.IO.In
	//	msg := in.Msg
	//
	//	fromWallet := msg.SenderAddr().String()
	//	toWallet := msg.DestAddr().String()
	//	amount := in.AsInternal().Amount.Nano()
	//	snake := in.AsInternal().Comment()
	//	//fmt.Println("snake", snake)
	//	r.handler(tx, fromWallet, toWallet, amount, snake)
	//
	//	//r.lastLt = tx.LT
	//	//r.lastHash = tx.Hash
	//}

}
