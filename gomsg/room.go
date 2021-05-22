package main

import (
	"net"

	// "go.mongodb.org/mongo-driver/bson"
	// "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type roomdatabase struct { //Define what a room is and who is in it
	ID			primitive.ObjectID		`bson:"_id,omitempty"`
	Name    	string					`bson:"name,omitempty"`
	
}
type room struct { //Define what a room is and who is in it
	
	name    	string					
	members 	map[net.Addr]*client
}

type privateRoom struct {
	name 		string
	members 	map[net.Addr]*client
	password 	string
}

type privateroomdatabase struct { //Define what a room is and who is in it
	ID			primitive.ObjectID		`bson:"_id,omitempty"`
	Name    	string					`bson:"name,omitempty"`
	Password 	string					`bson:"password,omitempty"`
}



	//database
	//for databse, collection server: store ids that reference rooms
	
//	Server      primitive.ObjectID 		`bson:"server,omitempty"`
//	Message      primitive.ObjectID 	`bson:"message,omitempty"`
//	Roomname	string			   		


func (r *room) broadcast(sender *client, msg string) { //Send message to all members of the room except sender
	for addr, m := range r.members {
		if addr != sender.conn.RemoteAddr() {
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
