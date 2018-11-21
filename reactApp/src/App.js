import React, { Component } from 'react';
import './App.css';
import LineChart from "./LineChart.js"
import sampleData from "./sampleData.json"

class App extends Component {
  render() {
    return (
      <div className="App">
        <LineChart data={sampleData}/>
      </div>
    );
  }
}

export default App;
