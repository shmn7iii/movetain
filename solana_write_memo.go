package main

import (
	"context"
	"log"

	"github.com/portto/solana-go-sdk/common"
	"github.com/portto/solana-go-sdk/program/memoprog"
	"github.com/portto/solana-go-sdk/types"
)

func writeMemo(content string) string {
	return sendTx(createMemoTX(content))
}

func fetchRecentBlockhash() string {
	recentBlockhashResponse, err := SOLANA_CLIENT.GetRecentBlockhash(context.Background())
	if err != nil {
		log.Fatalf("failed to get recent blockhash, err: %v", err)
	}
	return recentBlockhashResponse.Blockhash
}

func createMemoTX(content string) types.Transaction {
	tx, err := types.NewTransaction(types.NewTransactionParam{
		Signers: []types.Account{FEEPAYER, FEEPAYER},
		Message: types.NewMessage(types.NewMessageParam{
			FeePayer:        FEEPAYER.PublicKey,
			RecentBlockhash: fetchRecentBlockhash(),
			Instructions: []types.Instruction{
				memoprog.BuildMemo(memoprog.BuildMemoParam{
					SignerPubkeys: []common.PublicKey{FEEPAYER.PublicKey},
					Memo:          []byte(content),
				}),
			},
		}),
	})
	if err != nil {
		log.Fatalf("failed to new a transaction, err: %v", err)
	}
	return tx
}

func sendTx(tx types.Transaction) string {
	txhash, err := SOLANA_CLIENT.SendTransaction(context.Background(), tx)
	if err != nil {
		log.Fatalf("failed to send tx, err: %v", err)
	}
	log.Println("[Solana]  BOT has created a transaction")
	log.Println("[Solana]    Tx Hash:", txhash)
	return txhash
}
