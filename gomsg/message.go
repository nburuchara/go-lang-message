package main

import (
	// "bufio"
	// "fmt"
	// "net"
	// "strings"

	// "go.mongodb.org/mongo-driver/bson"
	// "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson/primitive"
)


type msgdatabase struct { //Define what a user is and who is in it
	ID			primitive.ObjectID		`bson:"_id,omitempty"`
	Content	    string 			 		`bson:"content,omitempty"`
	Room      	primitive.ObjectID		`bson:"room,omitempty"`
	
	//also here should include reference of room. 
	//here the content incude the nickname ... 
	//timestamp
}
type privatemsgdatabase struct { //Define what a user is and who is in it
	ID			primitive.ObjectID		`bson:"_id,omitempty"`
	Content	    string 			 		`bson:"content,omitempty"`
	Proom      	primitive.ObjectID		`bson:"room,omitempty"` //private room
	
	//also here should include reference of room. 
	//here the content incude the nickname ... 
	//timestamp
}
