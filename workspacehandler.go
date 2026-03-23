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

type eventSubscription interface {
	Next() bool
	Event() i3.Event
	Close() error
}

var subscribeWorkspaceEvents = func() eventSubscription {
	return i3.Subscribe(i3.WorkspaceEventType)
}

var runI3Command = i3.RunCommand

var shouldStoreEventInBackStack = true
var shouldStoreEventInFrontStack = false
var workspaceNavigationMutex sync.Mutex

func (wh *WorkspaceHandler) ListenEvents() {
	subscription := subscribeWorkspaceEvents()
	defer subscription.Close()
	for subscription.Next() {
		event := subscription.Event().(*i3.WorkspaceEvent)
		if event.Change == "focus" {
			workspaceNavigationMutex.Lock()
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

			workspaceNavigationMutex.Unlock()
		}
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
	goToWorkspaceName(*namePtr, false, true)
}

func (wh *WorkspaceHandler) treatFrontMessage() {
	namePtr := wh.ds.PopOnFront()
	if namePtr == nil {
		return
	}
	goToWorkspaceName(*namePtr, true, false)
}

func goToWorkspaceName(name string, storeOnBack, storeOnFront bool) {
	workspaceNavigationMutex.Lock()
	shouldStoreEventInBackStack = storeOnBack
	shouldStoreEventInFrontStack = storeOnFront
	workspaceNavigationMutex.Unlock()
	command, err := runI3Command("workspace " + name)
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
