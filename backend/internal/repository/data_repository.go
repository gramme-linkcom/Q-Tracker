package repository

import (
	"database/sql"
	"errors"
	"kfqt_backend/internal/model"
)

// RoomStatus は現在の部屋全体の状況を表す構造体
type RoomStatus struct {
	CurrentNumber int  `json:"currentNumber"`
	IsActive      bool `json:"isActive"`
}

func GetRoomStatus(db *sql.DB) (RoomStatus, error) {
	var room RoomStatus
	err := db.QueryRow("SELECT current_number, is_active FROM room_status WHERE id = 1").Scan(&room.CurrentNumber, &room.IsActive)
	return room, err
}

func SetRoomStatus(db *sql.DB, roomStatus RoomStatus) error {
	query := "UPDATE room_status SET is_active = ?, current_number = ? WHERE id = 1"
	isRoomActive := 0
	if ( roomStatus.IsActive == true ) {
		isRoomActive = 1
	} 

	_, err := db.Exec(query, isRoomActive, roomStatus.CurrentNumber)
	return err
}

func GetYoungerGroups(db *sql.DB) (*model.Ticket, error) {
	query := "SELECT number, device_id, status FROM tickets WHERE status = 'waiting' ORDER BY number ASC LIMIT 1;"
	var t model.Ticket
	// 1件だけの取得なので QueryRow を使い、Scan で構造体に流し込む
	err := db.QueryRow(query).Scan(&t.Number, &t.DeviceID, &t.Status)
	
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("現在待機中のチケットはありません")
		}
		return nil, err
	}

	return &t, nil
}

func GetAheadGroups(db *sql.DB, myNumberStr string) (int) {
	myAheadGroups := 0
	query := "SELECT COUNT(*) FROM tickets WHERE number < ? AND status = 'waiting'"
	err := db.QueryRow(query, myNumberStr).Scan(&myAheadGroups)
	if err != nil {
		// エラー時は安全に0にしておく
		return 0
	}
	return myAheadGroups 
}

// GetActiveTickets は待機中("waiting")のチケットの一覧を番号順にそのまま取得する
func GetWaitingTickets(db *sql.DB) ([]model.Ticket, error) {
	rows, err := db.Query("SELECT number, device_id, status FROM tickets WHERE status = 'waiting' ORDER BY number ASC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tickets []model.Ticket
	for rows.Next() {
		var t model.Ticket
		if err := rows.Scan(&t.Number, &t.DeviceID, &t.Status); err != nil {
			return nil, err
		}
		tickets = append(tickets, t)
	}

	if tickets == nil {
		tickets = []model.Ticket{}
	}
	return tickets, nil
}

// GetActiveTickets は待機中("waiting")のチケットの一覧を番号順にそのまま取得する
func GetActiveTickets(db *sql.DB) ([]model.Ticket, error) {
	rows, err := db.Query("SELECT number, device_id, status FROM tickets WHERE status = 'waiting' OR status = 'serving' ORDER BY number ASC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tickets []model.Ticket
	for rows.Next() {
		var t model.Ticket
		if err := rows.Scan(&t.Number, &t.DeviceID, &t.Status); err != nil {
			return nil, err
		}
		tickets = append(tickets, t)
	}

	if tickets == nil {
		tickets = []model.Ticket{}
	}
	return tickets, nil
}
