package model


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
	ReservedTime       string `json:"reservedTime"`  // 指定された予約時間
	ServeStartTime     string `json:"serveStartTime"` // 稼働開始時間 (HH:MM)
	ServeEndTime       string         `json:"serveEndTime"`   // 稼働終了時間 (HH:MM)
	SlotInterval       int            `json:"slotInterval"`   // 枠の粒度 (分)
	MaxBookingsPerSlot int            `json:"maxBookingsPerSlot"` // 1枠あたりの最大予約数
	SlotBookings       map[string]int `json:"slotBookings"`       // 時間枠ごとの現在の予約数
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
	Tickets            []Ticket `json:"tickets"`
}

type Ticket struct {
	Number 		int    `json:"number"`
	Uuid		string `json:"uuid"`
	DeviceID	string `json:"device_id"`
	Status 		string `json:"status"` // "waiting", "called", "canceled", "absent"
	ReservedTime string `json:"reserved_time"` // 指定された予約時間
}

type ResultTicket struct {
	Uuid	string		`json:"uuid"`
	TicketNumber int	`json:"ticketNumber"` 
}
