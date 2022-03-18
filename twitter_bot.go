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

	// 更新なし
	if mention_timeline_data.NewestID <= latest_replied_id {
		updated_latest_replied_id = latest_replied_id
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

		// 君らは親子かな？
		if tweet_data.AuthorUserName != parent_tweet_data.AuthorUserName {
			continue
		}

		// 内容をつなげる
		NFT_content := "[Movetain NFT] " + parent_tweet_data.AuthorName + "(@" + parent_tweet_data.AuthorUserName +
			")「" + parent_tweet_data.TweetText + "」 - " + parent_tweet_data.CreatedAt

		NFT_media_URL := parent_tweet_data.ImageURL

		// NFT発行
		nftAddress, err := mintNFT(NFT_content, NFT_media_URL)
		if err != nil {
			continue
		}

		// 返信
		reply_content := "🎉 Success!" +
			"\nI created a NFT on Solana (devnet)." +
			"\nYou can see your NFT on Solana Explorer:" +
			"\n https://explorer.solana.com/address/" + nftAddress + "?cluster=devnet"

		reply_id, err := reply2Tweet(tweet_id, reply_content)
		if err != nil {
			continue
		}

		log.Println("[Twitter] BOT replied:", reply_id)
	}

	// Newest IDを更新し返却
	updated_latest_replied_id = mention_timeline_data.NewestID
	return
}
