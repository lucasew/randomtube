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
)

func main() {
    defer CleanupPhase()
    Report("Processo de geração de vídeo iniciado. Limite %ds ou %d videos, 0 é sem limite", maxSeconds, maxVideos)
    ctx, cancel := context.WithCancel(context.Background())
    AddCleanupHook(cancel)
    Log("Starting up...")
    MustBinary("ffmpeg")
    joinedVideo, err := ConcatVideos(GetVideos(ctx)...)
    BailOutIfError(err)
    video, err := PostVideo(ctx, joinedVideo)
    BailOutIfError(err)
    MarkTelegramVideosAsProcessedCleanupHook()
    Report("Video postado em http://youtu.be/%s", video.Id)
}

func init() {
    flag.StringVar(&TELEGRAM_BOT, "tg", os.Getenv("TELEGRAM_BOT"), "Telegram bot token")
    flag.StringVar(&FETCH_ENDPOINT, "fe", os.Getenv("FETCH_ENDPOINT"), "Fetch endpoint (pipedream)")
    flag.StringVar(&tempdir, "tmp", os.TempDir(), "Temporary folder for file processing")
    flag.IntVar(&maxVideos, "mv", 0, "Max number of videos in a bundle")
    flag.IntVar(&maxSeconds, "ms", 0, "Stop adding videos when their lengths pass x seconds")
    flag.BoolVar(&dontMarkVideosAsProcessed, "dp", false, "Don't delete processed videos from the queue")
    flag.Parse()
}

func SetupVideoStream(ctx context.Context) chan(*Video) {
    endpointResult, err := FetchTelegramEndpoint(FETCH_ENDPOINT)
    BailOutIfError(err)
    videos := NewVideoStreamFromTelegramEndpoint(ctx, endpointResult, VideoStreamProps{
        Seconds: maxSeconds,
        Amount: maxVideos,
    })
    return videos
}

func GetVideos(ctx context.Context) ([]*Video) {
    ctx, cancel := context.WithCancel(ctx)
    defer cancel()
    videoStream := SetupVideoStream(ctx)
    downloadedVideos := make([]*Video, 0, 10)
    for video := range videoStream {
        downloadedVideos = append(downloadedVideos, video)
    }
    return downloadedVideos
}
