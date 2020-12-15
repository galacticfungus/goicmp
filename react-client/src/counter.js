import React, { Component } from "react";

class Counter extends Component {
    constructor(props) {
        super(props);
        this.state = { currentPings: 0 };
    }
  render() {
    return (
      <blockquote className="blockquote text-center">
        <p className="lead">Current Count: {this.state.currentPings}</p>
      </blockquote>
    );
  }

  updateCounter() {
    fetch("/data").then((res) => {
      res
        .json()
        .then((data) => {
          console.log("data is ", data);
          this.setState({ currentPings: data['count'] });
        })
        .catch(console.log);
    });
  }

  componentDidMount() {
    this.timerID = setInterval(() => this.updateCounter(), 1000);
  }

    componentWillUnmount() {
      clearInterval(this.timerID)
  }
}

export default Counter;

//class App extends Component {

//   state = {
//     contacts: []
//   }
//   ...
// }
