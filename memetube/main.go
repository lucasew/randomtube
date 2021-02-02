package main

import (
	"context"
	"flag"
	"os"
	"os/exec"
)

var (
    TELEGRAM_BOT string
    FETCH_ENDPOINT string
    tempdir string
    maxVideos int
    maxSeconds int
)

func init() {
    flag.StringVar(&TELEGRAM_BOT, "tg", os.Getenv("TELEGRAM_BOT"), "Telegram bot token")
    flag.StringVar(&FETCH_ENDPOINT, "fe", os.Getenv("FETCH_ENDPOINT"), "Fetch endpoint (pipedream)")
    flag.StringVar(&tempdir, "tmp", os.TempDir(), "Temporary folder for file processing")
    flag.IntVar(&maxVideos, "mv", 0, "Max number of videos in a bundle")
    flag.IntVar(&maxSeconds, "ms", 0, "Stop adding videos when their lengths pass x seconds")
    flag.Parse()
}

func main() {
    ctx := context.Background()
    MustBinary("ffmpeg")
    endpointResult, err := FetchTelegramEndpoint(FETCH_ENDPOINT)
    BailOutIfError(err)
    videoStream := NewVideoStreamFromTelegramEndpoint(ctx, endpointResult)
    if maxVideos != 0 {
        videoStream = VideoStreamLimitByAmount(videoStream, maxVideos)
    }
    if maxSeconds != 0 {
        videoStream = VideoStreamLimitByTotalLength(videoStream, maxSeconds)
    }
    downloadedVideos := []*Video{}
    defer videoStream.Close()
    for video := range videoStream.Chan() {
        downloadedVideos = append(downloadedVideos, video)
    }
    joinedVideo, err := ConcatVideos(downloadedVideos...)
    BailOutIfError(err)
    cmd := exec.Command("mv", joinedVideo.Filename, "/tmp/out_memetube.mp4")
    BailOutIfError(cmd.Run())
}
