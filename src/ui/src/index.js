import React from "react";
import ReactDOM from "react-dom";
import { BrowserRouter, Route, Switch } from "react-router-dom";
import Axios from "axios";
import Highcharts from "highcharts";
import HighchartsDrilldown from 'highcharts-drilldown';

import "./index.css";

const NotFound = () =>
	<div>
		<h3>404 page not found</h3>
	</div>



class UploadView extends React.Component {
	render() {
		return (
				<div>
					<h1>fitnotes</h1>
					<form action="http://localhost:8080/api/v1/import/fitnotes" method="post" encType="multipart/form-data">
					<input type="file" name="files" multiple />
					<input type="submit" value="Submit" />
					</form>
					<h1>strava</h1>
					<form action="http://localhost:8080/api/v1/import/strava" method="post" encType="multipart/form-data">
					<input type="file" name="files" multiple />
					<input type="submit" value="Submit" />
					</form>
				</div>
		       );
	}
}

class LineGraphView extends React.Component {
	constructor(props) {
		super(props);

		this.state = {
			data: [],
		};
	}
	componentWillMount() {
		//{"x": [<readable time string>,...], "y": [<value>,...]}
		Axios.get(this.props.host + this.props.url)
			.then(res => {
				this.setState({
					data: res.data.data,
				});
				this.renderAsync();
			});
	}
	componentDidMount() {
	}
	componentWillUnmount() {
		this.state.chart.destroy();
	}
	renderAsync() {
		// plotBands
		var plotBands = [{
			id: "0",
			from: this.props.bandOne,
			to: this.props.bandTwo,
			color: "rgba(244, 67, 54, 0.1) ",
		}, {
			id: "1",
			from: this.props.bandTwo,
			to: this.props.bandThree,
			color: "rgba(255, 235, 59, 0.1)",
		}, {
			id: "2",
			from: this.props.bandThree,
			to: 1000 * (this.props.bandThree - this.props.bandOne),
			color: "rgba(76, 175, 80, 0.1)",
		}];
		var config = {
			chart: {
				type: "line",
			},
			legend: {
				enabled: false,
			},
			series: [{
				name: this.props.titleX,
				data: this.state.data,
				pointPadding: 0,
				groupPadding: 0,
			}],
			title: {
				text: this.props.title + " - " + this.state.data[this.state.data.length - 1].y + this.props.unitsY
			},
			tooltip: {
				valueDecimals: 1,
				valueSuffix: " " + this.props.unitsY
			},
			xAxis: {
				showEmpty: false,
				type: "category",
			},
			yAxis: {
				min: 0,
				plotBands: JSON.parse(JSON.stringify(plotBands)),
				showEmpty: false,
				title: {
					text: this.props.titleY
				}
			},
		};
		this.setState({
			chart: Highcharts.chart(ReactDOM.findDOMNode(this), config)
		});
	}
	render() {
		return (
			<div className="graph-view line-graph-view" id={ this.state.id }></div>
	       );
	}
}

class ColumnGraphView extends React.Component {
	constructor(props) {
		super(props);

		this.state = {
			data: [],
			drilldown: [],
		};
	}
	componentWillMount() {
		HighchartsDrilldown(Highcharts);
		//{"x": [<readable time string>,...], "y": [<value>,...]}
		Axios.get(this.props.host + this.props.url)
			.then(res => {
				this.setState({
					data: res.data.data,
					drilldown: res.data.drilldown,
				});
				this.renderAsync();
			});
	}
	componentDidMount() {
	}
	componentWillUnmount() {
		this.state.chart.destroy();
	}
	renderAsync() {
		// drilldown series
		var drilldown = this.state.drilldown;
		for (var i = 0; i < drilldown.length; i++) {
			drilldown[i] = Object.assign({
				pointPadding: 0,
				groupPadding: 0,
			},drilldown[i])
		}
		// plotBands
		var plotBands = [{
			id: "0",
			from: this.props.bandOne,
			to: this.props.bandTwo,
			color: "rgba(244, 67, 54, 0.1) ",
		}, {
			id: "1",
			from: this.props.bandTwo,
			to: this.props.bandThree,
			color: "rgba(255, 235, 59, 0.1)",
		}, {
			id: "2",
			from: this.props.bandThree,
			to: 1000 * (this.props.bandThree - this.props.bandOne),
			color: "rgba(76, 175, 80, 0.1)",
		}];
		var config = {
			chart: {
				type: "column",
				events: {
					drilldown: function (e) {
						for (var i = 0; i < plotBands.length; i++) {
							this.yAxis[0].removePlotBand(plotBands[i].id);
						}
					},
					drillup: function (e) {
						for (var i = 0; i < plotBands.length; i++) {
							this.yAxis[0].addPlotBand(plotBands[i]);
						}
					},
				},
			},
			drilldown: {
				activeAxisLabelStyle: { "color": "#666666", "cursor": "default", "fontSize": "11px", "fontWeight": "normal", "textDecoration": "none"},
				activeDataLabelStyle: { "color": "#666666", "cursor": "default", "fontSize": "11px", "fontWeight": "normal", "textDecoration": "none"},
				animation: false,
				series: drilldown,
			},
			legend: {
				enabled: false,
			},
			series: [{
				name: this.props.titleX,
				data: this.state.data,
				pointPadding: 0,
				groupPadding: 0,
			}],
			title: {
				text: this.props.title
			},
			tooltip: {
				valueDecimals: 1,
				valueSuffix: " " + this.props.unitsY
			},
			xAxis: {
				showEmpty: false,
				type: "category",
			},
			yAxis: {
				min: 0,
				plotBands: JSON.parse(JSON.stringify(plotBands)),
				showEmpty: false,
				title: {
					text: this.props.titleY
				}
			},
		};
		this.setState({
			chart: Highcharts.chart(ReactDOM.findDOMNode(this), config)
		});
	}
	render() {
		return (
			<div className="graph-view column-graph-view" id={ this.state.id }></div>
	       );
	}
}

