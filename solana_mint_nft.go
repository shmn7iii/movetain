package main

import (
	"context"
	"log"
	"strings"

	"github.com/portto/solana-go-sdk/common"
	"github.com/portto/solana-go-sdk/pkg/pointer"
	"github.com/portto/solana-go-sdk/program/assotokenprog"
	"github.com/portto/solana-go-sdk/program/metaplex/tokenmeta"
	"github.com/portto/solana-go-sdk/program/sysprog"
	"github.com/portto/solana-go-sdk/program/tokenprog"
	"github.com/portto/solana-go-sdk/types"
)

func mintNFT(content string, imageURL string) (nftAddress string, err error) {
	// contentからjsonを作る
	jsonStringReader := generateJsonStringReader(content, imageURL)
	// jsonをIPFSにあげる
	jsonCID, err := uploadJson2ipfs(jsonStringReader)
	if err != nil {
		log.Println("[Solana] can't upload json to IPFS:", err)
		return
	}

	// mint
	nftAddress, sig, err := requestMintToken(jsonCID)
	if err != nil {
		log.Println("[Solana] can't mint NFT:", err)
		return
	}

	log.Println("[Solana ] 🪪 BOT has minted a NFT")
	log.Println("[Solana ]      Account:  ", nftAddress)
	log.Println("[Solana ]      Signature:", sig)

	return
}

func generateJsonStringReader(content string, imageURL string) (jsonStringReader *strings.Reader) {
	jsonStringReader = strings.NewReader(
		"{" +
			"\n  \"name\": \"Movetain Tweet Token\"," +
			"\n  \"description\": \"" + content + "\"," +
			"\n  \"image\": \"" + imageURL + "\"," +
			"\n  \"external_url\": \"https://www.shmn7iii.net/movetain\"" +
			"\n}" +
			"")
	return
}

// 以下直接は呼び出さない想定

// ミントをリクエスト
func requestMintToken(jsonCID string) (nftAddress string, sig string, err error) {
	mint := types.NewAccount()
	nftAddress = mint.PublicKey.ToBase58()

	ata, _, err := common.FindAssociatedTokenAddress(FEEPAYER.PublicKey, mint.PublicKey)
	if err != nil {
		log.Println("[Solana]  failed to find a valid ata, err:", err)
		return
	}

	tokenMetadataPubkey, err := tokenmeta.GetTokenMetaPubkey(mint.PublicKey)
	if err != nil {
		log.Println("[Solana]  failed to find a valid token metadata, err:", err)
		return
	}

	tokenMasterEditionPubkey, err := tokenmeta.GetMasterEdition(mint.PublicKey)
	if err != nil {
		log.Println("[Solana]  failed to find a valid master edition, err:", err)
		return
	}

	mintAccountRent, err := SOLANA_CLIENT.GetMinimumBalanceForRentExemption(context.Background(), tokenprog.MintAccountSize)
	if err != nil {
		log.Println("[Solana]  failed to get mint account rent, err:", err)
		return
	}

	recentBlockhashResponse, err := SOLANA_CLIENT.GetRecentBlockhash(context.Background())
	if err != nil {
		log.Println("[Solana]  failed to get recent blockhash, err:", err)
		return
	}

	tx, err := types.NewTransaction(types.NewTransactionParam{
		Signers: []types.Account{mint, FEEPAYER},
		Message: types.NewMessage(types.NewMessageParam{
			FeePayer:        FEEPAYER.PublicKey,
			RecentBlockhash: recentBlockhashResponse.Blockhash,
			Instructions: []types.Instruction{
				sysprog.CreateAccount(sysprog.CreateAccountParam{
					From:     FEEPAYER.PublicKey,
					New:      mint.PublicKey,
					Owner:    common.TokenProgramID,
					Lamports: mintAccountRent,
					Space:    tokenprog.MintAccountSize,
				}),
				tokenprog.InitializeMint(tokenprog.InitializeMintParam{
					Decimals: 0,
					Mint:     mint.PublicKey,
					MintAuth: FEEPAYER.PublicKey,
				}),
				tokenmeta.CreateMetadataAccount(tokenmeta.CreateMetadataAccountParam{
					Metadata:                tokenMetadataPubkey,
					Mint:                    mint.PublicKey,
					MintAuthority:           FEEPAYER.PublicKey,
					Payer:                   FEEPAYER.PublicKey,
					UpdateAuthority:         FEEPAYER.PublicKey,
					UpdateAuthorityIsSigner: true,
					IsMutable:               true,
					MintData: tokenmeta.Data{
						Name:                 "Movetain Tweet Token",
						Symbol:               "MTT",
						Uri:                  "http://" + HOST_IP + ":8080/ipfs/" + jsonCID,
						SellerFeeBasisPoints: 100,
						Creators: &[]tokenmeta.Creator{
							{
								Address:  FEEPAYER.PublicKey,
								Verified: true,
								Share:    100,
							},
						},
					},
				}),
				assotokenprog.CreateAssociatedTokenAccount(assotokenprog.CreateAssociatedTokenAccountParam{
					Funder:                 FEEPAYER.PublicKey,
					Owner:                  FEEPAYER.PublicKey,
					Mint:                   mint.PublicKey,
					AssociatedTokenAccount: ata,
				}),
				tokenprog.MintTo(tokenprog.MintToParam{
					Mint:   mint.PublicKey,
					To:     ata,
					Auth:   FEEPAYER.PublicKey,
					Amount: 1,
				}),
				tokenmeta.CreateMasterEdition(tokenmeta.CreateMasterEditionParam{
					Edition:         tokenMasterEditionPubkey,
					Mint:            mint.PublicKey,
					UpdateAuthority: FEEPAYER.PublicKey,
					MintAuthority:   FEEPAYER.PublicKey,
					Metadata:        tokenMetadataPubkey,
					Payer:           FEEPAYER.PublicKey,
					MaxSupply:       pointer.Uint64(0),
				}),
			},
		}),
	})
	if err != nil {
		log.Println("[Solana]  failed to new a tx, err:", err)
		return
	}

	sig, err = SOLANA_CLIENT.SendTransaction(context.Background(), tx)
	if err != nil {
		log.Println("[Solana]  failed to send tx, err:", err)
		return
	}

	return
}
