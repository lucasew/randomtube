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


func Report(message string, format ...interface{}) error {
    text := struct{
        Message string `json:"messsage"`
    }{
        Message: fmt.Sprintf(message, format...),
    }
    buf := bytes.NewBufferString("")
    err := json.NewEncoder(buf).Encode(text)
    if err != nil {
        return err
    }
    u, _ := url.Parse(fmt.Sprintf("%s/notify", FETCH_ENDPOINT))
    req := http.Request{
        Method: http.MethodGet,
        URL: u,
        Body: NewReadCloserWrapper(buf),
    }
    _, err = http.DefaultClient.Do(&req)
    if err != nil {
        log.Printf("Failed to report data upstream: %s", err)
    }
    return err
}

func Log(message string, format ...interface{}) {
    log.Printf(message, format...)
}
