package main

import (
	"errors"
	"log"
	"strings"

	"github.com/gorilla/websocket"
)

type Chatroom struct {
	Name    string
	Clients map[string]*websocket.Conn
}

type Message struct {
	From    string
	Room    string
	Message string
}

func (self *Chatroom) Send(msg Message) {
	for client, ws := range self.Clients {
		err := ws.WriteJSON(msg)

		if err != nil {
			log.Printf("error: %v", err)
			ws.Close()
			delete(self.Clients, client)
			self.SendChannelMsg(client + " left channel.")
		}
	}
}

func (self *Chatroom) SendChannelMsg(msg string) {
	self.Send(self.ChannelMsg(msg))
}

func (self *Chatroom) ChannelMsg(msg string) Message {
	return Message{From: "#" + self.Name, Message: msg, Room: self.Name}
}

func (self *Chatroom) Users() []string {
	users := make([]string, len(self.Clients))

	i := 0

	for c := range self.Clients {
		users[i] = c
		i++
	}

	return users
}

func joinOrCreateRoom(name string, user string, ws *websocket.Conn) *Chatroom {
	room, ok := rooms[name]

	if !ok {
		room = &Chatroom{Clients: make(map[string]*websocket.Conn), Name: name}
		rooms[name] = room
	}

	_, already_there := room.Clients[user]

	if already_there {
		ws.WriteJSON(room.ChannelMsg("User " + user + " is already in room."))
		return nil
	}

	room.Clients[user] = ws

	return room
}

func leaveRoom(name string, user string) error {
	room, ok := rooms[name]

	if !ok {
		return errors.New("Room " + name + " does not exist.")
	}

	_, ok = room.Clients[user]

	if !ok {
		return errors.New("User " + user + " not present in room " + name + ".")
	}

	delete(room.Clients, user)

	room.SendChannelMsg("User " + user + " left channel.")

	return nil
}

func handleMessages(room *Chatroom, usr string, ws *websocket.Conn) {
	defer ws.Close()

	room.SendChannelMsg("Welcome to " + room.Name + ", " + usr + "!")

	for {
		var msg Message

		err := ws.ReadJSON(&msg)

		if err != nil {
			log.Printf("error %v", err)
			delete(room.Clients, usr)
			room.SendChannelMsg(usr + " left channel.")
			break
		}

		cmd := strings.Split(msg.Message, " ")
		command, ok := COMMANDS[cmd[0]]

		if ok {
			command(room, ws, usr, cmd)
			continue
		}

		room.Send(msg)
	}
}
