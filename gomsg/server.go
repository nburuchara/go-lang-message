//Runchang's main task: 
//adding mongodb : 
// here are some collections: 
/*
1. (server)list of room -include: roomname, thread number(that the room is processing on)
		2. room -include: roomname, collection of msg
						3. msg (is an array in the room nickname, timestamp(maybe), content of msg
										4. user-info -include: nickname, IP address(maybe), room()
*/


//bug: right now, we have info both store in local and database, and it may influence,: e.g room with same name? 
//problem: how to deal with msg, store it in the database, may need the help of react par. 
//we need a funtion wich will print out all previous msg in the room

//05/16: 
//todo-adding more command: quit room, create a room/join an existing room, read previous msg
//also in the datebase, msg should include reference of rooms. 
//there is an issue about "ok", maybe we can avoid this issue by adding a new command


package main

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"net"
//	"math/rand"
	"time"
	// "regexp"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

//	"go.mongodb.org/mongo-driver/mongo"
)


type server struct {
	rooms    map[string]*room
	privateRooms map[string]*privateRoom
	commands chan command

	//for databse, collection server: store ids that reference rooms
	//maybe we do not need server collection
	// is the parent of room and user
//	ID			primitive.ObjectID		`bson:"_id,omitempty"`

}




func newServer() *server { //server builder
	return &server{
		rooms:    make(map[string]*room),
		privateRooms: make(map[string]*privateRoom),
		commands: make(chan command),
	}
}

func (s *server) run() { //constant function waiting for commands
	for cmd := range s.commands {
		switch cmd.id {
		case CMD_JOINP:
			s.joinPrivate(cmd.client, cmd.args)
		// case CMD_PASSWORD:
		// 	s.setPassword(cmd.client, cmd.args)
		case CMD_NICK:
			s.nick(cmd.client, cmd.args)
		case CMD_JOIN:
			s.join(cmd.client, cmd.args)
		case CMD_ROOMS:
			s.listrooms(cmd.client, cmd.args)
		case CMD_MSG:
			s.msg(cmd.client, cmd.args)
		case CMD_OMSG: 				//read previous msg 
			s.omsg(cmd.client, cmd.args)	
		case CMD_OPMSG: 				//read previous msg in private room 
		s.opmsg(cmd.client, cmd.args)
		case CMD_QUIT:
			s.quit(cmd.client, cmd.args)
		case CMD_PMSG:
			s.privateMsg(cmd.client, cmd.args)
		}
	}
}

func makeTimestamp() int64 {
    return time.Now().UnixNano() / (int64(time.Millisecond)/int64(time.Nanosecond))
}

// func (s *server) setPassword(c *client, args []string) {
// 	if len(args) < 2 {
// 		c.privateMsg("Password is required. usage /password PASSWORD")
// 		return
// 	}

// 	pword := args[1]
// 	if(c.privateRoom.password == "") {
// 		c.privateRoom.password = pword
// 		c.privateMsg(fmt.Sprintf("password is set to %s", c.privateRoom.password))
// 	} else if (pword != c.privateRoom.password) {
// 		c.privateMsg(fmt.Sprintf("incorrect password, %s, try again using /password PASSWORD", c.nick))
// 	} else if (pword == c.privateRoom.password) {
// 		c.privateMsg(fmt.Sprintf("Welcome, %s", c.nick))
// 	}
// }


