package main

type commandID int

const (
	CMD_NICK commandID = iota
	CMD_JOIN
	CMD_ROOMS
	CMD_MSG
	CMD_QUIT
	CMD_PASSWORD
	CMD_JOINPRIVATE
	CMD_PRIVATEMSG
)

type command struct {
	id     commandID
	client *client
	args   []string
}