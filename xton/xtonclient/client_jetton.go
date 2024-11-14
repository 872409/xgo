package tonclient

import (
	"context"
	"errors"
	"github.com/xssnick/tonutils-go/address"
	"github.com/xssnick/tonutils-go/tlb"
	"github.com/xssnick/tonutils-go/ton"
	"github.com/xssnick/tonutils-go/ton/jetton"
)

func NewJettonClient(workerCtx context.Context, api ton.APIClientWrapped, treasuryAddress *address.Address, jettonMaster *address.Address) (*JettonClient, error) {

	jettonMasterClient := jetton.NewJettonMasterClient(api, jettonMaster)

	// get our jetton wallet address
	treasuryJettonWallet, err := jettonMasterClient.GetJettonWallet(workerCtx, treasuryAddress)
	if err != nil {
		return nil, err
	}

	return &JettonClient{
		treasuryAddress:      treasuryAddress,
		treasuryJettonWallet: treasuryJettonWallet,
	}, nil

}

type JettonClient struct {
	treasuryAddress      *address.Address
	treasuryJettonWallet *jetton.WalletClient
}

type JettonTransfer struct {
	SrcAddr              *address.Address
	Comment              string
	TransferNotification *jetton.TransferNotification
}

func (receiver *JettonClient) HandleInTx(tx *tlb.Transaction) (*JettonTransfer, error) {
	if tx.IO.In != nil && tx.IO.In.MsgType == tlb.MsgTypeInternal {
		ti := tx.IO.In.AsInternal()
		src := ti.SrcAddr

		// verify that sender is our jetton wallet
		if ti.SrcAddr.Equals(receiver.treasuryJettonWallet.Address()) {
			var transfer jetton.TransferNotification
			err := tlb.LoadFromCell(&transfer, ti.Body.BeginParse())
			if err != nil {
				return nil, err
			}
			payloadSlice := transfer.ForwardPayload.BeginParse()
			payloadOpcode, _ := payloadSlice.LoadUInt(32)
			//fmt.Println("payloadOpcode", payloadOpcode, er)
			if payloadOpcode != 0 {
				return nil, errors.New("invalid payload opcode")
			}

			payloadComment := payloadSlice.MustLoadStringSnake()
			//fmt.Println("payloadComment", payloadComment)
			// convert decimals to 6 for USDT (it can be fetched from jetton details too), default is 9
			//amt := tlb.MustFromNano(transfer.Amount.Nano(), 6)

			// reassign sender to real jetton sender instead of its jetton wallet contract
			src = transfer.Sender
			//log.Println("received", amt.String(), "USDT from", src.String(), payloadComment)
			return &JettonTransfer{
				SrcAddr:              src,
				Comment:              payloadComment,
				TransferNotification: &transfer,
			}, nil

		}
	}
	return nil, errors.New("invalid transaction")
}
