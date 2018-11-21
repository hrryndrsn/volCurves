import React, { Component } from "react";
import "./App.css";
import * as d3 from "d3";

const width = 650;
const height = 400;
const margin = { top: 20, right: 5, bottom: 20, left: 35 };

class LineChart extends Component {

  state = {
    bars: []
  };

  calculateGraph(data) {
    const expr = data.sorted[Object.keys(data.sorted)[1]];
    // if (!data) return {};
    const putArr = expr.Puts;
    const callArr = expr.Calls;

    // X axis is the strike
    // get min and max
    // const extent = d3.extent(data, d => d)
    console.log("puts", putArr)
    console.log("calls", callArr)
    // return JSON.stringify(data.sorted);

    // Y axis is Implied vol
    // 1 line Call bid LV
    // 2 line Call ask LV
    // 3 line Put ask LV
  }

  render() {
    const data = this.props.data;
    if (data) {
      //we have data
      return (
        <div>
          <svg width={width} height={height} />
          <div>{console.log(this.calculateGraph(data))}</div>
        </div>
      );
    } else {
      return (
        //we have no / invalid props.data
        <div>
          <p>no data...</p>
        </div>
      );
    }
  }
}

export default LineChart;
