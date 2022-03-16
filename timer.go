package main

import (
	"log"
	"strconv"
	"time"
)

// タイマー本体
func timer(latest_replied_id string) {
	ticker := time.NewTicker(20 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			// 内容
			updated_latest_replied_id := timer_do(latest_replied_id)
			// 最終返信を更新
			latest_replied_id = updated_latest_replied_id
		}
	}
}

// 毎回やる処理
func timer_do(latest_replied_id string) string {
	// メンションタイムラインを取得
	mention_timeline_data := getMentionTimelineData()
	// ツイートディクショナリーを取得
	dictionary := mention_timeline_data.TweetDictionaries

	for tweet_id := range dictionary {
		// 起動前のツイート・返信済みのツイートは無視
		tweet_id_i, _ := strconv.Atoi(tweet_id)
		latest_replied_id_i, _ := strconv.Atoi(latest_replied_id)
		if tweet_id_i <= latest_replied_id_i {
			continue
		}

		// ツイートのデータを取得
		tweet_data := getTweetData(tweet_id)

		// 親からの呼び出しの場合は無視 // TODO: 暇だったら返信とかしてもいい
		tweet_conversation_id := tweet_data.ConversationID
		if tweet_id == tweet_conversation_id {
			continue
		}

		// NFTの発行
		// 親ツイートのデータを取得
		parent_tweet_data := getTweetData(tweet_conversation_id)

		// 内容をつなげる
		memo_content := "[Movetain MEMO]" +
			"\n " + parent_tweet_data.AuthorName + " @" + parent_tweet_data.AuthorUserName +
			"\n " + parent_tweet_data.TweetText +
			"\n  - " + parent_tweet_data.CreatedAt

		// メモ書く
		txhash := writeMemo(memo_content)

		// 返信
		reply_content := "🎉 Success!" +
			"\nI created a Memo Transaction on Solana (devnet)." +
			"\nYou can see your memo on Solana Explorer:" +
			"\n https://explorer.solana.com/tx/" + txhash + "?cluster=devnet"

		reply_id := reply2Tweet(tweet_id, reply_content)

		log.Println("[Twitter] BOT replied:", reply_id)
	}

	// 現在のNewest IDを返す
	return mention_timeline_data.NewestID
}
