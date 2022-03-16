package main

import (
	"context"
	"log"

	"github.com/g8rswimmer/go-twitter/v2"
)

// API問合せ
// 直接は呼び出さない想定
func requestUserMentionTimeline() *twitter.UserMentionTimelineResponse {
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

	timeline, err := TWITTER_CLIENT.UserMentionTimeline(context.Background(), *BOT_USER_ID, opts)
	if err != nil {
		// TODO: Panicでいい？
		log.Panicf("[Twitter] user mention timeline error: %v", err)
	}

	return timeline
}

// 必要なものだけ入れた構造体
type MentionTimelineData struct {
	TweetDictionaries map[string]*twitter.TweetDictionary
	NewestID          string
}

// データを取得
func getMentionTimelineData() MentionTimelineData {
	log.Println("[Twitter] Getting UserMentionTimeline...")

	timeline := requestUserMentionTimeline()

	mention_timeline_data := MentionTimelineData{
		TweetDictionaries: timeline.Raw.TweetDictionaries(),
		NewestID:          timeline.Meta.NewestID,
	}
	return mention_timeline_data
}
