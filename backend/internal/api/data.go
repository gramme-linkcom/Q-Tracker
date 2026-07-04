package api

import (
	"database/sql"
	"encoding/json"
	"kfqt_backend/internal/repository"
	"net/http"
)

type APIEnv struct {
	DB *sql.DB
}

// 💡 1グループあたりのプレイ時間設定（5分）
const PlayTimePerGroup = 5

type UserQueueResponse struct {
	WaitTime      int		`json:"waitTime"`
	WaitingGroups int		`json:"waitingGroups"`
	CurrentNumber int		`json:"currentNumber"`
	IsActive      bool		`json:"isActive"`
	NoticeMessage string	`json:"noticeMessage"` // 自動計算の混雑目安
	InfoMessage   string	`json:"infoMessage"`   // configから読み込んだ運営の手動メッセージ
}

// 💻 2. 管理者コンソール（WebSocket）用のフルデータレスポンス
type AdminQueueResponse struct {
	WaitTime      int					`json:"waitTime"`
	WaitingGroups int					`json:"waitingGroups"`
	CurrentNumber int					`json:"currentNumber"`
	IsActive      bool					`json:"isActive"`
	NoticeMessage string				`json:"noticeMessage"`
	InfoMessage   string 				`json:"infoMessage"`
	Tickets       []repository.Ticket	`json:"tickets"`
}

func (env *APIEnv) GetStatusHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// 1. リポジトリから純粋なDBデータを個別に取得
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

	// 2. 取得したデータをもとに、API側でロジック計算を行う
	waitingGroups := len(tickets)

	waitTime := 0
	if room.IsActive {
		waitTime = (waitingGroups + 1) * PlayTimePerGroup
	} else {
		waitTime = waitingGroups * PlayTimePerGroup
	}

	noticeMessage := "現在、すぐにご案内できます。"
	if waitTime >= 15 {
		noticeMessage = "現在、ご案内までに大幅にお時間がかかります。"
	} else if waitTime > 0 {
		noticeMessage = "現在、ご案内までに多少お時間がかかります。"
	}

	// 3. レスポンスデータを組み立てて送出
	response := UserQueueResponse{
		WaitTime:      waitTime,
		WaitingGroups: waitingGroups,
		CurrentNumber: room.CurrentNumber,
		IsActive:      room.IsActive,
		NoticeMessage: noticeMessage,
		InfoMessage:   "",
	}

	json.NewEncoder(w).Encode(response)
}
