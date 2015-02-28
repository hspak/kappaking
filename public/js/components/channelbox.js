/** @jsx React.DOM */

var ChannelTable = React.createClass({
  render: function() {
    return(
      <div className="channelTable">
        <ChannelCell />
      </div>
    );
  }
});

var ChannelCell = React.createClass({
  render: function() {
    return (
      <div className="channelCell">
        <ChannelStatic />
        <ChannelDynamic />
      </div>
    );
  }
});

var ChannelStatic = React.createClass({
  render: function() {
    return (
      <div className="channelStatic">
        <div className="displayName">Display Name</div>
        <div className="channelLogo">Logo</div>
      </div>
    );
  }
});

var ChannelDynamic = React.createClass({
  render: function() {
    return (
      <div className="channelDynamic">
        <div className="gameTitle">Game: Title</div>
        <div className="viewerCount">Viewer: Count</div>
        <div className="kappaCount">Kappa: Count</div>
      </div>
    );
  }
});

React.render(
  <ChannelTable />,
  document.getElementById('content')
);
