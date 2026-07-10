package service

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"kfqt_backend/internal/model"
	"kfqt_backend/internal/repository"
	"kfqt_backend/internal/system"
	"log"
	"net/http"
	"strings"
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
		if !cfg.AllowNoTimeSlot {
			http.Error(w, `{"error": "ただいま時間指定なしでの整理券発行を停止しております。お手数ですが、ドロップダウンより時間枠をご指定ください。"}`, http.StatusBadRequest)
			return
		}
		reservedTime = GetCurrentTimeSlot()
	} else {
		// 指定された時間枠が過去の枠（またはすでに開始している枠）でないかチェック
		if IsSlotPast(reservedTime) {
			http.Error(w, `{"error": "指定された時間帯はすでに受付を終了しています"}`, http.StatusBadRequest)
			return
		}
	}

	// 時間帯の予約上限チェック（時間指定あり・なし両方で、優先枠と当日枠の合算値を確認）
	count, err := GetSlotBookingCount(env.DB, reservedTime)
	maxBookings := cfg.MaxBookingsPerSlot
	if maxBookings <= 0 {
		maxBookings = 5 // 安全フォールバック
	}
	if err == nil && count >= maxBookings {
		http.Error(w, `{"error": "指定された時間帯は予約上限に達しています"}`, http.StatusBadRequest)
		return
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

func IsSlotPast(slot string) bool {
	if slot == "" {
		return false
	}
	parts := strings.Split(slot, " - ")
	if len(parts) < 2 {
		return false
	}
	startTimeStr := parts[0]

	jst, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		jst = time.UTC
	}
	now := time.Now().In(jst)
	todayStr := now.Format("2006-01-02 ")

	parsedStart, err := time.ParseInLocation("2006-01-02 15:04", todayStr+startTimeStr, jst)
	if err != nil {
		return false
	}

	return now.After(parsedStart) || now.Equal(parsedStart)
}

func GetSlotBookingCount(db *sql.DB, reservedTime string) (int, error) {
	baseSlot := strings.Replace(reservedTime, " (無指定)", "", 1)
	standbySlot := baseSlot + " (無指定)"

	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM tickets WHERE (reserved_time = ? OR reserved_time = ?) AND status IN ('waiting', 'serving')", baseSlot, standbySlot).Scan(&count)
	return count, err
}
