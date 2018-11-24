import React, { Component } from "react";
import * as d3 from "d3";

export default class Axis extends Component {
    
    constructor(props) {
        super(props)
        this.props = props
    }

    componentDidMount() { this.renderAxis() }
    componentDidUpdate() { this.renderAxis() }
    
    renderAxis() {
        
	    const scale = d3.scaleLinear(this.props.data.XExtent)
	                    .domain(this.props.data.XExtent)
	                    .range([0, this.props.width]);
    	const axis = d3.axisBottom(scale);

	    d3.select(this.refs.g)
	      .call(axis);  
    }

    render() {
    	return <g transform="translate(10, 400)" ref="g" />
    }
}
