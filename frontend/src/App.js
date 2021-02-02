import React, { Component } from 'react';
import { BrowserRouter, Route, Switch } from 'react-router-dom';
import './App.css' 
import Home from './components/Home/Home.js';
import Deploy from './components/Deploy/Deploy.js';
import Error from './components/Error/Error.js';
import Navigation from './components/Navigation/Navigation.js';

var functions = require('./components/functions');

class App extends Component {
  
  constructor(props) {
    super(props);
    this.state = {
      deployHistory: []
    }
  }

  componentDidMount() {
    functions.connect((msg) => {
      console.log("New Message", msg);
      this.setState(prevState => ({
         deployHistory: [...this.state.deployHistory, msg]
      }))
      console.log(this.state);
    });
  }
  
  render () {
    return (
      <BrowserRouter>
        <div>
          <Navigation />
            <Switch>
              <Route path="/" component={Home} exact/>
              <Route path="/deploy"render={(props) => (<Deploy {...props} deployHistory={this.state.deployHistory} />)} />
              <Route component={Error}/>
            </Switch>
        </div> 
      </BrowserRouter>
    );
  }
}
 
export default App;