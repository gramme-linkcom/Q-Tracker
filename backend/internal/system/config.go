package system

import (
	"encoding/json"
	"log"
	"os"
)

func ReadConfig() (*Config) {
	filePath := "./data/config.json"
	fileBytes, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("[ERROR] ファイルの読み込みに失敗しました: %v", err)
		return nil
	}

	var config Config
	err = json.Unmarshal(fileBytes, &config)
	if err != nil {
		log.Fatalf("[ERROR] JSONのパースに失敗しました: %v", err)
		return nil
	}
	return &config
}
