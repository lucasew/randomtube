package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

func GetYoutubeToken() (token *oauth2.Token, err error) {
    token = new(oauth2.Token)
    endpointResponse, err := http.Get(fmt.Sprintf("%s/yta", FETCH_ENDPOINT))
    if err != nil {
        return
    }
    var rawToken struct {
        AccessToken string `json:"oauth_access_token"`
        RefreshToken string `json:"oauth_refresh_token"`
    }
    err = json.NewDecoder(endpointResponse.Body).Decode(&rawToken)
    if err != nil {
        return nil, err
    }
    token.AccessToken = rawToken.AccessToken
    token.RefreshToken = rawToken.RefreshToken
    return
}

func GetYoutubeTokenBytes() ([]byte, error) {
    endpointResponse, err := http.Get(fmt.Sprintf("%s/yta", FETCH_ENDPOINT))
    if err != nil {
        return nil, err
    }
    buf := bytes.NewBuffer([]byte{})
    _, err = io.Copy(buf, endpointResponse.Body)
    if err != nil {
        return nil, err
    }
    return buf.Bytes(), nil
}


func GetYoutubeService(ctx context.Context) (*youtube.Service, error) {
    token, err := GetYoutubeToken()
    if err != nil {
        return nil, err
    }
    return youtube.NewService(ctx, option.WithTokenSource(oauth2.StaticTokenSource(token)))
}

func PostVideo(ctx context.Context, video *Video) (*youtube.Video, error) {
    log.Printf("Posting '%s' to youtube", video.Filename)
    videoFile, err := os.Open(video.Filename)
    defer videoFile.Close()
    if err != nil {
        return nil, err
    }
    service, err := GetYoutubeService(ctx)
    if err != nil {
        return nil, err
    }
    ytVideo := &youtube.Video{
        Snippet: &youtube.VideoSnippet{
            Title: "Random video compilation",
            Description: "Random videos found across the internet",
        },
        Status: &youtube.VideoStatus{PrivacyStatus: "unlisted"},
    }
    ytVideo, err = service.Videos.
        Insert([]string{"snippet,status"}, ytVideo).
        NotifySubscribers(true).
        Context(ctx).Media(videoFile).
        Do()
    if err != nil {
        return nil, err
    }
    log.Printf("Video posted successfully on https://youtu.be/%s", ytVideo.Id)
    return ytVideo, nil
}
