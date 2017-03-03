package main

import (
	"flag"
	"log"
	"strconv"
	"strings"

	"net/http"

	"github.com/gorilla/websocket"
)

type Chatroom struct {
	Clients map[string]*websocket.Conn
}

type Message struct {
	From    string
	Message string
}

var rooms = make(map[string]Chatroom)
var DEFAULT_ROOM = "general"

func (self *Chatroom) Send(msg Message) {
	for client, ws := range self.Clients {
		err := ws.WriteJSON(msg)

		if err != nil {
			log.Printf("error: %v", err)
			ws.Close()
			delete(self.Clients, client)
		}
	}
}

func validUser(user string) bool {
	return len(user) < 80 && !strings.HasPrefix(user, "#")
}

func makeWsHandler() func(http.ResponseWriter, *http.Request) {
	upgrader := websocket.Upgrader{}
	return func(w http.ResponseWriter, r *http.Request) {
		var usr string
		chatroom := r.URL.Query()["room"]
		username := r.URL.Query()["username"]
		roomname := DEFAULT_ROOM

		if len(chatroom) == 1 {
			roomname = chatroom[0]
		}

		if len(username) == 1 && validUser(username[0]) {
			usr = username[0]
		} else {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Need a valid username (shorter than 80 characters and does not start with a hash)"))
			return
		}

		room, ok := rooms[roomname]

		if !ok {
			rooms[roomname] = Chatroom{Clients: make(map[string]*websocket.Conn)}
			room = rooms[roomname]
		}

		ws, err := upgrader.Upgrade(w, r, nil)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Could not open websocket from connection."))
			return
		}

		defer ws.Close()

		room.Clients[usr] = ws

		room.Send(Message{From: "#channel",
			Message: strings.Join([]string{"Welcome to ", roomname, ", ", usr, "!"}, "")})

		for {
			var msg Message

			err := ws.ReadJSON(&msg)

			if err != nil {
				log.Printf("error %v", err)
				delete(room.Clients, usr)
				break
			}

			room.Send(msg)
		}
	}
}

func main() {
	var port int
	flag.IntVar(&port, "p", 8080, "the port to use")
	flag.Parse()

	fs := http.FileServer(http.Dir("./public"))
	http.Handle("/", fs)

	http.HandleFunc("/ws", makeWsHandler())

	err := http.ListenAndServe(strings.Join([]string{":", strconv.Itoa(port)}, ""),
		nil)

	if err != nil {
		log.Fatal("Server died with message:", err)
	}
}
