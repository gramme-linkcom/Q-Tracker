package service

import (
	"encoding/json"
	"kfqt_backend/internal/model"
	"kfqt_backend/internal/repository"
	"kfqt_backend/internal/system"
	"log"
	"net/http"
)

func (env *APIEnv) BookTicketHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

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

	if !system.ReadConfig().IsBookingAvailable {
		http.Error(w, `{"error": "ただいま整理券の新規発行を停止しております"}`, http.StatusBadRequest)
		return
	}

	bookingNumber, err := repository.CreateUserTicket(env.DB, req.PushToken)
	if err != nil {
		log.Printf("[ERROR] 整理券の発行失敗: %v", err)
		http.Error(w, `{"error": "Server error"}`, http.StatusInternalServerError)
		return
	}

	log.Printf("[INFO] 整理券を発行(発行者: ユーザー): 番号=%d", bookingNumber)

	// 管理者コンソールへ送出
	tickets, err := repository.GetActiveTickets(env.DB)
	if err == nil {
		var queueData []interface{}
		for _, t := range tickets {
			queueData = append(queueData, t)
		}
		select {
		case QueueUpdateChan <- queueData:
		default:
			log.Println("[WS] チャンネルが詰まっているためスキップしました")
		}
	}

	response := model.BookingResponse{
		BookingNumber: int(bookingNumber),
		Success:       true,
	}

	w.WriteHeader(http.StatusCreated) // 201 Created
	json.NewEncoder(w).Encode(response)
}
