var socket = new WebSocket("ws://192.168.1.48:8080/ws");
/*
var http = require('http');
var WebSocketServer = require('websocket').server;
var server = http.createServer();
server.listen(8080);

var socket = new WebSocketServer({
    httpServer: server
});
*/
module.exports.connect = function connect(cb) {
  console.log("Attempting Connection...");

  socket.onopen = () => {
    console.log("Successfully Connected");
  };

  socket.onmessage = msg => {
    //console.log(msg);
    cb(msg);
  };

  socket.onclose = event => {
    console.log("Socket Closed Connection: ", event);
  };

  socket.onerror = error => {
    console.log("Socket Error: ", error);
  };
};

module.exports.sendMsg = function sendMsg(msg) {
  console.log("sending msg: ", msg);
  socket.send(msg);
};
