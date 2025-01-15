package main

// TODO: Go ahead with double stack
// TODO: Is it better Unix or TCP socket?
// TODO: Think other way for ignore events sent by myself

var x int

func main() {
	// var mutex sync.Mutex
	//
	// go func() {
	// 	subscription := i3.Subscribe(i3.WorkspaceEventType)
	// 	defer subscription.Close()
	//
	// 	for subscription.Next() {
	// 		event := subscription.Event().(*i3.WorkspaceEvent)
	// 		if event.Change == "focus" {
	// 			mutex.Lock()
	// 			if event.Old.Name != "" {
	// 				if x == 0 {
	// 					s.Push(event.Old.Name)
	// 					log.Printf("[ %s ] - Workspace pushed to stack", event.Old.Name)
	// 				} else {
	// 					x--
	// 				}
	// 			}
	// 			mutex.Unlock()
	// 		}
	// 	}
	//
	// 	if err := subscription.Close(); err != nil {
	// 		log.Printf("Error closing subscription: %v", err)
	// 	}
	// }()
	//
	// }
}

// func handleConnection(conn net.Conn, s *stack.Stack, mutex *sync.Mutex) {
// 	defer conn.Close()
//
// 	buf := make([]byte, 4096)
// 	n, err := conn.Read(buf)
// 	if err != nil {
// 		log.Printf("Error reading from connection: %v", err)
// 		return
// 	}
//
// 	msg := strings.TrimSpace(string(buf[:n]))
// 	log.Printf("Message received: %s", msg)
//
// 	if msg == "back" {
// 		mutex.Lock()
// 		defer mutex.Unlock()
//
// 		if s.Len() == 0 {
// 			log.Println("Stack is empty, no workspace to go back to")
// 			return
// 		}
//
// 		workspaceName := s.Pop().(string)
// 		log.Printf("Switching to workspace: %s", workspaceName)
// 		x++
// 		if _, err := i3.RunCommand("workspace " + workspaceName); err != nil {
// 			log.Printf("Failed to switch workspace: %v", err)
// 		}
// 	} else {
// 		log.Printf("Invalid command: %s", msg)
// 	}
// }
