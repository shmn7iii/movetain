package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/g8rswimmer/go-twitter/v2"
	shell "github.com/ipfs/go-ipfs-api"
	"github.com/portto/solana-go-sdk/client"
	"github.com/portto/solana-go-sdk/rpc"
	"github.com/portto/solana-go-sdk/types"
)

// グローバル変数
var (
	JSON_KEYS      jsonKeys
	FEEPAYER       types.Account
	TWITTER_CLIENT *twitter.Client
	SOLANA_CLIENT  *client.Client
	IPFS_SHELL     *shell.Shell
)

type jsonKeys struct {
	ClientId       string `json:"clientId"`
	ClientIdSecret string `json:"clientIdSecret"`
	BotUserId      string `json:"botUserId"`
	FeePayerBase58 string `json:"feePayerBase58"`
	HostIP         string `json:"hostIp"`
}

func loadEnv() (jsonKeys jsonKeys, accessToken string, feePayer types.Account) {
	// Json読み込み
	raw, err := ioutil.ReadFile("secrets/keys.json")
	if err != nil {
		log.Fatalf("[Twitter] can't read secrets/keys.json: %v", err)
		return
	}
	json.Unmarshal(raw, &jsonKeys)

	feePayer, err = types.AccountFromBase58(jsonKeys.FeePayerBase58)
	if err != nil {
		log.Fatalf("[Solana]  can't load FeePayer: %v", err)
	}

	accessToken, err = requestAccessToken(jsonKeys)
	if err != nil {
		log.Fatalf("[Twitter] can't get ACCESS_TOKEN: %v", err)
	}

	return
}

func createClients(twitterAccessToken string) (twitterClient *twitter.Client, solanaClient *client.Client, ipfsShell *shell.Shell) {
	// twitter
	twitterClient = &twitter.Client{
		Authorizer: authorize{Token: twitterAccessToken},
		Client:     http.DefaultClient,
		Host:       "https://api.twitter.com",
	}

	// solana
	solanaClient = client.NewClient(rpc.DevnetRPCEndpoint)
	resp, err := solanaClient.GetVersion(context.TODO())
	if err != nil {
		log.Fatalf("[Solana]  Failed to version info, err: %v", err)
	}
	log.Println("[Solana]  Solana client has launched. version", resp.SolanaCore)

	// IPFS
	ipfsShell = shell.NewShell("ipfs:5001")

	return
}

func init() {
	// 各種読み込み
	var twitterAccessToken string
	JSON_KEYS, twitterAccessToken, FEEPAYER = loadEnv()
	// 各種クライアント作成
	TWITTER_CLIENT, SOLANA_CLIENT, IPFS_SHELL = createClients(twitterAccessToken)
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
