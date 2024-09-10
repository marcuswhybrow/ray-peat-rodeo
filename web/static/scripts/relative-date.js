function relativeDate(origStr) {
  let buildDate = Date.parse(origStr);
  let delta = Date.now() - buildDate;

  let hours = Math.floor(delta / (60 * 60 * 1000));
  if (hours < 24) {
    return 'today';
  }

  let days = Math.floor(delta / (24 * 60 * 60 * 1000));
  if (days == 1) {
    return 'today';
  } else if (days == 2) {
    return 'yesterday';
  } else if (days < 7) {
    return days + ' days ago';
  }

  let weeks = Math.floor(delta / (7 * 24 * 60 * 60 * 1000));
  if (weeks == 1) {
    return 'a week ago';
  } else if (weeks <= 16) {
    return weeks + ' weeks ago';
  }

  return origStr
}
