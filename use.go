package rpc

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// Allow connections from any origin
		return true
	},
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	// Upgrade the HTTP connection to a WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Error upgrading connection:", err)
		return
	}
	defer conn.Close()

	fmt.Println("Client connected!")

	// Echo messages back to the client
	for {
		// Read a message
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Error reading message:", err)
			break
		}

		// Print the received message
		fmt.Printf("Received: %s\n", message)

		// Write the message back to the client
		err = conn.WriteMessage(messageType, message)
		if err != nil {
			fmt.Println("Error writing message:", err)
			break
		}
	}
}

// func main() {
// 	Load_file("main.go")
// 	Load_file("parser.go")
// 	fmt.Println("functions_in_file: ", functions_in_file)

// 	os.WriteFile(typescriptRPCFileName, []byte(""), 0644)

// 	Setup_rpc(add)
// 	Setup_rpc(printNum)
// 	Setup_rpc(Load_file)

// 	http.HandleFunc("/ws", handleWebSocket)
// 	http.ListenAndServe(":8080", nil)

// }
