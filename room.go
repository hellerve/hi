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

var rooms = make(map[string]*Chatroom)
var DEFAULT_ROOM = "general"

func (self *Chatroom) Send(msg Message) {
	for client, ws := range self.Clients {
		err := ws.WriteJSON(msg)

		if err != nil {
			log.Printf("error: %v", err)
			ws.Close()
			delete(self.Clients, client)
			self.SendChannelMsg("User " + client + " left channel.")
		}
	}
	destroyIfEmpty(self)
}

func destroyIfEmpty(room *Chatroom) {
	if len(room.Clients) == 0 {
		delete(rooms, room.Name)
	}
}

func (self *Chatroom) Leave(user string) error {
	_, ok := self.Clients[user]

	if !ok {
		return errors.New("User " + user + " not present in room " + self.Name + ".")
	}

	delete(self.Clients, user)

	self.SendChannelMsg("User " + user + " left channel.")

	destroyIfEmpty(self)

	return nil
}

func (self *Chatroom) SendChannelMsg(msg string) {
	self.Send(self.ChannelMsg(msg))
}

func (self *Chatroom) ChannelMsg(msg string) Message {
	return Message{From: "#" + self.Name, Message: msg, Room: self.Name}
}

func SystemMsg(msg string) Message {
	return Message{From: "hi", Message: msg}
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

func RoomNames() []string {
	rs := make([]string, len(rooms))

	i := 0
	for r := range rooms {
		rs[i] = r
		i++
	}

	return rs
}

func joinOrCreateRoom(name string, user string, ws *websocket.Conn) *Chatroom {
	room, ok := rooms[name]

	if !ok {
		room = &Chatroom{Clients: make(map[string]*websocket.Conn), Name: name}
		rooms[name] = room
	}

	_, already_there := room.Clients[user]

	if already_there {
		ws.WriteJSON(SystemMsg("User " + user + " is already in room."))
		return nil
	}

	room.Clients[user] = ws

	room.SendChannelMsg("Welcome to " + room.Name + ", " + user + "!")

	return room
}

func leaveRoom(name string, user string) error {
	room, ok := rooms[name]

	if !ok {
		return errors.New("Room " + name + " does not exist.")
	}
	return room.Leave(user)
}

func handleMessages(usr string, ws *websocket.Conn) {
	defer ws.Close()

	for {
		var msg Message

		err := ws.ReadJSON(&msg)

		if err != nil {
			log.Printf("error %v", err)
			for _, room := range rooms {
				room.Leave(usr)
			}
			break
		}

		room, ok := rooms[msg.Room]

		if !ok {
			ws.WriteJSON(SystemMsg("room does not exist."))
			continue
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
