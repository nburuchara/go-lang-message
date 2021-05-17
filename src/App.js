import React,{Component} from 'react'
import logo from './logo.svg';
import './App.css';
import styled from 'styled-components'

const Styles = styled.div `


  // - - HEADER - - //

.header {
  text-align: center;
}

.header h1 {
  margin-top: 5%;
  font-family: Kanit;
}

.header img  {
  width: 200px;
}

  // - - LOADING - - //

.loading {
  text-align: center;
  margin-top: 10%;
}

.loading img {
  width: 20%;
}

  // - - CHAT - - //

.chat {
  text-align: center;
}

.chat p {
  font-family: Kanit;
  color: #F54631;
  font-size: 20.5px;
}

.chat input {
  width: 300px;
  height: 45px;
  border-radius: 8px;
  outline: none;
  text-align: center;
  border: 0.5px solid black;
  margin-bottom: 7px;
}

.chat button {
  margin-top: 20px;
  height: 40px;
  width: 120px;
  background-color: #42C7FC;
  color: white;
  font-family; Kanit;
  border-radius: 8px;
  outline: none;
  border: 2px solid #42C7FC;
}

  // - - CHAT PAGE - - //

.chatPage {
  text-align: center;
  
}

.chatPage h3 {
  margin-top: 10%;
  margin-bottom: 45px;
  font-family: Kanit;
}

.chatPage label {
  color: #42C7FC;
}

.chatPage h1 {
  margin-bottom: 45px;
  font-family: Kanit;
}

.chatPage button {
  height: 45px;
  width: 150px;
  background-color: #42C7FC;
  color: white;
  font-family; Kanit;
  border-radius: 8px;
  outline: none;
  border: 2px solid #42C7FC;
  margin-bottom: 50px;
}

  // - - NEW ROOM - - //

.newRoom {
  text-align: center;
}

.newRoom h1 {
  margin-top: 10%;
  font-family: Kanit;
  color: #42C7FC;
  margin-bottom: 30px;
}

.newRoom input {
  height: 42.5px;
  width: 315px;
  font-family: Kanit;
  text-align: center;
  border-radius: 8px;
  outline: none;
  border: 0.5px solid black;
  margin-bottom: 30px;
}

.newRoom button {
  height: 45px;
  width: 150px;
  background-color: #42C7FC;
  color: white;
  font-family; Kanit;
  border-radius: 8px;
  outline: none;
  border: 2px solid #42C7FC;
  margin-bottom: 50px;
  font-family: Kanit;
}

  // - - JOIN ROOM - - //

.joinRoom {
  text-align: center;
}

.joinRoom h1 {
  margin-top: 10%;
  font-family: Kanit;
  color: #42C7FC;
  margin-bottom: 40px;
}

.joinRoom button {
  height: 55px;
  width: 400px;
  background-color: #42C7FC;
  color: white;
  font-family; Kanit;
  border-radius: 8px;
  outline: none;
  border: 2px solid #42C7FC;
  margin-bottom: 50px;
  font-family: Kanit;
}


  // - - CHAT ROOM - - //

.crHeader {
  
}

.crHeader h1 {
  color: #42C7FC;
  text-align: center;
  font-family: Kanit;
  margin-top: 3%;
}

.chatRoom textarea {
  margin-top: 8px;
  margin-left: 27.5px;
  border-radius: 10px;
  border: 0.2px solid transparent;
  width: 93%;
  height: 60px;
  outline: none;
  padding: 7px;
  resize: none;
  margin-bottom: 5px;
}

.chatRoom button {
  background-color: white;
  width: 75px;
  height: 35px;
  margin-bottom: 20px;
  margin-left: 27.5px;
  border-radius: 8px;
  border: 2px solid transparent;
}


`

class App extends Component {
  constructor(props) {
    super()
    this.state = {
      home: true,
      errMsg: "",
      chat: false,
      nickname: "",
      newRoom: false,
      createdRoom : "",
      loadingScreen: false,
      joinRoom: false,
      chatRoom: false,
      roomMessage : "",
      selectedRoom: "",
      messageArray: []
    }
  }

  goToChat = () => {
    if (!this.state.nickname == "") {
      this.setState({
        home: false,
        chat: true
      })
    } else {
      this.setState({
        errMsg: "Please enter a name to proceed"
      })
    }
    
  }

