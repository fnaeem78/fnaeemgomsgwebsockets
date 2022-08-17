From the directory fnaeemgomsgwebsockets

* Run the server go build main.go 
* Run the client go build client.go
* The server can then be invoked via curl/postman to post a publish message which the server will handle

Client waits for input of the topic on the command line and
tries to register that topic with the server. 

The server is however closing the connection immediately after
the client registers and hence publishes are not getting 
through. 

At this point it will require deep knowledge of websockets and debugging 
as to why the server is closing the connection. 

