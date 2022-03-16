package main

import (
	"context"
	"log"

	"github.com/g8rswimmer/go-twitter/v2"
)

func requestCreateTweet(req twitter.CreateTweetRequest) *twitter.CreateTweetResponse {
	tweetResponse, err := TWITTER_CLIENT.CreateTweet(context.Background(), req)
	if err != nil {
		log.Panicf("create tweet error: %v", err)
	}
	return tweetResponse
}

// リプ
func reply2Tweet(tweet_id string, content string) string {
	req := twitter.CreateTweetRequest{
		Text: content,
		Reply: &twitter.CreateTweetReply{
			InReplyToTweetID: tweet_id,
		},
	}
	tweetResponse := requestCreateTweet(req)
	return tweetResponse.Tweet.ID
}

// 普通のツイート
func tweet(content string) string {
	req := twitter.CreateTweetRequest{
		Text: content,
	}
	tweetResponse := requestCreateTweet(req)
	return tweetResponse.Tweet.ID
}
