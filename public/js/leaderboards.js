/** @jsx React.DOM */

var Leaderboards = React.createClass({
  render: function() {
    return (
      <div>
        <Header />
        <KappaBoard />
        <KPMBoard />
        <AvgBoard />
      </div>
    );
  }
});

var KappaBoard = React.createClass({
  render: function() {
    return (
      <p>KappaBoard</p>
    );
  }
});

var KPMBoard = React.createClass({

  render: function() {
    return (
      <p>KPMBoard</p>
    );
  }
});

var AvgBoard = React.createClass({

  render: function() {
    return (
      <p>AvgBoard</p>
    );
  }
});

React.render(
  <Leaderboards />,
  document.getElementById('tables')
);
