/** @jsx React.DOM */

 var Header = React.createClass({
  render: function() {
    return (
      <div className="header">
        <a href="/">
          <div className="face">
            <img src="/img/kappa.png"></img>
            <div className="title">KappaKing</div>
          </div>
        </a>
        <div className="links">
          <a href="/leaderboards">leaderboards</a>
          <a href="/faq">faq</a>
          <a href="https://github.com/hspak/kappaking">source</a>
        </div>
      </div>
    );
  }
});
 
