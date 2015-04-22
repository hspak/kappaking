function convertMinutes(minutes, ago) {
  var ret = "";
  if (minutes == 0) {
    ret = "now";
  } else if (minutes < 60) {
    ret = minutes + "m";
  } else if (minutes < 1440) { // 60*24
    ret = Math.round(minutes/60) + "h " + (minutes%60) + "m";
  } else if (minutes < 10080) { // 60*24*7
    ret = Math.round(minutes/1440) + "d";
  } else {
    ret = Math.round(minutes/10080) + "w";
  }

  if (ago && minutes >= 60) {
    ret += " ago";
  }
  return ret;
}
