package system

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"kfqt_backend/internal/model"
	"log"
	"os"

	"github.com/google/uuid"
)

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

	defaultConfig := &model.Config{
		PageTitle:				"Q-Tracker",
		RoomName: 				"Room",
		TimeRequired: 				5,
		TimeRequiredRangeMin: 		1,
		TimeRequiredRangeMax: 		10,
		ServeStartTime: 		"09:00",
		ServeEndTime: 			"17:00",
		Infomation:				"",
		CallInAdvanceMessage: 	"まもなくご案内いたします。アトラクション付近でお待ちください。",
		CallCurrentMessage:		"お待たせいたしました。順番が来ましたので、窓口までお越しください。",
		MessageAvailable:		"現在、すぐにご案内できます。",
		MessageNormalDelay:		"現在、ご案内までに多少お時間がかかります。",
		MessageHeavyDelay:		"現在、ご案内までに大幅にお時間がかかります。",
		IsServiceAvailable:		false,
		IsBookingAvailable:		false,
		AdminConsoleAddress:	uuidStr,
		SlotInterval:           30,
		MaxBookingsPerSlot:     5,
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
	
	adminPsw := os.Getenv("ADMIN_CONSOLE_PSW")
	if (adminPsw == "") {
		fmt.Println("[ERROR] The administrator console password has not been set.")
		os.Exit(-1)
	}

	err := os.MkdirAll("./data", 0755)
	if err != nil {
		fmt.Println("Failed to create directory: ./data")
		os.Exit(-1)
	}
	
	_, err = os.Stat("./data/config.json")

	if err != nil {
		createConfig()
	}

	log.Printf("[LOG] 管理者コンソールURL: /console/admin/%s\n", ReadConfig().AdminConsoleAddress)
}

func InitDB(database *sql.DB) {
	if err := migrate(database); err != nil {
		log.Fatal("データベースマイグレーション失敗:", err)
	}
	log.Println("データベース初期化完了（SQLite）")
}
