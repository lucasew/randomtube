package main

import (
	"bytes"
	"fmt"
	"os/exec"
)

type Video struct {
    Filename string
    Length int
}

const ffmpegFilter = `
[0:v]scale=1920*2:-1,crop=1920:1080,boxblur=luma_radius=min(h\,w)/20:luma_power=1:chroma_radius=min(cw\,ch)/20:chroma_power=1[bg];
[0:v]scale=-1:1080[ov];
[bg][ov]overlay=(W-w)/2,crop=1920:1080[composed];
anullsrc=r=44100:cl=stereo[anull]
`

func NewVideoFromAnotherVideo(filename string) (*Video, error) {
    var ret Video
    ret.Filename = GetFilenameFromURL("output.mkv")
    Log("Ingesting video: '%s' as '%s'", filename, ret.Filename)
    sizebuf := bytes.NewBuffer([]byte{})
    cmd := exec.Command("ffprobe", "-v", "error", "-show_entries", "format=duration", "-of", "default=noprint_wrappers=1:nokey=1", filename)
    cmd.Stdout = sizebuf
    err := cmd.Run()
    if err != nil {
        return nil, err
    }
    var length float32
    fmt.Fscan(sizebuf, &length)
    ret.Length = int(length)
    err = Command("ffmpeg", "-i", filename,
        "-t", fmt.Sprintf("%d", ret.Length),
        "-c:a", "aac", "-c:v", "h264",
        "-filter_complex", ffmpegFilter,
        "-map", "[composed]",
        "-map", "0:a?",
        "-map", "[anull]",
        ret.Filename,
    )
    AddFileCleanupHook(ret.Filename)
    if err != nil {
        return nil, err
    }
    return &ret, nil
}

func ConcatVideos(videos ...*Video) (*Video, error) {
    var ret Video
    ret.Filename = GetFilenameFromURL("output.mkv")
    ret.Length = 0
    Log("Concating %d videos as '%s'", len(videos), ret.Filename)
    for _, video := range videos {
        ret.Length += video.Length
    }
    files := make([]string, len(videos))
    for i := 0; i < len(videos); i++ {
        files[i] = fmt.Sprintf("file '%s'", videos[i].Filename)
    }
    listFile, err := WriteLines(files...)
    if err != nil {
        return nil, err
    }
    err = Command("ffmpeg", "-f", "concat", "-safe", "0", "-i", listFile, ret.Filename)
    if err != nil {
        return nil, err
    }
    AddFileCleanupHook(ret.Filename)
    return &ret, nil
}
