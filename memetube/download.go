package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

func GetFilenameFromURL(url string) (string) {
    parts := strings.Split(url, ".")
    extension := ""
    if len(parts) > 1 {
        extension = parts[len(parts) - 1]
    }
    filename := uuid.New().String()
    if extension != "" {
        filename += "." + extension
    }
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
    log.Printf("Downloading %s (%d bytes)...", filename, resp.ContentLength)
    if err != nil {
        return "", err
    }
    _, err = io.Copy(f, resp.Body)
    return filename, err
}

