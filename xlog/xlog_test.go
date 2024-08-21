package xlog

import (
	"testing"
)

func TestInfo(t *testing.T) {
	//type args struct {
	//	channel string
	//	msg     string
	//	data    []interface{}
	//}
	//tests := []struct {
	//	name string
	//	args args
	//}{
	//	{name: "aaa", args: args{channel: DefaultChannel, msg: "DefaultChannel", data: []interface{}{"1", "0"}}},
	//	{name: "aaa", args: args{channel: UserEarningChannel, msg: "err", data: []interface{}{"UserEarningChannel", "1"}}},
	//	{name: "aaa", args: args{channel: UserEarningChannel, msg: "err", data: []interface{}{"UserEarningChannel", "2"}}},
	//	{name: "aaa", args: args{channel: UserEarningChannel, msg: "err", data: []interface{}{"UserEarningChannel", "3"}}},
	//	{name: "aaa", args: args{channel: UserEarningChannel, msg: "err", data: []interface{}{"UserEarningChannel", "4"}}},
	//	{name: "aaa", args: args{channel: UserEarningChannel, msg: "err", data: []interface{}{"UserEarningChannel", "5"}}},
	//	{name: "aaa", args: args{channel: UserEarningChannel, msg: "err", data: []interface{}{"UserEarningChannel", "6"}}},
	//}
	//for _, tt := range tests {
	//	t.Run(tt.name, func(t *testing.T) {
	//		go Default.Info(tt.args.msg, tt.args.data...)
	//		go Wallet.Error(tt.args.msg, tt.args.data...)
	//	})
	//}
	//
	//time.Sleep(2 * time.Second)
}
