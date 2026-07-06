package main

import (
	"kfqt_backend/internal"
	"kfqt_backend/internal/db"
	"kfqt_backend/internal/middleware"
	"kfqt_backend/internal/service"
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

	env := &service.APIEnv{DB: database}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /", middleware.SameSiteOnlyMiddleware(internal.GetIndexHandlerfunc))
	mux.HandleFunc("GET /api/data", middleware.SameSiteOnlyMiddleware(env.GetStatusHandler))
	mux.HandleFunc("POST /api/booking", middleware.SameSiteOnlyMiddleware(env.BookTicketHandler))
	mux.HandleFunc("POST /api/booking/cancel", middleware.SameSiteOnlyMiddleware(env.CancelBookingHandler))
	mux.HandleFunc("GET /console/admin/{admin_console_address}", func(w http.ResponseWriter, r *http.Request) {
		middleware.SameSiteOnlyMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			service.AdminConsoleHandler(w, r)
		})).ServeHTTP(w, r)
	})
	mux.HandleFunc("POST /console/admin/{admin_console_address}", func(w http.ResponseWriter, r *http.Request) {
		middleware.SameSiteOnlyMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			service.AdminConsoleHandler(w, r)
		})).ServeHTTP(w, r)
	})
	loggedMux := middleware.LoggerMiddleware(mux)
	mainMux := http.NewServeMux()
	mainMux.Handle("/", loggedMux)
	mainMux.HandleFunc("GET /console/ws", func(w http.ResponseWriter, r *http.Request) {
		service.WebSocketHandler(w, r)
	})

	log.Println("サーバー起動: http://localhost:8080")
	if err := http.ListenAndServe(":8080", mainMux); err != nil {
		log.Fatal(err)
	}
}
