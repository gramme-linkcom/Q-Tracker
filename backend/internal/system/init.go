package system

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/google/uuid"
)

type Config struct {
	PageTitle			string	`json:"page_title"`				// ヘッダー名前
	RoomName			string 	`json:"room_name"`
	
	WaitTime			int		`json:"wait_time"`				// 秒単位指定
	// 時間設定の人為的ミス防止
	WaitTimeRangeMin	int 	`json:"wait_time_range_min"`	// 秒単位指定
	WaitTimeRangeMax	int 	`json:"wait_time_range_max"`	// 秒単位指定

	ServeStartTime		string	`json:"serve_start_time"`		// HH:MM 形式
	ServeEndTime		string	`json:"serve_end_time"`			// HH:MM 形式

	Infomation			string	`json:"infomation"`				// アトラクションからのお知らせ
	CallInAdvanceMessage string	`json:"call_in_advice_message"` // 案内直前のユーザーに送られるメッセージ
	// 状況に応じたお知らせメッセージ
    MessageAvailable    string `json:"message_available"`		// すぐに案内できるとき
    MessageNormalDelay  string `json:"message_normal_delay"`	// 少し時間がかかるとき
    MessageHeavyDelay   string `json:"message_heavy_delay"`		// 大幅に時間がかかるとき

	IsEventAvailable	bool 	`json:"is_event_available"`		// システムが利用できるかどうか

	// Push通知
	PushNotification	bool	`json:"push_notification"`		// プッシュ通知を利用するかどうか

	//利用者がキューを入れられるかどうか
	IsBookingAvailable	bool	`json:"is_booking_available"`	// 予約を入れられるかどうか

	// adminコンソールの入口ランダム化
	AdminConsoleAddress	string	`json:"admin_console_address"`
}

func createConfig() {
	file, err := os.Create("./data/config.json")
	if err != nil {
		fmt.Println("ファイルの作成に失敗しました:", err)
		return
	}
	defer file.Close()

	newUUID, err := uuid.NewRandom()
	if err != nil {
		fmt.Println("UUIDの生成に失敗しました:", err)
		return
	}
	uuidStr := newUUID.String()

	defaultConfig := &Config{
		PageTitle:				"Q-Tracker",
		RoomName: 				"Room",
		WaitTime: 				3000,
		WaitTimeRangeMin: 		60,
		WaitTimeRangeMax: 		6000,
		ServeStartTime: 		"17:00",
		ServeEndTime: 			"09:00",
		Infomation:				"",
		CallInAdvanceMessage: 	"まもなくご案内いたします。<br />アトラクションの手前までお進みください。",
		MessageAvailable:		"現在、すぐにご案内できます。",
		MessageNormalDelay:		"現在、ご案内までに多少お時間がかかります。",
		MessageHeavyDelay:		"現在、ご案内までに大幅にお時間がかかります。",
		IsEventAvailable:		false,
		PushNotification:		true,
		IsBookingAvailable:		false,
		AdminConsoleAddress:	uuidStr,
	}

	jsonData, err := json.MarshalIndent(defaultConfig, "", "  ")
	if err != nil {
		return
	}

	_, err = file.WriteString(string(jsonData))
	if err != nil {
		fmt.Println("Failed to write:", err)
		return
	}

	fmt.Println("Created config file.")
}

func Init() {
	fmt.Println("Initialization in progress...")
	err := os.MkdirAll("./data", 0755)
	if err != nil {
		fmt.Println("Failed to create directory: ./data")
		os.Exit(-1)
	}
	
	_, err = os.Stat("./data/config.json")

	if err != nil {
		createConfig()
	}
}

func InitDB(database *sql.DB) {
	if err := migrate(database); err != nil {
		log.Fatal("データベースマイグレーション失敗:", err)
	}
	log.Println("データベース初期化完了（SQLite）")
}
