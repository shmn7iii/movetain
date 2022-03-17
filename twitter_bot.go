package main

import (
	"log"
	"strconv"
)

func botMain(latest_replied_id string) (updated_latest_replied_id string) {
	// メンションタイムラインを取得
	mention_timeline_data, err := getMentionTimelineData()
	if err != nil {
		return
	}

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
		tweet_data, err := getTweetData(tweet_id)
		if err != nil {
			continue
		}

		// 親からの呼び出しの場合は無視
		tweet_conversation_id := tweet_data.ConversationID
		if tweet_id == tweet_conversation_id {
			continue
		}

		// 親ツイートのデータを取得
		parent_tweet_data, err := getTweetData(tweet_conversation_id)
		if err != nil {
			continue
		}

		// 内容をつなげる
		memo_content := "[Movetain NFT]" +
			"\n " + parent_tweet_data.AuthorName + " @" + parent_tweet_data.AuthorUserName +
			"\n " + parent_tweet_data.TweetText +
			"\n  - " + parent_tweet_data.CreatedAt

		// メモ書く
		// txhash, err := writeMemo(memo_content)
		// if err != nil {
		// 	continue
		// }

		// NFT発行
		nftAddress, err := mintNFT(memo_content, "Null")
		if err != nil {
			continue
		}

		// 返信
		// reply_content := "🎉 Success!" +
		// 	"\nI created a Memo Transaction on Solana (devnet)." +
		// 	"\nYou can see your memo on Solana Explorer:" +
		// 	"\n https://explorer.solana.com/tx/" + txhash + "?cluster=devnet"

		reply_content := "🎉 Success!" +
			"\nI created a NFT on Solana (devnet)." +
			"\nYou can see your NFT on Solana Explorer:" +
			"\n https://explorer.solana.com/" + nftAddress + "?cluster=devnet"

		reply_id, err := reply2Tweet(tweet_id, reply_content)
		if err != nil {
			continue
		}

		log.Println("[Twitter] BOT replied:", reply_id)
	}

	// 現在のNewest IDを返す
	return mention_timeline_data.NewestID
}
