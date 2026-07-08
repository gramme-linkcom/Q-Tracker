package service

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func ManifestHandler(w http.ResponseWriter, r *http.Request) {
	// アクセスしてきた現在のプロトコル（http or https）とホスト名を自動取得！
	scheme := "http"
	if r.TLS != nil || r.Header.Get("X-Forwarded-Proto") == "https" {
		scheme = "https"
	}
	host := r.Host

	// 動的に start_url を組み立てる
	startURL := fmt.Sprintf("%s://%s/", scheme, host)

	// マニフェストの構造をマップで定義
	manifest := map[string]interface{}{
		"name":             "デジタル整理券システム",
		"short_name":       "Q-Tracker",
		"description":      "学校祭アトラクション順番待ちアプリ",
		"start_url":        startURL, // 🚀 ここが動的にngrokのURLになります！
		"display":          "standalone",
		"background_color": "#141416",
		"theme_color":      "#1e1e22",
	}

	// JSONとして出力
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*") // CORSブロックを完全に破壊するガード
	json.NewEncoder(w).Encode(manifest)
}
