package main

import (
	"context"
)

type VideoStream struct {
    ch chan(*Video)
    ctx context.Context
}

type VideoStreamProps struct {
    Seconds int
    Amount int
}

func (s *VideoStream) Chan() <-chan(*Video) {
    return s.ch
}

func (s *VideoStream) Close() {
    _, cancel := context.WithCancel(s.ctx)
    cancel()
}

func NewVideoStreamFromTelegramEndpoint(ctx context.Context, source *TelegramEndpointData, props VideoStreamProps) chan(*Video) {
    ch := make(chan(*Video))
    go func() {
        seconds := 0
        amount := 0
        defer close(ch)
        videos := source.Videos
        for i := 0; i < len(videos); i++ {
            if amount > props.Amount && props.Amount > 0 {
                return
            }
            if seconds > props.Seconds && props.Seconds > 0 {
                return
            }
            seconds += videos[i].Length
            amount++
            video, err := FetchVideoFromTelegram(videos[i].FileID)
            select {
            case <-ctx.Done():
                return
            default:
                if err != nil {
                    Log("Error when looking for '%s': %s", videos[i].FileID, err)
                    continue
                }
            }
            select {
            case ch <-video:
                continue
            case <-ctx.Done():
                return
            }
        }
    }()
    return ch
}
