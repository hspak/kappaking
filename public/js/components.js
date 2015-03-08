/** @jsx React.DOM */

function convertMinutes(minutes) {
  if (minutes == 0) {
    return "now";
  } else if (minutes < 60) {
    return minutes + "m ago";
  } else if (minutes < 1440) { // 60*24
    return Math.round(minutes/60) + "h " + (minutes%60) + "m ago";
  } else if (minutes < 10080) { // 60*24*7
    return Math.round(minutes/1440) + "d ago";
  } else {
    return Math.round(minutes/10080) + "w ago";
  }
}

var ChannelTable = React.createClass({
  mixins: [SetIntervalMixin],
  getInitialState: function() {
    return {streams: []};
  },
  componentDidMount: function() {
    this.dataUpdate();
    this.setInterval(this.dataUpdate, 5000);
  },
  dataUpdate: function() {
    var xhr = new XMLHttpRequest();
    xhr.open('get', document.URL + "api/get/data", true);
    xhr.onload = function() {
      var data = JSON.parse(xhr.responseText);
      this.setState({ streams: data.Streams });
    }.bind(this);
    xhr.send();
  },
  render: function() {
    var cells = [];
    var first = true;
    var firstCell;
    this.state.streams.sort(function(a, b) {
      return parseInt(b.currkpm) - parseInt(a.currkpm);
    });

    var i = 1;
    this.state.streams.forEach(function(stream) {
      var since = Math.round((Date.now() - Date.parse(stream.maxkpm_date))/60000);
      var sinceConvert = convertMinutes(since);
      if (stream.logo == "") {
        stream.logo = "http://static-cdn.jtvnw.net/jtv_user_pictures/xarth/404_user_150x150.png";
      }
      if (stream.maxkpm > 0) {
        stream.maxkpm_date = sinceConvert;
      } else {
        stream.maxkpm_date = 0;
      }
      if (first) {
        firstCell = <ChannelCell first={first} stream={stream} key={stream.display_name} />;
        first = false;
      } else {
        cells.push(<ChannelCell first={first} stream={stream} key={stream.display_name} />);
      }
      stream.display_name = i + ' ' + stream.display_name;
      i += 1;
    });
    return(
      React.createElement("div", {className: "cells"},
        firstCell,
        React.createElement("div", {className: "channelTable"},
          cells
        )
      )
    );
  }
});

var ChannelCell = React.createClass({
  render: function() {
    var channelType;
    var contentStatic = <ChannelStatic
            displayName={this.props.stream.display_name}
            logo={this.props.stream.logo} />;
    var contentDynamic =
          <ChannelDynamic
            game={this.props.stream.game}
            viewers={this.props.stream.viewers}
            minutes={this.props.stream.minutes}
            kappa={this.props.stream.kappa}
            maxkpm={this.props.stream.maxkpm}
            date={this.props.stream.maxkpm_date}
            currkpm={this.props.stream.currkpm} />;

    if (this.props.first) {
      channelType = <div key={this.props.key} className="channelCellFirst">
        {contentStatic}
        {contentDynamic}
        </div>;
    } else {
      channelType = <div key={this.props.key} className="channelCell">
        {contentStatic}
        {contentDynamic}
        </div>;
    }

    return (
      <a href={this.props.stream.url}>
        {channelType}
      </a>
    );
  }
});

var ChannelStatic = React.createClass({
  render: function() {
    return (
      <div className="channelStatic">
        <div className="displayName">{this.props.displayName}</div>
        <div className="channelLogo"><img src={this.props.logo}></img></div>
      </div>
    );
  }
});

var ChannelDynamic = React.createClass({
  render: function() {
    return (
      <div className="channelDynamic">
        <div className="kpm">
          <div className="currKpm">KPM: {this.props.currkpm}</div>
          <div className="maxKpm">MAX: {this.props.maxkpm}</div>
          <div className="maxKpmDate">set: {this.props.date}</div>
        </div>
        <div className="nonkpm">
          <div className="kappa">Kappa: {this.props.kappa}</div>
          <div className="minutes">Minutes Recorded: {this.props.minutes}</div>
          <div className="gameTitle">Game: {this.props.game}</div>
          <div className="viewerCount">Viewer: {this.props.viewers}</div>
        </div>
      </div>
    );
  }
});

React.render(
  <ChannelTable />,
  document.getElementById('content')
);
