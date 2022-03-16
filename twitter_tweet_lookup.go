package main

import (
	"context"
	"log"
	"strings"

	"github.com/g8rswimmer/go-twitter/v2"
)

// API問合せ
// 直接は呼び出さない想定
func requestTweetLookup(tweet_id string) map[string]*twitter.TweetDictionary {
	opts := twitter.TweetLookupOpts{
		Expansions: []twitter.Expansion{
			twitter.ExpansionEntitiesMentionsUserName,
			twitter.ExpansionAuthorID,
		},
		TweetFields: []twitter.TweetField{
			twitter.TweetFieldCreatedAt,
			twitter.TweetFieldConversationID,
			twitter.TweetFieldAttachments,
		},
	}

	tweetResponse, err := TWITTER_CLIENT.TweetLookup(context.Background(), strings.Split(tweet_id, ","), opts)
	if err != nil {
		log.Panicf("tweet lookup error: %v", err)
	}

	dictionaries := tweetResponse.Raw.TweetDictionaries()
	return dictionaries
}

type TweetData struct {
	ID             string
	ConversationID string
	AuthorName     string
	AuthorUserName string
	TweetText      string
	Image          string
	CreatedAt      string
}

func getTweetData(tweet_id string) TweetData {
	log.Println("[Twitter] Getting TweetLookup...")

	dictionaries := requestTweetLookup(tweet_id)

	tweet_data := TweetData{
		ID:             tweet_id,
		ConversationID: dictionaries[tweet_id].Tweet.ConversationID,
		AuthorName:     dictionaries[tweet_id].Author.Name,
		AuthorUserName: dictionaries[tweet_id].Author.UserName,
		TweetText:      dictionaries[tweet_id].Tweet.Text,
		CreatedAt:      dictionaries[tweet_id].Tweet.CreatedAt,
		// 今回は画像の実装は見送り
		// Image:       dictionaries[tweet_id].AttachmentMedia,
	}
	return tweet_data
}
