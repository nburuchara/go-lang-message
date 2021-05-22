//definitions of commands
package main

type commandID int

const (
	CMD_NICK commandID = iota
	CMD_JOIN
	CMD_ROOMS
	CMD_MSG
	CMD_QUIT
	CMD_OMSG
	CMD_OPMSG
	CMD_JOINP
	CMD_PMSG
	CMD_PASSWORD
)

type command struct {
	id     commandID
	client *client
	args   []string
}


