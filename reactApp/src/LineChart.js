import React, { Component } from "react";
import "./App.css";
import * as d3 from "d3";

const width = 650;
const height = 400;
const margin = { top: 20, right: 5, bottom: 20, left: 35 };

class LineChart extends Component {
  state = {
    bars: [],
    myRef: React.createRef()
  };

  constructor(props) {
    super(props)
    this.myRef = React.createRef(); 
  }


  calculateGraph(data) {
    const expr = data.sorted[Object.keys(data.sorted)[1]];
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
    const extent = d3.extent(orderedPuts, d => d.Instrument.strike);
    const xScale = d3
      .scaleLinear(extent)
      .domain(extent)
      .range([margin.left, width - margin.right]);

    //  map vol to y-position
    // get min/max of high temp
    const [min, max] = d3.extent(orderedPuts, d => d.OrderBook.askIv);
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

    //line generator
    var line = d3.line()
      .x(d => xScale(d.Instrument.strike))
      .y(d => yScale(d.OrderBook.askIv))

    var putLine = line(orderedPuts);
    var callLine = line(orderedCalls);

    
    console.log(this.myRef)

    return { bars, putLine, callLine, xScale, yScale };
  }


  render() {
    const data = this.props.data;
    const gridData = this.calculateGraph(data);
    const pd = gridData.putLine
    const cd = gridData.callLine


    if (data) {
      //we have data
      return (
        <div>
          <svg width={width} height={height}>
            {/* lines */}
            <path d={pd} fill="none" stroke="red" />
            <path d={cd} fill="none" stroke="blue" />
            {/* axis */}
            <g id="xAxis" ref={this.myRef} transform={`translate(0, ${height - margin.bottom})`} />
            <g id="yAxis" transform={`translate(${margin.left}, 0)`} />            
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
