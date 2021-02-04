package main

import (
	"context"
	"flag"
	"os"
)

var (
    TELEGRAM_BOT string
    FETCH_ENDPOINT string
    tempdir string
    maxVideos int
    maxSeconds int
    dontMarkVideosAsProcessed bool
    dontCleanup bool
    reportChat int64
)

func main() {
    Log("Starting up...")
    endpointResult, err := FetchTelegramEndpoint()
    BailOutIfError(err)
    Report("Processo de geração de vídeo iniciado. Limite %ds ou %d videos, 0 é sem limite", maxSeconds, maxVideos)
    defer CleanupPhase()
    ctx, cancel := context.WithCancel(context.Background())
    AddCleanupHook(cancel)

    downloadedVideos := make([]*Video, 0, 10)
    for video := range NewVideoStreamFromTelegramEndpoint(ctx, endpointResult, VideoStreamProps{
        Seconds: maxSeconds,
        Amount: maxVideos,
    }) {
        downloadedVideos = append(downloadedVideos, video)
    }
    MustBinary("ffmpeg")
    joinedVideo, err := ConcatVideos(downloadedVideos...)
    BailOutIfError(err)
    video, err := PostVideo(ctx, joinedVideo)
    BailOutIfError(err)
    Report("Video postado em http://youtu.be/%s", video.Id)
    MarkTelegramVideosAsProcessedCleanupHook()
}

func init() {
    flag.StringVar(&TELEGRAM_BOT, "tg", os.Getenv("TELEGRAM_BOT"), "Telegram bot token")
    flag.StringVar(&FETCH_ENDPOINT, "fe", os.Getenv("FETCH_ENDPOINT"), "Fetch endpoint (pipedream)")
    flag.StringVar(&tempdir, "tmp", os.TempDir(), "Temporary folder for file processing")
    flag.IntVar(&maxVideos, "mv", 0, "Max number of videos in a bundle")
    flag.IntVar(&maxSeconds, "ms", 0, "Stop adding videos when their lengths pass x seconds")
    flag.BoolVar(&dontMarkVideosAsProcessed, "dp", false, "Don't delete processed videos from the queue")
    flag.BoolVar(&dontCleanup, "dc", false, "Don't cleanup processed artifacts")
    flag.Parse()
}

