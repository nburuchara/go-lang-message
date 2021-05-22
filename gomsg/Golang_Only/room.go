package main

import (
	"net"
)

type room struct {
	name    string
	members map[net.Addr]*client
}

type privateRoom struct {
	name 		string
	members 	map[net.Addr]*client
	password 	string
}

func (r *room) broadcast(sender *client, msg string) {
	for addr, m := range r.members {
		if sender.conn.RemoteAddr() != addr {
			m.msg(msg)
		}
	}
}

func (pr *privateRoom) broadcastPriv(sender *client, msg string) {
	for addr, m := range pr.members {
		if sender.conn.RemoteAddr() != addr {
			m.msg(msg)
		}
	}
}