package model

import "kfqt_backend/internal/repository"

type UserQueueResponse struct {
	WaitTime           int    `json:"waitTime"`
	WaitingGroups      int    `json:"waitingGroups"`
	MyAheadGroups      int    `json:"myAheadGroups"`
	CurrentNumber      int    `json:"currentNumber"`
	NextNumber         int    `json:"nextNumber"`
	TimeRequired       int    `json:"timeRequired"`
	IsBookingAvailable bool   `json:"isBookingAvailable"`
	IsServiceAvailable bool   `json:"isServiceAvailable"`
	IsActive           bool   `json:"isActive"`
	NoticeMessage      string `json:"noticeMessage"` // 自動計算の混雑目安
	InfoMessage        string `json:"infoMessage"`   // configから読み込んだ運営の手動メッセージ
}

// 管理者コンソール（WebSocket）用のフルデータレスポンス
type AdminQueueResponse struct {
	WaitTime           int             `json:"waitTime"`
	WaitingGroups      int             `json:"waitingGroups"`
	CurrentNumber      int             `json:"currentNumber"`
	IsBookingAvailable bool            `json:"isBookingAvailable"`
	IsActive           bool            `json:"isActive"`
	NoticeMessage      string          `json:"noticeMessage"`
	InfoMessage        string          `json:"infoMessage"`
	Tickets            []repository.Ticket `json:"tickets"`
}
