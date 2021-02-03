package main

import (
	"log"
	"os"
	"sync"

	"github.com/google/uuid"
)

var cleanupHooks = map[string]func(){}
var cleanupHooksLock sync.Mutex

func AddCleanupHook(f func()) string {
    id := uuid.New().String()
    cleanupHooksLock.Lock()
    defer cleanupHooksLock.Unlock()
    cleanupHooks[id] = f
    return id
}

func AddFileCleanupHook(filename string) string {
    return AddCleanupHook(func() {
        log.Printf("cleanup: Deleting file %s", filename)
        os.Remove(filename)
    })
}

func CleanupPhase() {
    cleanupHooksLock.Lock()
    defer cleanupHooksLock.Unlock()
    log.Printf("Starting cleanup phase with %d items", len(cleanupHooks))
    for _, step := range cleanupHooks {
        step()
    }
}
