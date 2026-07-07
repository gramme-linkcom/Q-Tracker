package service

import (
	"database/sql"
	"encoding/json"
	"kfqt_backend/internal/model"
	"kfqt_backend/internal/repository"
	"kfqt_backend/internal/system"
	"log"
	"net/http"
	"time"
	_ "time/tzdata"
)

type APIEnv struct {
	DB *sql.DB
}

// IsWithinServeTime は現在時刻が config の稼働時間内（開始〜終了）にあるかを判定します
func IsWithinServeTime(startTimeStr, endTimeStr string) bool {
	// 1. 日本時間（JST）のロケーションを取得（サーバーが海外にあってもバグらせないため）
	jst, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		log.Printf("[ERROR] タイムゾーンの読み込みに失敗: %v", err)
		return false
	}

	// 2. 現在の日本時刻を取得
	now := time.Now().In(jst)

	// 3. 比較用に「今日の年・月・日」のフォーマット文字列を作る (例: "2026-07-07 ")
	todayStr := now.Format("2006-01-02 ")

	// 4. "2006-01-02 15:04" というレイアウト型紙を使って、configの文字列を今日のTime型にパースする
	layout := "2006-01-02 15:04"
	
	parsedStart, err := time.ParseInLocation(layout, todayStr+startTimeStr, jst)
	if err != nil {
		log.Printf("[ERROR] 開始時刻のパースに失敗 (%s): %v", startTimeStr, err)
		return false
	}

	parsedEnd, err := time.ParseInLocation(layout, todayStr+endTimeStr, jst)
	if err != nil {
		log.Printf("[ERROR] 終了時刻のパースに失敗 (%s): %v", endTimeStr, err)
		return false
	}

	// 5. time.Before() と time.After() を使って、現在時刻がその間にあるか比較！
	// now が parsedStart より後（または同時）、かつ parsedEnd より前（または同時）なら true
	if (now.After(parsedStart) || now.Equal(parsedStart)) && (now.Before(parsedEnd) || now.Equal(parsedEnd)) {
		return true
	}

	return false
}

func (env *APIEnv) GetStatusHandler(w http.ResponseWriter, r *http.Request) {
	config := system.ReadConfig()

	w.Header().Set("Content-Type", "application/json")

	myNumberStr := r.URL.Query().Get("myNumber")

	myAheadGroups := 0
	if myNumberStr != "" && myNumberStr != "0" {
		myAheadGroups = repository.GetAheadGroups(env.DB, myNumberStr)
	}

	// リポジトリから純粋なDBデータを個別に取得
	room, err := repository.GetRoomStatus(env.DB)
	if err != nil {
		http.Error(w, `{"error": "ルーム状況の取得に失敗しました"}`, http.StatusInternalServerError)
		return
	}

	tickets, err := repository.GetActiveTickets(env.DB)
	if err != nil {
		http.Error(w, `{"error": "チケット一覧の取得に失敗しました"}`, http.StatusInternalServerError)
		return
	}

	// 取得したデータをもとに、API側でロジック計算を行う
	waitingGroups := len(tickets)
	currentNumber := 0
	nextNumber	  := 0
	if len(tickets) > 0 {
		currentNumber = tickets[0].Number
	}
	if len(tickets) > 1 {
		nextNumber = tickets[1].Number
	}
	

	waitTime := 0
	waitTime = waitingGroups * config.TimeRequired

	noticeMessage := config.MessageAvailable
	if waitTime >= 15 {
		noticeMessage = config.MessageHeavyDelay
	} else if waitTime > 0 {
		noticeMessage = config.MessageNormalDelay
	}

	isServiceAvailable := false
	if (config.IsServiceAvailable == true) {
		isServiceAvailable = IsWithinServeTime(config.ServeStartTime, config.ServeEndTime)
	}

	// 3. レスポンスデータを組み立てて送出
	response := model.UserQueueResponse{
		WaitTime:      waitTime,
		TimeRequired : config.TimeRequired,
		WaitingGroups: waitingGroups,
		MyAheadGroups: myAheadGroups,
		CurrentNumber:	currentNumber,
		NextNumber: 	nextNumber,
		IsBookingAvailable: config.IsBookingAvailable,
		IsServiceAvailable: isServiceAvailable,
		IsActive:      room.IsActive,
		NoticeMessage: noticeMessage,
		InfoMessage:   config.Infomation,
	}

	json.NewEncoder(w).Encode(response)
}