func (s *server) joinPrivate(c *client, args []string) {
	if len(args) < 3 {
		c.privateMsg("Room name and password is required. usage /joinp ROOMNAME PASSWORD")
		return
	}

	allroom:= Getprivateroomslist()

	for _, name := range allroom {
		strname := name.(string)
		r2,found := s.privateRooms[strname]	
		if !found{
			r2 = &privateRoom{
				name:    strname,
				members: make(map[net.Addr]*client),
			}
			s.privateRooms[strname] = r2
		}
	}
	roomName := args[1]
	pword := args[2]

	ok2 := Check3(roomName)


	r,ok := s.privateRooms[roomName]
	fmt.Println(ok)
	fmt.Println(ok2)

	if !ok2 {
		r = &privateRoom{
			name:    roomName,
			members: make(map[net.Addr]*client),
			password:pword,
		}
		s.privateRooms[roomName] = r
		privateroomdb := privateroomdatabase{
			Name:	roomName,
			Password:	pword,
		}
		err := Addprivateroom(privateroomdb)
		if err != nil {
			fmt.Println(err)
			return
		}

	} else if (pword != r.password){
		c.privateMsg("Incorrect password, try again using /joinp ROOMNAME PASSWORD")
		return
	} else if (pword == r.password) {
		//user joins a private room in the database
		r.members[c.conn.RemoteAddr()] = c
		s.quitCurrentRoom(c)
		s.quitCurrentPrivateRoom(c)

		c.privateRoom = r
		//use roomname to find object id of room

		privateroomdb := privateroomdatabase{
			Name:	roomName,
		}
		rnid :=Findprivateroomid(privateroomdb)

		myroomid := rnid.(primitive.ObjectID)

		//upadate user, add room reference to the user. 
		userdb := userdatabase{
			//right now the user is not in any room, and nick is "anonym"
			Nickname: c.nick, 
			Number: c.idnumber,
			Roomname: myroomid, //question: should it be reference of the room(object id?) or just the roomname
		}
		err4 := Updateuser_room(userdb)
		if err4 != nil {
			fmt.Println(err4)
			return
		}
		r.broadcastPriv(c, fmt.Sprintf("%s joined the private room", c.nick))

		c.privateMsg(fmt.Sprintf("welcome to the private room %s", roomName))
		return
	}
	r.members[c.conn.RemoteAddr()] = c
	s.quitCurrentRoom(c)
	s.quitCurrentPrivateRoom(c)
	c.privateRoom = r

	privateroomdb := privateroomdatabase{
		Name:	roomName,
	}
	rnid :=Findprivateroomid(privateroomdb)

	myroomid := rnid.(primitive.ObjectID)

			//upadate user, add room reference to the user. 
	userdb := userdatabase{
		//right now the user is not in any room, and nick is "anonym"
		Nickname: c.nick, 
		Number: c.idnumber,
		Roomname: myroomid, //question: should it be reference of the room(object id?) or just the roomname
	}
	err4 := Updateuser_room(userdb)
	if err4 != nil {
		fmt.Println(err4)
		return
	}

	r.broadcastPriv(c, fmt.Sprintf("%s joined the private room", c.nick))

	c.privateMsg(fmt.Sprintf("welcome to the private room %s", roomName))
}




func (s *server) newClient(conn net.Conn) { //Server states that a new client has joined the server and nicknames them anonym
	//TO-DO: create a collection for user, maybe include the IP address of the user. (check client.go)
	//usernum := rand.Intn(100000)

	//usernum:=time.Now().UnixNano()/(1<<22)
	usernum := makeTimestamp()
	log.Printf("new client has connected: %s, userid = %d", conn.RemoteAddr().String(),usernum)
	//conn.RemoteAddr().String() is the IP address. 
	
	c := &client{
		conn:     conn,
		nick:     "anonym",
		idnumber: usernum,
		commands: s.commands,
	}
		//create a new user to the collection

		userdb := userdatabase{
			//right now the user is not in any room, and nick is "anonym"
			Nickname: "anonym", 
			Number: usernum, 
			IP:	conn.RemoteAddr().String(),
		}
		//print error if there is any
		err := Adduser(userdb)
		if err != nil {
			fmt.Println(err)
			return
		}


	c.readInput()
	//store the input in the User collection. 

}

func (s *server) nick(c *client, args []string) { //User can change their nickname
	c.nick = args[1]
	
	c.msg(fmt.Sprintf("all right,user = %d , I will call you %s", c.idnumber ,c.nick))
//	c.msg(fmt.Sprintf("all right, I  know  you %s", connstringIP))

	//here change the nick of the user. 
	//update nickname in the userdateabase collection 

	userdb := userdatabase{
		//right now the user is not in any room, and nick is "anonym"
		Nickname: c.nick, 
		Number: c.idnumber,
	}
	//just for testing
	//fmt.Sprintf("the roomname is %s", roomName) 

	//print error if there is any
	err := Updateuser_nick(userdb)
	if err != nil {
		fmt.Println(err)
		return
	}

}

