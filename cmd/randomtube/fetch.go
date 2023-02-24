package main

import (
	"bytes"
	"encoding/json"
	"fmt"
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
    QueueSize int
    QueueTime int
    ReportChatID int64 `json:"reportChatID"`
}

func FetchTelegramEndpoint() (*TelegramEndpointData, error) {
    Log("Fetching video list from endpoint...")
    resp, err := http.Get(fmt.Sprintf("%s/list", FETCH_ENDPOINT))
    if err != nil {
        return nil, err
    }
    var data TelegramEndpointData
    err = json.NewDecoder(resp.Body).Decode(&data)
    if err != nil {
        return nil, err
    }
    for _, video := range data.Videos {
        data.QueueSize++
        data.QueueTime += video.Length
    }
    reportChat = data.ReportChatID
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
        if (dontMarkVideosAsProcessed) {
            Log("warning: skipping mark processed videos as processed")
            return
        }
        Log("cleanup: mark processed videos as processed")
        payload := struct {
            Files []string `json:"files"`
        }{
            Files: ProcessedFileIDs,
        }
        buf := bytes.NewBufferString("")
        json.NewEncoder(buf).Encode(payload)
        u, err := url.Parse(fmt.Sprintf("%s/deleteIds", FETCH_ENDPOINT))
        if err != nil {
            Log("mark_processed: %s", err)
            return
        }
        req := http.Request{
            Method: http.MethodGet,
            URL: u,
            Body: NewReadCloserWrapper(buf),
        }
        res, err := http.DefaultClient.Do(&req)
        if err != nil || res.StatusCode != 200 {
            Report("Não foi possível marcar os vídeos processados como processados: err=%s statusCode=%d.\nClique no link para resolver: %s/deleteIds?ids=%s", err, res.StatusCode, FETCH_ENDPOINT, strings.Join(ProcessedFileIDs, ","))
            Log("mark_processed: %s. status_code: %d", err, res.StatusCode)
        }
    })
}

func FetchVideoFromTelegram(fileId string) (*Video, error) {
    Log("Downloading telegram video '%s'", fileId)
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
