package service

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"kfqt_backend/internal/model"
	"kfqt_backend/internal/system"

	"github.com/gorilla/websocket"
)

var ActiveAdminConn *websocket.Conn
var ConnMu sync.Mutex // 同時書き込みによるクラッシュを防ぐ安全弁
var QueueUpdateChan = make(chan []interface{}, 10)

// QueueUpdateChannnelを監視して、データが来たらWSにブロードキャストするループ関数
func StartQueueBroadcastLoop() {
	for queueData := range QueueUpdateChan {
		broadcastQueue(queueData)
	}
}

func broadcastQueue(queueData []interface{}) {
	ConnMu.Lock()
	defer ConnMu.Unlock()

	if ActiveAdminConn == nil {
		return
	}

	payload := map[string]interface{}{
		"type":  "queue_update",
		"queue": queueData,
	}

	// ログイン中の「唯一の1台」に直接送信！
	err := ActiveAdminConn.WriteJSON(payload)
	if err != nil {
		log.Println("[WS_WRITE_ERROR] 送信失敗、接続を切断します:", err)
		ActiveAdminConn.Close()
		ActiveAdminConn = nil
	}
}

type Response struct {
	Action    string `json:"action"`
	RequestID string `json:"request_id"` // どのリクエストへの返子か特定するため
	Status    string `json:"status"`     // "success" や "error"
	Message   string `json:"message"`
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// WebSocketHandler は最小限の接続維持と初期データ送信を行います
func WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	if !checkAdminAuth(r) {
		http.Error(w, "Unauthorized Session", http.StatusUnauthorized)
		log.Println("[WS_AUTH_ERROR] 無効なセッションからのWebSocket接続を拒否しました")
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket connection failed:", err)
		return
	}
	defer func() {
		conn.Close()

		cookie, err := r.Cookie("admin_session")
		if err == nil {
			sessionMu.Lock()
			delete(adminSessions, cookie.Value)
			log.Printf("[WS_DISCONNECT] セッション %s を名簿から削除し、ロックを解放しました\n", cookie.Value[:8])
			sessionMu.Unlock()
		}
	}()

	sendInitialState(conn)

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			break
		}

		if !checkAdminAuth(r) {
			log.Println("[WS_AUTH_ERROR] 操作中にセッションが無効化されたため、パケットを破棄しました")
			_ = conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "Session Expired"))
			break
		}

		var newConfigData model.Config
		err = json.Unmarshal(msg, &newConfigData)
		if err != nil {
			log.Printf("[ERROR] JSONの変換に失敗しました: %v", err)
			return
		}
		system.SaveConfig(newConfigData)
		
	}
}

func sendInitialState(conn *websocket.Conn) {
	cfg := system.ReadConfig()

	// admin.html の updateDOM がそのままパースできる器
	initialData := map[string]interface{}{
		"nextNumber":    1,
		"currentNumber": 0,
		"waitingGroups": 0,
		"tickets":       []interface{}{}, // 待ち列一覧（空）
		"config": map[string]interface{}{
			"page_title":              cfg.PageTitle,
			"room_name":               cfg.RoomName,
			"time_required":           cfg.TimeRequired,
			"time_required_range_min": cfg.TimeRequiredRangeMin,
			"time_required_range_max": cfg.TimeRequiredRangeMax,
			"serve_start_time":        cfg.ServeStartTime,
			"serve_end_time":          cfg.ServeEndTime,
			"infomation":              cfg.Infomation,
			"is_booking_available":   cfg.IsBookingAvailable,
			"admin_console_address":  cfg.AdminConsoleAddress,
		},
	}

	_ = conn.WriteJSON(initialData)
}
