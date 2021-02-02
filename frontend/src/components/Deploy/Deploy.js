import React, { Component } from 'react';
import './Deploy.css';
import Message from '../Messages/Message.jsx'


class Deploy extends Component {

  render() {
     //console.log("Assigned messsage", global.returnMsg);
     //const messages = this.props.deployHistory.map((msg, index) => (
     const messages = this.props.deployHistory.map(msg => <Message message={msg.data} />);

     return (
      <div className="Deploy">
         <h2>NGINX Unit Deployment</h2>
         {messages}
      </div>
      )
   }
}
 
export default Deploy;