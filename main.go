package main

// TODO: Go ahead with double stack
// TODO: Is it better Unix or TCP socket?
// TODO: Think other way for ignore events sent by myself

var ignoreEvent int

func main() {
	ds := NewDoubleStack[string](10)
	wh := NewWorkspaceHandler(ds)
	s := NewServer("tcp", ":43223")
	ready := make(chan struct{})
	go wh.ListenEvents()
	go s.Start(wh.Handle, ready)
	select {}
}
