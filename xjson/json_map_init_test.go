package xjson

import (
	"fmt"
	"testing"
)

func TestConvertToMap(t *testing.T) {
	r, ee := ConvertToMap(`{"order_no": "T111831719329815", "fiat_code": "USD", "crypto_code": "USDT-BSC", "fiat_amount": "666", "product_name": "Buy-USDT-BSC", "crypto_amount": "0", "merchant_order_no": "234124313333"}`)
	fmt.Println(r, ee)
}
