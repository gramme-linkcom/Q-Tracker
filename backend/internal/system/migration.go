package system

import (
	"database/sql"
	"fmt"
)

func migrate(db *sql.DB) error {
	// 1. 整理券管理テーブルの作成
	ticketsTable := `
	CREATE TABLE IF NOT EXISTS tickets (
		number INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT,
		status TEXT NOT NULL DEFAULT 'waiting',
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`
	if _, err := db.Exec(ticketsTable); err != nil {
		return fmt.Errorf("failed to create tickets table: %w", err)
	}

	// 2. ルーム状況管理テーブルの作成
	roomStatusTable := `
	CREATE TABLE IF NOT EXISTS room_status (
		id INTEGER PRIMARY KEY CHECK (id = 1), -- 常に1行しか存在させない制約
		current_number INTEGER DEFAULT 0,
		is_active INTEGER DEFAULT 0 -- 0: false, 1: true
	);`
	if _, err := db.Exec(roomStatusTable); err != nil {
		return fmt.Errorf("failed to create room_status table: %w", err)
	}

	// 3. 初期レコードの投入（まだ1行もなければ、初期値 0番/空室 で作成）
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM room_status").Scan(&count)
	if err != nil {
		return err
	}
	if count == 0 {
		_, err = db.Exec("INSERT INTO room_status (id, current_number, is_active) VALUES (1, 0, 0)")
		if err != nil {
			return fmt.Errorf("failed to insert initial room_status: %w", err)
		}
	}

	return nil
}
