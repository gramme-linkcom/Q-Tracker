// console/login.go
package console

import (
	"crypto/rand"
	"encoding/hex"
	"html/template"
	"kfqt_backend/internal/api"
	"kfqt_backend/internal/repository"
	"kfqt_backend/internal/system"
	"net/http"
	"os"
)

func generateSecureToken() string {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return hex.EncodeToString(b)
}

func checkAdminAuth(env *api.APIEnv, r *http.Request) bool {
	cookie, err := r.Cookie("admin_session")
	if err != nil {
		return false
	}

	env.SessionMu.RLock()
	defer env.SessionMu.RUnlock()
	return env.AdminSessions != nil && env.AdminSessions[cookie.Value]
}

func AdminConsoleHandler(env *api.APIEnv, w http.ResponseWriter, r *http.Request) {
	config := system.ReadConfig()
	adminAddress := r.PathValue("admin_console_address")

	expectedAddress := config.AdminConsoleAddress
	if adminAddress != expectedAddress {
		http.Error(w, "Unauthorized URL Key", http.StatusUnauthorized)
		return
	}

	if r.Method == http.MethodPost && !checkAdminAuth(env, r) {
		r.ParseForm()
		password := r.FormValue("password")

		if password == os.Getenv("ADMIN_CONSOLE_PSW") {
			sessionID := generateSecureToken()

            env.SessionMu.Lock()
            if env.AdminSessions == nil {
                env.AdminSessions = make(map[string]bool)
            }

            if len(env.AdminSessions) > 0 {
                env.SessionMu.Unlock()
                http.Error(w, "Other administrator is already logged in. Access Denied.", http.StatusForbidden)
                return
            }

            env.AdminSessions[sessionID] = true
            env.SessionMu.Unlock()

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

	if !checkAdminAuth(env, r) {
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
