package xton

import (
	"context"
	"github.com/xssnick/tonutils-go/liteclient"
	"github.com/xssnick/tonutils-go/ton"
)

func GetAPIClient(ctx context.Context) (*ton.APIClient, context.Context) {
	client := liteclient.NewConnectionPool()

	// 连接到测试网lite服务器
	err := client.AddConnectionsFromConfigUrl(ctx, "https://ton.org/global.config.json")
	if err != nil {
		panic(err)
	}

	ctxSticky := client.StickyContext(ctx)

	//seed := wallet.NewSeed()

	// 初始化ton api lite连接包装器
	api := ton.NewAPIClient(client)
	return api, ctxSticky
}
