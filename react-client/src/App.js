import React, {Component} from 'react'
import Header from './header'
import Counter from './counter'
class App extends Component {
  render() {
    return (
      <div>
        <Header></Header>
        <Counter></Counter>
      </div>
    );
  }
  
  componentDidMount() {
    console.log("Executing Mount");
    fetch("/data")
      .then((res) => {
        res.json()
          .then((data) => {
            console.log("data is ", data)
            this.setState({ currentPings: data });
            console.log(data);
          })
          .catch(console.log);
      },
      )
  }
}

export default App;

//class App extends Component {

//   state = {
//     contacts: []
//   }
//   ...
// }