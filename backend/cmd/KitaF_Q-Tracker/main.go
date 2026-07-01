package main

import (
	"log"
	"net/http"
)

func main() {
	// "public" フォルダの中身（後でNext.jsの画面をここに入れます）を配信する
	fs := http.FileServer(http.Dir("./public"))
	http.Handle("/", fs)

	log.Println("サーバー起動: http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
