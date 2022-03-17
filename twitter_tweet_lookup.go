package main

import (
	"context"
	"log"
	"strings"

	"github.com/g8rswimmer/go-twitter/v2"
)

// 必要なものだけ入れた構造体
type TweetData struct {
	ID             string
	ConversationID string
	AuthorName     string
	AuthorUserName string
	TweetText      string
	ImageURL       string
	CreatedAt      string
}

// データを取得
func getTweetData(tweet_id string) (tweet_data TweetData, err error) {
	log.Println("[Twitter] Getting TweetLookup...")

	tweetResponse, err := requestTweetLookup(tweet_id)
	if err != nil {
		log.Println("[Twitter] ERROR: can't get tweet data:", err)
		return
	}

	dictionaries := tweetResponse.Raw.TweetDictionaries()
	tweet_data = TweetData{
		ID:             tweet_id,
		ConversationID: dictionaries[tweet_id].Tweet.ConversationID,
		AuthorName:     dictionaries[tweet_id].Author.Name,
		AuthorUserName: dictionaries[tweet_id].Author.UserName,
		TweetText:      dictionaries[tweet_id].Tweet.Text,
		CreatedAt:      dictionaries[tweet_id].Tweet.CreatedAt,
	}

	if len(dictionaries[tweet_id].AttachmentMedia) != 0 {
		tweet_data.ImageURL = dictionaries[tweet_id].AttachmentMedia[0].URL
	}
	return
}

// 以下直接は呼び出さない想定

// APIへリクエストを送信
func requestTweetLookup(tweet_id string) (tweetResponse *twitter.TweetLookupResponse, err error) {
	opts := twitter.TweetLookupOpts{
		Expansions: []twitter.Expansion{
			twitter.ExpansionEntitiesMentionsUserName,
			twitter.ExpansionAttachmentsMediaKeys,
			twitter.ExpansionAuthorID,
		},
		TweetFields: []twitter.TweetField{
			twitter.TweetFieldCreatedAt,
			twitter.TweetFieldConversationID,
			twitter.TweetFieldAttachments,
		},
		MediaFields: []twitter.MediaField{
			twitter.MediaFieldURL,
		},
	}
	tweetResponse, err = TWITTER_CLIENT.TweetLookup(context.Background(), strings.Split(tweet_id, ","), opts)
	return
}
