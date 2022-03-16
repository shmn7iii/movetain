package main

import (
	"context"
	"log"

	"github.com/portto/solana-go-sdk/common"
	"github.com/portto/solana-go-sdk/program/memoprog"
	"github.com/portto/solana-go-sdk/rpc"
	"github.com/portto/solana-go-sdk/types"
)

// メモをトランザクションに埋め込み
func writeMemo(content string) (txhash string, err error) {
	// fetch recent block hash
	recentBlockhashResponse, err := fetchRecentBlockhash()
	if err != nil {
		log.Println("[Solana]  ERROR: can't fetch blockhash:", err)
		return
	}
	blockhash := recentBlockhashResponse.Blockhash

	// create transaction
	tx, err := createMemoTX(content, blockhash)
	if err != nil {
		log.Println("[Solana]  ERROR: can't create transaction:", err)
		return
	}

	// send transaction
	txhash, err = sendTx(tx)
	if err != nil {
		log.Println("[Solana]  ERROR: can't send transaction:", err)
		return
	}

	// OK!
	log.Println("[Solana]  BOT has created a transaction")
	log.Println("[Solana]   Tx Hash:", txhash)
	return
}

// 以下直接は呼び出さない想定

// fetch recent block hash
func fetchRecentBlockhash() (recentBlockhashResponse rpc.GetRecentBlockHashResultValue, err error) {
	recentBlockhashResponse, err = SOLANA_CLIENT.GetRecentBlockhash(context.Background())
	return
}

// create transaction
func createMemoTX(content string, blockhash string) (tx types.Transaction, err error) {
	tx, err = types.NewTransaction(types.NewTransactionParam{
		Signers: []types.Account{FEEPAYER, FEEPAYER},
		Message: types.NewMessage(types.NewMessageParam{
			FeePayer:        FEEPAYER.PublicKey,
			RecentBlockhash: blockhash,
			Instructions: []types.Instruction{
				memoprog.BuildMemo(memoprog.BuildMemoParam{
					SignerPubkeys: []common.PublicKey{FEEPAYER.PublicKey},
					Memo:          []byte(content),
				}),
			},
		}),
	})
	return
}

// send transaction
func sendTx(tx types.Transaction) (txhash string, err error) {
	txhash, err = SOLANA_CLIENT.SendTransaction(context.Background(), tx)
	return
}