  goToNewRoom = () => {
    this.setState({
      chat: false,
      newRoom: true
    })
  }
  handleChange = (event) => {
      this.setState({
          [event.target.id] : event.target.value
      })
  }

  createRoom = () => {
    this.setState({loadingScreen: true, newRoom: false}, () => {
        setTimeout(() => {
            this.setState({loadingScreen: false, joinRoom: true})
        }, 3000)
    })
  }

  sendMessage = () => {
    var obj = {}
    obj["01"] = this.state.nickname
    obj["02"] = this.state.roomMessage
    this.state.messageArray = this.state.messageArray.push(obj)
    console.log(obj)
  }

  goToChatRoom = () => {
    this.setState({loadingScreen: true, joinRoom: false, selectedRoom: this.state.createdRoom}, () => {
      setTimeout(() => {
          this.setState({loadingScreen: false, chatRoom: true})
      }, 3000)
  })
  }

  render () {

    let dialogueStyle = {
        width : "800px",
        maxWidth: "92.5%",
        margin: "0 auto",
        height: "650px",
        textAlignment: "center",
        zIndex: "999",
        backgroundColor: "#42C7FC",
        padding: "10px 20px 10px",
        borderRadius: "8px",
        flexDirection: "column",
        marginTop: "1%",
        marginBottom: "20px"
    }

    let dialogueStyle2 = {
      width : "99%",
      height: "75%",
      maxWidth: "92.5%",
      margin: "0 auto",
      textAlignment: "center",
      zIndex: "999",
      backgroundColor: "#fff",
      padding: "10px 20px 40px",
      borderRadius: "8px",
      flexDirection: "column",
      marginTop: "2%",
      marginBottom: "20px"
  }

    return (
      <Styles>
        {this.state.home && 
          <div>
            <div className="header">
                <h1>Go Language Web Messenger Tool</h1> 
                <img src="/assets/react.png"/>
                <img src="/assets/mongodb-leaf2.png"/>
                <img src="/assets/golang-gopher.png"/>
            </div>
            <div className="chat">
              <input
              id="nickname"
              onChange={this.handleChange}
              value={this.state.nickname}
              placeholder="Enter your name"
              /> 
              <p>{this.state.errMsg}</p>
              <button
              onClick={this.goToChat}
              >
                <b>Go to chat</b>
              </button>
            </div>
          </div>
        }
        {this.state.chat && 
          <div className="chatPage">
            <h3>Welcome <label>{this.state.nickname}</label></h3>
            <h1>Select an option</h1>
            <button
            onClick={this.goToNewRoom}
            ><b>New Room</b></button> <br/>
            <button><b>Join Room</b></button> <br/>
            <button><b>New Message</b></button> <br/>
            <button><b>Quit</b></button> <br/>
          </div>
        }
        {this.state.newRoom && 
          <div className="newRoom">
            <h1>Create a new room</h1>
            <input
            id="createdRoom"
            onChange={this.handleChange}
            value={this.state.createdRoom}
            placeholder="Enter the name of the new room"
            /> <br/>
            <button
            onClick={this.createRoom}
            >Create</button>
          </div>
        }
        {this.state.joinRoom && 
          <div className="joinRoom">
            <h1>Rooms</h1>
            <button
            onClick={this.goToChatRoom}
            >{this.state.createdRoom}</button>
          </div>
        }
        {this.state.chatRoom && 
          <div className="crHeader">
            <h1>{this.state.selectedRoom}</h1>
            <div className="chatRoom" style={dialogueStyle}>
              <div style={dialogueStyle2}>
                {this.state.messageArray.map(message => (
                  <div>
                    {message["01"]}
                  </div>
                ))}
              </div>
              <textarea
              id="roomMessage"
              value={this.state.roomMessage}
              onChange={this.handleChange}
              placeholder="Type your message here"
              />
              <button
              onClick={this.sendMessage}
              >Send</button>
            </div>
          </div>
        }
        {this.state.loadingScreen && 
            <div className="loading">
                <img src="/assets/gifLoader.gif"/>
            </div>
        }
      </Styles>
    );
  }
}

export default App;
