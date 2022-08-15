package main

import (
	"fmt"
	"net/http"
)

func main() {
	var subscribedTopics []string
	http.HandleFunc("/", writeResponse)
	http.HandleFunc("/subscribe", func(writer http.ResponseWriter, request *http.Request) {
		subscribedTopics = subscribeTopic(writer, request, &subscribedTopics)
	})

	http.ListenAndServe(":8080", nil)

}

func subscribeTopic(writer http.ResponseWriter, request *http.Request, subscribedTopics *[]string) []string {
	query := request.URL.Query()
	topic, present := query["topic"]
	currentTopic := topic[0]
	if !present || len(topic) == 0 {
		writer.WriteHeader(http.StatusBadRequest)
		return *subscribedTopics
	} else if contains(subscribedTopics, currentTopic) {
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write([]byte("topic exists"))
		return *subscribedTopics
	} else {
		subscribedTopics := append(*subscribedTopics, currentTopic)
		var topic string
		for i := range subscribedTopics {
			topic += subscribedTopics[i]
			topic += "\n"
		}

		writer.WriteHeader(http.StatusOK)
		writer.Write([]byte(topic))

		return subscribedTopics
	}
}
func contains(topics *[]string, newtopic string) bool {
	for _, value := range *topics {
		if value == newtopic {
			return true
		}
	}
	return false
}

func writeResponse(writer http.ResponseWriter, request *http.Request) {
	fmt.Fprintf(writer, "starting server")
}
