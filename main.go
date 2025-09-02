package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

// Upgrader is used to upgrade HTTP connections to WebSocket connections.
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	// Upgrade the HTTP connection to a WebSocket connection
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Error upgrading:", err)
		return
	}
	defer conn.Close()
	// Listen for incoming messages
	for {
		// Read message from the client
		_, message, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Error reading message:", err)
			break
		}
		// Try to parse as JSON for location payloads
		var payload map[string]interface{}
		if err := json.Unmarshal(message, &payload); err == nil {
			if t, ok := payload["type"].(string); ok && t == "location" {
				// Log a friendly line
				fmt.Printf("Location: lat=%v, lng=%v, acc=%v, ts=%v\n", payload["lat"], payload["lng"], payload["accuracy"], payload["timestamp"])
				// Echo back the same JSON
				if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
					fmt.Println("Error writing message:", err)
					break
				}
				continue
			}
		}
		fmt.Printf("Received: %s\\n", message)
		// Echo the message back to the client (non-JSON or non-location)
		if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
			fmt.Println("Error writing message:", err)
			break
		}
	}
}

func main() {
	http.HandleFunc("/ws", wsHandler)
	// Serve the static client (client.html) and any assets from current dir
	http.Handle("/", http.FileServer(http.Dir(".")))
	fmt.Println("WebSocket server started on :1234")
	err := http.ListenAndServe(":1234", nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
