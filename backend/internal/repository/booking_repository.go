package repository

import (
	"database/sql"
	"kfqt_backend/internal/model"

	"github.com/google/uuid"
)

func CreateUserTicket(db *sql.DB, pushToken string) (model.ResultTicket, error) {
	resultData := model.ResultTicket{
		Uuid: uuid.NewString(),
		TicketNumber: 0,
	}

	query := "INSERT INTO tickets (status, uuid, device_id) VALUES ('waiting', ?, ?)"
	
	result, err := db.Exec(query, resultData.Uuid, pushToken)
	if err != nil {
		return resultData, err
	}

	lastID, err := result.LastInsertId()
	if err != nil {
		return resultData, err
	}

	resultData.TicketNumber = int(lastID)

	return resultData, nil
}

// CancelUserTicket はユーザーが自分のスマホから整理券をキャンセルした時にステータスを書き換える
func CancelUserTicket(db *sql.DB, bookingNumber int) error {
	query := "UPDATE tickets SET status = 'canceled' WHERE number = ? AND status = 'waiting'"
	
	_, err := db.Exec(query, bookingNumber)
	return err
}

// AbsentUserTicket はユーザーが自分のスマホから整理券をキャンセルした時にステータスを書き換える
func AbsentUserTicket(db *sql.DB, bookingNumber int) error {
	query := "UPDATE tickets SET status = 'absent' WHERE number = ? AND status = 'waiting'"
	
	_, err := db.Exec(query, bookingNumber)
	return err
}

func MarkGroupAsServing(db *sql.DB, bookingNumber int) error {
	query := "UPDATE tickets SET status = 'serving' WHERE number = ? AND status = 'waiting'"
	
	_, err := db.Exec(query, bookingNumber)
	return err
}

func FinishServingGroup(db *sql.DB) error {
	query := "UPDATE tickets SET status = 'done' WHERE status = 'serving'"
	
	_, err := db.Exec(query)
	return err
}
