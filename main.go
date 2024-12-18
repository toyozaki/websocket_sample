package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// To manage WebSocket connections
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// From client
type Message struct {
	Type    string `json:"type"`
	Payload string `json:"payload"`
}

// From server
type Response struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	defer conn.Close()

	log.Println("Client connected")

	for {
		var msg Message
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Println("Read error:", err)
			break
		}

		log.Printf("Received message: %+v\n", msg)

		switch msg.Type {
		case "greeting":
			sendResponse(conn, "response", "Hello, client!")
		case "echo":
			sendResponse(conn, "response", msg.Payload)
		case "close":
			sendResponse(conn, "response", "Goodbye, client!")
			return
		default:
			sendResponse(conn, "error", "Unknown message type")
		}
	}

	log.Println("Client disconnected")
}

func sendResponse(conn *websocket.Conn, responseType string, message string) {
	response := Response{
		Type:    responseType,
		Message: message,
	}
	if err := conn.WriteJSON(response); err != nil {
		log.Println("Write error:", err)
	}
}

func main() {
	http.HandleFunc("/ws", handleWebSocket)

	fmt.Println("WebSocket server is running on ws://localhost:8080/ws")
	log.Fatal(http.ListenAndServe("127.0.0.1:8080", nil))
}
