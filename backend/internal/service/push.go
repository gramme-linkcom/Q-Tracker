package service

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	webpush "github.com/SherClockHolmes/webpush-go"
)

// クライアントに教えるVAPID公開キー
var VapidPublicKey = os.Getenv("PUSH_PUBLIC_KEY")
var vapidPrivateKey = os.Getenv("PUSH_PRIVATE_KEY")
var SubscriberStr = os.Getenv("SUBSCRIBER")

// VapidPublicKeyHandler はフロントに公開キーを渡します
func VapidPublicKeyHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte(VapidPublicKey))
}

// SendPushToUser は指定された宛先（JSON文字列）に対してピンポイントでPush通知を送信します
func SendPushToUser(subscriptionJSON string, message string) {
	if subscriptionJSON == "" || subscriptionJSON == "manual_issued_token" {
		// 手動発券やトークンがない場合はスキップ
		return
	}

	// 1. 保存されていたJSON文字列を WebPush の構造体にデコード
	var sub webpush.Subscription
	if err := json.Unmarshal([]byte(subscriptionJSON), &sub); err != nil {
		log.Printf("[PUSH_ERROR] 宛先JSONの解析に失敗しました: %v", err)
		return
	}

	// 2. 通知ペイロード（メッセージ）を作成
	payload, _ := json.Marshal(map[string]string{
		"title": "Q-Tracker 呼び出し",
		"body":  message,
	})

	// 3. 送信実行！
	resp, err := webpush.SendNotification(payload, &sub, &webpush.Options{
		Subscriber:      SubscriberStr,
		VAPIDPublicKey:  VapidPublicKey,
		VAPIDPrivateKey: vapidPrivateKey,
		TTL:             30, // 30秒間届かなければ消滅（古い通知が後から溜まるのを防ぐ）
	})

	if err != nil {
		log.Printf("[PUSH_ERROR] 送信に失敗しました: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		log.Printf("[PUSH_SUCCESS] ユーザーへのPush通知送信に成功しました (%d)", resp.StatusCode)
	} else {
		log.Printf("[PUSH_WARN] Pushサーバーがエラーを返しました status=%d", resp.StatusCode)
	}
}
