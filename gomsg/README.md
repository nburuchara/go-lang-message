# GoLang web messaging service

gomsg   Go team buruch1,zhou7,hutche1   GoLang web messaging service

for the mongodb part, please do
 go mod init {foldername}

go get go.mongodb.org/mongo-driver


# Usage
1. compile using go build
2. In seperate terminal, use command '%telnet localhost 8888' to connect to the server
3. The list of commands available are as follows:

| Command  | Arguments  | Action
| :---: | :---: | :---: | 
| /nick    | [nickname] | Changes display name of client | 
| /join    | [roomname] | Joins/creates public room | 
| /joinp   | [roomname] [password] | Joins/creates a private room that requires a password | 
| /rooms   |    none    | lists public rooms | 
| /msg     | [message]  | send message to public room you are in | 
| /pmsg    | [message]  | send message to private room you are in | 
| /omsg    |    none    | display old messages from the room | 
| /quit    |    none    | leave room | 