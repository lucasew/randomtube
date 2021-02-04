package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
)

func MustBinary(binary string) {
    _, err := exec.LookPath(binary)
    BailOutIfError(err)
}

func BailOutIfError(err error) {
    if err != nil {
        log.Fatal(err)
    }
}

func WriteLines(lines ...string) (string, error) {
    filename := GetFilenameFromURL("whatever.txt")
    f, err := os.Create(filename)
    defer f.Close()
    if err != nil {
        return "", err
    }
    for _, file := range lines {
        fmt.Fprintln(f, file)
    }
    AddFileCleanupHook(filename)
    return filename, nil
}

func Command(binary string, args ...string) error {
    Log("Run command: %s %s", binary, args)
    cmd := exec.Command(binary, args...)
    cmd.Stderr = os.Stderr
    return cmd.Run()
}

type ReadCloserWrapper struct {
    io.Reader
}
func (ReadCloserWrapper) Close() error {
    return nil
}

func NewReadCloserWrapper(w io.Reader) io.ReadCloser {
    return ReadCloserWrapper{w}
}


func Report(message string, format ...interface{}) {
    text := struct{
        ChatID int64 `json:"chat_id"`
        Message string `json:"text"`
    }{
        Message: fmt.Sprintf(message, format...),
        ChatID: reportChat,
    }
    buf := bytes.NewBufferString("")
    BailOutIfError(json.NewEncoder(buf).Encode(text))
    u, err := url.Parse(fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", TELEGRAM_BOT))
    BailOutIfError(err)
    req := http.Request{
        Method: http.MethodGet,
        URL: u,
        Body: NewReadCloserWrapper(buf),
        Header: http.Header{
            "Content-Type": []string{"application/json"},
        },
    }
    res, err := http.DefaultClient.Do(&req)
    BailOutIfError(err)
    if res.StatusCode != 200 {
        io.Copy(os.Stdout, res.Body)
        BailOutIfError(fmt.Errorf("Report returned status %d", res.StatusCode))
    }
    Log("Report sent!")
}

func Log(message string, format ...interface{}) {
    log.Printf(message, format...)
}
