package service

import (
	"encoding/json"
	"fmt"
	"kfqt_backend/internal/model"
	"kfqt_backend/internal/repository"
	"kfqt_backend/internal/system"
	"log"
	"net/http"
	"time"
)

func (env *APIEnv) BookTicketHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	cfg := system.ReadConfig()

	if r.Method != http.MethodPost {
		http.Error(w, `{"error": "Method Not Allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	// 1. フロントからのJSON（トークン）をデコード
	var req model.BookingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "不正なリクエストデータです"}`, http.StatusBadRequest)
		return
	}

	if (!cfg.IsBookingAvailable || !IsWithinServeTime(cfg.ServeStartTime, cfg.ServeEndTime) || !cfg.IsServiceAvailable) {
		http.Error(w, `{"error": "ただいま整理券の新規発行を停止しております"}`, http.StatusBadRequest)
		return
	}

	reservedTime := req.ReservedTime
	if reservedTime == "" {
		reservedTime = GetCurrentTimeSlot()
	} else {
		// 指定された時間枠の予約上限チェック（時間指定なしの場合は制限なし）
		var count int
		err := env.DB.QueryRow("SELECT COUNT(*) FROM tickets WHERE reserved_time = ? AND status IN ('waiting', 'serving')", reservedTime).Scan(&count)
		maxBookings := cfg.MaxBookingsPerSlot
		if maxBookings <= 0 {
			maxBookings = 5 // 安全フォールバック
		}
		if err == nil && count >= maxBookings {
			http.Error(w, `{"error": "指定された時間帯は予約上限に達しています"}`, http.StatusBadRequest)
			return
		}
	}

	bookingData, err := repository.CreateUserTicket(env.DB, req.PushToken, reservedTime)
	if err != nil {
		log.Printf("[ERROR] 整理券の発行失敗: %v", err)
		http.Error(w, `{"error": "Server error"}`, http.StatusInternalServerError)
		return
	}

	log.Printf("[INFO] 整理券を発行(発行者: ユーザー): 番号=%d", bookingData.TicketNumber)

	// 管理者コンソールへ送出
	tickets, err := repository.GetActiveTickets(env.DB)
	if err == nil {
		var queueData []interface{}
		for _, t := range tickets {
			queueData = append(queueData, t)
		}

		BroadcastQueue(BroadcastDatas{
			PushType: "queue_update",
			Queue:    queueData,
		})
	}

	response := model.BookingResponse{
		BookingNumber: bookingData.TicketNumber,
		Uuid:	bookingData.Uuid,
		Success:       true,
	}

	w.WriteHeader(http.StatusCreated) // 201 Created
	json.NewEncoder(w).Encode(response)
}

func GetCurrentTimeSlot() string {
	cfg := system.ReadConfig()
	interval := cfg.SlotInterval
	if interval <= 0 {
		interval = 30 // 安全フォールバック
	}

	jst, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		jst = time.UTC
	}
	now := time.Now().In(jst)
	
	// 総経過時間（分単位）で計算する
	currentTotalMinutes := now.Hour()*60 + now.Minute()
	
	startTotalMinutes := (currentTotalMinutes / interval) * interval
	endTotalMinutes := startTotalMinutes + interval
	
	startHour := (startTotalMinutes / 60) % 24
	startMin := startTotalMinutes % 60
	
	endHour := (endTotalMinutes / 60) % 24
	endMin := endTotalMinutes % 60
	
	return fmt.Sprintf("%02d:%02d - %02d:%02d (無指定)", startHour, startMin, endHour, endMin)
}
