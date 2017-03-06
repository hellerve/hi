package main

import (
	"strings"

	"github.com/gorilla/websocket"
)

type Command func(*Chatroom, *websocket.Conn, string, []string)

var COMMANDS map[string]Command

func listUsers(room *Chatroom, ws *websocket.Conn, user string, args []string) {
	ws.WriteJSON(room.ChannelMsg(strings.Join(room.Users(), ", ")))
}

func join(room *Chatroom, ws *websocket.Conn, user string, args []string) {
	if len(args) != 2 {
		ws.WriteJSON(room.ChannelMsg(WrongArgs(args)))
		return
	}

	roomname := args[1]

	joinOrCreateRoom(roomname, user, ws)
}

func leave(room *Chatroom, ws *websocket.Conn, user string, args []string) {
	if len(args) != 2 {
		ws.WriteJSON(room.ChannelMsg(WrongArgs(args)))
		return
	}

	roomname := args[1]

	err := leaveRoom(roomname, user)

	if err != nil {
		ws.WriteJSON(SystemMsg("Couldn't leave channel " + roomname + ": " + err.Error()))
	}

	ws.WriteJSON(SystemMsg("Left channel " + roomname + "."))
}

func WrongArgs(args []string) string {
	return "Malformatted arguments: " + strings.Join(args, " ")
}

func init() {
	// this avoids an initialization loop with joining
	COMMANDS = map[string]Command{
		"/list":  listUsers,
		"/join":  join,
		"/leave": leave,
	}
}
