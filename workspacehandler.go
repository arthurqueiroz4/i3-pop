package main

import (
	"go.i3wm.org/i3/v4"
	"log"
	"net"
	"strings"
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
				if ignoreEvent == 0 {
					wh.ds.PushOnBack(event.Old.Name)
					log.Printf("[ %s ] - Workspace pushed to stack", event.Old.Name)
				} else {
					ignoreEvent--
				}
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

func (wh *WorkspaceHandler) Handle(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte, 4096)
	n, err := conn.Read(buf)
	if err != nil {
		log.Printf("Error reading from connection: %v", err)
		return
	}

	msg := strings.TrimSpace(string(buf[:n]))
	switch msg {
	case BACK:
		wh.treatBackMessage()
	case FRONT:
		wh.treatFrontMessage()
	}
	_, err = conn.Write([]byte("Message received and processed"))
}

func (wh *WorkspaceHandler) treatBackMessage() {
	namePtr := wh.ds.PopOnBackAndPutOnFront()
	if namePtr == nil {
		return
	}
	goToWorkspaceName(*namePtr)
}

func (wh *WorkspaceHandler) treatFrontMessage() {
	// TODO: Implement this
	//name := wh.ds.PopOnFrontAndPutOnBack()
	//goToWorkspaceName(*name)
}

func goToWorkspaceName(name string) {
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
