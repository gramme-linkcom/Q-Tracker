package internal

import (
	"fmt"
	"kfqt_backend/internal/system"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func GetIndexHandlerfunc(w http.ResponseWriter, r *http.Request) {
	config := system.ReadConfig()

	publicDir := "./public"
	if r.URL.Path == "/" {
		tmplPath := filepath.Join(publicDir, "index.html")
		htmlBytes, err := os.ReadFile(tmplPath)
		if err != nil {
			log.Println("[ERROR] index.htmlの読み込みに失敗しました:", err)
			http.Error(w, "Read Error", http.StatusInternalServerError)
			return
		}

		htmlStr := string(htmlBytes)

		injectScript := fmt.Sprintf(`<script>
			window.__SERVER_CONFIG__ = {
				pageTitle: "%s",
				roomName: "%s"
			};
		</script>`, config.PageTitle, config.RoomName)

		htmlStr = strings.ReplaceAll(htmlStr, "<head>", "<head>"+injectScript)

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(htmlStr))
		return
	}
	http.FileServer(http.Dir(publicDir)).ServeHTTP(w, r)
}
