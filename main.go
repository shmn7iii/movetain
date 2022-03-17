package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/g8rswimmer/go-twitter/v2"
	"github.com/portto/solana-go-sdk/client"
	"github.com/portto/solana-go-sdk/rpc"
	"github.com/portto/solana-go-sdk/types"
)

// グローバル変数
var (
	ACCESS_TOKEN   string
	BOT_USER_ID    string
	FEEPAYER       types.Account
	TWITTER_CLIENT *twitter.Client
	SOLANA_CLIENT  *client.Client
)

type jsonKeys struct {
	ClientId       string `json:"clientId"`
	ClientIdSecret string `json:"clientIdSecret"`
	BotUserId      string `json:"botUserId"`
	FeePayerBase58 string `json:"feePayerBase58"`
}

func loadEnv() {
	// Json読み込み
	raw, err := ioutil.ReadFile("secrets/keys.json")
	if err != nil {
		log.Fatalf("[Twitter] can't read secrets/keys.json: %v", err)
		return
	}
	var jsonKeys jsonKeys
	json.Unmarshal(raw, &jsonKeys)

	// 変数定義
	BOT_USER_ID = jsonKeys.BotUserId
	FEEPAYER, err = types.AccountFromBase58(jsonKeys.FeePayerBase58)
	if err != nil {
		log.Fatalf("[Solana]  can't load FeePayer: %v", err)
	}
	ACCESS_TOKEN, err = requestAccessToken(jsonKeys)
	if err != nil {
		log.Fatalf("[Twitter] can't get ACCESS_TOKEN: %v", err)
	}
}

func createClients() {
	// twitter
	TWITTER_CLIENT = &twitter.Client{
		Authorizer: authorize{Token: ACCESS_TOKEN},
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

func init() {
	// 各種読み込み
	loadEnv()
	// 各種クライアント作成
	createClients()
}

func main() {
	// 起動時点のNewestID(=latest_replied_id)を取得
	mention_timeline_data, err := getMentionTimelineData()
	if err != nil {
		log.Fatalf("[Twitter] ERROR: can't get mention timeline on launch.")
	}

	// タイマー起動
	timer(mention_timeline_data.NewestID)
}
