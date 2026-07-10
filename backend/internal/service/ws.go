package service

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"kfqt_backend/internal/model"
	"kfqt_backend/internal/repository"
	"kfqt_backend/internal/system"

	"github.com/gorilla/websocket"
)

type BroadcastDatas struct {
	PushType	string
	Queue		[]interface{}
}

type AdminActionPacket struct {
	Action       string `json:"action"`
	Number       int    `json:"number"`
	ReservedTime string `json:"reserved_time"`
}

var ActiveAdminConn *websocket.Conn
var ConnMu sync.Mutex

func BroadcastQueue(data BroadcastDatas) {
	ConnMu.Lock()
	defer ConnMu.Unlock()

	if ActiveAdminConn == nil {
		log.Println("[LOG] 接続中の管理画面が存在しませんでした。")
		return
	}

	payload := map[string]interface{}{
		"type":  data.PushType,
		"queue": data.Queue,
	}

	// ログイン中の「唯一の1台」に直接送信！
	err := ActiveAdminConn.WriteJSON(payload)
	if err != nil {
		log.Println("[WS_WRITE_ERROR] 送信失敗、接続を切断します:", err)
		ActiveAdminConn.Close()
		ActiveAdminConn = nil
	}

	log.Println("BroadcastQueue を実行しました。")
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
func (env *APIEnv) WebSocketHandler(w http.ResponseWriter, r *http.Request) {
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
	
	ConnMu.Lock()
	ActiveAdminConn = conn
	ConnMu.Unlock()
	log.Println("[WS_CONNECT] 管理画面のWebSocketが正常に接続・登録されました。")

	defer func() {
		conn.Close()

		ConnMu.Lock()
		if ActiveAdminConn == conn { // 自分自身の接続ならnilにする
			ActiveAdminConn = nil
		}
		ConnMu.Unlock()
		log.Println("[WS_DISCONNECT] 管理画面のWebSocketが切断されました。")

		cookie, err := r.Cookie("admin_session")
		if err == nil {
			sessionID := cookie.Value
			// 10秒間の再接続猶予期間（グレースピリオド）を開始
			go func(sid string) {
				time.Sleep(10 * time.Second)
				
				ConnMu.Lock()
				defer ConnMu.Unlock()
				
				// 10秒後、ActiveAdminConnが空（再接続されていない）ならセッションを削除する
				if ActiveAdminConn == nil {
					sessionMu.Lock()
					delete(adminSessions, sid)
					log.Printf("[WS_DISCONNECT] 10秒間再接続がなかったため、セッション %s を破棄しました\n", sid[:8])
					sessionMu.Unlock()
				} else {
					log.Printf("[WS_DISCONNECT] 再接続が確認されたため、セッション %s の破棄をキャンセルしました\n", sid[:8])
				}
			}(sessionID)
		}
	}()

	env.sendInitialState(conn)

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

		var actionData AdminActionPacket
		_ = json.Unmarshal(msg, &actionData)

		if actionData.Action != "" {
			// 操作コマンドだった場合の処理
			log.Printf("[WS_ACTION] 操作を受信しました: Action=%s, Number=%d", actionData.Action, actionData.Number)
			
			switch actionData.Action {
			case "cancel_ticket":
				err := repository.CancelUserTicket(env.DB, actionData.Number)
				if err != nil {
					log.Println("[ERROR] 整理券をキャンセルできませんでした。")
					continue
				}
				tickets, err := repository.GetWaitingTickets(env.DB)
				if err == nil {
					var queueData []interface{}
					for _, t := range tickets {
						queueData = append(queueData, t)
					}
					BroadcastQueue(BroadcastDatas{PushType: "queue_update", Queue: queueData})
				}
				nextTicketData, err := repository.GetYoungerGroups(env.DB)
				if err != nil {
					log.Println("次のチケット情報の取得に失敗しました。")
					continue
				} else if (nextTicketData.DeviceID != ""){
					go SendPushToUser(nextTicketData.DeviceID, system.ReadConfig().CallCurrentMessage)
				}

			case "absent_ticket":
				err := repository.AbsentUserTicket(env.DB, actionData.Number)
				if err != nil {
					log.Println("[ERROR] 整理券を不在キャンセルできませんでした。")
					continue
				}
				tickets, err := repository.GetWaitingTickets(env.DB)
				if err == nil {
					var queueData []interface{}
					for _, t := range tickets {
						queueData = append(queueData, t)
					}
					BroadcastQueue(BroadcastDatas{PushType: "queue_update", Queue: queueData})
				}
				nextTicketData, err := repository.GetYoungerGroups(env.DB)
				if err != nil {
					log.Println("次のチケット情報の取得に失敗しました。")
					continue
				} else if (nextTicketData.DeviceID != ""){
					go SendPushToUser(nextTicketData.DeviceID, system.ReadConfig().CallCurrentMessage)
				}

			case "group_enter":
				youngTicketData, err := repository.GetYoungerGroups(env.DB)
				if err != nil {
					log.Println("チケット情報の取得に失敗しました。")
					continue
				}

				log.Printf("取得したID: %d\n", youngTicketData.Number)

				roomStatus := repository.RoomStatus {
					IsActive: true,
					CurrentNumber: youngTicketData.Number,
				}
				if res := repository.SetRoomStatus(env.DB, roomStatus); res != nil {
					log.Println("ルーム状況の更新に失敗しました")
					continue
				}
				if res := repository.MarkGroupAsServing(env.DB, youngTicketData.Number); res != nil {
					log.Println("ユーザーステータスの更新に失敗しました。")
					continue
				}
				log.Printf("入室処理を実行: %d\n", youngTicketData.Number)
				tickets, _ := repository.GetWaitingTickets(env.DB)
				if err == nil {
					var queueData []interface{}
					for _, t := range tickets {
						queueData = append(queueData, t)
					}
					BroadcastQueue(BroadcastDatas{PushType: "queue_update", Queue: queueData})
				}

				nextTicketData, err := repository.GetYoungerGroups(env.DB)
				if err != nil {
					log.Println("次のチケット情報の取得に失敗しました。")
					continue
				} else if (nextTicketData.DeviceID != ""){
					go SendPushToUser(nextTicketData.DeviceID, system.ReadConfig().CallInAdvanceMessage)
				}
			case "group_exit":
				roomStatus := repository.RoomStatus {
					IsActive: false,
					CurrentNumber: 0,
				}
				if res := repository.SetRoomStatus(env.DB, roomStatus); res != nil {
					log.Println("ルーム状況の更新に失敗しました")
					continue
				}
				if res := repository.FinishServingGroup(env.DB); res != nil {
					log.Println("ユーザーステータスの更新に失敗しました。")
					continue
				}
				log.Println("退出処理を実行")
				tickets, err := repository.GetWaitingTickets(env.DB)
				if err == nil {
					var queueData []interface{}
					for _, t := range tickets {
						queueData = append(queueData, t)
					}
					BroadcastQueue(BroadcastDatas{PushType: "queue_update", Queue: queueData})
				}

				nextTicketData, err := repository.GetYoungerGroups(env.DB)
				if err != nil {
					log.Println("次のチケット情報の取得に失敗しました。")
					continue
				} else if (nextTicketData.DeviceID != ""){
					go SendPushToUser(nextTicketData.DeviceID, system.ReadConfig().CallCurrentMessage)
				}
			
			case "clear_all":
				room, err := repository.GetRoomStatus(env.DB)
				if err == nil && room.IsActive {
					log.Println("[WS_WARN] 現在グループが入場中のため、リセットを拒否しました。")
					
					payload := map[string]interface{}{
						"type":    "operation_denied",
						"message": "現在アクティブな入場グループが存在するため、整理券番号のリセットは拒否されました。退場処理を先に行うか、空室時に実行してください。",
					}
					_ = conn.WriteJSON(payload)
					continue
				}

				cfg := system.ReadConfig()
				if (!cfg.IsBookingAvailable || !IsWithinServeTime(cfg.ServeStartTime, cfg.ServeEndTime) || !cfg.IsServiceAvailable) {
					remainTickets, err := repository.GetWaitingTickets(env.DB)
					if (len(remainTickets) != 0) {
						log.Println("[WS_WARN] 待機列が存在するため、リセットを拒否しました。")
						
						payload := map[string]interface{}{
							"type":    "operation_denied",
							"message": "待機列が存在するため、整理券番号のリセットは拒否されました。",
						}
						_ = conn.WriteJSON(payload)
						continue
					}

					_, err = env.DB.Exec("DELETE FROM tickets")
					if err != nil {
						log.Printf("[DB_ERROR] チケットデータの全削除に失敗しました: %v\n", err)
						payload := map[string]interface{}{
							"type":    "operation_denied",
							"message": "データベースエラーによりリセットに失敗しました。",
						}
						_ = conn.WriteJSON(payload)
						continue
					}

					_, err = env.DB.Exec("DELETE FROM sqlite_sequence WHERE name = 'tickets'")
					if err != nil {
						log.Printf("[DB_ERROR] 自動採番カウンターのリセットに失敗しました: %v\n", err)
					} else {
						log.Println("[DB_LOG] 整理券の採番カウンターを1番にリセットしました。")
					}

					// 4. ルーム状況（room_status）も初期状態（0番、空室）にリセット
					resetRoom := repository.RoomStatus{
						CurrentNumber: 0,
						IsActive:      false,
					}
					if err := repository.SetRoomStatus(env.DB, resetRoom); err != nil {
						log.Printf("[DB_ERROR] ルーム状況のリセットに失敗しました: %v\n", err)
					}

					log.Println("[WS_ACTION] 整理券番号のリセット処理（全データ削除＆採番初期化）が成功しました。")

					// 5. 空っぽになった最新の待機列を管理画面（および全ユーザー画面）へリアルタイム同期配信！
					tickets, err := repository.GetWaitingTickets(env.DB)
					if err == nil {
						var queueData []interface{}
						for _, t := range tickets {
							queueData = append(queueData, t)
						}
						BroadcastQueue(BroadcastDatas{PushType: "queue_update", Queue: queueData})
					}
				} else {
					log.Println("[WS_WARN] 整理券受付中のため、リセットを拒否しました。")
					
					payload := map[string]interface{}{
						"type":    "operation_denied",
						"message": "整理券受付中のため、整理券番号のリセットは拒否されました。リセットしたい場合は、整理券受付を停止してください。",
					}
					_ = conn.WriteJSON(payload)
					continue
				}

			case "reflesh":
				tickets, err := repository.GetWaitingTickets(env.DB)
				if err == nil {
					var queueData []interface{}
					for _, t := range tickets {
						queueData = append(queueData, t)
					}
					BroadcastQueue(BroadcastDatas{PushType: "queue_update", Queue: queueData})
				}
			
			case "issue_manual_ticket":
				cfg := system.ReadConfig()
				if (!cfg.IsBookingAvailable || !IsWithinServeTime(cfg.ServeStartTime, cfg.ServeEndTime) || !cfg.IsServiceAvailable) {
					log.Println("発券停止中のため、発券できませんでした。")
					payload := map[string]interface{}{
						"type":    "operation_denied",
						"message": "発券停止中のため、発券できませんでした。",
					}
					_ = conn.WriteJSON(payload)
					continue
				}
				reservedTime := actionData.ReservedTime
				if reservedTime == "" {
					reservedTime = GetCurrentTimeSlot()
				}
				bookingData, err := repository.CreateUserTicket(env.DB, "", reservedTime)
				if err != nil {
					log.Printf("[ERROR] 整理券の発行失敗: %v", err)
					continue
				}

				log.Printf("[INFO] 整理券を発行(発行者: 管理者): 番号=%d (指定枠: %s)", bookingData.TicketNumber, reservedTime)
				tickets, err := repository.GetWaitingTickets(env.DB)
				if err == nil {
					var queueData []interface{}
					for _, t := range tickets {
						queueData = append(queueData, t)
					}
					BroadcastQueue(BroadcastDatas{PushType: "queue_update", Queue: queueData})
				}
			}
		} else {
			// アクションが含まれていない場合は、従来通りの「設定データ」として処理
			var newConfigData model.Config
			err = json.Unmarshal(msg, &newConfigData)
			if err != nil {
				log.Printf("[ERROR] JSONの変換に失敗しました: %v", err)
				continue
			}
			system.SaveConfig(newConfigData)
			log.Println("[WS_CONFIG] 構成設定を更新しました")
		}
	}
}

func (env *APIEnv) sendInitialState(conn *websocket.Conn) {
	tickets, err := repository.GetActiveTickets(env.DB)
	if err != nil {
		log.Fatalln("DBからのデータ取得に失敗しました。DBが破損している可能性があります。")
		return
	}
	var ticketsData []interface{}
	for _, t := range tickets {
		ticketsData = append(ticketsData, t)
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

	cfg := system.ReadConfig()

	// admin.html の updateDOM がそのままパースできる器
	initialData := map[string]interface{}{
		"nextNumber":    nextNumber,
		"currentNumber": currentNumber,
		"waitingGroups": waitingGroups,
		"tickets":       ticketsData, // 待ち列一覧（空）
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
			"slot_interval":           cfg.SlotInterval,
			"max_bookings_per_slot":   cfg.MaxBookingsPerSlot,
			"allow_no_time_slot":      cfg.AllowNoTimeSlot,
		},
	}

	_ = conn.WriteJSON(initialData)
}
