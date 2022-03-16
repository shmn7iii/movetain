package main

import (
	"context"
	"log"

	"github.com/g8rswimmer/go-twitter/v2"
)

// 必要なものだけ入れた構造体
type MentionTimelineData struct {
	TweetDictionaries map[string]*twitter.TweetDictionary
	NewestID          string
}

// データを取得
func getMentionTimelineData() (mention_timeline_data MentionTimelineData, err error) {
	log.Println("[Twitter] Getting UserMentionTimeline...")

	timeline, err := requestUserMentionTimeline()
	if err != nil {
		log.Println("[Twitter] ERROR: can't get mention timeline:", err)
		return
	}

	mention_timeline_data = MentionTimelineData{
		TweetDictionaries: timeline.Raw.TweetDictionaries(),
		NewestID:          timeline.Meta.NewestID,
	}
	return
}

// 以下直接は呼び出さない想定

// APIへリクエストを送信
func requestUserMentionTimeline() (timeline *twitter.UserMentionTimelineResponse, err error) {
	opts := twitter.UserMentionTimelineOpts{
		TweetFields: []twitter.TweetField{
			twitter.TweetFieldCreatedAt,
			twitter.TweetFieldAuthorID,
			twitter.TweetFieldConversationID,
			twitter.TweetFieldPublicMetrics,
			twitter.TweetFieldContextAnnotations,
		},
		UserFields: []twitter.UserField{twitter.UserFieldUserName},
		Expansions: []twitter.Expansion{twitter.ExpansionAuthorID},
		MaxResults: 5,
	}
	timeline, err = TWITTER_CLIENT.UserMentionTimeline(context.Background(), *BOT_USER_ID, opts)
	return
}
