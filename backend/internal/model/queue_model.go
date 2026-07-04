package model

import "kfqt_backend/internal/repository"

type QueueStatusResponse struct {
	WaitTime      int                 `json:"waitTime"`       // 計算された待ち時間（分）
	WaitingGroups int                 `json:"waitingGroups"`   // 待機列にいるリアルな組数
	CurrentNumber int                 `json:"currentNumber"`   // 現在案内中の整理券番号
	IsActive      bool                `json:"isActive"`       // 部屋が稼働中かどうか
	NoticeMessage string              `json:"noticeMessage"`  // configから引っ張る動的なお知らせ
	Tickets       []repository.Ticket `json:"tickets"`       // 管理画面用の待機中のチケット一覧（空なら自動で[]になる）
}
