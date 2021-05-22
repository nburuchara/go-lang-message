package main

import (
	"fmt"
	"log"
	"net"
	"strings"
)

type server struct {
	rooms    	map[string]*room
	privateRooms map[string]*privateRoom
	commands 	chan command
}

func newServer() *server {
	return &server{
		rooms:    	make(map[string]*room),
		privateRooms: make(map[string]*privateRoom),
		commands: 	make(chan command),
	}
}

func (s *server) run() {
	for cmd := range s.commands {
		switch cmd.id {
		case CMD_JOINPRIVATE:
			s.joinPrivate(cmd.client, cmd.args)
		case CMD_NICK:
			s.nick(cmd.client, cmd.args)
		case CMD_JOIN:
			s.join(cmd.client, cmd.args)
		case CMD_ROOMS:
			s.listRooms(cmd.client)
		case CMD_MSG:
			s.msg(cmd.client, cmd.args)
		case CMD_QUIT:
			s.quit(cmd.client)
		case CMD_PRIVATEMSG:
			s.privateMsg(cmd.client, cmd.args)
		}
	}
}


/*func (s *server) setPassword(c *client, args []string) {
	if len(args) < 2 {
		c.privateMsg("Password is required. usage /password PASSWORD")
		return
	}

	pword := args[1]
	if(c.privateRoom.password == "") {
		c.privateRoom.password = pword
		c.privateMsg(fmt.Sprintf("password is set to %s", c.privateRoom.password))
	} else if (pword != c.privateRoom.password) {
		c.privateMsg(fmt.Sprintf("Please enter password, %s", c.nick))
	} else if (pword == c.privateRoom.password) {
		c.privateMsg(fmt.Sprintf("Welcome, %s", c.nick))
	}
}*/


func (s *server) joinPrivate(c *client, args []string) {
	if len(args) < 2 {
		c.privateMsg("Room name is required. usage /joinPrivate ROOMNAME PASSWORD")
		return
	}

	roomName := args[1]
	pword := args[2]

	r,ok := s.privateRooms[roomName]
	if !ok {
		r = &privateRoom{
			name:    roomName,
			members: make(map[net.Addr]*client),
			password:pword,
		}
		s.privateRooms[roomName] = r
	} else if (pword != r.password){
		c.privateMsg("Incorrect password, try again using /joinPrivate ROOMNAME PASSWORD")
		return
	} else if (pword == r.password) {
		r.members[c.conn.RemoteAddr()] = c
		s.quitCurrentRoom(c)
		c.privateRoom = r

		r.broadcastPriv(c, fmt.Sprintf("%s joined the room", c.nick))

		c.privateMsg(fmt.Sprintf("welcome to %s", roomName))
		return
	}
	r.members[c.conn.RemoteAddr()] = c
	s.quitCurrentRoom(c)
	c.privateRoom = r

	r.broadcastPriv(c, fmt.Sprintf("%s joined the room", c.nick))

	c.privateMsg(fmt.Sprintf("welcome to %s", roomName))
}


func (s *server) newClient(conn net.Conn) {
	log.Printf("new client has joined: %s", conn.RemoteAddr().String())

	c := &client{
		conn:     conn,
		nick:     "anonymous",
		commands: s.commands,
	}

	c.readInput()
}

func (s *server) nick(c *client, args []string) {
	if len(args) < 2 {
		c.msg("nick is required. usage: /nick NAME")
		return
	}

	c.nick = args[1]
	c.msg(fmt.Sprintf("all right, I will call you %s", c.nick))
}

func (s *server) join(c *client, args []string) {
	if len(args) < 2 {
		c.msg("room name is required. usage: /join ROOM_NAME")
		return
	}

	roomName := args[1]

	r, ok := s.rooms[roomName]
	if !ok {
		r = &room{
			name:    roomName,
			members: make(map[net.Addr]*client),
		}
		s.rooms[roomName] = r
	}
	r.members[c.conn.RemoteAddr()] = c

	s.quitCurrentRoom(c)
	c.room = r

	r.broadcast(c, fmt.Sprintf("%s joined the room", c.nick))

	c.msg(fmt.Sprintf("welcome to %s", roomName))
}

func (s *server) listRooms(c *client) {
	var rooms []string
	for name := range s.rooms {
		rooms = append(rooms, name)
	}

	c.msg(fmt.Sprintf("available rooms: %s", strings.Join(rooms, ", ")))
}

func (s *server) msg(c *client, args []string) {
	if len(args) < 2 {
		c.msg("message is required, usage: /msg MSG")
		return
	}

	msg := strings.Join(args[1:], " ")
	c.room.broadcast(c, c.nick+": "+msg)
}

func (s *server) privateMsg(c *client, args []string) {
	if len(args) < 2 {
		c.privateMsg("message is required, usage: /privateMsg MSG")
		return
	}

	msg := strings.Join(args[1:], " ")
	c.privateRoom.broadcastPriv(c, c.nick+": "+msg)
}

func (s *server) quit(c *client) {
	log.Printf("client has left the chat: %s", c.conn.RemoteAddr().String())

	s.quitCurrentRoom(c)

	c.msg("sad to see you go =(")
	c.conn.Close()
}

func (s *server) quitCurrentRoom(c *client) {
	if c.room != nil {
		oldRoom := s.rooms[c.room.name]
		delete(s.rooms[c.room.name].members, c.conn.RemoteAddr())
		oldRoom.broadcast(c, fmt.Sprintf("%s has left the room", c.nick))
	}
	if c.privateRoom != nil {
		oldRoom := s.privateRooms[c.privateRoom.name]
		delete(s.privateRooms[c.privateRoom.name].members, c.conn.RemoteAddr())
		oldRoom.broadcastPriv(c, fmt.Sprintf("%s has left the room", c.nick))
	}
}