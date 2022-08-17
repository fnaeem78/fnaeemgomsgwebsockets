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
	// if no topic was specified return false
	if !present || len(topic) == 0 {
		writer.WriteHeader(http.StatusBadRequest)
		return false
	}
	// if the topic was found we want to find
	// all the connections that are the value of that topic in our map
	// and write the message to them.
	if containsTopic(subscribedTopicsMap, currentTopic) {
		conns := subscribedTopicsMap[currentTopic]
		var topicName []string
		topicName = append(topicName, "[topic= ", currentTopic, "]:   ", string(messageToPublish))
		for _, currentConn := range conns {
			currentConn.WriteJSON(topicName)
		}
		return true
	} else {
		// if an attempt was made to publish a topic
		// that no one had subscribed to return an error.
		writer.WriteHeader(http.StatusNotFound)
		writer.Write([]byte("topic not found"))
		return false
	}

}

// subscribe to a topic, initially creating a map of topic and slice of connections with 1 element
// subsequent invocations of subscribe with same topic will grow the slice of connections for the
// same topic key.
func subscribeTopic(writer http.ResponseWriter, request *http.Request, subscribedTopicsMap *map[string][]websocket.Conn, conn websocket.Conn) *map[string][]websocket.Conn {
	query := request.URL.Query()
	topic, present := query["topic"]
	var currentTopic string = topic[0]
	// if no topic was specified return 400
	if !present || len(topic) == 0 {
		writer.WriteHeader(http.StatusBadRequest)
		return subscribedTopicsMap
	} else {
		// otherwise get the slice of connections for the topic
		// whether its empty or not, we append the current connection to it
		// a required improvement would be to make sure that the current connection
		// is not already in the slice of connections before doing the append.
		conns := (*subscribedTopicsMap)[currentTopic]
		(*subscribedTopicsMap)[currentTopic] = append(conns, conn)
		// this was commented out to see if this was an issue on why the connection
		// is being dropped, but it is probably something else.
		//writer.WriteHeader(http.StatusOK)
		//writer.Write([]byte("topic added"))
		return subscribedTopicsMap
	}

}

// invoked by publish to make sure that the topic must exist
// which will cause at least one connection in the topic-->connections map.
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