func (s *server) join(c *client, args []string) { //user can join a chat room, if it doesnt exist it is created
	if len(args) < 2 {
		c.msg("room name is required. usage: /join ROOM_NAME")
		return
	}

	//read all rooms from the database
	allroom:= Getroomslist()

	for _, name := range allroom {
		strname := name.(string)
		r2,found := s.rooms[strname]	
		if !found{
			r2 = &room{
				name:    strname,
				members: make(map[net.Addr]*client),
			}
			s.rooms[strname] = r2
		}
	}
	
	//TO-DO: add a user to the room list, also fetech info inside of the room. 
	roomName := args[1]
	r, ok := s.rooms[roomName]


	// roomdb := roomdatabase{
	// 	Name:	roomName,
	// }
	// ok2 := Checkthisroomname(roomName)
	ok2 := Check2(roomName)
	// try := Getthisroomname
	// fmt.Println(try)


	//listallroom := Getroomslist()
	//check := Checkroomname(listallroom, roomName)

	fmt.Println(ok)
	fmt.Println(ok2)
	//(bug)ok should check the database, instead of the server. 
	if !ok2 {
		r = &room{
			name:    roomName,
			members: make(map[net.Addr]*client),
		}
		s.rooms[roomName] = r
		//database sample input
		// { 
		// 	_id: "ro"
		// 	server_id: "webechat", // reference to server document 
		// 	name: "room1"
		//  }
		//create a new room to the collection
		roomdb := roomdatabase{
			Name:	roomName,
		}
		//just for testing
		//fmt.Sprintf("the roomname is %s", roomName) 

		//print error if there is any
		err := Addroom(roomdb)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
	
	r.members[c.conn.RemoteAddr()] = c

	s.quitCurrentRoom(c)

	c.room = r

//use roomname to find object id of room

roomdb := roomdatabase{
	Name:	roomName,
}
rnid :=Findroomid(roomdb)

myroomid := rnid.(primitive.ObjectID)
// rex := regexp.MustCompile(`\"([^)]+)\"`)
// out := rex.FindAllStringSubmatch(myroom, -1)
// for _, i := range out {
//     fmt.Println(i[1])
// }


// fmt.Println("rn is equal to")
// fmt.Println(myroomid)

// id := "60a1f60c41c0f629c000acb5"

// 	docID, errid := primitive.ObjectIDFromHex(id)
// 	if errid != nil {
// 		fmt.Println(errid)
// 		return
// 	}


//upadate user, add room reference to the user. 
	userdb := userdatabase{
		//right now the user is not in any room, and nick is "anonym"
		Nickname: c.nick, 
		Number: c.idnumber,
		Roomname: myroomid, //question: should it be reference of the room(object id?) or just the roomname
	}
	//just for testing
	//fmt.Sprintf("the roomname is %s", roomName) 

	//print error if there is any
	err4 := Updateuser_room(userdb)
	if err4 != nil {
		fmt.Println(err4)
		return
	}

	r.broadcast(c, fmt.Sprintf("%s has joined the room", c.nick))
	c.msg(fmt.Sprintf("welcome to %s", r.name)) 
	
}

func (s *server) listrooms(c *client, args []string) { //lists all current rooms
	//here should it be string or bson? 
	//TO-DO: call and fetch from the database, fetch the list of room. 

	// var rooms []string
	// for name := range s.rooms {
	// 	rooms = append(rooms, name)
	// }
	
	listallroom := Getroomslist()

	c.msg(fmt.Sprintf("available rooms are: %v", listallroom))
}
func (s *server) listprivaterooms(c *client, args []string) { //lists all current privaterooms
	//here should it be string or bson? 
	//TO-DO: call and fetch from the database, fetch the list of privateroom. 

	// var privaterooms []string
	// for name := range s.privaterooms {
	// 	privaterooms = append(privaterooms, name)
	// }
	
	listallprivateroom := Getprivateroomslist()

	c.privateMsg(fmt.Sprintf("available privaterooms are: %v", listallprivateroom))
}

func (s *server) omsg(c *client, args []string) {
	if c.room == nil {
		c.err(errors.New("you must join the room first"))
		return
	}
	userdb := userdatabase{
		Number:	c.idnumber,
	}
	myroom :=Findroom(userdb).(primitive.ObjectID)
	// listallmsg := Getallmsg()

	msgdb := msgdatabase{
		Room:	myroom,
	}
	listallmsg := Getroommsg(msgdb)
	c.msg(fmt.Sprintf("here are old message:\n  %v", listallmsg))

}
func (s *server) opmsg(c *client, args []string) {
	if c.privateRoom == nil {
		c.err(errors.New("you must join a private room first"))
		return
	}
	userdb := userdatabase{
		Number:	c.idnumber,
	}
	myroom :=Findroom(userdb).(primitive.ObjectID)
	// listallmsg := Getallmsg()

	privatemsgdb := privatemsgdatabase{
		Proom:	myroom,
	}
	listallmsg := Getroomprivatemsg(privatemsgdb)
	c.privateMsg(fmt.Sprintf("here are old message in this private room:\n  %v", listallmsg))

}

func (s *server) msg(c *client, args []string) { //sends a message to the other clients in a room
	if c.room == nil {
		c.err(errors.New("you must join the room first"))
		return
	}
	//TO-DO: store the msg in the database, along with info{nickname, msg}:, msg should be a part of the room list
	//should be a forloop args
	//create a new room to the collection

	// for  index, entered := range args {
	// 	fmt.Println("sentence number", index, " saying", entered)
	// }


	words := c.nick +": " + strings.Join(args[1:len(args)]," ") + "\n "

	fmt.Println(words)

	//read the roomname from user

	userdb := userdatabase{
		Number: c.idnumber,
	}


	myroom := Findroom(userdb).(primitive.ObjectID)


	// fmt.Println("rn is equal to")
	fmt.Println(myroom)

	// id := "60a1f60c41c0f629c000acb5"

	// docID, err := primitive.ObjectIDFromHex(id)

	msgdb := msgdatabase{
		Content:	words,
		Room: 		myroom,
	}
	//print error if there is any
	err2 := Addmsg(msgdb)
	if err2 != nil {
		fmt.Println(err2)
		return
	}
	c.room.broadcast(c, words)
}	

func (s *server) privateMsg(c *client, args []string) {
	if len(args) < 2 {
		c.privateMsg("message is required, usage: /privateMsg MSG")
		return
	}
	if c.privateRoom == nil {
		c.err(errors.New("you must join the private room first"))
		return
	}
	pmsg := c.nick +": " + strings.Join(args[1:len(args)]," ") + "\n "
	userdb := userdatabase{
		Number: c.idnumber,
	}
	myroom := Findroom(userdb).(primitive.ObjectID)
	fmt.Println(myroom)
	privatemsgdb := privatemsgdatabase{
		Content:	pmsg,
		Proom: 		myroom,
	}
	//print error if there is any
	err2 := Addprivatemsg(privatemsgdb)
	if err2 != nil {
		fmt.Println(err2)
		return
	}
	c.privateRoom.broadcastPriv(c, pmsg)
}


func (s *server) quit(c *client, args []string) { //quits the room and then the server
	log.Printf("client has disconnected: %s", c.conn.RemoteAddr().String())

	s.quitCurrentRoom(c)

	c.msg("sad to see you go")
	//let just delete the user once they leave
	userdb := userdatabase{
		//right now the user is not in any room, and nick is "anonym"
		Number: c.idnumber,
	}
	//just for testing
	//fmt.Sprintf("the roomname is %s", roomName) 

	//print error if there is any
	err := Removeuser(userdb)
	if err != nil {
		fmt.Println(err)
		return
	}
	//check if there any user in the server, if not remove all rooms
	// count := Countuser()
	// if count == 0 {
	// 	Removeallrooms()
	// }

//	c.conn.Close()
}

func (s *server) quitCurrentRoom(c *client) { //leave the current room
	if c.room != nil {
		delete(c.room.members, c.conn.RemoteAddr())
		c.room.broadcast(c, fmt.Sprintf("%s has left the room", c.nick))
	}
	//TO-DO: when the user quit, delte the user from the the room. 
	//it should have a list of room, and for each room there is a list of active uers in the room

}
func (s *server) quitCurrentPrivateRoom(c *client) { //leave the current room
	if c.room != nil {
		delete(c.privateRoom.members, c.conn.RemoteAddr())
		c.privateRoom.broadcastPriv(c, fmt.Sprintf("%s has left the private room", c.nick))
	}
	//TO-DO: when the user quit, delte the user from the the room. 
	//it should have a list of room, and for each room there is a list of active uers in the room

}




//database functions 
//Addroom creates a new room to the collection.
func Addroom(roomdb roomdatabase) error {
    coll := conn.Database(DBName).Collection(roomsCollection)
    roomdb.ID = primitive.NewObjectID()
    insertResult, err := coll.InsertOne(ctx, roomdb)
    if err != nil {
        fmt.Printf("Could not create new room. Id: %s\n", roomdb.ID)
        return err
    }
    fmt.Printf("Created new room. ID: %s\n", insertResult.InsertedID)
    return nil
}
//addprivate room 
func Addprivateroom(privateroomdb privateroomdatabase) error {
    coll := conn.Database(DBName).Collection(privateroomsCollection)
    privateroomdb.ID = primitive.NewObjectID()
    insertResult, err := coll.InsertOne(ctx, privateroomdb)
    if err != nil {
        fmt.Printf("Could not create new privateroom. Id: %s\n", privateroomdb.ID)
        return err
    }
    fmt.Printf("Created new privateroom. ID: %s\n", insertResult.InsertedID)
    return nil
}

//Adduser creates a new user to the collection.
func Adduser(userdb userdatabase) error {
    coll := conn.Database(DBName).Collection(usersCollection)
    userdb.ID = primitive.NewObjectID()
    insertResult, err := coll.InsertOne(ctx, userdb)
    if err != nil {
        fmt.Printf("Could not create new user. Id: %s\n", userdb.ID)
        return err
    }
    fmt.Printf("Created new user. ID: %s\n", insertResult.InsertedID)
    return nil
}

//Addmsg creates a new msg to the collection.
func Addmsg(msgdb msgdatabase) error {
    coll := conn.Database(DBName).Collection(msgsCollection)
    msgdb.ID = primitive.NewObjectID()
    insertResult, err := coll.InsertOne(ctx, msgdb)
    if err != nil {
        fmt.Printf("Could not create new msg. Id: %s\n", msgdb.ID)
        return err
    }
    fmt.Printf("Created new msg. ID: %s\n", insertResult.InsertedID)
    return nil
}
func Addprivatemsg(privatemsgdb privatemsgdatabase) error {
    coll := conn.Database(DBName).Collection(privatemsgsCollection)
    privatemsgdb.ID = primitive.NewObjectID()
    insertResult, err := coll.InsertOne(ctx, privatemsgdb)
    if err != nil {
        fmt.Printf("Could not create new privatemsg. Id: %s\n", privatemsgdb.ID)
        return err
    }
    fmt.Printf("Created new privatemsg. ID: %s\n", insertResult.InsertedID)
    return nil
}




//update user's nickname.
func Updateuser_nick(userdb userdatabase) error {

    coll := conn.Database(DBName).Collection(usersCollection)


	filter := bson.D{{"number", userdb.Number}}
    update := bson.D{{"$set",
        bson.D{
            {"nickname", userdb.Nickname},
        },
    }}

    nickResult, err := coll.UpdateOne(ctx, filter,update)
    if err != nil {
        fmt.Printf("Could not update the user's nickname. Id: %s\n", userdb.ID)
        return err
    }
	fmt.Printf("Updated %v Documents!\n", nickResult.ModifiedCount)
    return nil
}

//update user's roomname.
func Updateuser_room(userdb userdatabase) error {

    coll := conn.Database(DBName).Collection(usersCollection)

	filter := bson.D{{"number", userdb.Number}}
    update := bson.D{{"$set",
        bson.D{
            {"roomname", userdb.Roomname},
        },
    }}

    roomResult, err := coll.UpdateOne(ctx, filter,update)
    if err != nil {
        fmt.Printf("Could not update the user's room. Id: %s\n", userdb.ID)
        return err
    }
	fmt.Printf("Updated %v Documents!\n", roomResult.ModifiedCount)
    return nil
}
//read and check to see if we have such room in the database, search by the roomname
// func Readroomname(roomdb roomdatabase) error {
// 	coll := conn.Database(DBName).Collection(roomsCollection)

// 	insertResult, err := coll.InsertOne(ctx, roomdb)
// 	if err != nil {
// 		fmt.Printf("Could not create new room. Id: %s\n", roomdb.ID)
// 		return err
// 	}
// 	fmt.Printf("Crokeated new room. ID: %s\n", insertResult.InsertedID)
// 	return nil
// }

func Check2(str string) bool{
	listallroom := Getroomslist()

		for index, value := range listallroom {
			if value == str {
				fmt.Println(index)
				fmt.Println(value)
			 return true
			
			}
  
		}
	return false
}
func Check3(str string) bool{
	listallroom := Getprivateroomslist()

		for index, value := range listallroom {
			if value == str {
				fmt.Println(index)
				fmt.Println(value)
			 return true
			
			}
  
		}
	return false
}
// func Checkthisroomname(str string) bool{


// 	filter := bson.D{{"name", str}}
// 	coll := conn.Database(DBName).Collection(roomsCollection)
// 	err := coll.FindOne(ctx, filter)

// 	if err != nil {
// 		fmt.Println(nil)
// 		fmt.Println(err)
// 		// ErrNoDocuments means that the filter did not match any documents in the collection
// 		// fmt.Printf("Could not find the room in checkthisroomname. Id: %s\n", roomdb.ID)
// 		return false
// 	}
// 	fmt.Printf("found document in checkthisroomname" )	
// 	return true
// }
// // func Getthisroomname(roomdb roomdatabase) string{
// 	coll := conn.Database(DBName).Collection(roomsCollection)
// 	err := coll.FindOne(ctx, roomdb)

// 	if err != nil {
// 		// ErrNoDocuments means that the filter did not match any documents in the collection
// 		fmt.Printf("Could not find the room. Id: %s\n", roomdb.ID)
// 		log.Fatal(err)
	
// 	}
// 	id := roomdb.ID.Hex() 
	
// 	return id
// }
// func Getthisroomname(roomdb roomdatabase) string{

// 	coll := conn.Database(DBName).Collection(roomsCollection)
// 	err := coll.FindOne(ctx, roomdb)
// 	id := roomdb.ID.string() 
// 	if err != nil {
// 		// ErrNoDocuments means that the filter did not match any documents in the collection
// 		fmt.Printf("Could not find the room. Id: %s\n", roomdb.ID)
// 		log.Fatal(err)
		 
// 		return id
// 	}

// 	return id
// }






func Findroom(userdb userdatabase) interface {} {
	coll := conn.Database(DBName).Collection(usersCollection)

	var result bson.M
	err := coll.FindOne(ctx, userdb).Decode(&result)
	if err != nil {
		fmt.Printf("Could not find the user. Id: %s\n", userdb.ID)
       // return err
	}
	//fmt.Printf("found document in findroom%v", result)
	// for key, value := range result {
	// 	strvale:=value.(string)

	// 	// if (value == "roomname"){
	// 	// // value == "roomname"{
	// 	  fmt.Println(strvale)
	// 	// }
	// 	fmt.Println("no room name")
	// }

	for index, s := range result {
		if index == "roomname"{
			rn := s
			//fmt.Println(s)
			return rn
		}
	}
	return err
	// var raw bson.Raw = result
	// err := raw.Validate()
	// if err != nil { return err }
	// rn := raw.Lookup("roomname")
	// for key, value := range result {
	// 	if key == "Roomname"{
	// 	  listall := value
	// 	  return listall
	// 	  //fmt.Println(value)
	// 	}

	// }
 	// fmt.Println("wrong")
	// return userdb.ID
}
//find the room object id by roomname 
func Findroomid(roomdb roomdatabase) interface {} {

	coll := conn.Database(DBName).Collection(roomsCollection)

	var result bson.M
	err := coll.FindOne(ctx, roomdb).Decode(&result)
	if err != nil {
		fmt.Printf("Could not find the user. Id: %s\n", roomdb.ID)
       // return err
	}

	for index, s := range result {
		if index == "_id"{
			rn := s
			//fmt.Println(s)
			return rn
		}
	}
	return err
}

//find the privateroom object id by privateroomname 
func Findprivateroomid(privateroomdb privateroomdatabase) interface {} {

	coll := conn.Database(DBName).Collection(privateroomsCollection)

	var result bson.M
	err := coll.FindOne(ctx, privateroomdb).Decode(&result)
	if err != nil {
		fmt.Printf("Could not find the user. Id: %s\n", privateroomdb.ID)
       // return err
	}

	for index, s := range result {
		if index == "_id"{
			rn := s
			//fmt.Println(s)
			return rn
		}
	}
	return err
}


//remove the user once they quit
func Removeuser(userdb userdatabase)error{
	coll := conn.Database(DBName).Collection(usersCollection)

	filter := bson.D{{"number", userdb.Number}}

	result, err := coll.DeleteOne(ctx, filter)
	if err != nil {
		fmt.Printf("Could not delete the user. Id: %s\n", userdb.ID)
        return err
	}
	fmt.Printf("DeleteOne removed %v document(s)\n", result.DeletedCount)
	return nil
}

//count how many user in the database
func Countuser() int64{
	coll := conn.Database(DBName).Collection(usersCollection)
	itemCount, err := coll.CountDocuments(ctx, bson.D{})
	if err != nil {
		log.Fatal(err)
	}
	return itemCount
}

//remove all user 
func Removeallrooms() error {
	coll := conn.Database(DBName).Collection(roomsCollection)

	err := coll.Drop(ctx)
	if err != nil {
		return err
	}
	return nil
}
//fetch the whole collection of room 

func Getroomslist()([]interface{}) {
    coll := conn.Database(DBName).Collection(roomsCollection)
    
	cursor, err := coll.Find(ctx, bson.D{})
	if err != nil {
		log.Fatal(err)
	}
	
	var roomlist []bson.M
	if err = cursor.All(ctx, &roomlist); err != nil {
		fmt.Printf("Could list all rooms.")
		//return err
	}
	// fmt.Println(roomlist)
	
	var listall []interface{}
	// err := c.Find(bson.M{}).All(&data)
	// handle err 
	for _, doc := range roomlist {
	  for key, value := range doc {
		  if key == "name"{
			listall = append(listall, value)
			//fmt.Println(value)
		  }

	  }
	}  
	return listall 


}
//read all message 
func Getallmsg()([]interface{}) {
    coll := conn.Database(DBName).Collection(msgsCollection)
    
	cursor, err := coll.Find(ctx, bson.D{})
	if err != nil {
		log.Fatal(err)
	}
	
	var msglist []bson.M
	if err = cursor.All(ctx, &msglist); err != nil {
		fmt.Printf("Could list all msgs.")
		//return err
	}
	
	var listall []interface{}
	// err := c.Find(bson.M{}).All(&data)
	// handle err 
	for _, doc := range msglist {
	  for key, value := range doc {
		  if key == "content"{
			listall = append(listall, value)
			// fmt.Println(value)
		  }

	  }
	}  
	// fmt.Println("check mark here ")
	// fmt.Println(listall)
	
	return listall 
}

//read msg in this room. 
func Getroommsg(msgdb msgdatabase) []interface {} {
	coll := conn.Database(DBName).Collection(msgsCollection)

	filterCursor, err := coll.Find(ctx, msgdb)
	if err != nil {
		log.Fatal(err)
	}
	var msgsFiltered []bson.M
	if err = filterCursor.All(ctx, &msgsFiltered); err != nil {
		fmt.Printf("Could list all msgs in this room.")
	}
	//fmt.Println(msgsFiltered)

	var listall []interface{}
	// err := c.Find(bson.M{}).All(&data)
	// handle err 
	for _, doc := range msgsFiltered {
	  for key, value := range doc {
		  if key == "content"{
			listall = append(listall, value)
			// fmt.Println(value)
		  }

	  }
	}  
	return listall
	
}

//list all private room
func Getprivateroomslist()([]interface{}) {
    coll := conn.Database(DBName).Collection(privateroomsCollection)
    
	cursor, err := coll.Find(ctx, bson.D{})
	if err != nil {
		log.Fatal(err)
	}
	
	var privateroomlist []bson.M
	if err = cursor.All(ctx, &privateroomlist); err != nil {
		fmt.Printf("Could list all privaterooms.")
		//return err
	}
	// fmt.Println(privateroomlist)
	
	var listall []interface{}
	// err := c.Find(bson.M{}).All(&data)
	// handle err 
	for _, doc := range privateroomlist {
	  for key, value := range doc {
		  if key == "name"{
			listall = append(listall, value)
			//fmt.Println(value)
		  }

	  }
	}  
	return listall 
}
//read msg in this room. 
func Getroomprivatemsg(privatemsgdb privatemsgdatabase) []interface {} {
	coll := conn.Database(DBName).Collection(privatemsgsCollection)

	filterCursor, err := coll.Find(ctx, privatemsgdb)
	if err != nil {
		log.Fatal(err)
	}
	var privatemsgsFiltered []bson.M
	if err = filterCursor.All(ctx, &privatemsgsFiltered); err != nil {
		fmt.Printf("Could list all privatemsgs in this privateroom.")
	}
	//fmt.Println(privatemsgsFiltered)

	var listall []interface{}
	// err := c.Find(bson.M{}).All(&data)
	// handle err 
	for _, doc := range privatemsgsFiltered {
	  for key, value := range doc {
		  if key == "content"{
			listall = append(listall, value)
			// fmt.Println(value)
		  }

	  }
	}  
	return listall
	
}