import React from 'react';
//import ReactDOM from "react-dom";
import Form from "@rjsf/material-ui";
import logo from './nginx_unit.png';

const unitSchemas = require('./schemas.js');
const submitData = ({ formData }, e) => sendMsg(JSON.stringify(formData));

var socket = new WebSocket("ws://localhost:8080/ws");

let connect = cb => {
  console.log("Attempting Connection...");

  socket.onopen = () => {
    console.log("Successfully Connected");
  };

  socket.onmessage = msg => {
    console.log(msg);
    cb(msg);
  };

  socket.onclose = event => {
    console.log("Socket Closed Connection: ", event);
  };

  socket.onerror = error => {
    console.log("Socket Error: ", error);
  };
};

let sendMsg = msg => {
  console.log("sending msg: ", msg);
  socket.send(msg);
};

export { connect, sendMsg };

const Home = () => { 
  <div className="App">
    <header className="App-header">
      <img src={logo} className="App-logo" alt="logo" />
    </header>
    <div>
      <Form schema={unitSchemas.testSchema} onSubmit={submitData} uiSchema={unitSchemas.uiSchema} />
    </div>
  </div>
}

export default Home;