package main

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

func GetFilenameFromURL(url string) (string) {
    urlPart := ""
    if len(url) < 10 {
        urlPart = url
    } else {
        urlPart = url[len(url) - 10:]
    }
    parts := strings.Split(urlPart, ".")
    extension := ""
    if len(parts) > 1 {
        extension = parts[len(parts) - 1]
    }
    filename := uuid.New().String()
    if extension == "" {
        extension = "mp4"
    }
    filename += "."
    filename += extension
    return filepath.Join(tempdir, filename)
}


func Download(url string) (path string, err error) {
    filename := GetFilenameFromURL(url)
    f, err := os.Create(filename)
    defer f.Close()
    if err != nil {
        return "", err
    }
    resp, err := http.Get(url)
    Log("Downloading %s (%d bytes)...", filename, resp.ContentLength)
    if err != nil {
        return "", err
    }
    _, err = io.Copy(f, resp.Body)
    AddFileCleanupHook(filename)
    return filename, err
}

