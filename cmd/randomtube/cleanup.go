package main

import (
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
        Log("cleanup: Deleting file %s", filename)
        os.Remove(filename)
    })
}

func CleanupPhase() {
    if dontCleanup {
        return
    }
    cleanupHooksLock.Lock()
    defer cleanupHooksLock.Unlock()
    Log("Starting cleanup phase with %d items", len(cleanupHooks))
    for _, step := range cleanupHooks {
        step()
    }
}
