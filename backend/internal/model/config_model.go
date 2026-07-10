package model

type Config struct {
	PageTitle			string	`json:"page_title"`				// ヘッダー名前
	RoomName			string 	`json:"room_name"`
	
	TimeRequired		int		`json:"time_required"`			// 分単位指定
	// 時間設定の人為的ミス防止
	TimeRequiredRangeMin int 	`json:"time_required_range_min"`	// 分単位指定
	TimeRequiredRangeMax int 	`json:"time_required_range_max"`	// 分単位指定

	ServeStartTime		string	`json:"serve_start_time"`		// HH:MM 形式
	ServeEndTime		string	`json:"serve_end_time"`			// HH:MM 形式

	Infomation			string	`json:"infomation"`				// アトラクションからのお知らせ
	CallInAdvanceMessage string	`json:"call_in_advance_message"` // 案内直前のユーザーに送られるメッセージ
	CallCurrentMessage	string `json:"call_current_message"`		// 案内対象者へ送られるメッセージ
	// 状況に応じたお知らせメッセージ
    MessageAvailable    string `json:"message_available"`		// すぐに案内できるとき
    MessageNormalDelay  string `json:"message_normal_delay"`	// 少し時間がかかるとき
    MessageHeavyDelay   string `json:"message_heavy_delay"`		// 大幅に時間がかかるとき

	IsServiceAvailable	bool 	`json:"is_service_available"`		// システムが利用できるかどうか

	//利用者がキューを入れられるかどうか
	IsBookingAvailable	bool	`json:"is_booking_available"`	// 予約を入れられるかどうか

	// adminコンソールの入口ランダム化
	AdminConsoleAddress	string	`json:"admin_console_address"`

	SlotInterval        int     `json:"slot_interval"`          // 予約枠の間隔 (分)
	MaxBookingsPerSlot  int     `json:"max_bookings_per_slot"`  // 1枠あたりの最大予約数
}
