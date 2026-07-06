package service

import (
	"database/sql"
	"encoding/json"
	"kfqt_backend/internal/model"
	"kfqt_backend/internal/repository"
	"kfqt_backend/internal/system"
	"net/http"
)

type APIEnv struct {
	DB *sql.DB
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
	if room.IsActive {
		waitTime = (waitingGroups + 1) * config.TimeRequired
	} else {
		waitTime = waitingGroups * config.TimeRequired
	}

	noticeMessage := config.MessageAvailable
	if waitTime >= 15 {
		noticeMessage = config.MessageHeavyDelay
	} else if waitTime > 0 {
		noticeMessage = config.MessageNormalDelay
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
		IsServiceAvailable: config.IsServiceAvailable,
		IsActive:      room.IsActive,
		NoticeMessage: noticeMessage,
		InfoMessage:   config.Infomation,
	}

	json.NewEncoder(w).Encode(response)
}
