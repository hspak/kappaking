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
      <div>KappaBoard</div>
    );
  }
});

var KPMBoard = React.createClass({

  render: function() {
    return (
      <div>KPMBoard</div>
    );
  }
});

var AvgBoard = React.createClass({

  render: function() {
    return (
      <div>AvgBoard</div>
    );
  }
});

React.render(
  <Leaderboards />,
  document.getElementById('content')
);
