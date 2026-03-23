package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.i3wm.org/i3/v4"
)

type fakeEventSubscription struct {
	events []i3.Event
	index  int
	closed bool
}

func (f *fakeEventSubscription) Next() bool {
	if f.index >= len(f.events) {
		return false
	}
	f.index++
	return true
}

func (f *fakeEventSubscription) Event() i3.Event {
	return f.events[f.index-1]
}

func (f *fakeEventSubscription) Close() error {
	f.closed = true
	return nil
}

func resetWorkspaceNavigationState() {
	workspaceNavigationMutex.Lock()
	shouldStoreEventInBackStack = true
	shouldStoreEventInFrontStack = false
	workspaceNavigationMutex.Unlock()
}

func TestWorkspaceHandlerProcessBack(t *testing.T) {
	resetWorkspaceNavigationState()
	ds := NewDoubleStack[string](10)
	wh := NewWorkspaceHandler(ds)
	ds.PushOnBack("2")

	var executedCommand string
	prevRunCommand := runI3Command
	runI3Command = func(command string) ([]i3.CommandResult, error) {
		executedCommand = command
		return []i3.CommandResult{{Success: true}}, nil
	}
	t.Cleanup(func() {
		runI3Command = prevRunCommand
	})

	wh.Process(BACK)

	assert.Equal(t, "workspace 2", executedCommand)
	assert.Equal(t, false, shouldStoreEventInBackStack)
	assert.Equal(t, true, shouldStoreEventInFrontStack)
}

func TestWorkspaceHandlerProcessFront(t *testing.T) {
	resetWorkspaceNavigationState()
	ds := NewDoubleStack[string](10)
	wh := NewWorkspaceHandler(ds)
	ds.PushOnFront("9")

	var executedCommand string
	prevRunCommand := runI3Command
	runI3Command = func(command string) ([]i3.CommandResult, error) {
		executedCommand = command
		return []i3.CommandResult{{Success: true}}, nil
	}
	t.Cleanup(func() {
		runI3Command = prevRunCommand
	})

	wh.Process(FRONT)

	assert.Equal(t, "workspace 9", executedCommand)
	assert.Equal(t, true, shouldStoreEventInBackStack)
	assert.Equal(t, false, shouldStoreEventInFrontStack)
}

func TestWorkspaceHandlerListenEventsStoresBackHistoryOnFocus(t *testing.T) {
	resetWorkspaceNavigationState()
	ds := NewDoubleStack[string](10)
	wh := NewWorkspaceHandler(ds)

	sub := &fakeEventSubscription{
		events: []i3.Event{
			&i3.WorkspaceEvent{
				Change: "init",
				Old:    i3.Node{Name: "0"},
			},
			&i3.WorkspaceEvent{
				Change: "focus",
				Old:    i3.Node{Name: "1"},
			},
			&i3.WorkspaceEvent{
				Change: "focus",
				Old:    i3.Node{Name: "2"},
			},
		},
	}

	prevSubscribe := subscribeWorkspaceEvents
	subscribeWorkspaceEvents = func() eventSubscription {
		return sub
	}
	t.Cleanup(func() {
		subscribeWorkspaceEvents = prevSubscribe
	})

	wh.ListenEvents()

	assert.True(t, sub.closed)
	assert.Equal(t, 2, ds.BackLength())
	assert.Equal(t, 0, ds.FrontLength())
	assert.Equal(t, "2", *ds.PeekOnBack())
	assert.Equal(t, true, shouldStoreEventInBackStack)
	assert.Equal(t, false, shouldStoreEventInFrontStack)
}

func TestWorkspaceHandlerListenEventsStoresFrontWhenFlagIsEnabled(t *testing.T) {
	resetWorkspaceNavigationState()
	workspaceNavigationMutex.Lock()
	shouldStoreEventInBackStack = false
	shouldStoreEventInFrontStack = true
	workspaceNavigationMutex.Unlock()

	ds := NewDoubleStack[string](10)
	wh := NewWorkspaceHandler(ds)

	sub := &fakeEventSubscription{
		events: []i3.Event{
			&i3.WorkspaceEvent{
				Change: "focus",
				Old:    i3.Node{Name: "A"},
			},
		},
	}

	prevSubscribe := subscribeWorkspaceEvents
	subscribeWorkspaceEvents = func() eventSubscription {
		return sub
	}
	t.Cleanup(func() {
		subscribeWorkspaceEvents = prevSubscribe
	})

	wh.ListenEvents()

	assert.True(t, sub.closed)
	assert.Equal(t, 0, ds.BackLength())
	assert.Equal(t, 1, ds.FrontLength())
	assert.Equal(t, "A", *ds.PeekOnFront())
	assert.Equal(t, true, shouldStoreEventInBackStack)
	assert.Equal(t, false, shouldStoreEventInFrontStack)
}
