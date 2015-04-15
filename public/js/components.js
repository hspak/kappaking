/** @jsx React.DOM */

var ChannelTable = React.createClass({
  mixins: [SetIntervalMixin],

  getInitialState: function() {
    return {streams: [], downtime: 0};
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
      this.setState({ streams: data.Streams, downtime: data.DownTime });
    }.bind(this);
    xhr.send();
  },

  render: function() {
    var cells = [];
    var first = true;
    var firstCell;
    var avg = 0;
    this.state.streams.sort(function(a, b) {
      return parseInt(b.currkpm) - parseInt(a.currkpm);
    });

    var i = 1;
    var url = "";
    this.state.streams.forEach(function(stream) {
      var since = Math.round((Date.now() - Date.parse(stream.maxkpm_date))/60000);
      var sinceConvert = convertMinutes(since);
      if (stream.logo == "") {
        stream.logo = "http://static-cdn.jtvnw.net/jtv_user_pictures/xarth/404_user_150x150.png";
      }
      url = "http://www.twitch.tv/" + stream.display_name;

      if (stream.maxkpm > 0) {
        stream.maxkpm_date = sinceConvert;
      } else {
        stream.maxkpm_date = "never";
      }

      if (stream.minutes > 0) {
        avg = (stream.kappa / stream.minutes).toFixed(2);
      } else {
        avg = 0;
      }

      if (first) {
        firstCell =
          <ChannelCell first={first} url={url} avg={avg} stream={stream} key={stream.display_name} />;
        first = false;
      } else {
        cells.push(<ChannelCell first={first} url={url} avg={avg} stream={stream} key={stream.display_name} />);
      }

      stream.display_name = i + ' ' + stream.display_name;
      i += 1;
    });
    return (
      <div>
        <Header />
        <DownTime time={this.state.downtime}/>
        <div className="cells">
          {firstCell}
          <div className="bar"></div>
          <div className="channelTable">
            {cells}
          </div>
        </div>
      </div>
    );
  }
});

var DownTime = React.createClass({
  render: function() {
    if (this.props.time > 0) {
      return (
        <div className="deadtime">
          Kappa Counter has been dead for: {this.props.time}s
        </div>
      );
    } else {
      return <div></div>;
    }
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
            avg={this.props.avg}
            kappa={this.props.stream.kappa}
            maxkpm={this.props.stream.maxkpm}
            date={this.props.stream.maxkpm_date}
            currkpm={this.props.stream.currkpm} />;

    if (this.props.first) {
      channelType = <div key={this.props.key} className="firstCrowns">
        <img src="/img/crown.png"></img>
        <div className="channelCellFirst">
        {contentStatic}
        {contentDynamic}
        </div>
        <img src="/img/crown.png"></img>
        </div>;
    } else {
      channelType = <div key={this.props.key} className="channelCell">
        {contentStatic}
        {contentDynamic}
        </div>;
    }

    return (
      <a href={this.props.url}>
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
          <div className="avgKpm">AVG: {this.props.avg}</div>
          <div className="maxKpm">MAX: {this.props.maxkpm}</div>
          <div className="maxKpmDate">set: {this.props.date}</div>
        </div>
        <div className="nonkpm">
          <div className="kappa">Kappa: {this.props.kappa}</div>
          <div className="gameTitle">{this.props.game}</div>
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
