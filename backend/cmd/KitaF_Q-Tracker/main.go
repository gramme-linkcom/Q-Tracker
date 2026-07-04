package main

import (
	"kfqt_backend/internal"
	"kfqt_backend/internal/api"
	"kfqt_backend/internal/db"
	"kfqt_backend/internal/middleware"
	"kfqt_backend/internal/system"

	"log"
	"net/http"
)

type TemplateData struct {
	PageTitle string
	RoomName string
}

func main() {
	system.Init()

	database, err := db.Open()
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	system.InitDB(database)

	env := &api.APIEnv{DB: database}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /", middleware.SameSiteOnlyMiddleware(internal.GetIndexHandlerfunc))
	mux.HandleFunc("GET /api/data", middleware.SameSiteOnlyMiddleware(env.GetStatusHandler))
	// mux.HandleFunc("GET /api/")

	loggedMux := middleware.LoggerMiddleware(mux)

	log.Println("サーバー起動: http://localhost:8080")
	if err := http.ListenAndServe(":8080", loggedMux); err != nil {
		log.Fatal(err)
	}
}
