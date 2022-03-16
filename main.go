package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/g8rswimmer/go-twitter/v2"
	"github.com/portto/solana-go-sdk/client"
	"github.com/portto/solana-go-sdk/rpc"
	"github.com/portto/solana-go-sdk/types"
)

var (
	BEARER_TOKEN   = flag.String("token", "", "Twitter API token")
	BOT_USER_ID    = flag.String("user_id", "", "BOT's user id")
	FeePayerBase58 = flag.String("feepayer", "", "FeePayer no base58."+
		" *keypair no base 58, private key no base 58 dato error deru")
	TWITTER_CLIENT *twitter.Client
	SOLANA_CLIENT  *client.Client
	FEEPAYER       types.Account
)

type authorize struct {
	Token string
}

func (a authorize) Add(req *http.Request) {
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", a.Token))
}

func init() {
	flag.Parse()
	FEEPAYER, _ = types.AccountFromBase58(*FeePayerBase58)

	// twitter
	TWITTER_CLIENT = &twitter.Client{
		Authorizer: authorize{Token: *BEARER_TOKEN},
		Client:     http.DefaultClient,
		Host:       "https://api.twitter.com",
	}

	// solana
	SOLANA_CLIENT = client.NewClient(rpc.DevnetRPCEndpoint)
	resp, err := SOLANA_CLIENT.GetVersion(context.TODO())
	if err != nil {
		log.Fatalf("[Solana]  Failed to version info, err: %v", err)
	}
	log.Println("[Solana]  Solana client has launched. version", resp.SolanaCore)
}

func main() {
	// 起動時点のNewestID(=latest_replied_id)を取得
	latest_replied_id := getMentionTimelineData().NewestID
	// タイマー起動
	timer(latest_replied_id)
}
