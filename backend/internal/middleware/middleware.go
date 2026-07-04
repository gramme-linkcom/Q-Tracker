package middleware

import (
	"fmt"
	"net/http"
	"strings"
)

func SameSiteOnlyMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		currentHost := r.Header.Get("X-Forwarded-Host")
		if currentHost == "" {
			currentHost = r.Host
		}

		scheme := "http://"
		if r.Header.Get("X-Forwarded-Proto") == "https" {
			scheme = "https://"
		}

		myOrigin := fmt.Sprintf("%s%s", scheme, currentHost)

		w.Header().Set("Access-Control-Allow-Origin", myOrigin)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-Requested-With")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		origin := r.Header.Get("Origin")
		referer := r.Header.Get("Referer")

		isAllowed := false

		if origin != "" && origin == myOrigin {
			isAllowed = true
		} else if referer != "" && strings.HasPrefix(referer, myOrigin) {
			isAllowed = true
		} else if origin == "" && referer == "" {
			isAllowed = true
		}

		if !isAllowed {
			http.Error(w, "Access Denied: 不正なオリジンからのアクセスです", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	}
}
