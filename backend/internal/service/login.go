package service

import (
	"crypto/rand"
	"encoding/hex"
	"html/template"
	"kfqt_backend/internal/repository"
	"kfqt_backend/internal/system"
	"net/http"
	"os"
	"sync"
)

var (
	adminSessions = make(map[string]bool)
	sessionMu     sync.RWMutex
)

func generateSecureToken() string {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return hex.EncodeToString(b)
}

func checkAdminAuth(r *http.Request) bool {
	cookie, err := r.Cookie("admin_session")
	if err != nil {
		return false
	}

	sessionMu.RLock()
	defer sessionMu.RUnlock()
	return adminSessions[cookie.Value]
}

func AdminConsoleHandler(w http.ResponseWriter, r *http.Request) {
	config := system.ReadConfig()
	adminAddress := r.PathValue("admin_console_address")

	expectedAddress := config.AdminConsoleAddress
	if adminAddress != expectedAddress {
		http.Error(w, "Unauthorized URL Key", http.StatusUnauthorized)
		return
	}

	if r.Method == http.MethodPost && !checkAdminAuth(r) {
		r.ParseForm()
		password := r.FormValue("password")

		if password == os.Getenv("ADMIN_CONSOLE_PSW") {
			sessionID := generateSecureToken()

			sessionMu.Lock()
			if len(adminSessions) > 0 {
				sessionMu.Unlock()
				http.Error(w, "Other administrator is already logged in. Access Denied.", http.StatusForbidden)
				return
			}

			adminSessions[sessionID] = true
			sessionMu.Unlock()

			cookie := &http.Cookie{
				Name:     "admin_session",
				Value:    sessionID,
				Path:     "/",
				HttpOnly: true,
				Secure:   true,
				SameSite: http.SameSiteLaxMode,
				MaxAge:   3600 * 12,
			}
			http.SetCookie(w, cookie)

			http.Redirect(w, r, r.URL.Path, http.StatusSeeOther)
			return
		}
		http.Error(w, "Invalid Password", http.StatusForbidden)
		return
	}

	if !checkAdminAuth(r) {
		tmpl, err := template.ParseFiles("templates/login.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, nil)
		return
	}

	tmpl, err := template.ParseFiles("templates/admin.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	type CombinedAdminData struct {
		NextNumber    int
		CurrentNumber int
		WaitingGroups int
		Tickets       []repository.Ticket
		Config        map[string]interface{}
	}

	cfg := system.ReadConfig()

	configData := map[string]interface{}{
		"PageTitle":            cfg.PageTitle,
		"RoomName":             cfg.RoomName,
		"TimeRequired":         cfg.TimeRequired,
		"TimeRequiredRangeMin": cfg.TimeRequiredRangeMin,
		"TimeRequiredRangeMax": cfg.TimeRequiredRangeMax,
		"ServeStartTime":       cfg.ServeStartTime,
		"ServeEndTime":         cfg.ServeEndTime,
		"Infomation":           cfg.Infomation,
		"IsBookingAvailable":   cfg.IsBookingAvailable,
		"IsServiceAvailable":	cfg.IsServiceAvailable,
		"AdminConsoleAddress":  adminAddress,
	}

	data := CombinedAdminData{
		NextNumber:    1,
		CurrentNumber: 0,
		WaitingGroups: 0,
		Tickets:       []repository.Ticket{},
		Config:        configData,
	}

	tmpl.Execute(w, data)
}
