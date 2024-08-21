package tonclient

import (
	"context"
	"github.com/xssnick/tonutils-go/address"
	"github.com/xssnick/tonutils-go/liteclient"
	"github.com/xssnick/tonutils-go/tlb"
	"github.com/xssnick/tonutils-go/ton"
	"log"
)

func NewClient(ctx context.Context) (*liteclient.ConnectionPool, *ton.APIClient, error) {
	pool := liteclient.NewConnectionPool()
	api := ton.NewAPIClient(pool, ton.ProofCheckPolicyFast)
	err := pool.AddConnectionsFromConfigUrl(ctx, "https://ton.org/global.config.json")
	return pool, api, err
}

func GetAccount(ctx context.Context, api ton.APIClientWrapped, addr *address.Address) (*tlb.Account, error) {
	// 每次循环，获取一次账号数据
	b, err := api.CurrentMasterchainInfo(ctx)
	if err != nil {
		log.Fatalln("get block err:", err.Error())
		return nil, err
	}
	account, err := api.GetAccount(ctx, b, addr)
	if err != nil {
		log.Fatalln("get account err:", err.Error())
		return nil, err
	}
	return account, err
}
