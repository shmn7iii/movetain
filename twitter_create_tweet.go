package main

import (
	"context"
	"log"

	"github.com/g8rswimmer/go-twitter/v2"
)

// リプ
func reply2Tweet(reply_to_tweet_id string, content string) (tweet_id string, err error) {
	req := twitter.CreateTweetRequest{
		Text: content,
		Reply: &twitter.CreateTweetReply{
			InReplyToTweetID: reply_to_tweet_id,
		},
	}
	tweetResponse, err := requestCreateTweet(req)
	if err != nil {
		log.Println("[Twitter] ERROR: can't create reply tweet:", err)
		return
	}

	tweet_id = tweetResponse.Tweet.ID
	return
}

// 普通のツイート（たぶん使わん）
// func tweet(content string) (tweet_id string, err error) {
// 	req := twitter.CreateTweetRequest{
// 		Text: content,
// 	}
// 	tweetResponse, err := requestCreateTweet(req)
// 	if err != nil {
// 		log.Println("[Twitter] ERROR: can't create tweet:", err)
// 		return
// 	}

// 	tweet_id = tweetResponse.Tweet.ID
// 	return
// }

// 以下直接は呼び出さない想定

// APIへリクエストを送信
func requestCreateTweet(req twitter.CreateTweetRequest) (tweetResponse *twitter.CreateTweetResponse, err error) {
	tweetResponse, err = TWITTER_CLIENT.CreateTweet(context.Background(), req)
	return
}
