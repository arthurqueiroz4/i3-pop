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

var shouldStoreEventInBackStack = true
var shouldStoreEventInFrontStack = false

func (wh *WorkspaceHandler) ListenEvents() {
	var mutex sync.Mutex
	subscription := i3.Subscribe(i3.WorkspaceEventType)
	defer subscription.Close()
	for subscription.Next() {
		event := subscription.Event().(*i3.WorkspaceEvent)
		if event.Change == "focus" {
			mutex.Lock()
			if shouldStoreEventInBackStack {
				wh.ds.PushOnBack(event.Old.Name)
				log.Printf("[ %s ] - Workspace pushed to back stack", event.Old.Name)
			} else {
				shouldStoreEventInBackStack = true
			}

			if shouldStoreEventInFrontStack {
				wh.ds.PushOnFront(event.Old.Name)
				log.Printf("[ %s ] - Workspace pushed to front stack", event.Old.Name)
				shouldStoreEventInFrontStack = false
			}

			mutex.Unlock()
		}
	}

	if err := subscription.Close(); err != nil {
		log.Printf("Error closing subscription: %v", err)
	}
}

var (
	BACK  = "back"
	FRONT = "front"
)

func (wh *WorkspaceHandler) Process(msg string) {
	switch msg {
	case BACK:
		wh.treatBackMessage()
	case FRONT:
		wh.treatFrontMessage()
	}
}

func (wh *WorkspaceHandler) treatBackMessage() {
	namePtr := wh.ds.PopOnBack()
	if namePtr == nil {
		return
	}
	goToWorkspaceName(*namePtr)
}

func (wh *WorkspaceHandler) treatFrontMessage() {
	namePtr := wh.ds.PopOnFront()
	if namePtr == nil {
		return
	}
	goToWorkspaceName(*namePtr)
}

func goToWorkspaceName(name string) {
	shouldStoreEventInBackStack = false
	shouldStoreEventInFrontStack = true
	command, err := i3.RunCommand("workspace" + name)
	if err != nil {
		log.Printf("Error running command on i3: %v", err)
		return
	}
	for _, c := range command {
		if !c.Success {
			log.Printf("Command returned unsuccessful message: %v", c)
		}
	}
}
