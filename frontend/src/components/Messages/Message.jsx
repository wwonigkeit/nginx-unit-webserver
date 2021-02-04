import React, { Component } from "react";
import "./Message.css";

class Message extends Component {

  componentDidMount () {
    this.scrollToBottom()
  }
  componentDidUpdate () {
    this.scrollToBottom()
  }
  scrollToBottom = () => {
    this.el.scrollIntoView({ behavior: 'smooth' });
  }

  constructor(props) {
    super(props);
    let temp = JSON.parse(this.props.message);
    this.state = {
      message: temp
    };
  }

  render() {
    return <div className="Message">
      {this.state.message.body}
        <div ref={el => { this.el = el; }} />
    </div>;
  }
}

export default Message;