package repository

import (
	"database/sql"
)

type Ticket struct {
	Number int    `json:"number"`
	Name   string `json:"name"`
	Status string `json:"status"` // "waiting", "called", "canceled", "absent"
}

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

// GetActiveTickets は待機中("waiting")のチケットの一覧を番号順にそのまま取得する
func GetActiveTickets(db *sql.DB) ([]Ticket, error) {
	rows, err := db.Query("SELECT number, name, status FROM tickets WHERE status = 'waiting' ORDER BY number ASC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tickets []Ticket
	for rows.Next() {
		var t Ticket
		if err := rows.Scan(&t.Number, &t.Name, &t.Status); err != nil {
			return nil, err
		}
		tickets = append(tickets, t)
	}

	if tickets == nil {
		tickets = []Ticket{}
	}
	return tickets, nil
}

// func GetWaitStatusByDB(database *sql.DB) {
	
// }
