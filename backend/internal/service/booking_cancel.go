package service

import (
	"encoding/json"
	"kfqt_backend/internal/model"
	"kfqt_backend/internal/repository"
	"log"
	"net/http"
)

func (env *APIEnv) CancelBookingHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		http.Error(w, `{"error": "Method Not Allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	var req model.CancelRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "不正なリクエストデータです"}`, http.StatusBadRequest)
		return
	}

	// 💡 リポジトリの関数を呼び出して、DBのステータスを 'canceled' に変更
	err := repository.CancelUserTicket(env.DB, req.BookingNumber)
	if err != nil {
		log.Printf("[ERROR] 整理券のキャンセル処理に失敗しました: %v", err)
		http.Error(w, `{"error": "サーバーエラーによりキャンセルに失敗しました"}`, http.StatusInternalServerError)
		return
	}

	log.Printf("[INFO] ユーザーが整理券をキャンセルしました: 番号=%d", req.BookingNumber)

	// 成功レスポンスを返却
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"success": true}`))
}
