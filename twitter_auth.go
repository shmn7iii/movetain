package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/g8rswimmer/go-twitter/v2"
)

type authorize struct {
	Token string
}

func (a authorize) Add(req *http.Request) {
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", a.Token))
}

type responseJson struct {
	TokenType    string `json:"token_type"`
	ExpiresIn    string `json:"expires_in"`
	AccessToken  string `json:"access_token"`
	Scope        string `json:"scope"`
	RefreshToken string `json:"refresh_token"`
}

func regenerateTwitterClient(jsonKeys jsonKeys) (newClient *twitter.Client) {
	log.Println("[Twitter] Regenerating Twitter client.")

	// アクセストークンを再生成
	newToken, err := requestAccessToken(jsonKeys)
	if err != nil {
		log.Fatalf("[Twitter] can't get ACCESS_TOKEN: %v", err)
	}

	// クライアント再生成
	newClient = &twitter.Client{
		Authorizer: authorize{Token: newToken},
		Client:     http.DefaultClient,
		Host:       "https://api.twitter.com",
	}
	return
}

func requestAccessToken(jsonKeys jsonKeys) (accessToken string, err error) {
	// secrets/refreshtokenから読み込み
	bytes, err := ioutil.ReadFile("secrets/refreshtoken")
	if err != nil {
		log.Fatalf("[Twitter] can't read secrets/refreshtoken: %v", err)
		return
	}
	// 今回使うリフレッシュトークンを設定
	refresh_token := string(bytes)

	// 新しくアクセストークンを取得

	// data-urlencode
	values := url.Values{}
	values.Add("refresh_token", refresh_token)
	values.Add("grant_type", "refresh_token")

	// リクエストを作成
	req, err := http.NewRequest(
		"POST",
		"https://api.twitter.com/2/oauth2/token",
		strings.NewReader(values.Encode()),
	)
	if err != nil {
		log.Println("[Twitter] ERROR: can't create new http request:", err)
		return
	}

	// ヘッダーを設定
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Basic認証
	req.SetBasicAuth(jsonKeys.ClientId, jsonKeys.ClientIdSecret)

	// クライアント作成・実行
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("[Twitter] ERROR: can't do http request:", err)
		return
	}

	// レスポンスBodyを読み込み
	byteArray, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("[Twitter] ERROR: can't read http response body:", err)
		return
	}

	defer resp.Body.Close()

	// 構造体に落とし込む
	jsonBytes := ([]byte)(byteArray)
	var responseJson responseJson
	json.Unmarshal(jsonBytes, &responseJson)

	// 新しいリフレッシュトークンを保存
	f, err := os.Create("secrets/refreshtoken")
	data := []byte(responseJson.RefreshToken)
	_, err = f.Write(data)
	if err != nil {
		log.Println("[Twitter] ERROR: can't write secrets/refreshtoken:", err)
		log.Println("[Twitter] RefreshToken:", refresh_token)
		return
	}

	// 取得したアクセストークン
	access_token = responseJson.AccessToken
	return
}
