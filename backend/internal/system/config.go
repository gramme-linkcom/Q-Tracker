package system

import (
	"encoding/json"
	"kfqt_backend/internal/model"
	"log"
	"os"
)

func ReadConfig() (*model.Config) {
	filePath := "./data/config.json"
	fileBytes, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("[ERROR] ファイルの読み込みに失敗しました: %v", err)
		return nil
	}

	var config model.Config
	err = json.Unmarshal(fileBytes, &config)
	if err != nil {
		log.Fatalf("[ERROR] JSONのパースに失敗しました: %v", err)
		return nil
	}
	return &config
}

func SaveConfig(newConfigData model.Config) (error){
	filePath := "./data/config.json"
	fileBytes, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("[ERROR] ファイルの読み込みに失敗しました: %v", err)
		return err
	}
	var config model.Config
	err = json.Unmarshal(fileBytes, &config)
	if err != nil {
		log.Fatalf("[ERROR] JSONの解析に失敗しました: %v", err)
		return err
	}

	config.PageTitle 				= newConfigData.PageTitle
	config.RoomName  				= newConfigData.RoomName
	config.TimeRequired 			= newConfigData.TimeRequired
	config.TimeRequiredRangeMin 	= newConfigData.TimeRequiredRangeMin
	config.TimeRequiredRangeMax 	= newConfigData.TimeRequiredRangeMax
	config.ServeStartTime 			= newConfigData.ServeStartTime
	config.ServeEndTime   			= newConfigData.ServeEndTime
	config.Infomation     			= newConfigData.Infomation
	config.IsBookingAvailable 		= newConfigData.IsBookingAvailable
	config.IsServiceAvailable		= newConfigData.IsServiceAvailable
	config.CallCurrentMessage		= newConfigData.CallCurrentMessage
	config.CallInAdvanceMessage		= newConfigData.CallInAdvanceMessage
	config.SlotInterval             = newConfigData.SlotInterval
	config.MaxBookingsPerSlot       = newConfigData.MaxBookingsPerSlot


	updatedBytes, err := json.MarshalIndent(config, "", "    ")
	if err != nil {
		log.Fatalf("[ERROR] JSONへの変換に失敗しました: %v", err)
		return err
	}

	err = os.WriteFile(filePath, updatedBytes, 0644)
	if err != nil {
		log.Fatalf("[ERROR] ファイルの保存に失敗しました: %v", err)
		return err
	}

	log.Println("[LOG] Config が書き換えられました。")
	return nil
}
