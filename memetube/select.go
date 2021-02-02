package main

import (
	"context"
)

type VideoStream struct {
    ch chan(*Video)
    ctx context.Context
}

func (s *VideoStream) Chan() <-chan(*Video) {
    return s.ch
}

func (s *VideoStream) Close() {
    _, cancel := context.WithCancel(s.ctx)
    cancel()
}

func NewVideoStreamFromTelegramEndpoint(ctx context.Context, source *TelegramEndpointData) *VideoStream {
    var stream VideoStream
    stream.ch = make(chan(*Video))
    stream.ctx = ctx
    go func() {
        videos := source.Videos
        for i := 0; i < len(videos); i++ {
            video, err := FetchVideoFromTelegram(videos[i].FileID)
            if err != nil {
                continue
            }
            select {
            case stream.ch <- video:
                continue
            case <-ctx.Done():
                return
            }
        }
    }()
    return &stream
}

func VideoStreamLimitByTotalLength(stream *VideoStream, seconds int) *VideoStream {
    ctx, cancel := context.WithCancel(stream.ctx)
    var ret VideoStream
    ret.ch = make(chan(*Video))
    ret.ctx = ctx
    go func() {
        defer cancel()
        defer close(ret.ch)
        curLength := 0
        for video := range stream.Chan() {
            curLength += video.Length
            ret.ch <- video
            if (curLength > seconds) {
                return
            }
        }
    }()
    return &ret
}

func VideoStreamLimitByAmount(stream *VideoStream, amount int) *VideoStream {
    ctx, cancel := context.WithCancel(stream.ctx)
    var ret VideoStream
    ret.ch = make(chan(*Video))
    ret.ctx = ctx
    go func() {
        i := 0
        defer cancel()
        defer close(ret.ch)
        for video := range stream.Chan() {
            i++
            ret.ch <- video
            if (i > amount) {
                return
            }
        }
    }()
    return &ret
}
