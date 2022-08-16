package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"io"
	"log"
	"net/http"
)

func main() {

	var mapOfTopicAndConnections = make(map[string][]websocket.Conn)

	http.HandleFunc("/", writeResponse)
	http.HandleFunc("/subscribe", func(writer http.ResponseWriter, request *http.Request) {
		var upgrader = websocket.Upgrader{}
		var conn, err = upgrader.Upgrade(writer, request, nil)
		if err != nil {
			log.Print("Error during connection upgradation:", err)
			return
		}
		subscribeTopic(writer, request, &mapOfTopicAndConnections, *conn)
		defer conn.Close()
	})

	http.HandleFunc("/publish", func(writer http.ResponseWriter, request *http.Request) {
		body, err := io.ReadAll(request.Body)
		if err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			writer.Write([]byte("no message in body to publish"))
			return
		}
		publishMessage(writer, request, mapOfTopicAndConnections, body)
	})

	http.ListenAndServe(":8081", nil)

}

func publishMessage(writer http.ResponseWriter, request *http.Request, subscribedTopicsMap map[string][]websocket.Conn, messageToPublish []byte) bool {
	query := request.URL.Query()
	topic, present := query["topic"]
	var currentTopic string = topic[0]
	if !present || len(topic) == 0 {
		writer.WriteHeader(http.StatusBadRequest)
		return false
	}
	if containsTopic(subscribedTopicsMap, currentTopic) {
		conns := subscribedTopicsMap[currentTopic]
		var topicName []string
		topicName = append(topicName, "[topic= ", currentTopic, "]:   ", string(messageToPublish))
		for _, currentConn := range conns {
			currentConn.WriteJSON(topicName)
		}
		return true
	} else {
		writer.WriteHeader(http.StatusNotFound)
		writer.Write([]byte("topic not found"))
		return false
	}

}

func subscribeTopic(writer http.ResponseWriter, request *http.Request, subscribedTopicsMap *map[string][]websocket.Conn, conn websocket.Conn) *map[string][]websocket.Conn {
	query := request.URL.Query()
	topic, present := query["topic"]
	var currentTopic string = topic[0]
	if !present || len(topic) == 0 {
		writer.WriteHeader(http.StatusBadRequest)
		return subscribedTopicsMap
	} else {
		conns := (*subscribedTopicsMap)[currentTopic]
		(*subscribedTopicsMap)[currentTopic] = append(conns, conn)
		//writer.WriteHeader(http.StatusOK)
		//writer.Write([]byte("topic added"))
		return subscribedTopicsMap
	}

}

func containsTopic(topics map[string][]websocket.Conn, newtopic string) bool {
	for key, _ := range topics {
		if key == newtopic {
			return true
		}
	}
	return false
}

func writeResponse(writer http.ResponseWriter, request *http.Request) {
	fmt.Fprintf(writer, "starting server on 8081")
}
