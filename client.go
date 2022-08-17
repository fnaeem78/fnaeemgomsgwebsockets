package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
)

var done chan interface{}

func receiveHandler(connection *websocket.Conn) {
	defer close(done)
	for {
		_, msg, err := connection.ReadMessage()
		if err != nil {
			fmt.Println("Error in receive:", err)
			return
		}
		fmt.Printf("Received message: %s\n", msg)
	}
}

func main() {

	var topicToSubscribe string
	fmt.Scanf("%s", &topicToSubscribe)
	socketUrl := "ws://localhost:8081" + "/subscribe?topic=" + topicToSubscribe
	fmt.Printf("Connecting to %s", socketUrl)
	conn, _, err := websocket.DefaultDialer.Dial(socketUrl, nil)
	if err != nil {
		log.Fatal("Error connecting to Websocket Server:", err)
	}
	defer conn.Close()
	go receiveHandler(conn)
	// stopping prog from exiting while receiveHandler runs in the bck
	fmt.Scanln()

}
