package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type TelegramVideo struct {
    FileID string `json:"file_id"`
    Length int `json:"length"`
}

type TelegramEndpointData struct {
    Videos []TelegramVideo `json:"videos"`
}

func FetchTelegramEndpoint(endpoint string) (*TelegramEndpointData, error) {
    log.Printf("Fetching video list from endpoint...")
    resp, err := http.Get(fmt.Sprintf("%s/list", endpoint))
    if err != nil {
        return nil, err
    }
    var data TelegramEndpointData
    err = json.NewDecoder(resp.Body).Decode(&data)
    if err != nil {
        return nil, err
    }
    return &data, nil
}

type ResultData struct {
    FilePath string `json:"file_path"`
}
type TelegramFileResult struct {
    Result ResultData `json:"result"`
}

func FetchVideoFromTelegram(fileId string) (*Video, error) {
    log.Printf("Downloading telegram video '%s'", fileId)
    resp, err := http.Get(fmt.Sprintf("https://api.telegram.org/bot%s/getFile?file_id=%s", TELEGRAM_BOT, fileId))
    if err != nil {
        return nil, err
    }
    var tgRes interface{}
    err = json.NewDecoder(resp.Body).Decode(&tgRes)
    if err != nil {
        return nil, err
    }
    telegramUri := tgRes.(map[string]interface{})["result"].(map[string]interface{})["file_path"].(string)
    videoDownloadUrl := fmt.Sprintf("https://api.telegram.org/file/bot%s/%s", TELEGRAM_BOT, telegramUri)
    video, err := Download(videoDownloadUrl)
    if err != nil {
        return nil, err
    }
    return NewVideoFromAnotherVideo(video)
}
