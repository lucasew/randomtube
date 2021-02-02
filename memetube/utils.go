package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

func MustBinary(binary string) {
    _, err := exec.LookPath(binary)
    if err != nil {
        log.Fatal(err)
    }
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
    return filename, nil
}
