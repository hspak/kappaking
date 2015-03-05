/** @jsx React.DOM */

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
    this.state.streams.sort(function(a, b) {
      return parseInt(b.currkpm) - parseInt(a.currkpm);
    });
    this.state.streams.forEach(function(stream) {
      if (stream.logo == "") {
        stream.logo = "http://static-cdn.jtvnw.net/jtv_user_pictures/xarth/404_user_150x150.png";
      }
      cells.push(<ChannelCell stream={stream} key={stream.display_name} />);
    });
    return(
      <div className="channelTable">
        {cells}
      </div>
    );
  }
});

var ChannelCell = React.createClass({
  render: function() {
    return (
      <div key={this.props.key} className="channelCell">
        <ChannelStatic
          displayName={this.props.stream.display_name}
          logo={this.props.stream.logo}
          game={this.props.stream.game}
          viewers={this.props.stream.viewers} />
        <ChannelDynamic
          minutes={this.props.stream.minutes}
          kappa={this.props.stream.kappa}
          maxkpm={this.props.stream.maxkpm}
          currkpm={this.props.stream.currkpm} />
      </div>
    );
  }
});

var ChannelStatic = React.createClass({
  render: function() {
    return (
      <div className="channelStatic">
        <div className="displayName">{this.props.displayName}</div>
        <div className="channelLogo"><img src={this.props.logo}></img></div>
        <div className="gameTitle">Game: {this.props.game}</div>
        <div className="viewerCount">Viewer: {this.props.viewers}</div>
      </div>
    );
  }
});

var ChannelDynamic = React.createClass({
  render: function() {
    return (
      <div className="channelDynamic">
        <div className="currKpm">KPM: {this.props.currkpm}</div>
        <div className="maxKpm">MAX KPM: {this.props.maxkpm}</div>
        <div className="kappa">Kappa: {this.props.kappa}</div>
        <div className="minutes">Minutes Recorded: {this.props.minutes}</div>
      </div>
    );
  }
});

React.render(
  <ChannelTable />,
  document.getElementById('content')
);
