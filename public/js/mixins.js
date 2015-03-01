/** @jsx React.DOM */

var SetIntervalMixin = {
  componentWillMount: function() {
    this.intervals = [];
  },
  setInterval: function() {
    this.intervals.push(setInterval.apply(null, arguments));
  },
  componentWillUnmonut: function() {
    this.intervals.map(clearInterval);
  }
};
