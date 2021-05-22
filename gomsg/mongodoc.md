the whole database suppose to be like this: 

Embedded Document Pattern



// server document
{
   _id: "webchat",
   name: "somewhat"
}


// room document
{
   _id: "ro"
   server_id: "webechat", // reference to server document 
   roomname: "room1"
}


//user document ??? should it refers to server or room. 
{
   _id: "use1"
   server_id: "webchat"
   nickname: "jack"
}

//
//message document
{
   _id: "msg1"
   room_id: "ro"
   nickname: "jack"
   timestamp: 
   content: ""

}

