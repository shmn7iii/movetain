package main

import (
	"time"
)

// タイマー本体
func timer(latest_replied_id string) {
	ticker := time.NewTicker(20 * time.Second)
	defer ticker.Stop()
	count := 0
	for {
		select {
		case <-ticker.C:
			// 内容
			updated_latest_replied_id := botMain(latest_replied_id)
			// 最終返信を更新
			latest_replied_id = updated_latest_replied_id
			// カウント
			count++
			// 2時間(360x20)経ったらクライアントを再起動
			if count == 350 {
				TWITTER_CLIENT = regenerateTwitterClient(JSON_KEYS)
				count = 0
			}
		}
	}
}
