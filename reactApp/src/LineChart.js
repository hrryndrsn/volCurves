import React, { Component } from "react";
import "./App.css";
import * as d3 from "d3";
import Axis from "./Axis"

const width = 800;
const height = 500;
const margin = { top: 20, right: 5, bottom: 20, left: 35 };

class LineChart extends Component {
  state = {
    bars: [],
  };


  calculateGraph(data) {
    const expr = data.sorted[Object.keys(data.sorted)[1]];
    console.log(data)
    if (!data) return {};
    const putArr = expr.Puts;
    const callArr = expr.Calls;

    //order by strike price
    var orderedPuts = putArr.sort(function(a, b) {
      return a.Instrument.strike - b.Instrument.strike;
    });
    var orderedCalls = callArr.sort(function(a, b) {
      return a.Instrument.strike - b.Instrument.strike;
    });

    // X axis is the strike
    // get min and max
    const XExtent = d3.extent(orderedPuts, d => d.Instrument.strike);
    const xScale = d3
      .scaleLinear(XExtent)
      .domain(XExtent)
      .range([margin.left, width - margin.right]);

    //  map vol to y-position
    // get min/max of high temp
    const [min, max] = d3.extent(orderedPuts, d => d.OrderBook.askIv);
    const yScale = d3
      .scaleLinear()
      .domain([Math.min(min, 0), max])
      .range([height - margin.bottom, margin.top]);


    //Ask line generator
    var askLine = d3.line()
      .x(d => xScale(d.Instrument.strike))
      .y(d => yScale(d.OrderBook.askIv))

    //ask lines 
    var askPutLine = askLine(orderedPuts);
    var askCallLine = askLine(orderedCalls);

    //Ask line generator
    const bidLine = d3.line()
      .x(d => xScale(d.Instrument.strike))
      .y(d => yScale(d.OrderBook.bidIv))

     //ask lines 
     var bidPutLine = bidLine(orderedPuts);
     var bidCallLine = bidLine(orderedCalls);

    return { askPutLine, askCallLine, bidPutLine, bidCallLine, XExtent, xScale, yScale };
  }


  render() {
    const data = this.props.data;
    const gridData = this.calculateGraph(data);
    const askPut = gridData.askPutLine
    const askCall = gridData.askCallLine
    const bidPut = gridData.bidPutLine
    const bidCall = gridData.bidCallLine


    if (data) {
      //we have data
      return (
        <div>
          <svg width={width+40} height={height}>
            {/* ask lines */}
            <path d={askPut} fill="none" stroke="red" />
            <path d={askCall} fill="none" stroke="blue" />
            {/* bid lines */}
            <path d={bidPut} fill="none" stroke="green" />
            <path d={bidCall} fill="none" stroke="orange" />
            {/* Hardcarded Legend */}
            <rect x={50} y={50} width={10} height={10} fill="blue"/>
            <text x={70} y={60}>Ask calls</text>
            <rect x={50} y={80} width={10} height={10} fill="red"/>
            <text x={70} y={90}> Ask puts</text>
            {/* axis */}
            <Axis data={gridData} width={width} height={height}/>     
          </svg>

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
