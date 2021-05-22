
//TO-DO: create a collection for user, maybe include the IP address of the user. (check client.go)

package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"

	// "go.mongodb.org/mongo-driver/bson"
	// "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type client struct {
	conn     net.Conn
	nick     string
	idnumber int64
	room     *room
	privateRoom *privateRoom
	commands chan<- command
}

type userdatabase struct { //Define what a user is and who is in it
	ID			primitive.ObjectID		`bson:"_id,omitempty"`
	//database
	Number		int64					`bson:"number,omitempty"`
	Roomname    primitive.ObjectID 		`bson:"roomname,omitempty"`
	Nickname    string 			 		`bson:"nickname,omitempty"`
	IP			string			   		`bson:"ip,omitempty"`
}


func (c *client) readInput() {
	for {
		msg, err := bufio.NewReader(c.conn).ReadString('\n')
		if err != nil {
			return
		}

		msg = strings.Trim(msg, "\r\n")

		args := strings.Split(msg, " ")   //breaks the input into chunks as an id, the client, and the arguments
		cmd := strings.TrimSpace(args[0]) //finds what the /cmd is

		switch cmd { //define the different parts of each command and store it in thread safe channel for the server
		case "/nick":
			c.commands <- command{
				id:     CMD_NICK,
				client: c,
				args:   args,
			}
		case "/joinp":
			c.commands <- command{
				id:		CMD_JOINP,
				client: c,
				args:	args,
			}
		// case "/password":
		// 	c.commands <- command{
		// 		id:		CMD_PASSWORD,
		// 		client:	c,
		// 		args:	args,
		// 	}
		case "/join":
			c.commands <- command{
				id:     CMD_JOIN,
				client: c,
				args:   args,
			}
		case "/rooms":
			c.commands <- command{
				id:     CMD_ROOMS,
				client: c,
			}
		case "/omsg":
			c.commands <- command{
				id:     CMD_OMSG,
				client: c,
			}	
		case "/opmsg":
			c.commands <- command{
				id:     CMD_OPMSG,
				client: c,
			}		
		case "/msg":
			c.commands <- command{
				id:     CMD_MSG,
				client: c,
				args:   args,
			}
		case "/pmsg":
			c.commands <- command{
				id:		CMD_PMSG,
				client:	c,
				args:	args,
			}
		case "/quit":
			c.commands <- command{
				id:     CMD_QUIT,
				client: c,
			}
		default:
			c.err(fmt.Errorf("unknown command: %s", cmd))
		}
	}
}

func (c *client) err(err error) {
	c.conn.Write([]byte("err: " + err.Error() + "\n"))
}

func (c *client) msg(msg string) { //print message to user
	c.conn.Write([]byte("> " + msg + "\n"))
}
func (c *client) privateMsg(msg string) {
	c.conn.Write([]byte("> " + msg + "\n"))
}
