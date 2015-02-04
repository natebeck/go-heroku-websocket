package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
)

var connections map[*websocket.Conn]bool

func main() {

	// Start the Web Server
	StartWebServer()

	// count := 0
	// for {
	// 	count++
	// 	sendAll([]byte("Test: " + strconv.Itoa(count)))
	// 	time.Sleep(5000 * time.Millisecond)
	// }
}

func sendAll(msg []byte) {
	for conn := range connections {
		if err := conn.WriteMessage(websocket.TextMessage, msg); err != nil {
			delete(connections, conn)
			conn.Close()
		}
	}
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	// Taken from gorilla's website
	conn, err := websocket.Upgrade(w, r, nil, 1024, 1024)
	if _, ok := err.(websocket.HandshakeError); ok {
		http.Error(w, "Not a websocket handshake", 400)
		return
	} else if err != nil {
		log.Println(err)
		return
	}
	log.Println("Succesfully upgraded connection")
	connections[conn] = true
}

func StartWebServer() {
	// Setup Web Interface
	port := os.Getenv("PORT")

	connections = make(map[*websocket.Conn]bool)

	// handle all requests by serving a file of the same name
	fs := http.Dir("public/")
	fileHandler := http.FileServer(fs)
	http.Handle("/", fileHandler)
	http.HandleFunc("/feed", wsHandler)

	log.Printf("Web Service Running on port %s\n", port)

	addr := fmt.Sprintf(":%s", port)
	// this call blocks -- the progam runs here forever
	err := http.ListenAndServe(addr, nil)
	fmt.Println(err.Error())
}
