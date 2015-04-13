/** @jsx React.DOM */

var Leaderboards = React.createClass({
  getInitialState: function() {
    return {leaders: []};
  },

  componentDidMount: function() {
    var xhr = new XMLHttpRequest();
    xhr.open('get', '//' + location.host + "/api/get/leaders", true);
    xhr.onload = function() {
      var data = JSON.parse(xhr.responseText);
      this.setState({ leaders: data });
    }.bind(this);
    xhr.send(); 
  },

  render: function() {
    if (Array.isArray(this.state.leaders.highest_avg)) {
      this.state.leaders.highest_avg.sort(function(a, b) {
        return parseInt(b.avg) - parseInt(a.avg);
      });

      this.state.leaders.highest_kpm.forEach(function(kpm) {
        kpm.kpm_date = kpm.kpm_date.split("T")[0];
      })
    }

    return (
      <div>
        <Header />
        <div className="tables">
          <KappaBoard kappas={this.state.leaders.most_kappa}/>
          <KPMBoard kpm={this.state.leaders.highest_kpm}/>
          <AvgBoard avg={this.state.leaders.highest_avg}/>
        </div>
      </div>
    );
  }
});

var KappaBoard = React.createClass({
  render: function() {
    var rows = [];
    if (Array.isArray(this.props.kappas)) {
      this.props.kappas.forEach(function(kappa) {
        rows.push(<tr key={kappa.name}>
          <td>{kappa.name}</td>
          <td>{kappa.kappas}</td></tr>);
      });
    }
    return (
      <div className="kappa-table">
        <div className="table-name">Most Kappas</div>
        <table>
          <thead>
            <td>Streamer</td>
            <td>Kappas</td>
          </thead>
          {rows}
        </table>
      </div>
    );
  }
});

var KPMBoard = React.createClass({
  render: function() {
    var rows = [];
    if (Array.isArray(this.props.kpm)) {
      console.log(this.props.kpm);
      this.props.kpm.forEach(function(kpm) {
        rows.push(<tr key={kpm.name}>
          <td>{kpm.name}</td>
          <td>{kpm.kpm}</td>
          <td>{kpm.kpm_date}</td></tr>);
      });
    }
    return (
      <div className="kappa-table">
        <div className="table-name">Highest KPM</div>
        <table>
          <thead>
            <td>Streamer</td>
            <td>KPM</td>
            <td>Date</td>
          </thead>
          {rows}
        </table>
      </div>
    );
  } 
});

var AvgBoard = React.createClass({
  render: function() {
    var rows = [];
    if (Array.isArray(this.props.avg)) {
      this.props.avg.forEach(function(avg) {
        rows.push(<tr key={avg.name}>
          <td>{avg.name}</td>
          <td>{avg.avg.toFixed(2)}</td>
          <td>{avg.minutes}</td></tr>);
      });
    }
    return (
      <div className="kappa-table">
        <div className="table-name">Highest Average KPM</div>
        <table>
          <thead>
            <td>Streamer</td>
            <td>Average</td>
            <td>Minutes</td>
          </thead>
          {rows}
        </table>
      </div>
    );
  } 
});

React.render(
  <Leaderboards />,
  document.getElementById('content')
);