class DashboardView extends React.Component {
	render() {
		return (
			<div id="dashboard">
				<div id="powerlifting-graphs">
					<LineGraphView
						host="http://localhost:8080"
						url='/api/v1/data/maxes?fields=[{"name":"Flat Barbell Bench Press"},{"name":"Barbell Back Squat"},{"name":"Conventional Barbell Deadlift"}]&increment=3'
						title="Total"
						titleY="Max (lbs)"
						titleX="total"
						unitsY="lbs"
						bandOne="0"
						bandTwo="750"
						bandThree="1000"
					/>
					<LineGraphView
						host="http://localhost:8080"
						url='/api/v1/data/maxes?fields=[{"name":"Flat Barbell Bench Press"}]&increment=3'
						title="Bench"
						titleY="Max (lbs)"
						titleX="max bench"
						unitsY="lbs"
						bandOne="0"
						bandTwo="165"
						bandThree="315"
					/>
					<LineGraphView
						host="http://localhost:8080"
						url='/api/v1/data/maxes?fields=[{"name":"Barbell Back Squat"}]&increment=3'
						title="Squat"
						titleY="Max (lbs)"
						titleX="max squat"
						unitsY="lbs"
						bandOne="0"
						bandTwo="225"
						bandThree="405"
					/>
					<LineGraphView
						host="http://localhost:8080"
						url='/api/v1/data/maxes?fields=[{"name":"Conventional Barbell Deadlift"}]&increment=3'
						title="Deadlift"
						titleY="Max (lbs)"
						titleX="max deadlift"
						unitsY="lbs"
						bandOne="0"
						bandTwo="315"
						bandThree="495"
					/>
				</div>
				<div id="training-graphs">
					<ColumnGraphView
						host="http://localhost:8080"
						url='/api/v1/data/drilldown?fields=[{"name":"Brazilian Jiu-Jitsu","attr":2},{"name":"Judo","attr":2},{"name":"Wrestling","attr":2},{"name":"Boxing","attr":2},{"name":"Kickboxing","attr":2},{"name":"MMA","attr":2}]&increment=3'
						title="training"
						titleY="Time training"
						titleX="training"
						unitsY="hours"
						bandOne="0"
						bandTwo="16"
						bandThree="32"
					/>
					<ColumnGraphView
						host="http://localhost:8080"
						url='/api/v1/data/drilldown?fields=[{"name":"Road Cycling","attr":1}]&increment=3'
						title="biking"
						titleY="Distance biked"
						titleX="biking"
						unitsY="miles"
						bandOne="0"
						bandTwo="144"
						bandThree="280"
					/>
				</div>
			</div>
		);
	}
}

const Main = () => (
	<main>
	<Switch>
	<Route exact path="/" component={UploadView} />
	<Route path="/dashboard" component={DashboardView} />
	<Route path="/graph/training" component={() => (<ColumnGraphView
		host="http://localhost:8080"
		url='/api/v1/data/drilldown?fields=[{"name":"Brazilian%20Jiu-Jitsu","attr":2},{"name":"Judo","attr":2},{"name":"Wrestling","attr":2},{"name":"Boxing","attr":2},{"name":"Kickboxing","attr":2},{"name":"MMA","attr":2}]&increment=3' 
		title="training"
		titleY="Time training"
		titleX="training"
		unitsY="hours"
		bandOne="0"
		bandTwo="16"
		bandThree="32"
	/>)} />
	<Route path="/graph/biking" component={() => (<ColumnGraphView
		host="http://localhost:8080"
		url='/api/v1/data/drilldown?fields=[{"name":"Road Cycling","attr":1}]&increment=3'
		title="biking"
		titleY="Distance biked"
		unitsY="miles"
		titleX="biking"
		bandOne="0"
		bandTwo="144"
		bandThree="280"
	/>)} />
	<Route path="*" component={NotFound} />
	</Switch>
	</main>
)

const App = () => (
	<div>
	<Main />
	</div>
)

ReactDOM.render((
	<BrowserRouter>
	<App />
	</BrowserRouter>
), document.getElementById("root"))
