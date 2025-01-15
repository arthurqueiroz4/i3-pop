package main

import (
	"go.i3wm.org/i3/v4"
	"log"
	"sync"
)

type WorkspaceHandler struct {
	ds *DoubleStack[string]
}

func NewWorkspaceHandler(ds *DoubleStack[string]) *WorkspaceHandler {
	return &WorkspaceHandler{ds: ds}
}

func (wh *WorkspaceHandler) ListenEvents() {
	var mutex sync.Mutex
	subscription := i3.Subscribe(i3.WorkspaceEventType)
	defer subscription.Close()
	for subscription.Next() {
		event := subscription.Event().(*i3.WorkspaceEvent)
		if event.Change == "focus" {
			mutex.Lock()
			if event.Old.Name != "" {
				if x == 0 {
					wh.ds.PushOnBack(event.Old.Name)
					log.Printf("[ %s ] - Workspace pushed to stack", event.Old.Name)
				} else {
					x--
				}
			}
			mutex.Unlock()
		}
	}

	if err := subscription.Close(); err != nil {
		log.Printf("Error closing subscription: %v", err)
	}
}
