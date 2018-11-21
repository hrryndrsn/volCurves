import React, { Component } from 'react';
import './App.css';

const width = 650;
const height = 400;
const margin = {top: 20, right: 5, bottom: 20, left: 35};

class Chart extends Component {
  state = {
    bars: []
  }

  render() {
    const data = this.props.data
    if (data) {
      //we have data
      return (
        <svg width={width} height={height} />
      );
    } else {
      return (
        //we have no / invalid props.data 
        <div><p>no data...</p></div>
      );
    }

  }
}

export default Chart;


