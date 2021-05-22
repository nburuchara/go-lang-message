package main

import (
	"log"
	"net"
//	"context"
	"fmt"
//	"time"
	// "go.mongodb.org/mongo-driver/bson"
	// "go.mongodb.org/mongo-driver/mongo"
	// "go.mongodb.org/mongo-driver/mongo/options"
	// "go.mongodb.org/mongo-driver/mongo/readpref"
)

// func returnIP(conn net.Conn) string{
// 	var strRemoteAddr string =  conn.RemoteAddr().String()
// 	return strRemoteAddr
// }

// var(
// 	ipaddress := returnIP
// )

func main() {
	s := newServer()
	go s.run() //waits for commands on a different thread

	listener, err := net.Listen("tcp", ":8888")
	if err != nil {
		log.Fatalf("unable to start serverL: %s", err.Error())
	}

	defer listener.Close()
	log.Printf("started server on :8888")
	
	//connect the database
	if err := createDBSession(); err != nil {
        fmt.Println(err)
        return
    }
    fmt.Println("Connected")  // Connected

/* //replace with mongo.go

	//add mongodb connection 
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb+srv://user_zhou:AB1Ck2Wg5MS3MfCq@cluster0.rktlj.mongodb.net/myFirstDatabase?retryWrites=true&w=majority"))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}

	//buid a database for gomsg
	database := client.Database("webmsg")
	//serversCollection := database.Collection("servers")
	roomsCollection := database.Collection("rooms")
	//messagesCollection := database.Collection("messages")
	usersCollection := database.Collection("users")


	//print the database name: for testing concern 
	databases, err := client.ListDatabaseNames(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(databases)
 */


	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("unable to accept connection: %s", err.Error())
			continue
		}

		go s.newClient(conn) // Creates new client on different thread
	}

	
	
}
