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
    dontPost bool
    reportChat int64
    youtubePrivacyStatus string
)

func main() {
    Log("Starting up...")
    endpointResult, err := FetchTelegramEndpoint()
    BailOutIfError(err)
    host, err := os.Hostname()
    BailOutIfError(err)
    Report("Processo de geração de vídeo iniciado em '%s'. Limite %ds ou %d videos, 0 é sem limite", host, maxSeconds, maxVideos)
    Report("Tamanho da fila: %d vídeos (%ds)", endpointResult.QueueSize, endpointResult.QueueTime)
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
    Log("Generated video: %s", joinedVideo.Filename)
    Command("ffprobe", joinedVideo.Filename)
    if (!dontPost) {
        video, err := PostVideo(ctx, joinedVideo)
        BailOutIfError(err)
        Report("Video postado em http://youtu.be/%s", video.Id)
        MarkTelegramVideosAsProcessedCleanupHook()
    } else {
        Report("Vídeo configurado para não ser postado")
    }
}

func init() {
    flag.StringVar(&TELEGRAM_BOT, "tg", os.Getenv("TELEGRAM_BOT"), "Telegram bot token")
    flag.StringVar(&FETCH_ENDPOINT, "fe", os.Getenv("FETCH_ENDPOINT"), "Fetch endpoint (pipedream)")
    flag.StringVar(&tempdir, "tmp", os.TempDir(), "Temporary folder for file processing")
    flag.StringVar(&youtubePrivacyStatus, "ps", "public", "Youtube video privacy status, can be public, unlisted or private")
    flag.IntVar(&maxVideos, "mv", 0, "Max number of videos in a bundle")
    flag.IntVar(&maxSeconds, "ms", 0, "Stop adding videos when their lengths pass x seconds")
    flag.BoolVar(&dontMarkVideosAsProcessed, "dd", false, "Don't delete processed videos from the queue")
    flag.BoolVar(&dontCleanup, "dc", false, "Don't cleanup processed artifacts")
    flag.BoolVar(&dontPost, "dp", false, "Don't post the video, implies dd")
    flag.Parse()
}

