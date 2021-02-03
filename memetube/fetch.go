package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
    "strings"
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

var ProcessedFileIDs = []string{}

func MarkTelegramVideosAsProcessedCleanupHook() {
    AddCleanupHook(func() {
        log.Printf("cleanup: mark processed videos as processed")
        payload := struct {
            Files []string `json:"files"`
        }{
            Files: ProcessedFileIDs,
        }
        buf := bytes.NewBufferString("")
        json.NewEncoder(buf).Encode(payload)
        u, err := url.Parse(fmt.Sprintf("%s/deleteIds", FETCH_ENDPOINT))
        if err != nil {
            log.Printf("mark_processed: %s", err)
            return
        }
        req := http.Request{
            Method: http.MethodGet,
            URL: u,
            Body: NewReadCloserWrapper(buf),
        }
        _, err = http.DefaultClient.Do(&req)
        if err != nil {
            Report("Não foi possível marcar os vídeos processados como processados, intervenção manual requerida: %s\nVídeos processados:\n%s", err, strings.Join(ProcessedFileIDs, "\n"))
            log.Printf("mark_processed: %s", err)
        }
    })
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
    videoFile, err := NewVideoFromAnotherVideo(video)
    if err != nil {
        return nil, err
    }
    ProcessedFileIDs = append(ProcessedFileIDs, fileId)
    return videoFile, nil
}
