package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"os"
	"os/signal"
	"time"
)

var done chan interface{}

func receiveHandler(connection *websocket.Conn) {
	defer close(done)
	for {
		_, msg, err := connection.ReadMessage()
		if err != nil {
			log.Println("Error in receive:", err)
			return
		}
		log.Printf("Received message: %s\n", msg)
	}
}

func main() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(
		interrupt,
	)

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
	// copied from gorilla's github since the connection was being closed.
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return
		case t := <-ticker.C:
			err := conn.WriteMessage(websocket.TextMessage, []byte(t.String()))
			if err != nil {
				log.Println("write:", err)
				return
			}
		case <-interrupt:
			log.Println("interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}
