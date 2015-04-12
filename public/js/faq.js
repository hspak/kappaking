/** @jsx React.DOM */

var Faq = React.createClass({
  render: function() {
    return (
      <div>
        <Header />
      </div>
    );
  }
});

React.render(
  <Faq />,
  document.getElementById('content')
);
