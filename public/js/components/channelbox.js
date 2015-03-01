/** @jsx React.DOM */

var ChannelTable = React.createClass({
  render: function() {
    var cells = [];
    this.props.streams.forEach(function(stream) {
      cells.push(<ChannelCell stream={stream} key={stream.name} />);
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
          kappa={this.props.stream.kappa}
        />
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
        <div className="kappaCount">Kappa: {this.props.kappa}</div>
      </div>
    );
  }
});

var STREAMS =
// {{{
{  
  "Streams":[  
    {  
      "display_name":"Riot Games",
      "game":"League of Legends",
      "viewers":298644,
      "kappa":9001,
      "logo":"http://static-cdn.jtvnw.net/jtv_user_pictures/riotgames-profile_image-4be3ad99629ac9ba-300x300.jpeg",
      "status":"NA LCS Spring - Week 6 Day 1",
      "url":"http://www.twitch.tv/riotgames"
    },
    {  
      "display_name":"Tempo_Storm",
      "game":"Hearthstone: Heroes of Warcraft",
      "viewers":46989,
      "kappa":9001,
      "logo":"http://static-cdn.jtvnw.net/jtv_user_pictures/tempo_storm-profile_image-ec2dcc7d4debd0b7-300x300.png",
      "status":"Lord of the Arena 3 presented by G2A.com feat. Trump, TotalBiscuit, TheOddOne, GiantWaffle, SivHD, Massan, Guardsman Bob, and Yogscast Zylus",
      "url":"http://www.twitch.tv/tempo_storm"
    },
    {  
      "display_name":"iBUYPOWER",
      "game":"Counter-Strike: Global Offensive",
      "viewers":26445,
      "kappa":9001,
      "logo":"http://static-cdn.jtvnw.net/jtv_user_pictures/ibuypower-profile_image-3cc840e53f9e5387-300x300.jpeg",
      "status":"SKDC vs Luminosity - iBP Invitational #1",
      "url":"http://www.twitch.tv/ibuypower"
    },
    {  
      "display_name":"LIRIK",
      "game":"DragonBall Xenoverse",
      "viewers":22577,
      "kappa":9001,
      "logo":"http://static-cdn.jtvnw.net/jtv_user_pictures/lirik-profile_image-4bbb8e826c7240bb-300x300.png",
      "status":"The Usual - @DatGuyLirik",
      "url":"http://www.twitch.tv/lirik"
    },
    {  
      "display_name":"The Creatures",
      "game":"Brink",
      "viewers":15904,
      "kappa":9001,
      "logo":"http://static-cdn.jtvnw.net/jtv_user_pictures/thecreatures-profile_image-01616d1c509d02d5-300x300.png",
      "status":"Creature Talk Episode 119",
      "url":"http://www.twitch.tv/thecreatures"
    },
    {  
      "display_name":"RiotGamesBrazil",
      "game":"League of Legends",
      "viewers":15236,
      "kappa":9001,
      "logo":"http://static-cdn.jtvnw.net/jtv_user_pictures/riotgamesbrazil-profile_image-1c8aec8985d294ef-300x300.jpeg",
      "status":"LCS NA – S6 D1",
      "url":"http://www.twitch.tv/riotgamesbrazil"
    },
    {  
      "display_name":"WCS",
      "game":"StarCraft II: Heart of the Swarm",
      "viewers":15016,
      "kappa":9001,
      "logo":"http://static-cdn.jtvnw.net/jtv_user_pictures/wcs-profile_image-4270cab0f4148e12-300x300.jpeg",
      "status":"WCS Premier League - Ro32 - Group H",
      "url":"http://www.twitch.tv/wcs"
    },
    {  
      "display_name":"sodapoppin",
      "game":"Journey",
      "viewers":12549,
      "kappa":9001,
      "logo":"http://static-cdn.jtvnw.net/jtv_user_pictures/sodapoppin-profile_image-10049b6200f90c14-300x300.png",
      "status":" ( ° ͜ʖ͡°)╭∩╮<@sodapoppintv>━╤デ╦︻(▀̿̿Ĺ̯̿̿▀̿ ̿) Sh*t Show Saturday #75. Subs recommend games for me to play. ",
      "url":"http://www.twitch.tv/sodapoppin"
    },
    {  
      "display_name":"summit1g",
      "game":"Counter-Strike: Global Offensive",
      "viewers":12272,
      "kappa":9001,
      "logo":"http://static-cdn.jtvnw.net/jtv_user_pictures/summit1g-profile_image-1cf38b42b1545fd7-300x300.jpeg",
      "status":"CS w/ @summit1g. Salt is inevitable, I will be one with it.",
      "url":"http://www.twitch.tv/summit1g"
    },
    {  
      "display_name":"Giantwaffle",
      "game":"Minecraft",
      "viewers":10362,
      "kappa":9001,
      "logo":"http://static-cdn.jtvnw.net/jtv_user_pictures/giantwaffle-profile_image-4c53a2c25f94c6a9-300x300.png",
      "status":"Magical Mysteries! - FTB Infinity on The Build Guild Server",
      "url":"http://www.twitch.tv/giantwaffle"
    },
    {  
      "display_name":"SCGLive",
      "game":"Magic: The Gathering",
      "viewers":10008,
      "kappa":9001,
      "logo":"http://static-cdn.jtvnw.net/jtv_user_pictures/scglive-profile_image-cb37b7e8e06f1882-300x300.jpeg",
      "status":"SCG Open Series Baltimore, MD - February 28-March 1",
      "url":"http://www.twitch.tv/scglive"
    },
    {  
      "display_name":"mEclipse",
      "game":"Counter-Strike: Global Offensive",
      "viewers":9067,
      "kappa":9001,
      "logo":"http://static-cdn.jtvnw.net/jtv_user_pictures/meclipse-profile_image-0373ffcdea347024-300x300.jpeg",
      "status":"OH MY, IM AWAKE? THE SUN!? aHHHHH | @C9shroud",
      "url":"http://www.twitch.tv/meclipse"
    },
    {  
      "display_name":"OgamingLoL",
      "game":"League of Legends",
      "viewers":8240,
      "kappa":9001,
      "logo":"http://static-cdn.jtvnw.net/jtv_user_pictures/ogaminglol-profile_image-8871bc8519315c6a-300x300.jpeg",
      "status":"[FR] LCS NA - TSM vs DIGNITAS - W6D1 ",
      "url":"http://www.twitch.tv/ogaminglol"
    },
    {  
      "display_name":"HearthStats",
      "game":"Hearthstone: Heroes of Warcraft",
      "viewers":7952,
      "kappa":9001,
      "logo":"http://static-cdn.jtvnw.net/jtv_user_pictures/hearthstats-profile_image-b3c576b858ba8d8e-300x300.png",
      "status":"Gaara vs ARee | HearthStats Champion's League RO4",
      "url":"http://www.twitch.tv/hearthstats"
    },
    {  
      "display_name":"Itmejp",
      "game":"Dungeons & Dragons",
      "viewers":7541,
      "kappa":9001,
      "logo":"http://static-cdn.jtvnw.net/jtv_user_pictures/itmejp-profile_image-64703923f21827e3-300x300.png",
      "status":"RollPlay: The West Marches w/ DM @Silent0siris and players @itmeJP @CohhCarnage @dexbonus and @Ezekiel_III",
      "url":"http://www.twitch.tv/itmejp"
    },
    {  
      "display_name":"Mori09TV",
      "game":"League of Legends",
      "viewers":7260,
      "kappa":9001,
      "logo":"http://static-cdn.jtvnw.net/jtv_user_pictures/mori09tv-profile_image-a28291922a2e9442-300x300.png",
      "status":"[GER] LCS NA - W6D1",
      "url":"http://www.twitch.tv/mori09tv"
    },
    {  
      "display_name":"ESL_spain",
      "game":"League of Legends",
      "viewers":7175,
      "kappa":9001,
      "logo":"http://static-cdn.jtvnw.net/jtv_user_pictures/esl_spain-profile_image-33a4567a83dffa24-300x300.png",
      "status":"[Español] LCS NA S6D1 #SomosLCS",
      "url":"http://www.twitch.tv/esl_spain"
    },
    {  
      "display_name":"OMGitsfirefoxx",
      "game":"Minecraft",
      "viewers":6456,
      "kappa":9001,
      "logo":"http://static-cdn.jtvnw.net/jtv_user_pictures/omgitsfirefoxx-profile_image-f7261473917e00a2-300x300.jpeg",
      "status":"♥ MIANITE SEASON 2 Join4BootyTouches",
      "url":"http://www.twitch.tv/omgitsfirefoxx"
    },
    {  
      "display_name":"EternaLEnVyy",
      "game":"Dota 2",
      "viewers":6364,
      "kappa":9001,
      "logo":"http://static-cdn.jtvnw.net/jtv_user_pictures/eternalenvyy-profile_image-6b20e441fd6342b1-300x300.png",
      "status":"C9 EternaLEnVy Tryhard DotA (songs: songs to be tested)",
      "url":"http://www.twitch.tv/eternalenvyy"
    },
    {  
      "display_name":"Forsenlol",
      "game":"Hearthstone: Heroes of Warcraft",
      "viewers":5417,
      "kappa":9001,
      "logo":"http://static-cdn.jtvnw.net/jtv_user_pictures/forsenlol-profile_image-48b43e1e4f54b5c8-300x300.png",
      "status":"Forsen, ☑ 5 hours of sleep ☑ Drinking saturday ☐ Forsen RIP",
      "url":"http://www.twitch.tv/forsenlol"
    },
    {  
      "display_name":"Kolento",
      "game":"Hearthstone: Heroes of Warcraft",
      "viewers":5305,
      "kappa":9001,
      "logo":"http://static-cdn.jtvnw.net/jtv_user_pictures/kolento-profile_image-b750c756ada12684-300x300.jpeg",
      "status":"C9 Kolento constructed",
      "url":"http://www.twitch.tv/kolento"
    },
    {  
      "display_name":"Nick_28T",
      "game":"FIFA 15",
      "viewers":5200,
      "kappa":9001,
      "logo":"http://static-cdn.jtvnw.net/jtv_user_pictures/nick_28t-profile_image-c16da3de3ebcc6ec-300x300.png",
      "status":"",
      "url":""
    },
    {  
      "display_name":"RocketBeansTV",
      "game":"Gaming Talk Shows",
      "viewers":5089,
      "kappa":9001,
      "logo":"http://static-cdn.jtvnw.net/jtv_user_pictures/rocketbeanstv-profile_image-389ff6f2d5ae9804-300x300.jpeg",
      "status":"Wochenende",
      "url":"http://www.twitch.tv/rocketbeanstv"
    },
    {  
      "display_name":"xlg",
      "game":"League of Legends",
      "viewers":4271,
      "kappa":9001,
      "logo":"http://static-cdn.jtvnw.net/jtv_user_pictures/xlg-profile_image-2857629621043ccb-300x300.jpeg",
      "status":"Pré-Temporada Xtreme League: Circuito Desafiante - Tabela: http://uol.com/byd8x9",
      "url":"http://www.twitch.tv/xlg"
    },
    {  
      "display_name":"MANvsGAME",
      "game":"Dark Souls",
      "viewers":3809,
      "kappa":9001,
      "logo":"http://static-cdn.jtvnw.net/jtv_user_pictures/manvsgame-profile_image-b90a89f4bc2beeda-300x300.png",
      "status":"MAN vs DARK SOULS (PC) Quest for 100% v3.0!",
      "url":"http://www.twitch.tv/manvsgame"
    }
  ]
}
// }}}

React.render(
  <ChannelTable streams={STREAMS.Streams} />,
  document.getElementById('content')
);
