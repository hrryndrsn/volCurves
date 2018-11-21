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
    if (!data) return {};
    const putArr = expr.Puts;
    // const callArr = expr.Calls;

    // X axis is the strike
    // get min and max
    const extent = d3.extent(putArr, d => d.Instrument.strike);
    const xScale = d3
      .scaleLinear(extent)
      .domain(extent)
      .range([margin.left, width - margin.right]);

    // 2. map vol to y-position
    // get min/max of high temp
    const [min, max] = d3.extent(putArr, d => d.OrderBook.askIv);
    const yScale = d3
      .scaleLinear()
      .domain([Math.min(min, 0), max])
      .range([height - margin.bottom, margin.top]);

    // array of objects: x: y., height
    const bars = putArr.map(d => {
      return {
        x: xScale(d.Instrument.strike),
        y: yScale(d.OrderBook.askIv),
        height: 10,
        fill: "black"
      };
    });

    // console.log("puts", putArr)
    // console.log("Put extent:", extent)
    // console.log("askIvextent:", min, max)
    // console.log("bar data:", bars)

    return { bars, xScale, yScale };
  }

  render() {
    const data = this.props.data;
    const gridData = this.calculateGraph(data);
    const bars = gridData.bars;
    if (data) {
      //we have data
      return (
        <div>
          <svg width={width} height={height}>
            {bars.map(d => (
              <rect
                x={d.x}
                y={d.y}
                width={10}
                height={d.height}
                fill={d.fill}
              />
            ))}
          </svg>
          <div>{console.log(bars)}</div>
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
