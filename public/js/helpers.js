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

