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
    this.state.streams.forEach(function(stream) {
      cells.push(<ChannelCell key={stream.name} stream={stream} />);
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
      <div className="channelCell">
        <ChannelStatic
          displayName={this.props.stream.display_name}
          logo={this.props.stream.logo} />
        <ChannelDynamic
          game={this.props.stream.game}
          viewers={this.props.stream.viewers}
          kappa={this.props.stream.kappa} />
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
      </div>
    );
  }
});

var ChannelDynamic = React.createClass({
  render: function() {
    return (
      <div className="channelDynamic">
        <div className="gameTitle">Game: {this.props.game}</div>
        <div className="viewerCount">Viewer: {this.props.viewers}</div>
        <div className="kappaCount">KPM: {this.props.kappa}</div>
      </div>
    );
  }
});

React.render(
  <ChannelTable />,
  document.getElementById('content')
);
