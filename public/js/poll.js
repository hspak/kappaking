console.log("poll loaded");
setInterval(function() {
  console.log(document.URL);
  xmlHttp = new XMLHttpRequest();
  xmlHttp.open("GET", document.URL + "api/get/data", false);
  xmlHttp.send(null);
  // TODO: handle error
  Data = JSON.parse(xmlHttp.responseText);
}, 2000);
